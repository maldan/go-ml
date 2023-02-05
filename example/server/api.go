package main

import (
	ms_handler "github.com/maldan/go-ml/server/core/handler"
)

type User struct{}
type Template struct{}

type ArgsX struct {
	A int    `json:"a"`
	B int    `json:"b"`
	X string `json:"x"`
}

func (u User) GetIndex(ctx *ms_handler.Context, x ArgsX) int {
	return x.A + x.B
}
