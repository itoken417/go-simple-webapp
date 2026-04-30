package session

import (
	"context"
	"net/http"
)

const cookieName = "session_id"

type contextKey struct{}

// Session はリクエストに紐づくセッションを表す。
type Session struct {
	id    string
	entry *entry
}

// FromContext はContextからSessionを取り出す。
// ミドルウェアを通っていない場合はpanicする。
func FromContext(ctx context.Context) *Session {
	s, ok := ctx.Value(contextKey{}).(*Session)
	if !ok || s == nil {
		panic("session: middleware not applied")
	}
	return s
}

// Load はリクエストからセッションを復元、または新規作成する。
func Load(r *http.Request) *Session {
	c, err := r.Cookie(cookieName)
	if err == nil {
		if e, ok := globalStore.get(c.Value); ok {
			return &Session{id: c.Value, entry: e}
		}
	}
	id, e := globalStore.create()
	return &Session{id: id, entry: e}
}

// Inject はSessionをContextに埋め込む。
func (s *Session) Inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey{}, s)
}

// SetCookie はセッションCookieをレスポンスにセットする。
func (s *Session) SetCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    s.id,
		Path:     "/",
		MaxAge:   int(idleTimeout.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Get はセッションからキーに対応する値を返す。
func (s *Session) Get(key string) (any, bool) {
	v, ok := s.entry.data[key]
	return v, ok
}

// Set はセッションにキーと値を保存する。
func (s *Session) Set(key string, val any) {
	s.entry.data[key] = val
}

// Delete はセッションから指定キーを削除する。
func (s *Session) Delete(key string) {
	delete(s.entry.data, key)
}

// Regenerate はセッションIDを再生成する（ログイン後に必ず呼ぶ）。
func (s *Session) Regenerate(w http.ResponseWriter) {
	globalStore.delete(s.id)

	newID, newEntry := globalStore.create()
	for k, v := range s.entry.data {
		newEntry.data[k] = v
	}
	s.id = newID
	s.entry = newEntry

	s.SetCookie(w)
}

// Destroy はセッションを破棄し、Cookieを削除する（ログアウト時に呼ぶ）。
func (s *Session) Destroy(w http.ResponseWriter) {
	globalStore.delete(s.id)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
