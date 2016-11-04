package httproute

import (
	"context"
	"net/http"
	"strings"
)

// Router represents an HTTP router instance.
type Router struct {
	endpoints []*Endpoint
}

// NewRouter creates a new HTTP router instance.
func NewRouter() *Router {
	return &Router{
		endpoints: []*Endpoint{},
	}
}

// ServeHTTP satisfies 'http.Handler' interface requirements.
func (rt *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	endpoint, ctx := rt.match(path)
	if endpoint == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	endpoint.handle(ctx, rw, r)
}

func (rt *Router) match(path string) (*Endpoint, context.Context) {
	for _, endpoint := range rt.endpoints {
		if ctx, ok := endpoint.pattern.match(path); ok {
			return endpoint, ctx
		}
	}

	return nil, nil
}
