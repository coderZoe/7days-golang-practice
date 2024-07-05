package gee

import "errors"

// 改了HandlerFunc入参为Context了
type HandlerFunc func(*Context)

type Router struct {
	router map[string]HandlerFunc
}

//如下整个按gin aoi 做了个封装，可以看出gin的Engine其实就是个Handler，然后用net库起了个http server

func newRouter() *Router {
	return &Router{router: make(map[string]HandlerFunc)}
}

func (r *Router) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.router[key] = handler
}

func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.addRouter("GET", pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.addRouter("POST", pattern, handler)
}

func (r *Router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.router[key]; ok {
		handler(c)
	} else {
		c.errorCode(404, errors.New("404 NOT FOUND: "+c.Path+"\n"))
	}
}
