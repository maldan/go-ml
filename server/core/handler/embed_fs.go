package ms_handler

import (
	"embed"
	"fmt"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type EmbedFS struct {
	Root string
	Fs   embed.FS
}

func (e EmbedFS) Handle(args *Args) {
	// Prepare path
	pathWithoutKey := strings.Replace(args.Path, args.Route, "", 1)

	// Path inside vfs
	pathInsideVfs := strings.Replace(pathWithoutKey, e.Root, "", 1)
	if pathInsideVfs == "" {
		pathInsideVfs = "/"
	}
	if len(pathInsideVfs) > 0 && pathInsideVfs[0] != '/' {
		pathInsideVfs = "/" + pathInsideVfs
	}

	if pathInsideVfs == "/" {
		pathInsideVfs = "/index.html"
	}

	// Read file
	data, err := e.Fs.ReadFile(e.Root + pathInsideVfs)
	ms_error.FatalIfError(err)

	// Write to temp dir
	p2 := os.TempDir() + "/rapi_vfs/" + fmt.Sprintf("%v", os.Getpid()) + "/" + pathInsideVfs
	err = os.MkdirAll(filepath.Dir(p2), 0777)
	ms_error.FatalIfError(err)
	err = ml_file.New(p2).Write(data)
	ms_error.FatalIfError(err)

	// Serve file
	http.ServeFile(args.Response, args.Request, p2)
}
