package handler

import (
	"html/template"
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/member"
	"github.com/itoken417/go-simple-webapp/internal/middleware"
	"github.com/itoken417/go-simple-webapp/internal/secret"
	"github.com/itoken417/go-simple-webapp/internal/session"
	"github.com/itoken417/goutils/logger"
)

type loginData struct {
	CSRFToken string
	Email     string
	Error     string
}

func (h Hdl) LoginForm() {
	sess := session.FromContext(h.R.Context())
	if _, ok := sess.Get("member_id"); ok {
		http.Redirect(h.W, h.R, "/", http.StatusSeeOther)
		return
	}
	renderLoginForm(h.W, h.R, "", "")
}

func (h Hdl) Login() {
	email := h.R.FormValue("email")
	password := h.R.FormValue("password")

	if email == "" || password == "" {
		renderLoginForm(h.W, h.R, email, "メールアドレスとパスワードを入力してください")
		return
	}

	key, err := secret.LoadKey()
	if err != nil {
		logger.Log("シークレットキー読み込み失敗:", err)
		http.Error(h.W, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	members, err := member.Load()
	if err != nil {
		logger.Log("メンバー読み込み失敗:", err)
		http.Error(h.W, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, m := range members {
		if m.Email != email {
			continue
		}
		ok, err := member.VerifyPassword(password, m.Passwd, m.Salt, key)
		if err != nil {
			logger.Log("パスワード検証エラー:", err)
			http.Error(h.W, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if ok {
			sess := session.FromContext(h.R.Context())
			sess.Set("member_id", m.ID)
			sess.Regenerate(h.W)
			http.Redirect(h.W, h.R, "/", http.StatusSeeOther)
			return
		}
		break
	}

	renderLoginForm(h.W, h.R, email, "メールアドレスまたはパスワードが正しくありません")
}

func (h Hdl) Logout() {
	sess := session.FromContext(h.R.Context())
	sess.Destroy(h.W)
	http.Redirect(h.W, h.R, "/login", http.StatusSeeOther)
}

func renderLoginForm(w http.ResponseWriter, r *http.Request, email, errMsg string) {
	t, err := template.ParseFiles("web/tmpl/login.html")
	if err != nil {
		logger.ErrCheck(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "login.html", loginData{
		CSRFToken: middleware.CSRFToken(r),
		Email:     email,
		Error:     errMsg,
	})
}
