package gosn

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

const T_STRING = 8
const T_SLICE = 9

const T_STRUCT = 10
const T_MAP = 11

const T_CUSTOM = 12

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
	case T_STRING:
		return "string"
	case T_SLICE:
		return "slice"
	case T_STRUCT:
		return "struct"
	case T_MAP:
		return "map"
	case T_CUSTOM:
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
