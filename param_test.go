package httproute

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextParam(t *testing.T) {
	r, _ := http.NewRequest("", "", nil)
	r = r.WithContext(context.WithValue(nil, contextKey{"key"}, "value"))
	assert.Equal(t, "value", ContextParam(r, "key"))
}

func TestQueryParam(t *testing.T) {
	r, _ := http.NewRequest("", "?key=value", nil)
	assert.Equal(t, "value", QueryParam(r, "key"))
}
