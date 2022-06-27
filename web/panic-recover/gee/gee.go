package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// Handler function type
type HandlerFunc func(*Context)

// Define a router group for group control
type RouterGroup struct {
	prefix      string
	parent      *RouterGroup  //support group nesting
	middlewares []HandlerFunc // suppoort middleware
	engine      *Engine       // all groups in an engine can refer to the same engine instance
}

// Define the gee engine.
type Engine struct {
	router *router
	*RouterGroup
	groups        []*RouterGroup     // store all router groups
	htmlTemplates *template.Template // for html render: a struct for HTML template
	funcMap       template.FuncMap   // for html render: template rendering functions
}

// Constructor of Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}  // engine itself is a group
	engine.groups = []*RouterGroup{engine.RouterGroup} // so engine itself should be added to groups
	return engine
}

// set function map for the engine.
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// set html template for the engine.
// TODO: support group-level template? First, separate router group from gee...
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// Constructor of RouterGroup
// Create a new group under a group or engine.
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		engine: engine,
		prefix: group.prefix + prefix,
		parent: group,
	}
	engine.groups = append(engine.groups, newGroup) // add new group to engine
	return newGroup
}

// Add router for this group.
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Add a route mapping
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// Add a new GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// Add a new POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Start http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// apply middleware(s) to a group
// middleware will be called in the sequence of input arguments
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// ServeHTTP method implementation, to actually run service
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	middlewares := make([]HandlerFunc, 0)
	// TODO: find a better way to search groups
	for _, group := range engine.groups {
		// TODO: support wildcard or other group matching
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// create static handler
// relativePath is the relative path from root URL
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Params["filepath"]
		// Check if file exist or have access to it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
// relativePath is the relative path from root URL
// root is the absolute path to server local file
// TODO: separate static handler to other files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handler
	group.GET(urlPattern, handler)
}
