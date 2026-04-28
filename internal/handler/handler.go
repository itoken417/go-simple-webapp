package handler

import (
	"fmt"
	"net/http"
)

type Hdl struct {
	W http.ResponseWriter
	R *http.Request
}

func StaticHandle(path string) {
	dir := "./web" + path
	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

func (h Hdl) HelloHandle() {
	fmt.Fprintf(h.W, "Hello World")
}

func (h Hdl) TestHandle() {
	fmt.Fprintf(h.W, "test")
}
