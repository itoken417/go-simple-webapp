package main

import (
	"net/http"

	config "github.com/itoken417/go-simple-webapp/configs"
	"github.com/itoken417/go-simple-webapp/internal/handler"
	"github.com/itoken417/go-simple-webapp/internal/router"
	"github.com/itoken417/goutils/logger"
	"github.com/itoken417/goutils/mailsender"
	"github.com/joho/godotenv"
)

func init() {
	router.Add("/", (*handler.Hdl).HelloHandle)
	router.AddPost("/test", (*handler.Hdl).TestHandle)
	router.AddGet("/template", (*handler.Hdl).TemplateHandle)
	router.AddPost("/mail", (*handler.Hdl).MailHandle)
}

const addr = ":8080"

func main() {
	godotenv.Load()

	cfg := config.Get()
	handler.InitMailSender(mailsender.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	}, cfg.SMTPTo)

	logger.Init(true, false)
	defer func() {
		logger.Log("server stopped")
		logger.Destory()
	}()

	handler.StaticHandle("/static/")
	http.Handle("/", handler.RecoveryMiddleware(http.HandlerFunc(router.Router)))

	logger.Log("server starting on " + addr)
	err := http.ListenAndServe(addr, nil)
	logger.ErrCheck(err)
}
