package httprouter

import (
	"net/http"
	"strings"
)

type handler struct {
	endpoint *Endpoint
}

func newHandler(endpoint *Endpoint) *handler {
	return &handler{
		endpoint: endpoint,
	}
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// Serve endpoint if not on root one
	if h.endpoint != h.endpoint.root {
		h.endpoint.serve(rw, r)
		return
	}

	path := r.URL.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	for _, e := range h.endpoint.root.endpoints {
		ctx, ok := e.pattern.match(r.Context(), path)
		if ok {
			e.chain.ServeHTTP(rw, r.WithContext(ctx))
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
}
