package gee

import (
	"net/http"
)

type Engine struct {
	router *Router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	e.router.addRouter(method, pattern, handler)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRouter("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRouter("POST", pattern, handler)
}

// Engine还作为整体的Handler 不同的是我们抽出了一个中间层router，借助router可以抽出了处理函数HandlerFunc
// 可以不再依赖于http.HandlerFunc，而是用自己封装的context
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := newContext(writer, request)
	e.router.handle(context)
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
