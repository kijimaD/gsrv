package gsrv

import (
	"bufio"
	"bytes"
	"io"
	"path"
	"strconv"
	"strings"
	"syscall"
)

func readRequest(in io.Reader) HTTPRequest {
	req := HTTPRequest{}

	// io.Readerは再利用できない
	// TODO: io.Copyを使うのをやめて1つのio.Readerでパースする
	bbuf := new(bytes.Buffer)
	abuf := io.TeeReader(in, bbuf)

	readRequestLine(abuf, &req)
	readHeaderField(bbuf, &req)
	contentLength(&req)

	return req
}

func readRequestLine(in io.Reader, out *HTTPRequest) error {
	scanner := bufio.NewScanner(in)
	scanner.Scan()
	raw := scanner.Text()

	methodIdx := strings.Index(raw, " ")
	method := raw[:methodIdx]
	out.method = method

	raw = raw[methodIdx+1:]
	pathIdx := strings.Index(raw, " ")
	path := raw[:pathIdx]
	out.path = path

	raw = raw[pathIdx+1:]
	minorVersionIdx := strings.Index(raw, ".")
	minorVersion := raw[minorVersionIdx+1:]
	i, err := strconv.Atoi(minorVersion)
	if err != nil {
		return err
	}
	out.protoMinorVersion = i

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

	for {
		if strings.Index(raw, "\n") == -1 {
			break
		}

		keyIdx := strings.Index(raw, ": ")
		key := raw[:keyIdx]

		raw = raw[keyIdx+2:]                 // ": "の分
		valueIdx := strings.Index(raw, "\n") // TODO: HTTPでは \r\n を使うので直す
		value := raw[:valueIdx]

		h := HTTPHeaderField{
			name:  key,
			value: value,
		}
		out.header = append(out.header, h)
		raw = raw[valueIdx+1:]
	}

	return nil
}

// 構造体の中のヘッダーから、長さを取り出す
func contentLength(req *HTTPRequest) error {
	for _, r := range req.header {
		if r.name == "Content-Length" {
			result, err := strconv.Atoi(r.value)
			if err != nil {
				return err
			}
			req.length = result
			break
		}
	}
	return nil
}

func getFileInfo(docroot string, urlpath string) FileInfo {
	info := FileInfo{}
	st := syscall.Stat_t{}

	info.path = buildFSpath(docroot, urlpath)
	info.ok = false
	if err := syscall.Lstat(info.path, &st); err != nil {
		return info // 失敗
	}
	// TODO: ファイルじゃないなら失敗条件を入れる

	info.ok = true
	info.size = int(st.Size)
	return info
}

func buildFSpath(docroot string, urlpath string) string {
	path := path.Join(docroot, urlpath)
	return path
}
