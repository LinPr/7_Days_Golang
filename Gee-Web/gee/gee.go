package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

// Engine 作为调度全局的对象
type Engine struct {
	router       *router
	*RouterGroup                // 继承RouterGroup的方法
	groups       []*RouterGroup // stroe all groups
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine interface
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

/*
	The NewEngine function

1. creates a new instance of the Engine struct, initializes the router,
2. and sets up the RouterGroup as the root group of the engine.
3. It then adds the root group to the list of groups managed by the engine.
*/
func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = append(engine.groups, engine.RouterGroup)

	return engine
}

/*
	The NewGroup method of the RouterGroup struct

1. creates a new instance of the RouterGroup struct, which inherits the middleware stack and prefix of its parent group.
2. The new group is added to the list of groups managed by the engine.
*/
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 做一个继承
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, component string, handler HandlerFunc) {
	pattern := group.prefix + component
	log.Printf("Route %4s-%s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
