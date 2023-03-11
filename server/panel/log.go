package ms_panel

import (
	ms_error "github.com/maldan/go-ml/server/error"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"os"
	"strings"
)

type Log struct {
	Path string
}

func (l Log) GetList() []string {
	f, err := os.OpenFile(l.Path, os.O_RDONLY, 0777)
	ms_error.FatalIfError(err)

	b := make([]byte, 1024*16)
	f.Read(b)
	ss := string(b)
	lines := strings.Split(ss, "\n")

	return ml_slice.Reverse(lines)
}
