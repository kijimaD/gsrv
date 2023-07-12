package gsrv

type HTTPHeaderFields []HTTPHeaderField

type HTTPHeaderField struct {
	name  string
	value string
}

type HTTPRequest struct {
	protoMinorVersion int
	method            string
	path              string
	header            HTTPHeaderFields
	body              string
	length            int
}

type FileInfo struct {
	path string
	size int
	ok   bool
}
