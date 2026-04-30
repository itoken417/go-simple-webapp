//go:build !release

package middleware

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/itoken417/goutils/logger"
)

func handleError(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte) {
	logger.Log("panic recovered:", err)
	w.WriteHeader(http.StatusInternalServerError)
	t, parseErr := template.ParseFiles("web/tmpl/error_debug.html")
	if parseErr != nil {
		http.Error(w, fmt.Sprintf("Error: %v\n\n%s", err, stack), http.StatusInternalServerError)
		return
	}
	t.Execute(w, map[string]interface{}{
		"Error": fmt.Sprintf("%v", err),
		"Stack": string(stack),
	})
}
