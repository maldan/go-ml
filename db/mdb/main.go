package mdb

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	ml_fs "github.com/maldan/go-ml/util/io/fs"
	"os"
	"sync"
	"time"
)

type DBMapper interface {
	Map([]byte)
}

type DBContainer interface {
	Prepare(v any)
	Marshal(a any) []byte
	Unmarshall(b []byte, out any)

	GetMapper(fieldList string, tp any) any

	GetHeader() []byte
	SetHeader([]byte)

	GetStruct() map[string]string
}

type DataTable[T any] struct {
	mem      mmap.MMap
	file     *os.File
	rwLock   sync.RWMutex
	fileSize uint64

	Container DBContainer

	Header Header[T]

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

func New[T any](path string, name string, container DBContainer) *DataTable[T] {
	table := DataTable[T]{Path: path, Name: name}
	table.Header.table = &table

	table.Container = container
	table.Container.Prepare(*new(T))

	table.open()
	table.remap()
	table.readHeader()
	table.writeHeader()

	return &table
}
