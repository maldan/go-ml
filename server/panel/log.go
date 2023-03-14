package ms_panel

import (
	"github.com/maldan/go-ml/db/mdb"
	ms_log "github.com/maldan/go-ml/server/log"
)

type Log struct {
	Path string
}

type ArgsOffset struct {
	Page int `json:"page"`
}

func (l Log) GetList(args ArgsOffset) any {
	r := ms_log.LogDB.FindBy(mdb.ArgsFind[ms_log.LogBody]{
		Where: func(t *ms_log.LogBody) bool {
			return true
		},
	})
	return r
	/*// Open file
	f, err := os.OpenFile(l.Path, os.O_RDONLY, 0777)
	ms_error.FatalIfError(err)

	// Get size
	info, err := f.Stat()
	ms_error.FatalIfError(err)

	blockSize := 1024 * 4
	if int(info.Size()) < blockSize {
		blockSize = int(info.Size())
	}

	// Read lines
	b := make([]byte, blockSize)
	f.ReadAt(b, info.Size()-int64((args.Page)*blockSize))
	ss := string(b)
	lines := strings.Split(ss, "\n")

	return map[string]any{
		"lines": ml_slice.Reverse(lines),
		"total": info.Size() / int64(blockSize),
	}*/
}
