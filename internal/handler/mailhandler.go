package handler

import (
	"fmt"
	"net/http"

	"github.com/itoken417/goutils/logger"
	"github.com/itoken417/goutils/mailsender"
)

var mailSender *mailsender.Sender
var defaultMailTo string

// InitMailSender はメール送信クライアントを初期化する。main から呼ぶ。
func InitMailSender(cfg mailsender.Config, defaultTo string) {
	mailSender = mailsender.New(cfg)
	defaultMailTo = defaultTo
}

// MailHandle はフォームの subject / body を受け取り、設定済み宛先にメールを送信する。
func (h Hdl) MailHandle() {
	subject := h.R.FormValue("subject")
	body := h.R.FormValue("body")

	if err := mailSender.Send([]string{defaultMailTo}, subject, body); err != nil {
		logger.Log("メール送信エラー:", err)
		http.Error(h.W, "メール送信に失敗しました", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(h.W, "送信完了")
}
