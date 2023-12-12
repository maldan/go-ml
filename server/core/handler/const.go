package ms_handler

import (
	"bytes"
	"net/http"
)

type VirtualResponseWriter struct {
	StatusCode *int
	Buffer     *bytes.Buffer
	Response   http.ResponseWriter
}

func (v VirtualResponseWriter) Header() http.Header {
	return v.Response.Header()
}

func (v VirtualResponseWriter) Write(b []byte) (int, error) {
	v.Buffer.Write(b)
	return v.Response.Write(b)
}

func (v VirtualResponseWriter) WriteHeader(status int) {
	*v.StatusCode = status
	v.Response.WriteHeader(status)
}

func (v VirtualResponseWriter) AddHeader(name string, value string) {
	v.Response.Header().Add(name, value)
}

type Args struct {
	Route    string
	Path     string
	Response VirtualResponseWriter
	Request  *http.Request
	Body     []byte
}

type Handler interface {
	Handle(args *Args)
}

type RouteHandler struct {
	Path    string
	Handler Handler
}

type Context struct {
	AccessToken string
	Response    http.ResponseWriter
	Request     *http.Request
	RemoteIP    string
}
