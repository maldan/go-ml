package mdb_goson

import (
	"errors"
	"github.com/maldan/go-ml/db/goson/goson"
)

func (d *DataTable[T]) GenerateId() uint64 {
	d.rwLock.Lock()
	id := uint64(0)
	d.Header.AutoIncrement += 1
	id = d.Header.AutoIncrement
	d.rwLock.Unlock()
	d.writeAI()
	return id
}

func (d *DataTable[T]) Insert(v T) Record[T] {
	// Lock table
	d.rwLock.Lock()
	defer d.rwLock.Unlock()
	defer d.remap()

	bytes := goson.Marshal(v, d.Header.NameToId)

	bytes = wrap(bytes)

	// Get file size
	stat, err := d.file.Stat()
	if err != nil {
		panic(err)
	}

	// Write at end of file
	n, err := d.file.WriteAt(bytes, stat.Size())
	if err != nil {
		panic(err)
	}
	if n != len(bytes) {
		panic(errors.New("incomplete writing"))
	}

	// Return record info
	return Record[T]{
		offset: uint64(stat.Size()),
		size:   uint32(len(bytes)),
		table:  d,
	}
}
