package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
	csrfFieldName  = "csrf_token"
)

type csrfKey struct{}

// CSRF は Double Submit Cookie パターンで CSRF 対策を行う。
// 安全でないメソッド（POST/PUT/DELETE/PATCH）では Cookie とフォームフィールドまたは
// X-CSRF-Token ヘッダーのトークンが一致しない場合 403 を返す。
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			cookie, err := r.Cookie(csrfCookieName)
			if err != nil {
				token, err = generateCSRFToken()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     csrfCookieName,
					Value:    token,
					Path:     "/",
					HttpOnly: false,
					SameSite: http.SameSiteStrictMode,
				})
			} else {
				token = cookie.Value
			}
		default:
			cookie, err := r.Cookie(csrfCookieName)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			token = r.Header.Get(csrfHeaderName)
			if token == "" {
				token = r.FormValue(csrfFieldName)
			}
			if token == "" || token != cookie.Value {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		ctx := context.WithValue(r.Context(), csrfKey{}, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CSRFToken は context 経由でミドルウェアが設定したトークンを返す。
// テンプレートへの埋め込み時に使う。
func CSRFToken(r *http.Request) string {
	if token, ok := r.Context().Value(csrfKey{}).(string); ok {
		return token
	}
	return ""
}

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
