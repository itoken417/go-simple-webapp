package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itoken417/go-simple-webapp/internal/middleware"
	"github.com/itoken417/go-simple-webapp/internal/session"
)

func withSession(r *http.Request, loggedIn bool) *http.Request {
	sess := session.Load(r)
	if loggedIn {
		sess.Set("member_id", "1")
	}
	return r.WithContext(sess.Inject(r.Context()))
}

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func TestRequireLogin(t *testing.T) {
	h := middleware.RequireLogin([]string{"/"})(okHandler)

	tests := []struct {
		name       string
		path       string
		loggedIn   bool
		wantCode   int
		wantLocate string
	}{
		{"未ログイン・保護パス→リダイレクト", "/", false, http.StatusSeeOther, "/login"},
		{"ログイン済み・保護パス→通過", "/", true, http.StatusOK, ""},
		{"未ログイン・保護外パス→通過", "/login", false, http.StatusOK, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tc.path, nil)
			r = withSession(r, tc.loggedIn)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			if w.Code != tc.wantCode {
				t.Errorf("status: got %d, want %d", w.Code, tc.wantCode)
			}
			if tc.wantLocate != "" {
				if loc := w.Header().Get("Location"); loc != tc.wantLocate {
					t.Errorf("Location: got %q, want %q", loc, tc.wantLocate)
				}
			}
		})
	}
}
