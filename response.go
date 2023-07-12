package gsrv

import (
	"fmt"
	"io"
	"os"
	"time"
)

func respondTo(req HTTPRequest, out io.Writer, docroot string) {
	if req.method == "GET" {
		_ = doFileResponse(req, out, docroot)
	} else if req.method == "HEAD" {
		_ = doFileResponse(req, out, docroot)
	} else if req.method == "POST" {
		methodNotAllowed(req, out)
	} else {
		notImplemented(req, out)
	}
}

func doFileResponse(req HTTPRequest, out io.Writer, docroot string) error {
	info := getFileInfo(docroot, req.path)
	if !info.ok {
		notFound(req, out)
		return nil
	}
	if req.method == "GET" {
		f, err := os.Open(info.path)
		defer f.Close()
		if err != nil {
			return err
		}
		b := make([]byte, 1024)
		n, err := f.Read(b)
		if os.IsNotExist(err) {
			notFound(req, out)
		} else {
			outputCommonHeaderFields(req, out, "200 OK")
			fmt.Fprintf(out, "Content-Length: %d\r\n", info.size)
			fmt.Fprintf(out, "Content-Type: %s\r\n", guessContentType(info))
			fmt.Fprintf(out, "\r\n")
			out.Write(b[:n])
		}
	} else if req.method == "HEAD" {
		outputCommonHeaderFields(req, out, "200 OK")
		fmt.Fprintf(out, "Content-Length: %d\r\n", info.size)
		fmt.Fprintf(out, "Content-Type: %s\r\n", guessContentType(info))
		fmt.Fprintf(out, "\r\n")
	}
	return nil
}

func methodNotAllowed(req HTTPRequest, out io.Writer) {
	outputCommonHeaderFields(req, out, "405 Method Not Allowed")
	fmt.Fprintf(out, "Content-Type: text/html\r\n")
	fmt.Fprintf(out, "\r\n")
	fmt.Fprintf(out, "<html>\r\n")
	fmt.Fprintf(out, "<header>\r\n")
	fmt.Fprintf(out, "<title>405 Method Not Allowed</title>\r\n")
	fmt.Fprintf(out, "<header>\r\n")
	fmt.Fprintf(out, "<body>\r\n")
	fmt.Fprintf(out, "<p>The request method %s is not allowed</p>\r\n", req.method)
	fmt.Fprintf(out, "</body>\r\n")
	fmt.Fprintf(out, "</html>\r\n")

}

func notImplemented(req HTTPRequest, out io.Writer) {
	outputCommonHeaderFields(req, out, "501 Not Implemented")
	fmt.Fprintf(out, "Content-Type: text/html\r\n")
	fmt.Fprintf(out, "\r\n")
	fmt.Fprintf(out, "<html>\r\n")
	fmt.Fprintf(out, "<header>\r\n")
	fmt.Fprintf(out, "<title>501 Not Implemented</title>\r\n")
	fmt.Fprintf(out, "<header>\r\n")
	fmt.Fprintf(out, "<body>\r\n")
	fmt.Fprintf(out, "<p>The request method %s is not implemented</p>\r\n", req.method)
	fmt.Fprintf(out, "</body>\r\n")
	fmt.Fprintf(out, "</html>\r\n")
}

func notFound(req HTTPRequest, out io.Writer) {
	outputCommonHeaderFields(req, out, "404 Not Found")
	fmt.Fprintf(out, "Content-Type: text/html\r\n")
	fmt.Fprintf(out, "\r\n")
	if req.method == "HEAD" {
		fmt.Fprintf(out, "<html>\r\n")
		fmt.Fprintf(out, "<header><title>Not Found</title><header>\r\n")
		fmt.Fprintf(out, "<body><p>File not found</p></body>\r\n")
		fmt.Fprintf(out, "</html>\r\n")
	}
}

var (
	HTTP_MINOR_VERSION = 0
	SERVER_NAME        = "gsrv"
	SERVER_VERSION     = "1.0.0"
)

func outputCommonHeaderFields(req HTTPRequest, out io.Writer, status string) {
	n := time.Now()
	fmt.Fprintf(out, "HTTP/1.%d %s\r\n", HTTP_MINOR_VERSION, status)
	fmt.Fprintf(out, "Date: %s\r\n", n.Format(time.RFC1123))
	fmt.Fprintf(out, "Server: %s/%s\r\n", SERVER_NAME, SERVER_VERSION)
	fmt.Fprintf(out, "Connection: close\r\n")
}

func guessContentType(info FileInfo) string {
	return "text/plain"
}
