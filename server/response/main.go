package ms_response

import "io"

type Custom struct {
	Headers map[string]string
	Body    []byte
	Reader  io.Reader
}

type File struct {
	Headers map[string]string
	Path    string
}
