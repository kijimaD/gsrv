package gsrv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequest(t *testing.T) {
	r := strings.NewReader(`GET /dummy.txt HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)
	req := readRequest(r)
	assert.Equal(t, "GET", req.method)
	assert.Equal(t, "/dummy.txt", req.path)
	assert.Equal(t, 0, req.protoMinorVersion)

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
	r := strings.NewReader(`GET /dummy.txt HTTP/1.0
Connection: Close
Content-Type: text/plain
Content-Length: 100
`)

	req := HTTPRequest{}
	readRequestLine(r, &req)
	assert.Equal(t, "GET", req.method)
	assert.Equal(t, "/dummy.txt", req.path)
	assert.Equal(t, 0, req.protoMinorVersion)
}

func TestReadHeaderField(t *testing.T) {
	r := strings.NewReader(`GET /dummy.txt HTTP/1.0
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

func TestReadHeaderFieldReal(t *testing.T) {
	r := strings.NewReader(`Host: localhost:7777
Connection: keep-alive
Cache-Control: max-age=0
sec-ch-ua: "Not/A)Brand";v="99", "Google Chrome";v="115", "Chromium";v="115"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "Linux"
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Sec-Fetch-Site: none
Sec-Fetch-Mode: navigate
Sec-Fetch-User: ?1
Sec-Fetch-Dest: document
Accept-Encoding: gzip, deflate, br
Accept-Language: en-US,en;q=0.9
`)
	req := HTTPRequest{
		header: HTTPHeaderFields{},
	}
	readHeaderField(r, &req)
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
	result := getFileInfo(".", "dummy.txt")
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
