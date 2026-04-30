package middleware

import (
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/session"
)

// RequireLogin は指定パスへのアクセスにログインを要求するミドルウェア。
// 未ログインは /login にリダイレクトする。
func RequireLogin(paths []string) func(http.Handler) http.Handler {
	protected := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		protected[p] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := protected[r.URL.Path]; ok {
				sess := session.FromContext(r.Context())
				if _, ok := sess.Get("member_id"); !ok {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
