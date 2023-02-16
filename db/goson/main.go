package mdb_goson

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-ml/db/goson/core"
	"github.com/maldan/go-ml/db/goson/goson"
	ml_fs "github.com/maldan/go-ml/io/fs"
	"os"
	"sync"
	"time"
)

type DataTable[T any] struct {
	mem      mmap.MMap
	file     *os.File
	rwLock   sync.RWMutex
	fileSize uint64

	Header Header

	Path string
	Name string
}

func (d *DataTable[T]) SetBackupSchedule(dst string, each time.Duration) {
	go (func() {
		// File size
		fmt.Printf("Initializing backup: %v - %v\n", d.Name, d.fileSize)

		tt := time.Now()
		today := time.Now().Format("2006_01_02")
		timestamp := time.Now().Format("15_04_05")
		err := ml_fs.Copy(d.Path+"/"+d.Name, dst+"/"+today+"/"+d.Name+"_"+timestamp+".tmp")
		if err != nil {
			fmt.Printf("Backup Error: %v\n", err)
		} else {
			// Rename without .tmp, it means the file was successfully copied
			ml_fs.Rename(dst+"/"+today+"/"+d.Name+"_"+timestamp+".tmp", dst+"/"+today+"/"+d.Name+"_"+timestamp)
		}

		fmt.Printf("Backup done: %v - %v\n", d.Name, time.Since(tt))

		time.Sleep(each)
	})()
}

func New[T any](path string, name string) *DataTable[T] {
	d := DataTable[T]{Path: path, Name: name}

	// Check if type possible to serialize
	nid := core.NameToId{}
	nid.FromStruct(*new(T))
	_ = goson.Marshal(*new(T), nid)

	d.open()
	d.remap()
	d.readHeader()
	d.writeHeader()

	return &d
}
