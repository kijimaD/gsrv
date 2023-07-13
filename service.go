package gsrv

import (
	"io"
)

func Service(in io.Reader, out io.Writer, docroot string) {
	req := readRequest(in)
	respondTo(req, out, docroot)
}
