package core

const TypeBool = 1

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
	case TypeBool:
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
	case TypeSlice:
		return "slice"
	case TypeTime:
		return "time"
	case TypeStruct:
		return "struct"
	default:
		return "unknown"
	}
}
