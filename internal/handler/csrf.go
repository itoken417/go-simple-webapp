package handler

import (
	"html/template"
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/middleware"
	"github.com/itoken417/goutils/logger"
)

type csrfSampleData struct {
	CSRFToken string
	Result    string
	Error     string
}

func (h Hdl) CsrfForm() {
	renderCsrf(h.W, h.R, "", "")
}

func (h Hdl) Csrf() {
	message := h.R.FormValue("message")
	if message == "" {
		renderCsrf(h.W, h.R, "", "メッセージを入力してください")
		return
	}
	renderCsrf(h.W, h.R, message, "")
}

func renderCsrf(w http.ResponseWriter, r *http.Request, result, errMsg string) {
	t, err := template.ParseFiles("web/tmpl/csrf.htm")
	if err != nil {
		logger.ErrCheck(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "csrf.htm", csrfSampleData{
		CSRFToken: middleware.CSRFToken(r),
		Result:    result,
		Error:     errMsg,
	})
}
