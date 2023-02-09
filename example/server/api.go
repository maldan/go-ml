package main

import (
	ml_console "github.com/maldan/go-ml/io/console"
	ms_handler "github.com/maldan/go-ml/server/core/handler"
	ml_time "github.com/maldan/go-ml/util/time"
)

type User struct{}
type Template struct{}

type ArgsX struct {
	A int    `json:"a"`
	B int    `json:"b"`
	X string `json:"x"`
}

type Args2 struct {
	A       string       `json:"a"`
	B       int          `json:"b"`
	Created ml_time.Time `json:"created"`
}

func (u User) GetIndex(ctx *ms_handler.Context, x ArgsX) int {
	return x.A + x.B
}

func (u User) PostIndex(x Args2) {
	ml_console.PrettyPrint(x)
	// fmt.Printf("%+v\n", x.Created)
}
