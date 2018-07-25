package httproute

import (
	"net/http"
)

// Router represents an HTTP router instance.
type Router struct {
	root *Endpoint
}

// NewRouter creates a new HTTP router instance.
func NewRouter() *Router {
	return &Router{
		root: newEndpoint(""),
	}
}

// Endpoint creates a new HTTP router endpoint.
func (r *Router) Endpoint(pattern string) *Endpoint {
	return r.root.Endpoint(pattern)
}

// ServeHTTP satisfies the http.Handler interface.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.root.chain.ServeHTTP(rw, req)
}

// Use registers a new middleware in the HTTP handlers chain.
func (r *Router) Use(f func(http.Handler) http.Handler) *Router {
	r.root.Use(f)
	return r
}
