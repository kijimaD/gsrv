package gsrv

import (
	"io"
)

func service(in *io.Reader, out *io.Writer, docroot string) {
	// req := readRequest(in)
	// respondTo(req, out, docroot)

	// fmt.Println(out)
}

type HTTPHeaderFields []HTTPHeaderField

type HTTPHeaderField struct {
	name  string
	value string
}

type HTTPRequest struct {
	protoMinorVersion *int
	method            *string
	path              *string
	header            *HTTPHeaderFields
	body              *string
	length            int
}
