package httprouter

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextParam(t *testing.T) {
	r, err := http.NewRequest("", "", nil)
	assert.Nil(t, err)

	r = r.WithContext(context.WithValue(context.TODO(), contextKey{"key"}, "value"))
	assert.Equal(t, "value", ContextParam(r, "key"))
}

func TestQueryParam(t *testing.T) {
	r, err := http.NewRequest("", "?key=value", nil)
	assert.Nil(t, err)
	assert.Equal(t, "value", QueryParam(r, "key"))
}
