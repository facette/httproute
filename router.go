package httproute

import (
	"net/http"
)

// Router represents an HTTP router instance.
type Router struct {
	middlewares []func(http.Handler) http.Handler
	endpoints   []*Endpoint
	handler     *handler
	chain       http.Handler
}

// NewRouter creates a new HTTP router instance.
func NewRouter() *Router {
	router := &Router{}
	router.handler = newHandler(router)
	return router
}

// Endpoint creates a new HTTP router endpoint.
func (r *Router) Endpoint(pattern string) *Endpoint {
	endpoint := newEndpoint(pattern, r)
	r.endpoints = append(r.endpoints, endpoint)
	return endpoint
}

// ServeHTTP satisfies the http.Handler interface.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.chain.ServeHTTP(rw, req)
}

// Use registers a new middleware in the HTTP handlers chain.
func (r *Router) Use(f func(http.Handler) http.Handler) *Router {
	r.chain = r.handler
	r.middlewares = append(r.middlewares, f)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		r.chain = r.middlewares[i](r.chain)
	}
	return r
}
