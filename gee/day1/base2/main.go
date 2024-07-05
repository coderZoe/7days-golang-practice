package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

// 实现一个Handler 相当于自己实现了ServeMux
func (e *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(writer, "URL.Path=%q\n", req.URL.Path)
	case "/hello":
		for key, value := range req.Header {
			fmt.Fprintf(writer, "Header[%q] = %q\n", key, value)
		}
	default:
		fmt.Fprintf(writer, "404 NOT FOUND %s\n", req.URL)
	}
}

func main() {
	http.ListenAndServe(":9999", &Engine{})
}
