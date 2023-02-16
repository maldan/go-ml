package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-ml/db/goson/core"
	"log"
	"reflect"
	"unsafe"
)

type ValueMapper[T any] struct {
	Container T
	NameToId  core.NameToId
	MapOffset []unsafe.Pointer

	// Copy bytes and parse it and assign to pointer
	SetCustomContent []func([]byte) error
}

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func NewMapper[T any](nameToId core.NameToId) *ValueMapper[T] {
	mapper := ValueMapper[T]{
		Container:        *new(T),
		NameToId:         nameToId,
		MapOffset:        make([]unsafe.Pointer, 255),
		SetCustomContent: make([]func([]byte) error, 255),
	}

	typeOf := reflect.TypeOf(mapper.Container)
	start := unsafe.Pointer(&mapper.Container)

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
			// Convert unsafe pointer to reflect pointer with type hinting
			// rf := reflect.NewAt(typeOf.Field(i).Type, mapper.MapOffset[fieldId])

			realFn := fromBytes.Func.Interface()
			rfnPtr := unsafe.Pointer(&realFn)
			iface := (*emptyInterface)(rfnPtr)

			// De interface pointer. It means cast any to real type
			*(*emptyInterface)(rfnPtr) = emptyInterface{
				typ: iface.ptr,
			}
			//fmt.Printf("%v\n", realFn)

			// fmt.Printf("%x\n", mapper.MapOffset[fieldId])
			/*hdr := *(*IHeader)(unsafe.Pointer(&rf))
			fmt.Printf("%v\n", hdr)

			realFn := fromBytes.Func.Interface()
			rfnPtr := unsafe.Pointer(&realFn)
			castFn := *(*func(*IHeader, []byte) error)(rfnPtr)
			castFn(&hdr, []byte{1, 2, 3})*/

			//zzz := fromBytes.Func.Interface().(func(unsafe.Pointer, []byte) error)
			//zzz(mapper.MapOffset[fieldId], []byte{1, 2, 3})

			/*zyz := unsafe.Pointer(&zzz)

			zzx := *(*func(IHeader, []byte) error)(zyz)
			zzx(IHeader{}, []byte{1, 2, 3})*/

			//fmt.Printf("%v\n", yy)

			//fnPtr := *(*func(any, []byte) error)(unsafe.Pointer(fromBytes.Func.Pointer()))

			//fnPtr(uintptr(mapper.MapOffset[fieldId]), []byte{1})

			// fnPtr([]byte{})

			// Create callback
			mapper.SetCustomContent[fieldId] = func(bytes []byte) error {
				castFn := *(*func(unsafe.Pointer, []byte) error)(rfnPtr)
				return castFn(mapper.MapOffset[fieldId], bytes)

				/*ret := fromBytes.Func.Call([]reflect.Value{rf, reflect.ValueOf(bytes)})

				if len(ret) > 0 {
					err := ret[0].Interface()
					if err != nil {
						return err.(error)
					} else {
						return nil
					}
				}
				return nil*/
			}
		}
	}

	return &mapper
}

func typeSize(bytes []byte) int {
	switch bytes[0] {
	case core.T_BOOL, core.T_8:
		return 1
	case core.T_16:
		return 2
	case core.T_32, core.TypeF32:
		return 4
	case core.T_64, core.TypeF64:
		return 8
	case core.TypeString, core.T_CUSTOM:
		// Field size
		fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))
		return 2 + fieldSize // 2 - is length size info
	default:
		panic(fmt.Sprintf("unknown type %v", bytes[0]))
		return 0
	}
}

func applyType[T any](v *ValueMapper[T], bytes []byte, offset int, fieldId uint8) {
	off := v.MapOffset[fieldId]

	const offType = 1 // offset type
	const offLen = 2  // offset length of field

	switch bytes[offset] {
	case core.T_8:
		*(*uint8)(off) = bytes[offset]
		break
	case core.T_16:
		*(*uint16)(off) = binary.LittleEndian.Uint16(bytes[offset+offType:])
		break
	case core.T_32:
		*(*uint32)(off) = binary.LittleEndian.Uint32(bytes[offset+offType:])
		break
	case core.T_64:
		*(*uint64)(off) = binary.LittleEndian.Uint64(bytes[offset+offType:])
		break
	case core.TypeString:
		fieldSize := int(binary.LittleEndian.Uint16(bytes[offset+offType:]))
		bts := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))

		hh := (*reflect.StringHeader)(off)
		hh.Data = bts.Data + uintptr(offset) + offType + offLen
		hh.Len = fieldSize
		break
	case core.T_CUSTOM:
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

func handleStruct[T any](v *ValueMapper[T], bytes []byte, offset int, searchField uint8) int {
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

func (v *ValueMapper[T]) Map(bytes []byte, fieldList []uint8) {
	offset := 0

	for i := 0; i < len(fieldList); i++ {
		searchField := fieldList[i]

		if bytes[0] == core.TypeStruct {
			offset += handleStruct[T](v, bytes, offset, searchField)
		} else {
			panic(fmt.Sprintf("unmapped type %v", bytes[0]))
		}
	}
}
