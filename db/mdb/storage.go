package mdb

import (
	"encoding/binary"
	"errors"
	"github.com/edsrzf/mmap-go"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (d *DataTable[T]) open() {
	finalPath, err := filepath.Abs(d.Path)
	if err != nil {
		panic(err)
	}

	finalPath += "/" + d.Name

	// Check if file exists
	if _, err = os.Stat(finalPath); errors.Is(err, fs.ErrNotExist) {
		// Create path for file
		err = os.MkdirAll(filepath.Dir(finalPath), 0777)
		if err != nil {
			panic(err)
		}

		// Init file, because 0 length file fails with memory mapping
		if err = ioutil.WriteFile(finalPath, d.Header.ToBytes(), 0777); err != nil {
			panic(err)
		}
	}

	// Open file
	f, err := os.OpenFile(finalPath, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	// Get size
	s, _ := os.Stat(finalPath)

	d.fileSize = uint64(s.Size())
	d.file = f
}

func (d *DataTable[T]) remap() {
	// Unmap previous
	if d.mem != nil {
		err := d.mem.Unmap()
		if err != nil {
			panic(err)
		}
	}

	// Map new
	mem, err := mmap.Map(d.file, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	d.mem = mem
}

func (d *DataTable[T]) readHeader() {
	// Thread safe operation
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	d.Header.FromBytes(d.mem)
}

func (d *DataTable[T]) writeHeader() {
	// Thread safe operation
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	_, err := d.file.WriteAt(d.Header.ToBytes(), 0)
	if err != nil {
		panic(err)
	}
}

func (d *DataTable[T]) writeAI() {
	// Thread safe operation
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	// Prepare ai
	ai := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(ai, d.Header.AutoIncrement)

	_, err := d.file.WriteAt(ai, 8+1) // GOSON + Version
	if err != nil {
		panic(err)
	}
}
