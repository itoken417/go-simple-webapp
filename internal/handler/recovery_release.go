//go:build release

package handler

import (
	"fmt"
	"net/http"
	"net/smtp"

	config "github.com/itoken417/go-simple-webapp/configs"
	"github.com/itoken417/goutils/logger"
)

func handleError(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte) {
	logger.Log("panic recovered:", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	go notifyError(r, err, stack)
}

func notifyError(r *http.Request, err interface{}, stack []byte) {
	cfg := config.Get()
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)
	subject := fmt.Sprintf("Subject: [ERROR] %s %s\r\n", r.Method, r.URL.Path)
	body := fmt.Sprintf("Error: %v\r\nPath: %s\r\n\r\nStack Trace:\r\n%s", err, r.URL.Path, string(stack))
	msg := []byte(subject + "\r\n" + body)
	addr := cfg.SMTPHost + ":" + cfg.SMTPPort
	smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{cfg.SMTPTo}, msg)
}
