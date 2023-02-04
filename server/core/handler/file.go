package handler

import (
	ms "github.com/maldan/go-ml/server"
	ms_error "github.com/maldan/go-ml/server/error"
	"net/http"
	"os"
	"strings"
)

type FS struct {
	ContentPath string
}

func (f FS) Handle(args ms.HandlerArgs) {
	// Get current path
	cwd, err := os.Getwd()
	ms_error.FatalIfError(err)

	// Pure path without route // example /data/test -> /test
	routePath := strings.Replace(args.Request.URL.Path, args.Path, "", 1)

	path := strings.ReplaceAll(f.ContentPath, "@", cwd) + routePath
	path = strings.ReplaceAll(path, "\\", "/")

	// rapi_core.DisableCors(args.RW)
	http.ServeFile(args.Response, args.Request, path)
}
