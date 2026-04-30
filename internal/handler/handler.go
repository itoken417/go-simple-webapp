package handler

import (
	"fmt"
	"net/http"
)

type Hdl struct {
	W http.ResponseWriter
	R *http.Request
}

func Static(path string) {
	dir := "./web" + path
	http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

func (h Hdl) Hello() {
	fmt.Fprintf(h.W, "Hello World")
}

func (h Hdl) Test() {
	fmt.Fprintf(h.W, "test")
}

func (h Hdl) Panic() {
	panic("テスト用の意図的なパニックです")
}
