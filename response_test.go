package gsrv

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondToGet(t *testing.T) {
	r := strings.NewReader(`GET /dummy.txt HTTP/1.0`)
	req := readRequest(r)
	out := bytes.Buffer{}

	respondTo(req, &out, ".")
	assert.Contains(t, out.String(), "HTTP/1.0 200 OK")
	assert.Contains(t, out.String(), "Date: ")
	assert.Contains(t, out.String(), "Server: gsrv/1.0.0")
	assert.Contains(t, out.String(), "Connection: close")
	assert.Contains(t, out.String(), "Content-Length: 19")
	assert.Contains(t, out.String(), "Content-Type: text/plain")
	assert.Contains(t, out.String(), "this is dummy text")
}

func TestRespondToNotFound(t *testing.T) {
	r := strings.NewReader(`GET /this_is_not_exist HTTP/1.0`)
	req := readRequest(r)
	out := bytes.Buffer{}

	respondTo(req, &out, ".")
	assert.Contains(t, out.String(), "HTTP/1.0 404 Not Found")
	assert.Contains(t, out.String(), "Date: ")
	assert.Contains(t, out.String(), "Server: gsrv/1.0.0")
	assert.Contains(t, out.String(), "Connection: close")
	assert.NotContains(t, out.String(), "Content-Length: 7")
	assert.NotContains(t, out.String(), "Content-Type: text/plain")
	assert.NotContains(t, out.String(), "# gsrv")
}

func TestRespondToHEAD(t *testing.T) {
	r := strings.NewReader(`HEAD /dummy.txt HTTP/1.0`)
	req := readRequest(r)
	out := bytes.Buffer{}

	respondTo(req, &out, ".")
	assert.Contains(t, out.String(), "HTTP/1.0 200 OK")
}

func TestRespondToPOST(t *testing.T) {
	r := strings.NewReader(`POST /dummy.txt HTTP/1.0`)
	req := readRequest(r)
	out := bytes.Buffer{}

	respondTo(req, &out, ".")
	assert.Contains(t, out.String(), "<title>405 Method Not Allowed</title>")
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
