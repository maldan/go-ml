package ml_gosn

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

func Marshal(v any) []byte {
	ir := IR{}
	BuildIR(&ir, v, nil)
	return ir.Build()
}

func MarshalExt(v any, nameToId NameToId) []byte {
	ir := IR{}
	BuildIR(&ir, v, nameToId)
	return ir.Build()
}

func Unmarshall(bytes []byte, out any) {
	UnmarshallExt(bytes, out, nil)
}

func UnmarshallExt(bytes []byte, out any, idToName IdToName) {
	valueOf := reflect.ValueOf(out)
	typeByte := TypeStringToTypeByte(valueOf.Elem().Kind().String())
	unpack(bytes, valueOf.UnsafePointer(), typeByte, valueOf.Elem().Interface(), idToName)
}

func unpack(bytes []byte, ptr unsafe.Pointer, ptrType uint8, typeHint any, idToName IdToName) int {
	offset := 0

	// Read type
	tp := bytes[offset]
	offset += 1

	// fmt.Printf("ptrType %v - %v\n", ptrType, core.TypeToString(ptrType))

	switch tp {
	case T_BOOL:
		if ptrType == T_BOOL {
			*(*bool)(ptr) = bytes[offset] != 0
		}
		offset += 1
		break
	case T_8:
		if ptrType == T_8 {
			*(*uint8)(ptr) = bytes[offset]
		}
		offset += 1
		break
	case T_16:
		if ptrType == T_16 {
			*(*uint16)(ptr) = binary.LittleEndian.Uint16(bytes[offset:])
		}
		offset += 2
		break
	case T_32:
		if ptrType == T_32 {
			*(*uint32)(ptr) = binary.LittleEndian.Uint32(bytes[offset:])
		}
		offset += 4
		break
	case T_64:
		if ptrType == T_64 {
			*(*uint64)(ptr) = binary.LittleEndian.Uint64(bytes[offset:])
		}
		offset += 8
		break
	case T_F32:
		if ptrType == T_F32 {
			*(*uint32)(ptr) = binary.LittleEndian.Uint32(bytes[offset:])
		}
		offset += 4
		break
	case T_F64:
		if ptrType == T_F64 {
			*(*uint64)(ptr) = binary.LittleEndian.Uint64(bytes[offset:])
		}
		offset += 8
		break
	case T_STRING, T_SHORT_STRING, T_BIG_STRING:
		size := 0
		if tp == T_STRING {
			size = int(bytes[offset])
			offset += 1
		}
		if tp == T_SHORT_STRING {
			size = int(binary.LittleEndian.Uint16(bytes[offset:]))
			offset += 2
		}
		if tp == T_BIG_STRING {
			size = int(binary.LittleEndian.Uint32(bytes[offset:]))
			offset += 4
		}

		// Get blob and offset
		blob := bytes[offset : offset+size]
		offset += size

		if ptrType == T_STRING {
			// Eliminate memory leaks
			copyOfString := make([]byte, len(blob))
			copy(copyOfString, blob)

			*(*string)(ptr) = string(copyOfString)
		}
		break
	case T_SLICE, T_SHORT_SLICE, T_BIG_SLICE:
		// size := 0
		if tp == T_SLICE {
			// size = int(bytes[offset])
			offset += 1
		}
		if tp == T_SHORT_SLICE {
			// size = int(binary.LittleEndian.Uint16(bytes[offset:]))
			offset += 2
		}
		if tp == T_BIG_SLICE {
			// size = int(binary.LittleEndian.Uint32(bytes[offset:]))
			offset += 4
		}

		// Amount
		amount := int(binary.LittleEndian.Uint32(bytes[offset:]))
		offset += 4

		typeOf := reflect.TypeOf(typeHint).Elem()
		typeHint = reflect.New(typeOf).Elem().Interface()

		elemSlice := reflect.MakeSlice(reflect.SliceOf(typeOf), amount, amount)

		// Get pointer type of each element
		elementPtrType := TypeStringToTypeByte(typeOf.String())

		for i := 0; i < amount; i++ {
			offset += unpack(
				bytes[offset:],
				unsafe.Pointer(elemSlice.Index(i).Addr().Pointer()),
				elementPtrType,
				typeHint,
				idToName,
			)
		}

		g := elemSlice.Pointer()

		*(*reflect.SliceHeader)(ptr) = reflect.SliceHeader{
			Data: g,
			Len:  amount,
			Cap:  amount,
		}

		break
	case T_STRUCT, T_SHORT_STRUCT, T_BIG_STRUCT:
		// Size
		// size := 0
		if tp == T_STRUCT {
			// size = int(bytes[offset])
			offset += 1
		}
		if tp == T_SHORT_STRUCT {
			// size = int(binary.LittleEndian.Uint16(bytes[offset:]))
			offset += 2
		}
		if tp == T_BIG_STRUCT {
			// size = int(binary.LittleEndian.Uint32(bytes[offset:]))
			offset += 4
		}
		// fmt.Printf("READ SIZE: %v\n", size)

		// Amount
		amount := int(bytes[offset])
		offset += 1

		typeOf := reflect.TypeOf(typeHint)

		for i := 0; i < amount; i++ {
			fieldName := ""

			if idToName == nil {
				nameSize := int(bytes[offset])
				offset += 1
				fieldName = string(bytes[offset : offset+nameSize])
				offset += nameSize
			} else {
				// field name
				name, ok := idToName[bytes[offset]]
				if !ok {
					panic(fmt.Sprintf("field for id %v not found", bytes[offset]))
				}
				fieldName = name
				offset += 1
			}

			// Get real field
			field, _ := typeOf.FieldByName(fieldName)

			if field.Type.Kind() == reflect.Slice {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					T_SLICE,
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else if field.Type.Kind() == reflect.Struct {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					T_STRUCT,
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					TypeStringToTypeByte(field.Type.String()),
					typeHint,
					idToName,
				)
			}
		}
		break
	case T_CUSTOM, T_SHORT_CUSTOM, T_BIG_CUSTOM:
		// Size
		size := 0
		if tp == T_CUSTOM {
			size = int(bytes[offset])
			offset += 1
		}
		if tp == T_SHORT_CUSTOM {
			size = int(binary.LittleEndian.Uint16(bytes[offset:]))
			offset += 2
		}
		if tp == T_BIG_CUSTOM {
			size = int(binary.LittleEndian.Uint32(bytes[offset:]))
			offset += 4
		}

		// Read data
		blob := bytes[offset : offset+size]
		offset += size

		// Get pointer to type hint
		pp := reflect.PointerTo(reflect.TypeOf(typeHint))

		// Find method from bytes
		fromBytes, ok := pp.MethodByName("FromBytes")
		if ok {
			// Convert unsafe pointer to reflect pointer with type hinting
			rf := reflect.NewAt(reflect.TypeOf(typeHint), ptr)

			// Call
			ret := fromBytes.Func.Call([]reflect.Value{rf, reflect.ValueOf(blob)})
			if len(ret) > 0 {
				err := ret[0].Interface()
				if err != nil {
					panic(err)
				}
			}
		} else {
			panic(fmt.Sprintf("custom type %T doesn't have FromBytes method", typeHint))
		}

		break
	case T_BLOB:
		// Read size of data
		size := int(binary.LittleEndian.Uint32(bytes[offset:]))
		offset += 4

		// Get blob
		blob := bytes[offset : offset+size]
		offset += size

		if ptrType == T_SLICE {
			// Eliminate memory leaks
			copyOfBlob := make([]byte, len(blob))
			copy(copyOfBlob, blob)

			*(*[]byte)(ptr) = copyOfBlob
		}
		break
	default:
		panic(fmt.Sprintf("uknown type %v", tp))
	}

	return offset
}
