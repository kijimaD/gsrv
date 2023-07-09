package gsrv

import (
	"io"
	"strconv"
	"strings"
)

func readRequest(in io.Reader) HTTPRequest {
	req := HTTPRequest{}
	readRequestLine(in, &req)
	return req
}

func readRequestLine(in io.Reader, out *HTTPRequest) error {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, in)
	if err != nil {
		return err
	}
	raw := buf.String()

	methodIdx := strings.Index(raw, " ")
	method := raw[:methodIdx]
	out.method = &method

	raw = raw[methodIdx+1:]
	pathIdx := strings.Index(raw, " ")
	path := raw[:pathIdx]
	out.path = &path

	raw = raw[pathIdx+1:]
	minorVersionIdx := strings.Index(raw, ".")
	minorVersion := raw[minorVersionIdx+1:]
	i, err := strconv.Atoi(minorVersion)
	if err != nil {
		return err
	}
	out.protoMinorVersion = &i

	return nil
}
