package mdb

import (
	"encoding/binary"
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

type Header struct {
	FileId        string
	Version       uint8
	AutoIncrement uint64
	TotalRecords  uint64
	table         *DataTable
}

func (h *Header) FromBytes(bytes []byte) {
	offset := 0

	// File id
	h.FileId = string(bytes[0:8])
	if string(h.FileId) != "MEGADBGO" {
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
}

func (h *Header) ToBytes() []byte {
	offset := 0
	bytes := make([]byte, HEADER_SIZE)

	// File id
	copy(bytes, "MEGADBGO")
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
	data := h.table.Container.GetHeader()
	copy(bytes[offset:], data)

	return bytes
}
