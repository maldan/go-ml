package ms

import (
	ms_handler "github.com/maldan/go-ml/server/core/handler"
)

type PanelConfig struct {
	Login    string
	Password string

	HasRequestTab bool
	HasMethodTab  bool
	HasTestTab    bool
	HasDbTab      bool
	HasControlTab bool
}

type Config struct {
	Host   string
	Router []ms_handler.RouteHandler
	Panel  PanelConfig
}
