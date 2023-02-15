package mdb_goson

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"sync"
)

type DataTable[T any] struct {
	mem  mmap.MMap
	file *os.File
	// structInfo core.StructInfo
	rwLock sync.RWMutex

	Header Header

	Path string
	Name string
}

func New[T any](path string, name string) *DataTable[T] {
	d := DataTable[T]{Path: path, Name: name}
	//d.structInfo.FieldNameToId = map[string]int{}

	d.open()
	d.remap()
	d.readHeader()
	d.writeHeader()

	return &d
}
