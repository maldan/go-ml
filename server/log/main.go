package ms_log

import (
	"encoding/json"
	"fmt"
	gosn_driver "github.com/maldan/go-ml/db/driver/gosn"
	"github.com/maldan/go-ml/db/mdb"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ms_error "github.com/maldan/go-ml/server/error"
	ml_time "github.com/maldan/go-ml/util/time"
	"net/http"
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

type RequestBody struct {
	Id         int              `json:"id"`
	HttpMethod string           `json:"httpMethod"`
	Header     string           `json:"header"`
	Body       string           `json:"body"`
	Response   string           `json:"response"`
	Url        string           `json:"url"`
	StatusCode int              `json:"statusCode"`
	Created    ml_time.DateTime `json:"created"`
}

var LogDB *mdb.DataTable[LogBody]
var RequestDB *mdb.DataTable[RequestBody]

func Log(kind string, message any) {
	_, file, line, _ := runtime.Caller(1)

	b, _ := json.Marshal(message)

	LogDB.Insert(LogBody{
		Kind:    kind,
		File:    file,
		Line:    line,
		Body:    string(b),
		Created: ml_time.Now().UTC(),
	})
}

func LogRequest(request *http.Request, body []byte, response ms_handler.VirtualResponseWriter) {
	// Read header
	header := ""
	for k, v := range request.Header {
		header += fmt.Sprintf("%v: %v\n", k, v)
	}

	// Add to db
	RequestDB.Insert(RequestBody{
		Id:         int(RequestDB.GenerateId()),
		HttpMethod: request.Method,
		Header:     header,
		Body:       string(body),
		Response:   string(response.Buffer.Bytes()),
		StatusCode: *response.StatusCode,
		Url:        request.URL.RequestURI(),
		Created:    ml_time.Now().UTC(),
	})
}

func Init(logFile string) {
	LogDB = mdb.New[LogBody]("./db", "logs", &gosn_driver.Container{})
	RequestDB = mdb.New[RequestBody]("./db", "requests", &gosn_driver.Container{})

	r, w, err := os.Pipe()
	// mw := io.MultiWriter(os.Stdout, w)

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
					Created: ml_time.Now().UTC(),
				})
			}
			time.Sleep(time.Millisecond)
		}
	}()
}
