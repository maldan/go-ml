package ms_response

type Custom struct {
	Headers map[string]string
	Body    []byte
}

type File struct {
	Path string
}
