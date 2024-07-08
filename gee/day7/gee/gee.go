package gee

import (
	"net/http"
	"strings"
)

//这里有个比较大的核心思想，Engine其实是最大的路由组，也即/路由组
//golang下没有继承，因此使用组合代替继承

type Engine struct {
	*RouterGroup
	router *Router
}

func New() *Engine {
	baseRouter := newRouter()
	return &Engine{
		router: baseRouter,
		RouterGroup: &RouterGroup{
			prefix:   "",
			router:   baseRouter,
			children: map[string]*RouterGroup{},
		},
	}
}

// Engine还作为整体的Handler 不同的是我们抽出了一个中间层router，借助router可以抽出了处理函数HandlerFunc
// 可以不再依赖于http.HandlerFunc，而是用自己封装的context
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := newContext(writer, request)
	//将middleware设置进来 middle的原则是前缀能和路由匹配上
	context.middleware = e.getMiddleware(context.Path)
	e.router.handle(context)
}

// 递归的获取middleware
func (rg *RouterGroup) getMiddleware(path string) []HandlerFunc {
	result := make([]HandlerFunc, 0)
	if rg.prefix == "" || strings.HasPrefix(path, rg.prefix) {
		result = append(result, rg.middleware...)
	}
	for _, child := range rg.children {
		result = append(result, child.getMiddleware(path)...)
	}
	return result
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}
