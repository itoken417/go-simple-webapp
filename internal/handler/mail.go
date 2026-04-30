package handler

import (
	"html/template"
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/middleware"
	"github.com/itoken417/goutils/logger"
	"github.com/itoken417/goutils/mailsender"
)

var mailSender *mailsender.Sender
var defaultMailTo string

type mailFormData struct {
	CSRFToken string
	Subject   string
	Body      string
	Result    string
	Error     string
}

func InitMailSender(cfg mailsender.Config, defaultTo string) {
	mailSender = mailsender.New(cfg)
	defaultMailTo = defaultTo
	middleware.InitErrorNotifier(mailSender)
}

func (h Hdl) MailForm() {
	renderMailForm(h.W, h.R, "", "", "", "")
}

func (h Hdl) Mail() {
	subject := h.R.FormValue("subject")
	body := h.R.FormValue("body")

	if subject == "" || body == "" {
		renderMailForm(h.W, h.R, subject, body, "", "件名と本文を入力してください")
		return
	}

	if err := mailSender.Send([]string{defaultMailTo}, subject, body); err != nil {
		logger.Log("メール送信エラー:", err)
		renderMailForm(h.W, h.R, subject, body, "", "メール送信に失敗しました")
		return
	}
	renderMailForm(h.W, h.R, "", "", "送信しました", "")
}

func renderMailForm(w http.ResponseWriter, r *http.Request, subject, body, result, errMsg string) {
	t, err := template.ParseFiles("web/tmpl/mail.htm")
	if err != nil {
		logger.ErrCheck(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "mail.htm", mailFormData{
		CSRFToken: middleware.CSRFToken(r),
		Subject:   subject,
		Body:      body,
		Result:    result,
		Error:     errMsg,
	})
}
