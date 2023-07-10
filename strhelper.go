package gsrv

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

func readRequest(in io.Reader) HTTPRequest {
	req := HTTPRequest{}

	// io.Readerは再利用できない
	// TODO: io.Copyを使うのをやめて1つのio.Readerでパースする
	bbuf := new(bytes.Buffer)
	abuf := io.TeeReader(in, bbuf)

	readRequestLine(abuf, &req)
	readHeaderField(bbuf, &req)
	return req
}

func readRequestLine(in io.Reader, out *HTTPRequest) error {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	raw := scanner.Text()

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

func readHeaderField(in io.Reader, out *HTTPRequest) error {
	// TODO: copyを使うとio.Readerのカウンタがすべて走るので、scannerで1行ずつ処理するように書き直す
	buf := new(strings.Builder)
	_, err := io.Copy(buf, in)
	if err != nil {
		return err
	}
	raw := buf.String()
	headerLineIdx := strings.Index(raw, "\n")
	raw = raw[headerLineIdx+1:]

	headers := out.header
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
		headers = append(headers, h)

		raw = raw[valueIdx+1:]
		if strings.Index(raw, "\n") == -1 {
			break
		}
	}
	out.header = headers

	return nil
}

// 構造体の中のヘッダーから、長さを取り出す
func contentLength(req *HTTPRequest) (error, int) {
	var result int

	for _, r := range req.header {
		if r.name == "Content-Length" {
			result, _ = strconv.Atoi(r.value)
		}
		break
	}

	return nil, result
}
