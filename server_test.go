package gsrv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequestLine(t *testing.T) {
	r := strings.NewReader("GET /README.md HTTP/1.0")

	req := HTTPRequest{}
	readRequestLine(r, &req)
	assert.Equal(t, "GET", *req.method)
	assert.Equal(t, "/README.md", *req.path)
	assert.Equal(t, 0, *req.protoMinorVersion)
}
