package middleware

import "net/http"

// Chain はミドルウェアを外側から順に適用する。
// 例: Chain(h, A, B) → A(B(h))
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
