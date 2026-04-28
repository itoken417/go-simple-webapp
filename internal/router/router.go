package router

import (
	"net/http"

	"github.com/itoken417/go-simple-webapp/internal/handler"
	"github.com/itoken417/goutils/logger"
)

var route = make(map[string]map[string]func(*handler.Hdl), 0)

func init() {
	route["ANY"] = make(map[string]func(*handler.Hdl), 0)
	route["GET"] = make(map[string]func(*handler.Hdl), 0)
	route["POST"] = make(map[string]func(*handler.Hdl), 0)
}

func Add(p string, h func(*handler.Hdl)) {
	route["ANY"][p] = h
}

func AddGet(p string, h func(*handler.Hdl)) {
	Set("GET", p, h)
}
func AddPost(p string, h func(*handler.Hdl)) {
	Set("POST", p, h)
}

func Set(m, p string, h func(*handler.Hdl)) {
	if m != "" {
		route[m][p] = h
	}
}

func Router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var method string = r.Method
	logger.Log(method, path)
	handle := new(handler.Hdl)
	handle.W = w
	handle.R = r
	if rt, ok := route[method][path]; ok {
		rt(handle)
		return
	}
	if rt, ok := route["ANY"][path]; ok {
		rt(handle)
		return
	}
	logger.Log("404 not found:", path)
	http.NotFound(w, r)
}
