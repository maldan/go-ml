package mdb_goson

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maldan/go-ml/db/goson/core"
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
	endOfFile := stat.Size()
	n, err := d.file.WriteAt(bytes, endOfFile)
	if err != nil {
		panic(err)
	}
	if n != len(bytes) {
		panic(errors.New("incomplete writing"))
	}

	// Return record info
	return Record[T]{
		offset: uint64(endOfFile),
		size:   uint32(len(bytes)),
		table:  d,
	}
}

func (d *DataTable[T]) ForEach(fn func(offset uint64, size uint32) bool) {
	offset := uint64(core.HeaderSize)

	for {
		// Read size and flags
		size := binary.LittleEndian.Uint32(d.mem[offset+core.SIZE_OF_RECORD_START:])
		flags := int(d.mem[offset+core.SIZE_OF_RECORD_START+core.RecordSize])

		// If field not deleted
		if flags&core.MaskDeleted != core.MaskDeleted {
			if !fn(offset, size) {
				break
			}
		}

		// Go to next value
		offset += uint64(size)
		if offset >= uint64(len(d.mem)) {
			break
		}
	}
}

type ArgsFind[T any] struct {
	FieldList string
	Limit     int
	Where     func(*T) bool
}

type ArgsUpdate[T any] struct {
	FieldList string
	Limit     int
	Where     func(*T) bool
	Change    func(*T)
}

func (d *DataTable[T]) FindBy(args ArgsFind[T]) SearchResult[T] {
	// Lock table
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	// Return
	searchResult := SearchResult[T]{}

	// Create mapper for capturing values from bytes
	mapper := goson.NewMapper[T](d.Header.NameToId)

	// Field list
	fieldIdList := core.NameListToIdList(args.FieldList, d.Header.NameToId)

	// Go through each record
	d.ForEach(func(offset uint64, size uint32) bool {
		mapper.Map(d.mem[offset+core.SIZE_OF_RECORD_START+core.RecordSize+core.RecordFlags:], fieldIdList)

		// Collect values
		if args.Where(&mapper.Container) {
			searchResult.table = d
			searchResult.IsFound = true
			searchResult.Result = append(searchResult.Result, Record[T]{
				offset: offset,
				size:   size,
				table:  d,
			})

			// Check limit
			if args.Limit > 0 && len(searchResult.Result) >= args.Limit {
				return false
			}
		}

		return true
	})

	return searchResult
}

func (d *DataTable[T]) DeleteBy(args ArgsFind[T]) {
	// Lock table
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	// Create mapper for capturing values from bytes
	mapper := goson.NewMapper[T](d.Header.NameToId)

	// Field list
	fieldIdList := core.NameListToIdList(args.FieldList, d.Header.NameToId)

	// Go through each record
	counter := 0
	d.ForEach(func(offset uint64, size uint32) bool {
		mapper.Map(d.mem[offset+core.SIZE_OF_RECORD_START+core.RecordSize+core.RecordFlags:], fieldIdList)

		// Collect values
		if args.Where(&mapper.Container) {
			counter += 1

			// Read flags
			b := []byte{0}
			d.file.ReadAt(b, int64(offset+core.SIZE_OF_RECORD_START+core.RecordSize))
			fmt.Printf("FB: %v\n", b[0])
			b[0] |= core.MaskDeleted
			fmt.Printf("FAPL: %v\n", b[0])

			// Write back
			d.file.WriteAt(b, int64(offset+core.SIZE_OF_RECORD_START+core.RecordSize))

			// Check limit
			if args.Limit > 0 && counter >= args.Limit {
				return false
			}
		}

		return true
	})
}

func (d *DataTable[T]) UpdateBy(args ArgsFind[T], fields map[string]any) {
	// Lock table
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	// Create mapper for capturing values from bytes
	mapper := goson.NewMapper[T](d.Header.NameToId)
	fieldIdList := core.NameListToIdList(args.FieldList, d.Header.NameToId)

	// Go through each record
	counter := 0
	d.ForEach(func(offset uint64, size uint32) bool {
		mapper.Map(d.mem[offset+core.SIZE_OF_RECORD_START+core.RecordSize+core.RecordFlags:], fieldIdList)

		// Collect values
		if args.Where(&mapper.Container) {
			counter += 1

			// Check limit
			if args.Limit > 0 && counter >= args.Limit {
				return false
			}
		}

		return true
	})
}

func (d *DataTable[T]) Close() error {
	err := d.mem.Unmap()
	if err != nil {
		return err
	}

	err = d.file.Close()
	if err != nil {
		return err
	}

	return nil
}
