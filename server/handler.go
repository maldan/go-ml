package ms

import "net/http"

type HandlerArgs struct {
	Route    string
	Path     string
	Response http.ResponseWriter
	Request  *http.Request
}

type Handler interface {
	Handle(args HandlerArgs)
}

type RouteHandler struct {
	Path    string
	Handler Handler
}
