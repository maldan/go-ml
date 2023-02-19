package ml_gosn

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strings"
)

type IR struct {
	Type    int
	Id      uint8
	Name    string
	Content []byte
	List    []*IR
}

func (r *IR) Len() int {
	outSize := 0

	// Name Id, 0 means value doesn't have named
	// For example number or string doesn't have named, but struct field does
	if r.Id > 0 {
		outSize += 1
	}
	// Same as id. Builder can work in 2 mode. Field with id or fields with name.
	if r.Name != "" {
		outSize += 1
		outSize += len(r.Name)
	}

	// Type
	outSize += 1

	switch r.Type {
	case T_BOOL, T_8:
		outSize += 1
		break
	case T_16:
		outSize += 2
		break
	case T_32, T_F32:
		outSize += 4
		break
	case T_64, T_F64:
		outSize += 8
		break
	case T_STRING:
		outSize += 2
		outSize += len(r.Content)
		break
	case T_SLICE, T_BIG_SLICE:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 2

		if r.Type == T_BIG_SLICE {
			outSize += 2 // + 2 bytes for size
			outSize += 2 // + 2 bytes for amount of elements
		}

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case T_STRUCT:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 1

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case T_CUSTOM:
		outSize += 2
		outSize += len(r.Content)
		break
	case T_BLOB:
		outSize += 4
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

	// Name Id, 0 means value doesn't have named
	// For example number or string doesn't have named, but struct field does
	if r.Id > 0 {
		s = append(s, r.Id)
	}

	// Same as id. Builder can work in 2 mode. Field with id or fields with name.
	if r.Name != "" {
		s = append(s, uint8(len(r.Name)))
		s = append(s, r.Name...)
	}

	// Type
	s = append(s, uint8(r.Type))

	switch r.Type {
	case T_BOOL, T_8, T_16, T_32, T_64, T_F32, T_F64:
		// Content
		s = append(s, r.Content...)
		break
	case T_STRING:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Content
		s = append(s, r.Content...)
		break
	case T_SLICE, T_BIG_SLICE:
		// Size of slice
		l := r.Len()
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// 4 bytes for size
		if r.Type == T_BIG_SLICE {
			s = append(s, uint8(l>>16))
			s = append(s, uint8(l>>24))
		}

		// Amount of elements
		l = len(r.List)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// 4 bytes for amount
		if r.Type == T_BIG_SLICE {
			s = append(s, uint8(l>>16))
			s = append(s, uint8(l>>24))
		}

		// Elements
		for i := 0; i < len(r.List); i++ {
			s = append(s, r.List[i].Build()...)
		}
		break
	case T_STRUCT:
		// Size of struct
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
	case T_CUSTOM:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Content
		s = append(s, r.Content...)
		break
	case T_BLOB:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))
		s = append(s, uint8(l>>16))
		s = append(s, uint8(l>>24))

		// Content
		s = append(s, r.Content...)
		break
	default:
		panic("unknown type " + fmt.Sprintf("%v", r.Type))
	}

	return s
}

// BuildIR convert v any type to IR tree
func BuildIR(ir *IR, v any, nameToId NameToId) {
	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)

	switch typeOf.Kind() {
	case reflect.Bool:
		ir.Type = T_BOOL
		ir.Content = []byte{0}
		if valueOf.Bool() {
			ir.Content[0] = 1
		}
		break

	// Int 8
	case reflect.Int8:
		ir.Type = T_8
		ir.Content = []byte{uint8(valueOf.Int())}
		break
	case reflect.Uint8:
		ir.Type = T_8
		ir.Content = []byte{uint8(valueOf.Uint())}
		break

	// Int 16
	case reflect.Int16:
		b := []byte{0, 0}
		binary.LittleEndian.PutUint16(b, uint16(valueOf.Int()))
		ir.Type = T_16
		ir.Content = b
		break
	case reflect.Uint16:
		b := []byte{0, 0}
		binary.LittleEndian.PutUint16(b, uint16(valueOf.Uint()))
		ir.Type = T_16
		ir.Content = b
		break

	// Int 32
	case reflect.Int, reflect.Int32:
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, uint32(valueOf.Int()))
		ir.Type = T_32
		ir.Content = b
		break
	case reflect.Uint32:
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, uint32(valueOf.Uint()))
		ir.Type = T_32
		ir.Content = b
		break

	// Int 64
	case reflect.Int64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(b, uint64(valueOf.Int()))
		ir.Type = T_64
		ir.Content = b
		break
	case reflect.Uint64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(b, valueOf.Uint())
		ir.Type = T_64
		ir.Content = b
		break

	// Float
	case reflect.Float32:
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, math.Float32bits(float32(valueOf.Float())))
		ir.Type = T_F32
		ir.Content = b
		break
	case reflect.Float64:
		b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint64(b, math.Float64bits(valueOf.Float()))
		ir.Type = T_F64
		ir.Content = b
		break

	case reflect.String:
		ir.Type = T_STRING
		ir.Content = []byte(valueOf.String())

		// Big strings
		if len(ir.Content) > int(uint16(65535)) {
			ir.Type = T_BIG_STRING
		}

		break
	case reflect.Slice:
		ir.Type = T_SLICE

		// Blob types for []byte slice
		if valueOf.Type().Elem().String() == "uint8" {
			ir.Type = T_BLOB
			for i := 0; i < valueOf.Len(); i++ {
				ir.Content = append(ir.Content, uint8(valueOf.Index(i).Uint()))
			}
			break
		}

		// Any slice with any values
		for i := 0; i < valueOf.Len(); i++ {
			tr := IR{}
			ir.List = append(ir.List, &tr)
			BuildIR(&tr, valueOf.Index(i).Interface(), nameToId)
		}

		// Big slice
		if len(ir.List) > 0xFFFF {
			ir.Type = T_BIG_SLICE
		}

		break
	case reflect.Struct:
		// Check to custom encode method
		tb, ok := typeOf.MethodByName("ToBytes")
		if ok {
			// Call
			ret := tb.Func.Call([]reflect.Value{valueOf})
			if len(ret) > 0 {
				ir.Type = T_CUSTOM
				ir.Content = ret[0].Interface().([]byte)

				// Big custom
				if len(ir.Content) > 0xFFFF {
					ir.Type = T_BIG_CUSTOM
				}

				// Check an error
				if len(ret) > 1 {
					if ret[1].Interface() != nil {
						panic(ret[1])
					}
				}
			}
		} else {
			ir.Type = T_STRUCT
			for i := 0; i < typeOf.NumField(); i++ {
				if typeOf.Field(i).Name == strings.ToLower(typeOf.Field(i).Name) {
					panic(fmt.Sprintf("struct %v with private fields impossible to serialize", typeOf))
				}

				tr := IR{}

				// Mode without name id
				if nameToId == nil {
					tr.Name = typeOf.Field(i).Name
				} else {
					tr.Id, ok = nameToId[typeOf.Field(i).Name]
					if !ok {
						panic("name not found")
					}
				}

				ir.List = append(ir.List, &tr)
				BuildIR(&tr, valueOf.Field(i).Interface(), nameToId)
			}
			fmt.Printf("Size of struct: %v\n", ir.Len())
		}
		break
	default:
		panic("unsupported kind " + typeOf.Kind().String())
	}
}
