package mdb

const HEADER_SIZE = 2048

// SIZE_OF_RECORD_START is 2 bytes 0x1234 mark, stands for start of record
const SIZE_OF_RECORD_START = 2

// RecordSize is size of each record
const RecordSize = 4

// RecordFlags is deleted
const RecordFlags = 1

// SIZE_OF_RECORD_END is 2 bytes 0x5678  mark, stands for end of record
const SIZE_OF_RECORD_END = 2

const RecordStartMark = 0x12
const RecordEndMark = 0x34

const MaskDeleted = 0b1000_0000
const MaskTotalFields = 0b0011_1111

type StructInfo struct {
	FieldCount    int
	FieldNameToId map[string]int
	FieldType     []int
	FieldName     []string
	FieldOffset   []uintptr
}
