package ms

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

type LogBody struct {
	Kind    string    `json:"kind"`
	File    string    `json:"file"`
	Line    int       `json:"line"`
	Body    any       `json:"body"`
	Created time.Time `json:"created"`
}

func Log(kind string, message any) {
	/*if len(v) > 0 {
		fmt.Printf("[%v] [%v] "+message+"\n", kind, time.Now().Format("2006-01-02 15:04:05.999999"), v)
	} else {
		fmt.Printf("[%v] [%v] "+message+"\n", kind, time.Now().Format("2006-01-02 15:04:05.999999"))
	}*/
	_, file, line, _ := runtime.Caller(1)

	l := LogBody{
		Kind:    kind,
		File:    file,
		Line:    line,
		Body:    message,
		Created: time.Now(),
	}
	b, _ := json.Marshal(l)
	fmt.Println(string(b))
}

func startLog() {
	logfile := `logfile`
	f, _ := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	// mw := io.MultiWriter(os.Stdout, f)
	os.Stdout = f
	os.Stderr = f
}
