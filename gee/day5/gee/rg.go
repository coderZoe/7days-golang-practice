package gee

type RouterGroup struct {
	prefix     string                  //路由组的完整前缀 如/   /admin  /api这种
	router     *Router                 //这个router是全局路由，而非当前rg下的路由
	children   map[string]*RouterGroup //路由组支持创建子路由组 其中key为子路由的完整前缀 value是子路由实例
	middleware []HandlerFunc           //支持middleware 功能
}

// 以当前routerGroup为base衍生出子路由组
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	if group, ok := rg.children[prefix]; ok {
		return group
	}
	subRg := &RouterGroup{
		prefix:   rg.prefix + prefix, //以自己的prefix作为base衍生出子group
		router:   rg.router,          //全局router
		children: map[string]*RouterGroup{},
	}
	rg.children[subRg.prefix] = subRg
	return subRg
}

// 在路由组下添加路由 如当前路由组是/admin，则可以添加路由 GET info 那就不形成了GET /admin/info的路由
func (rg *RouterGroup) AddRouter(method string, subPattern string, handler HandlerFunc) {
	rg.router.addRouter(method, rg.prefix+subPattern, handler)
}

func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.AddRouter("GET", pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.AddRouter("POST", pattern, handler)
}

func (rg *RouterGroup) Use(middleware ...HandlerFunc) {
	rg.middleware = append(rg.middleware, middleware...)
}
