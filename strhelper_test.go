package gsrv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequest(t *testing.T) {
	r := strings.NewReader(`GET /README.md HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)
	req := readRequest(r)
	assert.Equal(t, "GET", *req.method)
	assert.Equal(t, "/README.md", *req.path)
	assert.Equal(t, 0, *req.protoMinorVersion)

	for i, h := range req.header {
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

	assert.Equal(t, 100, req.length)
}

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
		header: HTTPHeaderFields{},
	}
	readHeaderField(r, &req)

	for i, h := range req.header {
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

func TestContentLength(t *testing.T) {
	req := HTTPRequest{
		header: HTTPHeaderFields{
			{
				name:  "Connection",
				value: "Close",
			},
			{
				name:  "Content-Length",
				value: "100",
			}},
	}

	contentLength(&req)
	assert.Equal(t, 100, req.length)
}

func TestGetFileInfo(t *testing.T) {
	result := getFileInfo(".", "README.md")
	assert.Equal(t, true, result.ok)
	result = getFileInfo(".", "NotExists")
	assert.Equal(t, false, result.ok)
}

func TestFSpath(t *testing.T) {
	actual := buildFSpath("dir/dir", "index.html")
	assert.Equal(t, "dir/dir/index.html", actual)

	tests := []struct {
		docroot string
		urlpath string
		expect  string
	}{
		{
			docroot: "dir/dir",
			urlpath: "index.html",
			expect:  "dir/dir/index.html",
		},
		{
			docroot: ".",
			urlpath: "index.html",
			expect:  "index.html",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := buildFSpath(tt.docroot, tt.urlpath)
			if got != tt.expect {
				t.Errorf("got %s want %s", got, tt.expect)
			}
		})
	}
}
