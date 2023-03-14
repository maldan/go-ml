package mdb

import (
	"math"
)

/**
Record struct

[12 34] - record start
[0 0 0 0] - size of full record, include start + end
[0] - record flags, is deleted for example
[56 78] - record end
*/

type Record[T any] struct {
	offset uint64
	size   uint32
	table  *DataTable[T]
}

type SearchResult[T any] struct {
	IsFound bool `json:"isFound"`
	Count   int  `json:"count"`
	Result  []T  `json:"result"`
	// recordList []Record[T]
	table *DataTable[T]
}

func (s *SearchResult[T]) First() (T, bool) {
	if s.IsFound {
		return s.Result[0], true
	}
	return *new(T), false
}

func (s *Record[T]) Unpack() T {
	realData := unwrap(s.table.mem[s.offset : s.offset+uint64(s.size)])
	v := new(T)
	s.table.Container.Unmarshall(realData, v)
	return *v
}

func unwrap(bytes []byte) []byte {
	if bytes[0] != 0x12 {
		panic("non package")
	}
	hh := SIZE_OF_RECORD_START + RecordSize + RecordFlags
	return bytes[hh : len(bytes)-1]
}

func wrap(bytes []byte) []byte {
	fullSize := len(bytes) + SIZE_OF_RECORD_START + RecordSize + RecordFlags + SIZE_OF_RECORD_END

	// Calculate aligned size
	alignBy := 1
	alignedSize := math.Ceil(float64(fullSize)/float64(alignBy)) * float64(alignBy)
	zeroPadding := int(alignedSize) - fullSize
	fullSize = int(alignedSize)

	// Create array
	fullPackage := make([]byte, 0, fullSize)

	// Start
	fullPackage = append(fullPackage, 0x12)
	fullPackage = append(fullPackage, 0x34)

	// Size
	fullPackage = append(fullPackage, uint8(fullSize))
	fullPackage = append(fullPackage, uint8(fullSize>>8))
	fullPackage = append(fullPackage, uint8(fullSize>>16))
	fullPackage = append(fullPackage, uint8(fullSize>>24))

	// Flags
	fullPackage = append(fullPackage, 0)

	// Body
	fullPackage = append(fullPackage, bytes...)

	// Zero padding
	zero := make([]byte, zeroPadding)
	fullPackage = append(fullPackage, zero...)

	// End
	fullPackage = append(fullPackage, 0x56)
	fullPackage = append(fullPackage, 0x78)

	return fullPackage
}

/*func (s *Record[T]) Delete() bool {
	// Lock table
	s.table.rwLock.Lock()
	defer s.table.rwLock.Unlock()

	// Incorrect record
	if s.table.mem[s.offset] != core.RecordStartMark {
		return false
	}

	// Read flags
	b := []byte{0}
	s.table.file.ReadAt(b, int64(s.offset+core.RecordStart+core.RecordSize))
	fmt.Printf("FB: %v\n", b[0])
	b[0] |= core.MaskDeleted
	fmt.Printf("FAPL: %v\n", b[0])

	// Write back
	s.table.file.WriteAt(b, int64(s.offset+core.RecordStart+core.RecordSize))

	return true
}

func (s *Record[T]) Update(fields map[string]any) bool {
	// Unpack current
	unpack := s.Unpack()

	// Fill fields
	valueOf := reflect.ValueOf(&unpack).Elem()
	for name, value := range fields {
		field := valueOf.FieldByName(name)
		if field.Kind() == reflect.Invalid {
			continue
		}
		if !field.CanSet() {
			continue
		}
		switch value.(type) {
		case string:
			field.SetString(value.(string))
			break
		default:
			panic(fmt.Sprintf("can't set field for type %T", value))
		}
	}

	// Delete old
	s.Delete()

	// Create new
	s.offset = s.table.Insert(unpack).offset

	return true
}*/
