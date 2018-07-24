package httproute

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint(t *testing.T) {
	var e *Endpoint

	r := NewRouter()

	e = r.Endpoint("/")
	assert.Equal(t, []string{"OPTIONS"}, e.Methods())

	e = r.Endpoint("/").Any(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"}, e.Methods())

	e = r.Endpoint("/").Delete(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"DELETE", "OPTIONS"}, e.Methods())

	e = r.Endpoint("/").Get(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"GET", "HEAD", "OPTIONS"}, e.Methods())

	e = r.Endpoint("/").Head(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"HEAD", "OPTIONS"}, e.Methods())

	e = r.Endpoint("/").Patch(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"OPTIONS", "PATCH"}, e.Methods())

	e = r.Endpoint("/").Post(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"OPTIONS", "POST"}, e.Methods())

	e = r.Endpoint("/").Put(func(http.ResponseWriter, *http.Request) {})
	assert.Equal(t, []string{"OPTIONS", "PUT"}, e.Methods())
}
