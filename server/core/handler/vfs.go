package handler

import (
	"embed"
	ms "github.com/maldan/go-ml/server"
)

type VFS struct {
	Root string
	Fs   embed.FS
}

func (V VFS) Handle(args ms.HandlerArgs) {
	//TODO implement me
	panic("implement me")
}
