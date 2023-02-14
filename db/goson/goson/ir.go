package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-ml/db/goson/core"
	"reflect"
	"time"
)

type IR struct {
	Type    int
	Id      uint8
	Content []byte
	List    []*IR
}

func (r *IR) Len() int {
	outSize := 0

	// Name Id, 0 means value doesn't have name
	// For example number or string doesn't have name, but struct field does
	if r.Id > 0 {
		outSize += 1
	}

	// Type
	outSize += 1

	switch r.Type {
	case core.Type8, core.TypeBool:
		outSize += 1
		break
	case core.Type16:
		outSize += 2
		break
	case core.Type32:
		outSize += 4
		break
	case core.Type64:
		outSize += 8
		break
	case core.TypeStruct:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 1

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case core.TypeSlice:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 2

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case core.TypeString:
		outSize += 2
		outSize += len(r.Content)
		break
	case core.TypeTime:
		outSize += 1
		outSize += len(r.Content)
		break
	default:
		panic("unknown type " + fmt.Sprintf("%v", r.Type))
	}

	return outSize
}

// Build converts IR to bytes
func (r *IR) Build() []byte {
	s := make([]byte, 0, r.Len())

	// Name Id, 0 means value doesn't have name
	// For example number or string doesn't have name, but struct field does
	if r.Id > 0 {
		s = append(s, r.Id)
	}

	// Type
	s = append(s, uint8(r.Type))

	switch r.Type {
	case core.TypeBool, core.Type8, core.Type16, core.Type32, core.Type64, core.TypeF32, core.TypeF64:
		// Content
		s = append(s, r.Content...)
		break
	case core.TypeStruct:
		// Len of struct
		l := r.Len()
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Amount of elements
		l = len(r.List)
		s = append(s, uint8(l))

		// Elements
		for i := 0; i < len(r.List); i++ {
			s = append(s, r.List[i].Build()...)
		}
		break
	case core.TypeSlice:
		// Len of struct
		l := r.Len()
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Amount of elements
		l = len(r.List)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Elements
		for i := 0; i < len(r.List); i++ {
			s = append(s, r.List[i].Build()...)
		}
		break
	case core.TypeString:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Content
		s = append(s, r.Content...)
		break
	case core.TypeTime:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))

		// Content
		s = append(s, r.Content...)
		break
	default:
		panic("unknown type " + fmt.Sprintf("%v", r.Type))
	}

	return s
}

// BuildIR convert any type to IR tree
func BuildIR(ir *IR, v any, nameToId core.NameToId) {
	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)

	switch typeOf.Kind() {
	case reflect.Bool:
		ir.Content = []byte{0}
		if valueOf.Bool() {
			ir.Content[0] = 1
		}
		ir.Type = core.TypeBool
		break

	// Int 8
	case reflect.Int8:
		ir.Content = []byte{uint8(valueOf.Int())}
		ir.Type = core.Type8
		break
	case reflect.Uint8:
		ir.Content = []byte{uint8(valueOf.Uint())}
		ir.Type = core.Type8
		break

	// Int 16
	case reflect.Int16:
		b := []byte{0, 0}
		binary.LittleEndian.PutUint16(b, uint16(valueOf.Int()))
		ir.Content = b
		ir.Type = core.Type16
		break
	case reflect.Uint16:
		b := []byte{0, 0}
		binary.LittleEndian.PutUint16(b, uint16(valueOf.Uint()))
		ir.Content = b
		ir.Type = core.Type16
		break

	// Int 32
	case reflect.Int, reflect.Int32:
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, uint32(valueOf.Int()))
		ir.Content = b
		ir.Type = core.Type32
		break
	case reflect.Uint32:
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, uint32(valueOf.Uint()))
		ir.Content = b
		ir.Type = core.Type32
		break

	// Int 64
	case reflect.Int64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(b, uint64(valueOf.Int()))
		ir.Content = b
		ir.Type = core.Type64
		break
	case reflect.Uint64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(b, valueOf.Uint())
		ir.Content = b
		ir.Type = core.Type64
		break

	case reflect.String:
		ir.Type = core.TypeString
		ir.Content = []byte(valueOf.String())
		break
	case reflect.Slice:
		ir.Type = core.TypeSlice

		for i := 0; i < valueOf.Len(); i++ {
			tr := IR{}
			ir.List = append(ir.List, &tr)
			BuildIR(&tr, valueOf.Index(i).Interface(), nameToId)
		}
		break
	case reflect.Struct:
		if typeOf.Name() == "Time" {
			ir.Type = core.TypeTime
			ir.Content = []byte(valueOf.Interface().(time.Time).Format("2006-01-02T15:04:05.999-07:00"))
		} else {
			ir.Type = core.TypeStruct
			for i := 0; i < typeOf.NumField(); i++ {
				id, ok := nameToId[typeOf.Field(i).Name]
				if !ok {
					panic("name not found")
				}
				tr := IR{
					Id: id,
				}
				ir.List = append(ir.List, &tr)
				BuildIR(&tr, valueOf.Field(i).Interface(), nameToId)
			}
		}
		break
	default:
		panic("unsupported kind " + typeOf.Kind().String())
	}
}
