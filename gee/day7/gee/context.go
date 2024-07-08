package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	middleware []HandlerFunc //当前路由需执行的中间件
	index      int           //当前执行中间件的进度
}

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  writer,
		Request: req,
		Path:    req.URL.Path,
		Method:  req.Method,
		index:   -1,
	}
}

// 如果前端是form表单提交，通过key 拿到form表单提交的value
func (c *Context) GetFormByKey(key string) string {
	return c.Request.FormValue(key)
}

// 按URL根据key查询 param 如localhost:8080/user?name=tom 这里key就是name 返回就是tom
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) setHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) JSON(code int, result any) {
	c.setHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(result); err != nil {
		c.error(err)
	}
}

func (c *Context) String(code int, format string, values ...any) {
	c.setHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) error(err error) {
	http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
}

func (c *Context) errorCode(code int, err error) {
	http.Error(c.Writer, err.Error(), code)
}
func (c *Context) Fail(code int, err string) {
	c.index = len(c.middleware)
	c.JSON(code, map[string]any{"message": err})
}

// 通过next来控制转移 串行的执行中间件

//这块代码的逻辑挺深的，详细解释参见../Next.md解释
func (c *Context) Next() {
	c.index++
	s := len(c.middleware)
	for ; c.index < s; c.index++ {
		c.middleware[c.index](c)
	}
}
