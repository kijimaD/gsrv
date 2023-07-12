package gsrv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	r := strings.NewReader(`GET /README.md HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)
	buf := bytes.Buffer{}
	Service(r, &buf, ".")
	assert.Contains(t, buf.String(), "HTTP/1.0 200 OK")
	assert.Contains(t, buf.String(), "# gsrv")
}
