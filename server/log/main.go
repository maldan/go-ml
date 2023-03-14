package ms_log

import (
	"encoding/json"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_time "github.com/maldan/go-ml/util/time"
	"os"
	"runtime"
	"time"
)

type LogBody struct {
	Kind    string           `json:"kind"`
	File    string           `json:"file"`
	Line    int              `json:"line"`
	Body    string           `json:"body"`
	Created ml_time.DateTime `json:"created"`
}

var LogDB *mdb.DataTable[LogBody]

func Log(kind string, message any) {
	_, file, line, _ := runtime.Caller(1)

	b, _ := json.Marshal(message)

	LogDB.Insert(LogBody{
		Kind:    kind,
		File:    file,
		Line:    line,
		Body:    string(b),
		Created: ml_time.Now(),
	})
}

func Init(logFile string) {
	LogDB = mdb.New[LogBody]("./db", "logs", &gosn_driver.Container{})

	r, w, err := os.Pipe()
	ms_error.FatalIfError(err)
	os.Stdout = w
	os.Stderr = w

	go func() {
		for {
			x := make([]byte, 512)
			n, _ := r.Read(x)
			if n > 0 {
				LogDB.Insert(LogBody{
					Kind:    "raw",
					Body:    string(x[0:n]),
					Created: ml_time.Now(),
				})
			}
			time.Sleep(time.Millisecond)
		}
	}()
}
