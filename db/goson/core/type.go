package core

import (
	"strings"
)

const T_BOOL = 1

const T_8 = 2
const T_16 = 3
const T_32 = 4
const T_64 = 5

const TypeF32 = 6
const TypeF64 = 7

const TypeString = 8
const TypeSlice = 9

const TypeStruct = 10
const TypeMap = 11

const T_CUSTOM = 12

// Custom types
// const TypeTime = 12 // 0001-01-01T00:00:00+00:00

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
	case TypeF32:
		return "float32"
	case TypeF64:
		return "float64"
	case TypeString:
		return "string"
	case TypeStruct:
		return "struct"
	case TypeSlice:
		return "slice"
	case TypeMap:
		return "map"
	case T_CUSTOM:
		return "custom"
	default:
		return "unknown"
	}
}

func TypeStringToByteType(tp string) uint8 {
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
	case "string":
		return TypeString
	case "struct", "custom":
		return TypeStruct
	case "slice":
		return TypeSlice
	case "map":
		return TypeMap
	default:
		return 0
	}
}
