package gsrv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceOK(t *testing.T) {
	r := strings.NewReader(`GET /dummy.txt HTTP/1.0`)
	buf := bytes.Buffer{}
	Service(r, &buf, ".")
	assert.Contains(t, buf.String(), "HTTP/1.0 200 OK")
	assert.Contains(t, buf.String(), "this is dummy text")
}

func TestServiceNotFound(t *testing.T) {
	r := strings.NewReader(`GET /this_is_not_exists HTTP/1.0`)
	buf := bytes.Buffer{}
	Service(r, &buf, ".")
	assert.Contains(t, buf.String(), "HTTP/1.0 404 Not Found")
}
