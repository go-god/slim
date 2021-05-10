package slim

import (
	"net/http"
	"path"
	"regexp"
)

var (
	// reg match english letters for http method name
	regEnLetter = regexp.MustCompile("^[A-Z]+$")
)

// IRouter defines all router handle interface includes single and group router.
type IRouter interface {
	IRoutes
	Group(prefix string, h ...HandlerFunc) *RouterGroup
}

// IRoutes defines all router handle interface
type IRoutes interface {
	Use(middlewares ...HandlerFunc)

	Handle(httpMethod, relativePath string, handler HandlerFunc)
	Any(string, HandlerFunc)
	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	DELETE(string, HandlerFunc)
	PATCH(string, HandlerFunc)
	PUT(string, HandlerFunc)
	HEAD(string, HandlerFunc)
	CONNECT(pattern string, handler HandlerFunc)
	OPTIONS(pattern string, handler HandlerFunc)
	TRACE(pattern string, handler HandlerFunc)

	Static(relativePath string, root string)
}

// 判断是否实现了IRouter接口
var _ IRouter = &RouterGroup{}

// RouterGroup router group
type RouterGroup struct {
	prefix   string
	handlers HandlersChain // support middleware
	engine   *Engine       // all groups share a Engine instance
	parent   *RouterGroup  // support nesting
}

// Group creates a new router group. You should add all the routes
// that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:   group.calculateAbsolutePath(relativePath),
		engine:   engine,
		parent:   group,
		handlers: handlers,
	}

	engine.groups = append(engine.groups, newGroup)

	debugPrintf("engine groups len: %d", len(engine.groups))

	return newGroup
}

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.handlers = append(group.handlers, middlewares...)
}

// Handle registers a new request handle and middleware with the given path and method.
func (group *RouterGroup) Handle(httpMethod, relativePath string, handler HandlerFunc) {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}

	group.handle(httpMethod, relativePath, handler)
}

func (group *RouterGroup) handle(method string, relativePath string, handler HandlerFunc) {
	absolutePath := group.calculateAbsolutePath(relativePath)
	debugPrintf("Route %4s - %s", method, absolutePath)
	group.engine.router.addRoute(method, absolutePath, handler)
}

// anyMethods any method
var anyMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

// Any add any method router
func (group *RouterGroup) Any(pattern string, handler HandlerFunc) {
	for _, method := range anyMethods {
		group.handle(method, pattern, handler)
	}
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.handle(http.MethodGet, pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.handle("POST", pattern, handler)
}

// HEAD defines the method to add HEAD request
func (group *RouterGroup) HEAD(pattern string, handler HandlerFunc) {
	group.handle(http.MethodHead, pattern, handler)
}

// PUT defines the method to add PUT request
func (group *RouterGroup) PUT(pattern string, handler HandlerFunc) {
	group.handle(http.MethodPut, pattern, handler)
}

// PATCH defines the method to add PATCH request
func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.handle(http.MethodPatch, pattern, handler)
}

// DELETE defines the method to add DELETE request
func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.handle(http.MethodDelete, pattern, handler)
}

// CONNECT defines the method to add CONNECT request
func (group *RouterGroup) CONNECT(pattern string, handler HandlerFunc) {
	group.handle(http.MethodConnect, pattern, handler)
}

// OPTIONS defines the method to add OPTIONS request
func (group *RouterGroup) OPTIONS(pattern string, handler HandlerFunc) {
	group.handle(http.MethodOptions, pattern, handler)
}

// TRACE defines the method to add OPTIONS request
func (group *RouterGroup) TRACE(pattern string, handler HandlerFunc) {
	group.handle(http.MethodTrace, pattern, handler)
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.handlers) + len(handlers)
	if finalSize >= abortIndex {
		panic("too many handlers")
	}

	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.handlers)
	copy(mergedHandlers[len(group.handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.prefix, relativePath)
}

// createStaticHandler create static handler
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// Static serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

// ReturnObj 返回接口IRoutes
func (group *RouterGroup) ReturnObj() IRoutes {
	if group.parent != nil {
		return group.engine
	}

	return group
}
