package ms_handler

import (
	"errors"
	ms_error "github.com/maldan/go-ml/server/error"
	"net/http"
	"os"
	"strings"
)

type FS struct {
	ContentPath string
	Middleware  func(args *Args, filePath string) string
}

func (f FS) Handle(args *Args) {
	// Get current path
	cwd, err := os.Getwd()
	ms_error.FatalIfError(err)

	// Pure path without route // example /data/test -> /test
	routePath := strings.Replace(args.Request.URL.Path, args.Route, "", 1)

	path := strings.ReplaceAll(f.ContentPath, "@", cwd) + routePath
	path = strings.ReplaceAll(path, "\\", "/")

	if _, err3 := os.Stat(path); errors.Is(err3, os.ErrNotExist) {
		ms_error.Fatal(ms_error.Error{Code: 404, Description: "File not found"})
	}

	// Set middleware
	if f.Middleware != nil {
		path = f.Middleware(args, path)
	}

	http.ServeFile(args.Response, args.Request, path)
}
