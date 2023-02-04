package main

import (
	ms "github.com/maldan/go-ml/server"
)

type User struct{}
type Template struct{}

type ArgsX struct {
	A int    `json:"a"`
	B int    `json:"b"`
	X string `json:"x"`
}

func (u User) GetIndex(ctx *ms.Context, x ArgsX) int {
	return x.A + x.B
}
