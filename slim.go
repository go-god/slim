// Package slim for go restful api framework
package slim

import (
	"html/template"
	"net/http"
	"strings"
)

// HandlerFunc defines the handler used by slim middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc array
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}

	return nil
}

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router        *router
	noRoute       HandlersChain      // router not found chain
	groups        []*RouterGroup     // store all groups
	htmlTemplates *template.Template // for html render
	funcMap       template.FuncMap   // for html render func map
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.noRoute = []HandlerFunc{NoRoute()}
	return engine
}

// Default new engine
// use Logger() & Recovery handlers
func Default() *Engine {
	engine := New()
	engine.Use(AccessLog(), Recovery())
	return engine
}

// for custom render function
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载templates文件
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string, opts ...HTTPServerOption) (err error) {
	server := InitHTTPServer(opts...)
	server.Handler = engine
	server.Addr = addr

	debugPrintf("server run address: %s", addr)

	return server.ListenAndServe()
}

// ServeHTTP implement http ServeHTTP
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 合并所有的中间件
	var middlewares HandlersChain
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.handlers...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
}
