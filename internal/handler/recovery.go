package handler

import (
	"net/http"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				handleError(w, r, err, stack)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
