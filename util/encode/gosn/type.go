package ml_gosn

import (
	"strings"
)

const T_BOOL = 1

const T_8 = 2
const T_16 = 3
const T_32 = 4
const T_64 = 5

const T_F32 = 6
const T_F64 = 7

// By default has 2 byte length, max size is 65535
const T_STRING = 8
const T_SLICE = 9
const T_MAP = 10
const T_STRUCT = 11

// By default has 2 byte length, max size is 65535
const T_CUSTOM = 12

// Extended big types, each type has 4 byte length
const T_BIG_STRING = 20
const T_BIG_SLICE = 21
const T_BIG_MAP = 22
const T_BIG_STRUCT = 23
const T_BIG_CUSTOM = 24

// Extended tiny types, each type has 1 byte length
const T_SHORT_STRING = 25
const T_SHORT_SLICE = 26
const T_SHORT_MAP = 27
const T_SHORT_STRUCT = 28
const T_SHORT_CUSTOM = 29

// Blob, exactly the same as T_BIG_SLICE for []byte. But much faster
const T_BLOB = 30

func TypeToString(tp uint8) string {
	switch tp {
	case T_BOOL:
		return "bool"
	case T_8:
		return "uint8"
	case T_16:
		return "uint16"
	case T_32:
		return "uint32"
	case T_64:
		return "uint64"
	case T_F32:
		return "float32"
	case T_F64:
		return "float64"
	case T_STRING, T_BIG_STRING:
		return "string"
	case T_SLICE, T_BIG_SLICE:
		return "slice"
	case T_STRUCT:
		return "struct"
	case T_MAP:
		return "map"
	case T_CUSTOM, T_BIG_CUSTOM:
		return "any"
	default:
		return "unknown"
	}
}

func TypeStringToTypeByte(tp string) uint8 {
	switch strings.ToLower(tp) {
	case "bool":
		return T_BOOL
	case "int8", "uint8":
		return T_8
	case "int16", "uint16":
		return T_16
	case "int", "int32", "uint32":
		return T_32
	case "int64", "uint64":
		return T_64
	case "float32":
		return T_F32
	case "float64":
		return T_F64
	case "string":
		return T_STRING
	case "struct", "custom":
		return T_STRUCT
	case "slice":
		return T_SLICE
	case "map":
		return T_MAP
	default:
		return 0
	}
}
