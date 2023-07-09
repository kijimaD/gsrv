package gsrv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequestLine(t *testing.T) {
	r := strings.NewReader(`GET /README.md HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)

	req := HTTPRequest{}
	readRequestLine(r, &req)
	assert.Equal(t, "GET", *req.method)
	assert.Equal(t, "/README.md", *req.path)
	assert.Equal(t, 0, *req.protoMinorVersion)
}

func TestReadHeaderField(t *testing.T) {
	r := strings.NewReader(`GET /README.md HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)
	req := HTTPRequest{
		header: &HTTPHeaderFields{},
	}
	readHeaderField(r, &req)

	for i, h := range *req.header {
		if i == 0 {
			assert.Equal(t, "Connection", h.name)
			assert.Equal(t, "Close", h.value)
		}
		if i == 1 {
			assert.Equal(t, "Content-Type", h.name)
			assert.Equal(t, "text/plain", h.value)
		}
		if i == 2 {
			assert.Equal(t, "Content-Length", h.name)
			assert.Equal(t, "100", h.value)
		}
	}
}