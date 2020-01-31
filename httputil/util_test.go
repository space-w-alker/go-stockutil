package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ghetzel/testify/require"
)

func mt(ct string) *http.Request {
	req := httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, ct)
	return req
}

func TestMediaType(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, MediaType(mt(``)))
	assert.Equal(`text/plain`, MediaType(mt(`text/plain`)))
	assert.Equal(`text/plain`, MediaType(mt(`text/plain; charset=utf-8`)))
}

func TestIsMediaType(t *testing.T) {
	assert := require.New(t)

	req := mt(`text/plain; charset=utf-8`)

	assert.True(IsMediaType(req, `text/plain`))
	assert.True(IsMediaType(req, `text/plain`, `text/html`))
	assert.True(IsMediaType(req, `text/`))
	assert.False(IsMediaType(req))
	assert.False(IsMediaType(req, `text/html`))
}

func ExampleIsMediaType_o1() {
	req := httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/html`))
	// Output: true
}

func ExampleIsMediaType_o2() {
	req := httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/html`, `text/plain`))
	// Output: true
}

func ExampleIsMediaType_o3() {
	req := httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `text/`))
	// Output: true
}

func ExampleIsMediaType_o4() {
	req := httptest.NewRequest(`GET`, `/`, nil)
	req.Header.Set(`Content-Type`, `text/html; charset=utf-8`)

	fmt.Println(IsMediaType(req, `video/`))
	// Output: false
}
