package gee

import (
	"errors"
	"strings"
)

// 改了HandlerFunc入参为Context了
type HandlerFunc func(*Context)

type Router struct {
	roots map[string]*node
}

//如下整个按gin aoi 做了个封装，可以看出gin的Engine其实就是个Handler，然后用net库起了个http server

func newRouter() *Router {
	return &Router{roots: make(map[string]*node)}
}

func (r *Router) addRouter(method string, pattern string, handler HandlerFunc) {
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = NewTrie()
	}
	r.roots[method].insert(pattern, handler)
}

func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.addRouter("GET", pattern, handler)
}

func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.addRouter("POST", pattern, handler)
}

func (r *Router) getRouter(method string, pattern string) *node {
	root := r.roots[method]
	if root == nil {
		return nil
	} else {
		return root.search(pattern)
	}
}

func setParams(n *node, c *Context) {
	params := make(map[string]string)
	parts := parsePattern(n.pattern)
	searchParts := parsePattern(c.Path)
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	c.Params = params
}

func (r *Router) handle(c *Context) {
	if node := r.getRouter(c.Method, c.Path); node != nil && node.handler != nil {
		//在这一步将params解析并设置进来
		setParams(node, c)
		//将本函数的处理也作为middleware一环，加到最末环
		c.middleware = append(c.middleware, node.handler)
		//链式调用
		c.Next()
	} else {
		c.errorCode(404, errors.New("404 NOT FOUND: "+c.Path+"\n"))
	}
}
