package httprouter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	var (
		rr  *httptest.ResponseRecorder
		req *http.Request
		err error
	)

	r := New().
		Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("X-Test-Root", "root")
				h.ServeHTTP(rw, r)
			})
		})

	r.Endpoint("/").Get(func(rw http.ResponseWriter, r *http.Request) {})

	e := r.Endpoint("/a").
		Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				rw.Header().Add("X-Test-A", "a")
				h.ServeHTTP(rw, SetContextParam(r, "foo", 123))
			})
		}).
		Get(func(rw http.ResponseWriter, r *http.Request) {
			assert.Equal(t, 123, ContextParam(r, "foo"))
		}).
		Options(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusNoContent)
		})

	e.Endpoint("/b").Get(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	})

	r.Endpoint("/c").
		Delete(func(rw http.ResponseWriter, r *http.Request) {}).
		Patch(func(rw http.ResponseWriter, r *http.Request) {}).
		Put(func(rw http.ResponseWriter, r *http.Request) {}).
		Post(func(rw http.ResponseWriter, r *http.Request) {})

	r.Endpoint("/d/").Any(func(rw http.ResponseWriter, r *http.Request) {})

	r.Endpoint("/e").Head(func(rw http.ResponseWriter, r *http.Request) {})

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("OPTIONS", "/a", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "GET, HEAD, OPTIONS", rr.Header().Get("Allow"))
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))
	assert.Equal(t, "a", rr.Header().Get("X-Test-A"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/a", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))
	assert.Equal(t, "a", rr.Header().Get("X-Test-A"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("HEAD", "/a", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))
	assert.Equal(t, "a", rr.Header().Get("X-Test-A"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/a/", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusPermanentRedirect, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))
	assert.Equal(t, "a", rr.Header().Get("X-Test-A"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/a/b", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))
	assert.Equal(t, "a", rr.Header().Get("X-Test-A"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("OPTIONS", "/c", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "DELETE, OPTIONS, PATCH, POST, PUT", rr.Header().Get("Allow"))
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/c", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/c", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("OPTIONS", "/d/", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "DELETE, GET, HEAD, OPTIONS, PATCH, POST, PUT", rr.Header().Get("Allow"))
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("DELETE", "/d/", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("POST", "/d", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusPermanentRedirect, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("OPTIONS", "/e", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "HEAD, OPTIONS", rr.Header().Get("Allow"))
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("HEAD", "/e", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "root", rr.Header().Get("X-Test-Root"))

	rr = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/unknown", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
