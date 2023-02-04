package main

import (
	ms "github.com/maldan/go-ml/server"
	"github.com/maldan/go-ml/server/core/handler"
)

func main() {
	ms.Start(ms.Config{
		Host: "127.0.0.1:16000",
		Router: []ms.RouteHandler{
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
		Panel: ms.PanelConfig{},
	})
}
