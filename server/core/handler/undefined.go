package ms_handler

import (
	"fmt"
	ms_error "github.com/maldan/go-ml/server/error"
)

type Undefined struct {
}

func (r Undefined) Handle(args Args) {
	ms_error.Fatal(ms_error.Error{
		Code: 404,
		Description: fmt.Sprintf(
			"Resource for '%v' route not found",
			args.Path,
		),
	})
}
