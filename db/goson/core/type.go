package core

import (
	"strings"
)

const T_BOOL = 1

const Type8 = 2
const Type16 = 3
const Type32 = 4
const Type64 = 5

const TypeF32 = 6
const TypeF64 = 7

const TypeString = 8
const TypeTime = 9

const TypeStruct = 10

const TypeSlice = 11
const TypeMap = 12

func TypeToString(tp uint8) string {
	switch tp {
	case T_BOOL:
		return "bool"
	case Type8:
		return "i8"
	case Type16:
		return "i16"
	case Type32:
		return "i32"
	case Type64:
		return "i64"
	case TypeString:
		return "string"
	case TypeTime:
		return "time"
	case TypeStruct:
		return "struct"
	case TypeSlice:
		return "slice"
	case TypeMap:
		return "map"
	default:
		return "unknown"
	}
}

func TypeStringToByteType(tp string) uint8 {
	switch strings.ToLower(tp) {
	case "bool":
		return T_BOOL
	case "int8", "uint8":
		return Type8
	case "int16", "uint16":
		return Type16
	case "int32", "uint32":
		return Type32
	case "int64", "uint64":
		return Type64
	case "string":
		return TypeString
	case "time":
		return TypeTime
	case "struct":
		return TypeStruct
	case "slice":
		return TypeSlice
	case "map":
		return TypeMap
	default:
		return 0
	}
}
