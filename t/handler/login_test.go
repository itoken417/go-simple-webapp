package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/itoken417/go-simple-webapp/internal/handler"
	"github.com/itoken417/go-simple-webapp/internal/session"
)

func TestMain(m *testing.M) {
	// テンプレートファイル（web/tmpl/）はプロジェクトルート基準のパスで読むため
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..")
	if err := os.Chdir(root); err != nil {
		panic("chdir失敗: " + err.Error())
	}
	os.Exit(m.Run())
}

func withSession(r *http.Request, loggedIn bool) *http.Request {
	sess := session.Load(r)
	if loggedIn {
		sess.Set("member_id", "1")
	}
	return r.WithContext(sess.Inject(r.Context()))
}

func TestLoginForm_LoggedIn_Redirect(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	r = withSession(r, true)
	w := httptest.NewRecorder()

	handler.Hdl{W: w, R: r}.LoginForm()

	if w.Code != http.StatusSeeOther {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusSeeOther)
	}
	if loc := w.Header().Get("Location"); loc != "/" {
		t.Errorf("Location: got %q, want %q", loc, "/")
	}
}

func TestLoginForm_NotLoggedIn_ShowForm(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	r = withSession(r, false)
	w := httptest.NewRecorder()

	handler.Hdl{W: w, R: r}.LoginForm()

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "ログイン") {
		t.Error("レスポンスにログインフォームが含まれていない")
	}
}

func TestLogin_EmptyFields(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(""))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r = withSession(r, false)
	w := httptest.NewRecorder()

	handler.Hdl{W: w, R: r}.Login()

	if w.Code != http.StatusOK {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "メールアドレスとパスワードを入力してください") {
		t.Error("バリデーションエラーが表示されていない")
	}
}

func TestLogout_Redirect(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/logout", nil)
	r = withSession(r, true)
	w := httptest.NewRecorder()

	handler.Hdl{W: w, R: r}.Logout()

	if w.Code != http.StatusSeeOther {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusSeeOther)
	}
	if loc := w.Header().Get("Location"); loc != "/login" {
		t.Errorf("Location: got %q, want %q", loc, "/login")
	}
}
