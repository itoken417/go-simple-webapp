package main

import (
	"net/http"

	config "github.com/itoken417/go-simple-webapp/configs"
	"github.com/itoken417/go-simple-webapp/internal/handler"
	"github.com/itoken417/go-simple-webapp/internal/middleware"
	"github.com/itoken417/go-simple-webapp/internal/router"
	"github.com/itoken417/goutils/logger"
	"github.com/itoken417/goutils/mailsender"
	"github.com/joho/godotenv"
)

func init() {
	router.Add("/", (*handler.Hdl).Hello)
	router.AddPost("/test", (*handler.Hdl).Test)
	router.AddGet("/template", (*handler.Hdl).Template)
	router.AddGet("/mail", (*handler.Hdl).MailForm)
	router.AddPost("/mail", (*handler.Hdl).Mail)
	router.Add("/panic", (*handler.Hdl).Panic)
	router.AddGet("/csrf", (*handler.Hdl).CsrfForm)
	router.AddPost("/csrf", (*handler.Hdl).Csrf)
	router.AddGet("/login", (*handler.Hdl).LoginForm)
	router.AddPost("/login", (*handler.Hdl).Login)
	router.AddPost("/logout", (*handler.Hdl).Logout)
}

const addr = ":8080"

func main() {
	logger.Init(true, false)
	defer func() {
		logger.Log("server stopped")
		logger.Destory()
	}()

	godotenv.Load()

	cfg, err := config.Get()
	if err != nil {
		logger.Log("設定の読み込みに失敗:", err)
		return
	}
	handler.InitMailSender(mailsender.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		User:     cfg.SMTPUser,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
		Encoding: cfg.MailEncoding,
	}, cfg.SystemTo)

	handler.Static("/static/")
	http.Handle("/", middleware.Chain(
		http.HandlerFunc(router.Router),
		middleware.Exception,
		middleware.Session,
		middleware.CSRF,
		middleware.RequireLogin([]string{"/"}),
	))

	logger.Log("server starting on " + addr)
	err = http.ListenAndServe(addr, nil)
	logger.ErrCheck(err)
}
