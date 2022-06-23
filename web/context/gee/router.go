package gee

import (
	"log"
	"net/http"
)

// router interface
type router struct {
	handlers map[string]HandlerFunc
}

// constructor of router
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// Add a route mapping
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// actually run service
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
