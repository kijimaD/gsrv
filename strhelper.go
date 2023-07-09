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
	newLineIdx := strings.Index(raw, "\n")
	minorVersion := raw[minorVersionIdx+1 : newLineIdx]
	i, err := strconv.Atoi(minorVersion)
	if err != nil {
		return err
	}
	out.protoMinorVersion = &i

	return nil
}

func readHeaderField(in io.Reader, out *HTTPRequest) error {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, in)
	if err != nil {
		return err
	}
	raw := buf.String()
	headerLineIdx := strings.Index(raw, "\n")
	raw = raw[headerLineIdx+1:]

	for {
		keyIdx := strings.Index(raw, ":")
		key := raw[:keyIdx]

		raw = raw[keyIdx+2:]                 // ": "の分
		valueIdx := strings.Index(raw, "\n") // TODO: HTTPでは \r\n を使うので直す
		value := raw[:valueIdx]

		h := HTTPHeaderField{
			name:  key,
			value: value,
		}
		*out.header = append(*out.header, h)

		raw = raw[valueIdx+1:]
		if strings.Index(raw, "\n") == -1 {
			break
		}
	}

	return nil
}

// 構造体の中のヘッダーから、長さを取り出す
func contentLength(req *HTTPRequest) (error, int) {
	var result int

	for _, r := range *req.header {
		if r.name == "Content-Length" {
			result, _ = strconv.Atoi(r.value)
		}
		break
	}

	return nil, result
}
