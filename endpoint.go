package httproute

import (
	"net/http"
	"sort"
	"strings"
)

// Endpoint represents an HTTP router endpoint.
type Endpoint struct {
	pattern   *pattern
	handlers  map[string]http.HandlerFunc
	endpoints []*Endpoint
	router    *Router
}

func newEndpoint(pattern string, router *Router) *Endpoint {
	return &Endpoint{
		pattern:  newPattern(pattern),
		handlers: map[string]http.HandlerFunc{},
		router:   router,
	}
}

// Any registers a handler for any method.
func (e *Endpoint) Any(f http.HandlerFunc) *Endpoint {
	return e.register("", f)
}

// Delete registers a DELETE method handler.
func (e *Endpoint) Delete(f http.HandlerFunc) *Endpoint {
	return e.register("DELETE", f)
}

// Get registers a GET method handler.
func (e *Endpoint) Get(f http.HandlerFunc) *Endpoint {
	return e.register("GET", f)
}

// Head registers a HEAD method handler.
func (e *Endpoint) Head(f http.HandlerFunc) *Endpoint {
	return e.register("HEAD", f)
}

// Methods returns the list of methods available from the HTTP router endpoint.
func (e *Endpoint) Methods() []string {
	var hasGet, hasHead bool

	methods := []string{"OPTIONS"}
	for method := range e.handlers {
		methods = append(methods, method)

		if method == "GET" {
			hasGet = true
		} else if method == "HEAD" {
			hasHead = true
		}
	}
	if hasGet && !hasHead {
		methods = append(methods, "HEAD")
	}
	sort.Strings(methods)

	return methods
}

// Patch registers a PATCH method handler.
func (e *Endpoint) Patch(f http.HandlerFunc) *Endpoint {
	return e.register("PATCH", f)
}

// Post registers a POST method handler.
func (e *Endpoint) Post(f http.HandlerFunc) *Endpoint {
	return e.register("POST", f)
}

// Put registers a PUT method handler.
func (e *Endpoint) Put(f http.HandlerFunc) *Endpoint {
	return e.register("PUT", f)
}

func (e *Endpoint) register(method string, f http.HandlerFunc) *Endpoint {
	e.handlers[method] = f
	return e
}

func (e *Endpoint) serve(rw http.ResponseWriter, r *http.Request) {
	// Handle slash redirections
	if !e.pattern.wildcard {
		if e.pattern.slash && !strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path += "/"
			http.Redirect(rw, r, r.URL.String(), http.StatusPermanentRedirect)
			return
		} else if !e.pattern.slash && strings.HasSuffix(r.URL.Path, "/") {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			http.Redirect(rw, r, r.URL.String(), http.StatusPermanentRedirect)
			return
		}
	}

	handler, ok := e.handlers[r.Method]
	if !ok {
		if _, ok = e.handlers[""]; ok {
			// Use "Any" handler
			handler = e.handlers[""]
		} else {
			switch r.Method {
			case "HEAD":
				// Use GET method handler when HEAD is requested
				handler, ok = e.handlers["GET"]

			case "OPTIONS":
				rw.Header().Add("Allow", strings.Join(e.Methods(), ", "))
				rw.WriteHeader(http.StatusNoContent)
				return
			}
		}

		if !ok {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}

	handler(rw, r)
}
