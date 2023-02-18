package mdb_goson

import (
	"encoding/binary"
	"github.com/maldan/go-ml/db/goson/core"
)

/**
Header struct

[G O S O N D B $] - file id
[1] - version

[0 0 0 0 0 0 0 0] - auto increment
[0 0 0 0 0 0 0 0] - total records

[...] - name to id map
*/

/**
name to id map
[0] - total fields
[
	[0] - name length
	[...] - name
	[0] - id
] * totalFields
*/

type Header[T any] struct {
	FileId        string
	Version       uint8
	AutoIncrement uint64
	TotalRecords  uint64
	//NameToId      core.NameToId
	//IdToName      core.IdToName
	table *DataTable[T]
}

func (h *Header[T]) FromBytes(bytes []byte) {
	offset := 0

	// File id
	h.FileId = string(bytes[0:8])
	if string(h.FileId) != "GOSONDB$" {
		panic("non db")
	}
	offset += 8

	// Version
	h.Version = bytes[offset]
	offset += 1

	// AI
	h.AutoIncrement = binary.LittleEndian.Uint64(bytes[offset:])
	offset += 8

	// Total
	h.TotalRecords = binary.LittleEndian.Uint64(bytes[offset:])
	offset += 8

	// Set header
	h.table.Container.SetHeader(bytes[offset:])

	// Amount
	//amount := int(bytes[offset])
	//offset += 1

	/*// Init maps
	h.NameToId = map[string]uint8{}
	h.IdToName = map[uint8]string{}

	// Fill maps
	for i := 0; i < amount; i++ {
		// Name length
		nameLen := int(bytes[offset])
		offset += 1

		// Name
		name := string(bytes[offset : offset+nameLen])
		offset += nameLen

		// Fill map
		h.NameToId[name] = bytes[offset]
		offset += 1
	}

	// Fill id to name
	for name, id := range h.NameToId {
		h.IdToName[id] = name
	}*/
}

func (h *Header[T]) ToBytes() []byte {
	offset := 0
	bytes := make([]byte, core.HeaderSize)

	// File id
	copy(bytes, "GOSONDB$")
	offset += 8

	// Version
	if h.Version == 0 {
		h.Version = 1
	}
	bytes[offset] = h.Version
	offset += 1

	// AI and Total
	binary.LittleEndian.PutUint64(bytes[offset:], h.AutoIncrement)
	offset += 8
	binary.LittleEndian.PutUint64(bytes[offset:], h.TotalRecords)
	offset += 8

	// Get header data
	hdata := h.table.Container.GetHeader()

	copy(bytes[offset:], hdata)

	// Num of fields
	/*bytes[offset] = uint8(len(h.NameToId))
	offset += 1

	// Write name to id
	for name, id := range h.NameToId {
		// Name length
		bytes[offset] = uint8(len(name))
		offset += 1

		// Name
		copy(bytes[offset:], name)
		offset += len(name)

		// Id
		bytes[offset] = id
		offset += 1
	}*/

	return bytes
}
