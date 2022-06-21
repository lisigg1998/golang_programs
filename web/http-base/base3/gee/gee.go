package gee

import (
	"fmt"
	"net/http"
)

// Handler function type
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

// Define the interface of ServeHTTP
type Engine struct {
	router map[string]HandlerFunc
}

// Constructor of Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// Add a route mapping
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
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
	method := req.Method
	url := req.URL.Path
	key := method + "-" + url

	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 Not Found: %q\n", url)
	}
}
