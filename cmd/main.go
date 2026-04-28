package main

import (
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/handler"
	"github.com/itoken417/go-simple-webapp/internal/router"
	"github.com/itoken417/goutils/logger"
	"github.com/joho/godotenv"
)

func init() {
	router.Add("/", (*handler.Hdl).HelloHandle)
	router.AddPost("/test", (*handler.Hdl).TestHandle)
	router.AddGet("/template", (*handler.Hdl).TemplateHandle)
}

const addr = ":8080"

func main() {
	godotenv.Load()

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
