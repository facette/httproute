package httproute

import (
	"net/http"
	"strings"
)

type handler struct {
	router *Router
}

func newHandler(router *Router) *handler {
	return &handler{
		router: router,
	}
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}

	for _, e := range h.router.endpoints {
		if ctx, ok := e.pattern.match(r.Context(), path); ok {
			e.serve(rw, r.WithContext(ctx))
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
}
