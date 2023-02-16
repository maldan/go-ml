package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-ml/db/goson/core"
	"reflect"
	"unsafe"
)

func Marshal[T any](v T, nameToId core.NameToId) []byte {
	ir := IR{}
	BuildIR(&ir, v, nameToId)

	return ir.Build()
}

func Unmarshall[T any](bytes []byte, idToName core.IdToName) T {
	t := new(T)
	unpack(bytes, unsafe.Pointer(t), core.TypeStringToByteType(reflect.TypeOf(*t).Kind().String()), *t, idToName)
	return *t
}

func unpack(bytes []byte, ptr unsafe.Pointer, ptrType uint8, typeHint any, idToName core.IdToName) int {
	offset := 0

	// Read type
	tp := bytes[offset]
	offset += 1

	fmt.Printf("ptrType %v - %v\n", ptrType, core.TypeToString(ptrType))

	switch tp {
	case core.T_BOOL:
		if ptrType == core.T_BOOL {
			*(*bool)(ptr) = bytes[offset] != 0
		}
		offset += 1
		break
	case core.T_8:
		if ptrType == core.T_8 {
			*(*uint8)(ptr) = bytes[offset]
		}
		offset += 1
		break
	case core.T_16:
		if ptrType == core.T_16 {
			*(*uint16)(ptr) = binary.LittleEndian.Uint16(bytes[offset:])
		}
		offset += 2
		break
	case core.T_32:
		if ptrType == core.T_32 {
			*(*uint32)(ptr) = binary.LittleEndian.Uint32(bytes[offset:])
		}
		offset += 4
		break
	case core.T_64:
		if ptrType == core.T_64 {
			*(*uint64)(ptr) = binary.LittleEndian.Uint64(bytes[offset:])
		}
		offset += 8
		break
	case core.TypeString:
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		if ptrType == core.TypeString {
			*(*string)(ptr) = string(blob)
		}
		break
	case core.TypeStruct:
		// Size
		offset += 2

		// Amount
		amount := int(bytes[offset])
		offset += 1

		typeOf := reflect.TypeOf(typeHint)

		for i := 0; i < amount; i++ {
			// field name
			fieldName, ok := idToName[bytes[offset]]
			if !ok {
				panic(fmt.Sprintf("field for id %v not found", bytes[offset]))
			}
			offset += 1

			field, _ := typeOf.FieldByName(fieldName)

			if field.Type.Kind() == reflect.Slice {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					core.TypeSlice,
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else if field.Type.Kind() == reflect.Struct {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					core.TypeStruct,
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					core.TypeStringToByteType(field.Type.String()),
					typeHint,
					idToName,
				)
			}
		}
		break
	case core.TypeSlice:
		// Size
		offset += 2

		// Amount
		amount := int(bytes[offset])
		offset += 2

		typeOf := reflect.TypeOf(typeHint).Elem()
		typeHint = reflect.New(typeOf).Elem().Interface()

		elemSlice := reflect.MakeSlice(reflect.SliceOf(typeOf), amount, amount)
		arr := make([]any, amount, amount)

		for i := 0; i < amount; i++ {
			offset += unpack(
				bytes[offset:],
				unsafe.Pointer(elemSlice.Index(i).Addr().Pointer()),
				ptrType,
				typeHint,
				idToName,
			)
			arr[i] = elemSlice.Index(i).Interface()
		}

		g := elemSlice.Pointer()

		*(*reflect.SliceHeader)(ptr) = reflect.SliceHeader{
			Data: g,
			Len:  amount,
			Cap:  amount,
		}
		break
	case core.T_CUSTOM:
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
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
	default:
		panic(fmt.Sprintf("uknown type %v", tp))
	}

	return offset
}
