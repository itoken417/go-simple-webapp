package middleware

import (
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/session"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.Load(r)
		sess.SetCookie(w)
		next.ServeHTTP(w, r.WithContext(sess.Inject(r.Context())))
	})
}
