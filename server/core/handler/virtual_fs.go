package ms_handler

import (
	ml_vfs "github.com/maldan/go-ml/util/io/vfs"
	"net/http"
	"strings"
	"time"
)

type VirtualFS struct {
	Root string
	Fs   ml_vfs.IVFS
}

func (v VirtualFS) Handle(args *Args) {
	// Prepare path
	pathWithoutKey := strings.Replace(args.Path, args.Route, "", 1)

	// Path inside vfs
	pathInsideVfs := strings.Replace(pathWithoutKey, v.Root, "", 1)
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
	//data, err := v.Fs.ReadFile(v.Root + pathInsideVfs)
	//ms_error.FatalIfError(err)

	// Write to temp dir
	/*p2 := os.TempDir() + "/rapi_vfs/" + fmt.Sprintf("%v", os.Getpid()) + "/" + pathInsideVfs
	err = os.MkdirAll(filepath.Dir(p2), 0777)
	ms_error.FatalIfError(err)
	err = ml_file.New(p2).Write(data)
	ms_error.FatalIfError(err)*/

	// Serve file
	// http.ServeFile(args.Response, args.Request, p2)

	http.ServeContent(args.Response, args.Request, "", time.Now(), nil)
}
