package mdb

import (
	"encoding/binary"
	"errors"
	"fmt"
	ml_expression "github.com/maldan/go-ml/util/expression"
	"reflect"
)

type ArgsFind struct {
	// FieldList       string
	Offset          int
	Limit           int
	WhereExpression string
	Where           func(any2 any) bool
}

type ArgsUpdate struct {
	FieldList string
	Limit     int
	Where     func(any2 any) bool
	Change    func(any2 any)
}

func __insertToTable(table *DataTable, value any, isLock bool, isRemap bool) Record {
	// Lock table
	if isLock {
		table.rwLock.Lock()
		defer table.rwLock.Unlock()
	}
	if isRemap {
		defer table.remap()
	}

	bytes := table.Container.Marshal(value)
	// bytes := goson.Marshal(value, table.Header.NameToId)
	bytes = wrap(bytes)

	// Get file size
	stat, err := table.file.Stat()
	if err != nil {
		panic(err)
	}

	// Write at end of file
	endOfFile := stat.Size()
	n, err := table.file.WriteAt(bytes, endOfFile)
	if err != nil {
		panic(err)
	}
	if n != len(bytes) {
		panic(errors.New("incomplete writing"))
	}

	// Return record info
	return Record{
		offset: uint64(endOfFile),
		size:   uint32(len(bytes)),
		table:  table,
	}
}

func __markDeleted(table *DataTable, offset uint64) {
	// Read flags
	b := []byte{0}
	_, err := table.file.ReadAt(b, int64(offset+SIZE_OF_RECORD_START+RecordSize))
	if err != nil {
		panic(err)
	}
	b[0] |= MaskDeleted

	// Write back
	_, err = table.file.WriteAt(b, int64(offset+SIZE_OF_RECORD_START+RecordSize))
	if err != nil {
		panic(err)
	}
}

func __offsetUntil(slice []byte, seq ...uint8) (uint64, bool) {
	for i := 0; i < len(slice)-(len(seq)-1); i++ {
		isFound := true

		// Compare sequence
		for j := 0; j < len(seq); j++ {
			if slice[i+j] != seq[j] {
				isFound = false
				break
			}
		}

		if isFound {
			return uint64(i), true
		}
	}

	return 0, false
}

func (d *DataTable) GenerateId() uint64 {
	d.rwLock.Lock()
	id := uint64(0)
	d.Header.AutoIncrement += 1
	id = d.Header.AutoIncrement
	d.rwLock.Unlock()
	d.writeAI()
	return id
}

func (d *DataTable) Insert(v any) Record {
	return __insertToTable(d, v, true, true)
}

func (d *DataTable) InsertMany(v []any) {
	d.rwLock.Lock()
	defer d.rwLock.Unlock()
	defer d.remap()

	for i := 0; i < len(v); i++ {
		__insertToTable(d, v[i], false, false)
	}
}

func (d *DataTable) ForEach(fn func(offset uint64, size uint32) bool) {
	offset := uint64(HEADER_SIZE)

	// Empty table
	if len(d.mem) <= HEADER_SIZE {
		return
	}

	for {
		// Check header
		if !(d.mem[offset] == 0x12 && d.mem[offset+1] == 0x34) {
			fmt.Printf("corrupted at %v\n", offset)
			// error
			next, ok := __offsetUntil(d.mem[offset:], 0x12, 0x34)
			if ok {
				offset += next
			} else {
				break
			}
		}

		// Read size and flags
		size := binary.LittleEndian.Uint32(d.mem[offset+SIZE_OF_RECORD_START:])
		flags := int(d.mem[offset+SIZE_OF_RECORD_START+RecordSize])

		// Check size, it may be out of memory if record is corrupted
		if uint64(size)+offset > uint64(len(d.mem)) {
			offset += 1 // offset by 1 byte
			if offset > uint64(len(d.mem)) {
				break
			}
			continue
		}

		// Check end
		if !(d.mem[offset+uint64(size)-2] == 0x56 && d.mem[offset+uint64(size)-1] == 0x78) {
			offset += 1 // offset by 1 byte
			if offset > uint64(len(d.mem)) {
				break
			}
			continue
		}

		// If field not deleted
		if flags&MaskDeleted != MaskDeleted {
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

func (d *DataTable) Count(args ArgsFind) int {
	// Lock table
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	// Create mapper for capturing values from bytes
	mapperContainer := reflect.New(d.Type).Interface()
	// mapper := d.Container.GetMapper(args.FieldList, mapperContainer).(DBMapper)

	// Compile expression
	if args.WhereExpression != "" {
		expr, _ := ml_expression.Parse(args.WhereExpression)
		expr.Bind(mapperContainer)
		args.Where = func(any2 any) bool {
			return expr.Execute().(bool)
		}
	}

	offsetCounter := args.Offset
	counter := 0

	// Go through each record
	d.ForEach(func(offset uint64, size uint32) bool {
		// mapper.Map(d.mem[offset+SIZE_OF_RECORD_START+RecordSize+RecordFlags:])
		d.Container.Unmarshall(
			d.mem[offset+SIZE_OF_RECORD_START+RecordSize+RecordFlags:],
			mapperContainer,
		)

		// Collect values
		if args.Where(mapperContainer) {
			if offsetCounter > 0 {
				offsetCounter -= 1
				return true
			}

			counter += 1

			// Check limit
			if args.Limit > 0 && counter >= args.Limit {
				return false
			}
		}

		return true
	})

	return counter
}

func (d *DataTable) FindBy(args ArgsFind) SearchResult {
	// Lock table
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	// Return
	searchResult := SearchResult{}

	// Create mapper for capturing values from bytes
	mapperContainer := reflect.New(d.Type).Interface()
	// mapper := d.Container.GetMapper(args.FieldList, mapperContainer).(DBMapper)

	// Compile expression
	if args.WhereExpression != "" {
		expr, _ := ml_expression.Parse(args.WhereExpression)
		expr.Bind(mapperContainer)
		args.Where = func(any2 any) bool {
			return expr.Execute().(bool)
		}
	}

	offsetCounter := args.Offset

	// Go through each record
	d.ForEach(func(offset uint64, size uint32) bool {
		d.Container.Unmarshall(
			d.mem[offset+SIZE_OF_RECORD_START+RecordSize+RecordFlags:],
			mapperContainer,
		)

		// mapper.Map()

		// Collect values
		if args.Where(mapperContainer) {
			if offsetCounter > 0 {
				offsetCounter -= 1
				return true
			}

			searchResult.table = d
			searchResult.IsFound = true
			record := Record{
				offset: offset,
				size:   size,
				table:  d,
			}
			unpacked := record.Unpack()
			searchResult.Result = append(searchResult.Result, unpacked)
			searchResult.Count += 1

			// Check limit
			if args.Limit > 0 && len(searchResult.Result) >= args.Limit {
				return false
			}
		}

		return true
	})

	if len(searchResult.Result) == 0 {
		searchResult.Result = make([]any, 0)
	}

	return searchResult
}

/*func (d *DataTable) DeleteBy(args ArgsFind[T]) {
	// Lock table
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	// Create mapper for capturing values from bytes
	mapperContainer := new(T)
	mapper := d.Container.GetMapper(args.FieldList, mapperContainer).(DBMapper)

	// Go through each record
	counter := 0
	d.ForEach(func(offset uint64, size uint32) bool {
		mapper.Map(d.mem[offset+SIZE_OF_RECORD_START+RecordSize+RecordFlags:])

		// Collect values
		if args.Where(mapperContainer) {
			counter += 1

			__markDeleted(d, offset)

			// Check limit
			if args.Limit > 0 && counter >= args.Limit {
				return false
			}
		}

		return true
	})
}*/

/*func (d *DataTable) UpdateBy(args ArgsUpdate[T]) {
	// Lock table
	d.rwLock.Lock()
	defer d.rwLock.Unlock()
	defer d.remap()

	// Create mapper for capturing values from bytes
	mapperContainer := new(T)
	mapper := d.Container.GetMapper(args.FieldList, mapperContainer).(DBMapper)

	// Go through each record
	counter := 0
	d.ForEach(func(offset uint64, size uint32) bool {
		mapper.Map(d.mem[offset+SIZE_OF_RECORD_START+RecordSize+RecordFlags:])

		// Check condition
		if args.Where(mapperContainer) {
			counter += 1

			// Get record
			record := Record[T]{offset: offset, size: size, table: d}

			// Unpack value
			raw := record.Unpack()

			// Change data
			args.Change(&raw)

			// Insert new
			__insertToTable(d, raw, false, false)

			// Delete old
			__markDeleted(d, offset)

			// Check limit
			if args.Limit > 0 && counter >= args.Limit {
				return false
			}
		}

		return true
	})
}*/

func (d *DataTable) Close() error {
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
