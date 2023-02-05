package ms_handler

import "net/http"

type Args struct {
	Route    string
	Path     string
	Response http.ResponseWriter
	Request  *http.Request
}

type Handler interface {
	Handle(args Args)
}

type RouteHandler struct {
	Path    string
	Handler Handler
}

type Context struct {
	AccessToken string
	Headers     map[string]string
	Response    http.ResponseWriter
	Request     *http.Request
}
