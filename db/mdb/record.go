package mdb

import (
	"math"
	"reflect"
)

/**
Record struct

[12 34] - record start
[0 0 0 0] - size of full record, include start + end
[0] - record flags, is deleted for example
[56 78] - record end
*/

type Record struct {
	offset uint64
	size   uint32
	table  *DataTable
}

type SearchResult struct {
	IsFound bool  `json:"isFound"`
	Count   int   `json:"count"`
	Result  []any `json:"result"`

	Page    int `json:"page"`
	Total   int `json:"total"`
	PerPage int `json:"perPage"`
	// recordList []Record[T]
	table *DataTable
}

func (s *SearchResult) First() (any, bool) {
	if s.IsFound {
		return s.Result[0], true
	}
	return nil, false
}

func (s *Record) Unpack() any {
	realData := unwrap(s.table.mem[s.offset : s.offset+uint64(s.size)])
	v := reflect.New(s.table.Type).Interface()
	s.table.Container.Unmarshall(realData, v)
	return v
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
