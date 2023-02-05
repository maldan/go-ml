package main

import (
	ms "github.com/maldan/go-ml/server"
	"github.com/maldan/go-ml/server/core/handler"
)

func main() {
	ms.Start(ms.Config{
		Host: "127.0.0.1:16000",
		Router: []ms_handler.RouteHandler{
			{
				Path: "/api",
				Handler: ms_handler.API{
					ControllerList: []any{User{}, Template{}},
				},
			},
			{
				Path: "/data",
				Handler: ms_handler.FS{
					ContentPath: "example",
				},
			},
			{
				Path:    "/",
				Handler: ms_handler.EmbedFS{},
			},
		},
		Panel: ms.PanelConfig{},
	})
}
