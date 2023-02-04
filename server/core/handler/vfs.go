package handler

import (
	"embed"
	ms_config "github.com/maldan/go-ml/server/config"
)

type VFS struct {
	Root string
	Fs   embed.FS
}

func (V VFS) Handle(args ms_config.HandlerArgs) {
	//TODO implement me
	panic("implement me")
}
