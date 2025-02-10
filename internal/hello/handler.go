package hello

import (
	"fmt"
	"net/http"
)

type HelloHandler struct{}

func NewHelloHandler(router *http.ServeMux) {
	handler := &HelloHandler{}
	router.HandleFunc("/hello", handler.hello())
}

func (handler *HelloHandler) hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!!!"))
		fmt.Println("Hello!!!")
	}
}
