package main

import (
	"github.com/maldan/go-ml/db/mdb"
	ms "github.com/maldan/go-ml/server"
	"github.com/maldan/go-ml/server/core/handler"
)

var DataBase map[string]*mdb.DataTable

type ImageDescription struct {
	Path        string
	Hash        string
	Width       int
	Height      int
	Description string
}

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
					ContentPath: "db",
				},
			},
			/*{
				Path:    "/",
				Handler: ms_handler.EmbedFS{},
			},*/
		},
		Panel: ms.PanelConfig{},
		DataBase: ms.DataBaseConfig{
			Path:     "./db",
			DataBase: &DataBase,
			TableList: []ms.TableConfig{
				{Name: "x", Type: Gasofeal{}},
				{Name: "tags", Type: ImageDescription{}},
			},
		},
	})
}
