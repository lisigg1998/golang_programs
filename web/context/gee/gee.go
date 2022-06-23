package gee

import (
	"net/http"
)

// Handler function type
type HandlerFunc func(*Context)

// Define the interface of ServeHTTP
type Engine struct {
	router *router
}

// Constructor of Engine
func New() *Engine {
	return &Engine{router: newRouter()}
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

// ServeHTTP method implementation, to actually run service
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
