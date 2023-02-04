package handler

import (
	"fmt"
	ms_config "github.com/maldan/go-ml/server/config"
	ms_error "github.com/maldan/go-ml/server/error"
)

type Undefined struct {
}

func (r Undefined) Handle(args ms_config.HandlerArgs) {
	ms_error.Fatal(ms_error.Error{
		Code: 404,
		Description: fmt.Sprintf(
			"Resource for '%v' route not found",
			1, //args.Route,
		),
	})
}
