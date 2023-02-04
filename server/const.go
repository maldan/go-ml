package ms

import "net/http"

type DebugConfig struct {
	IsEnabled bool
}

type PanelConfig struct {
	Login    string
	Password string
}

type Context struct {
	AccessToken string
	Headers     map[string]string
	Response    http.ResponseWriter
	Request     *http.Request
}

type Config struct {
	Host string

	Router []RouteHandler

	Debug DebugConfig
	Panel PanelConfig

	EnableJsonWrapper bool
}
