package ms_handler

import (
	"embed"
)

type VFS struct {
	Root string
	Fs   embed.FS
}

func (V VFS) Handle(args Args) {
	//TODO implement me
	panic("implement me")
}
