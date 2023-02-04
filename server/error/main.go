package ms_error

import (
	"runtime"
	"strings"
	"time"
)

type Error struct {
	Status      bool   `json:"-"`
	Code        int    `json:"-"`
	Type        string `json:"type"`
	Field       string `json:"field,omitempty"`
	Description string `json:"description"`
	File        string `json:"-"`
	Line        int    `json:"-"`
	// Stack       string    `json:"stack,omitempty"`
	Created time.Time `json:"-"`
}

func Fatal(err Error) {
	_, file, line, _ := runtime.Caller(1)
	ff := strings.Split(file, "/")

	if err.Code == 0 {
		err.Code = 500
	}
	if err.Type == "" {
		err.Type = "unknown"
	}
	err.File = strings.Join(ff[len(ff)-2:], "/")
	err.Line = line
	err.Created = time.Now()

	// rapi_log.Error(err.Description)

	panic(err)
}

func FatalIfError(err error) {
	if err != nil {
		Fatal(Error{
			Description: err.Error(),
		})
	}
}

func FatalIfTrue(ok bool, err Error) {
	if ok {
		Fatal(err)
	}
}

func FatalIfFalse(ok bool, err Error) {
	if !ok {
		Fatal(err)
	}
}
