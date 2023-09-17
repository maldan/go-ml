package mmath_la

import (
	"encoding/binary"
	"golang.org/x/exp/constraints"
	"math"
)

type Vector4[T constraints.Float] struct {
	X T
	Y T
	Z T
	W T
}

func (v *Vector4[T]) Clone() Vector4[T] {
	return Vector4[T]{v.X, v.Y, v.Z, v.W}
}

func (v Vector4[T]) FromBytes(data []byte) Vector4[T] {
	v.X = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0:4])))
	v.Y = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0+4 : 4+4])))
	v.Z = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0+4+4 : 4+4+4])))
	v.W = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0+4+4+4 : 4+4+4+4])))
	return v
}
