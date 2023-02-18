package gosn_driver

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-ml/util/encode/gosn"
	"log"
	"reflect"
	"unsafe"
)

type Mapper struct {
	Container any
	NameToId  ml_gosn.NameToId
	MapOffset []unsafe.Pointer

	// Copy bytes and parse it and assign to pointer
	SetCustomContent []func([]byte) error

	// FieldIdList
	FieldIdList []uint8
}

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func NewMapper(nameToId ml_gosn.NameToId, fieldList string, out any) *Mapper {
	mapper := Mapper{
		Container:        out,
		NameToId:         nameToId,
		MapOffset:        make([]unsafe.Pointer, 255),
		SetCustomContent: make([]func([]byte) error, 255),
	}

	// Convert field name to field id
	mapper.FieldIdList = ml_gosn.NameListToIdList(fieldList, nameToId)

	typeOf := reflect.TypeOf(mapper.Container).Elem()
	start := reflect.ValueOf(out).UnsafePointer()
	//start := unsafe.Pointer(&mapper.Container)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldId, ok := nameToId[typeOf.Field(i).Name]
		if !ok {
			log.Fatalf("field id %v not found", typeOf.Field(i).Name)
		}
		mapper.MapOffset[fieldId] = unsafe.Add(start, typeOf.Field(i).Offset)

		// Get pointer to type hint
		fieldPointer := reflect.PointerTo(typeOf.Field(i).Type)

		// From bytes
		fromBytes, ok := fieldPointer.MethodByName("FromBytes")
		if ok {
			// Get real function but it's in any type
			realFn := fromBytes.Func.Interface()
			rfnPtr := unsafe.Pointer(&realFn)
			iface := (*emptyInterface)(rfnPtr)

			// De interface pointer. It means cast any to real type
			*(*emptyInterface)(rfnPtr) = emptyInterface{
				typ: iface.ptr,
			}

			// Create callback
			mapper.SetCustomContent[fieldId] = func(bytes []byte) error {
				castFn := *(*func(unsafe.Pointer, []byte) error)(rfnPtr)
				return castFn(mapper.MapOffset[fieldId], bytes)
			}
		}
	}

	return &mapper
}

func typeSize(bytes []byte) int {
	switch bytes[0] {
	case ml_gosn.T_BOOL, ml_gosn.T_8:
		return 1
	case ml_gosn.T_16:
		return 2
	case ml_gosn.T_32, ml_gosn.T_F32:
		return 4
	case ml_gosn.T_64, ml_gosn.T_F64:
		return 8
	case ml_gosn.T_STRING, ml_gosn.T_CUSTOM:
		// Field size
		fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))
		return 2 + fieldSize // 2 - is length size info
	default:
		panic(fmt.Sprintf("unknown type %v", bytes[0]))
		return 0
	}
}

func applyType(v *Mapper, bytes []byte, offset int, fieldId uint8) {
	off := v.MapOffset[fieldId]

	const offType = 1 // offset type
	const offLen = 2  // offset length of field

	switch bytes[offset] {
	case ml_gosn.T_8:
		*(*uint8)(off) = bytes[offset]
		break
	case ml_gosn.T_16:
		*(*uint16)(off) = binary.LittleEndian.Uint16(bytes[offset+offType:])
		break
	case ml_gosn.T_32:
		*(*uint32)(off) = binary.LittleEndian.Uint32(bytes[offset+offType:])
		break
	case ml_gosn.T_64:
		*(*uint64)(off) = binary.LittleEndian.Uint64(bytes[offset+offType:])
		break
	case ml_gosn.T_STRING:
		fieldSize := int(binary.LittleEndian.Uint16(bytes[offset+offType:]))
		bts := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))

		hh := (*reflect.StringHeader)(off)
		hh.Data = bts.Data + uintptr(offset) + offType + offLen
		hh.Len = fieldSize
		break
	case ml_gosn.T_CUSTOM:
		fieldSize := int(binary.LittleEndian.Uint16(bytes[offset+offType:]))

		if v.SetCustomContent[fieldId] == nil {
			panic(fmt.Sprintf("set custom content not found for type %v", bytes[offset]))
		} else {
			err := v.SetCustomContent[fieldId](bytes[offset+offType+offLen : offset+offType+offLen+fieldSize])
			if err != nil {
				panic(err)
			}
		}
		break
	default:
		panic(fmt.Sprintf("can't apply type %v", bytes[offset]))
	}
}

func handleStruct(v *Mapper, bytes []byte, offset int, searchField uint8) int {
	// Type
	offset += 1

	// Size
	offset += 2
	size := int(binary.LittleEndian.Uint16(bytes[1:]))

	// Field amount
	amount := int(bytes[offset])
	offset += 1

	for i := 0; i < amount; i++ {
		// Read field id
		fieldId := bytes[offset]
		offset += 1

		// Field matches
		if fieldId == searchField {
			applyType(v, bytes, offset, fieldId)
			return size
		}

		// Go next
		fieldSize := typeSize(bytes[offset:])
		offset += 1 // field type
		offset += fieldSize
	}

	return size
}

func (v *Mapper) Map(bytes []byte) {
	offset := 0

	for i := 0; i < len(v.FieldIdList); i++ {
		searchField := v.FieldIdList[i]

		if bytes[0] == ml_gosn.T_STRUCT {
			offset += handleStruct(v, bytes, offset, searchField)
		} else {
			panic(fmt.Sprintf("unmapped type %v", bytes[0]))
		}
	}
}
