package main

import (
	ms "github.com/maldan/go-ml/server"
	ms_config "github.com/maldan/go-ml/server/config"
	"github.com/maldan/go-ml/server/core/handler"
)

func main() {
	ms.Start(ms_config.Config{
		Host: "127.0.0.1:16000",
		Router: []ms_config.RouteHandler{
			{
				Path: "/api",
				Handler: handler.API{
					ControllerList: []any{User{}, Template{}},
				},
			},
			{
				Path: "/data",
				Handler: handler.FS{
					ContentPath: "example",
				},
			},
			{
				Path:    "/",
				Handler: handler.VFS{},
			},
		},
	})
}
