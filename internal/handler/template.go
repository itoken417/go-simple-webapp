package handler

import (
	"html/template"

	"github.com/itoken417/goutils/logger"
)

func (h Hdl) Template() {
	t := template.Must(template.ParseFiles("web/tmpl/template.htm"))
	str := "template test"
	if err := t.ExecuteTemplate(h.W, "template.htm", str); err != nil {
		logger.ErrCheck(err)
	}
}
