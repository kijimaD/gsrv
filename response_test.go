package gsrv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondTo(t *testing.T) {
	r := strings.NewReader(`GET /README.md HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)
	req := readRequest(r)
	out := bytes.Buffer{}

	respondTo(req, &out, ".")
}

func TestDoFileResponse(t *testing.T) {

}

func TestNotImplemented(t *testing.T) {
	r := HTTPRequest{
		protoMinorVersion: 1,
		method:            "GET",
		path:              ".",
		header:            HTTPHeaderFields{},
		body:              "",
		length:            0,
	}
	out := bytes.Buffer{}
	notImplemented(r, &out)
	assert.Contains(t, out.String(), "HTTP/1.0 501 Not Implemented")
	// 日付日時があるので完全一致にできない
	assert.Contains(t, out.String(), "Date: ")
	assert.Contains(t, out.String(), "Server: gsrv/1.0.0")
	assert.Contains(t, out.String(), "Connection: close")
	assert.Contains(t, out.String(), "Content-Type: text/html")
	// 改行が\r\nなのでヒアドキュメント複数行一致検査ができない
	assert.Contains(t, out.String(), "<title>501 Not Implemented</title>")
}

func TestMethodNotAllowed(t *testing.T) {
	r := HTTPRequest{
		protoMinorVersion: 1,
		method:            "GET",
		path:              ".",
		header:            HTTPHeaderFields{},
		body:              "",
		length:            0,
	}
	out := bytes.Buffer{}
	methodNotAllowed(r, &out)
	assert.Contains(t, out.String(), "HTTP/1.0 405 Method Not Allowed")
	// 日付日時があるので完全一致にできない
	assert.Contains(t, out.String(), "Date: ")
	assert.Contains(t, out.String(), "Server: gsrv/1.0.0")
	assert.Contains(t, out.String(), "Connection: close")
	assert.Contains(t, out.String(), "Content-Type: text/html")
	// 改行が\r\nなのでヒアドキュメント複数行一致検査ができない
	assert.Contains(t, out.String(), "<title>405 Method Not Allowed</title>")
}

func TestOutputCommonHeaderFields(t *testing.T) {
	r := HTTPRequest{
		protoMinorVersion: 1,
		method:            "GET",
		path:              ".",
		header:            HTTPHeaderFields{},
		body:              "",
		length:            0,
	}
	out := bytes.Buffer{}
	outputCommonHeaderFields(r, &out, ".")
	assert.Contains(t, out.String(), "HTTP/1.0 .")
	// 日付日時があるので完全一致にできない
	assert.Contains(t, out.String(), "Date: ")
	assert.Contains(t, out.String(), "Server: gsrv/1.0.0")
	assert.Contains(t, out.String(), "Connection: close")
}
