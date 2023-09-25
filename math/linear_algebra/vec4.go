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

func (v Vector4[T]) Sub(v2 Vector4[T]) Vector4[T] {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	v.W -= v2.W
	return v
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

func (v Vector4[T]) TransformMatrix4x4(m Matrix4x4[T]) Vector4[T] {
	result := Vector4[T]{}

	result.X = v.X*m.Raw[0] + v.Y*m.Raw[4] + v.Z*m.Raw[8] + v.W*m.Raw[12]
	result.Y = v.X*m.Raw[1] + v.Y*m.Raw[5] + v.Z*m.Raw[9] + v.W*m.Raw[13]
	result.Z = v.X*m.Raw[2] + v.Y*m.Raw[6] + v.Z*m.Raw[10] + v.W*m.Raw[14]
	result.W = v.X*m.Raw[3] + v.Y*m.Raw[7] + v.Z*m.Raw[11] + v.W*m.Raw[15]

	return result
}

func (v Vector4[T]) MultiplyMatrix4x4(m Matrix4x4[T]) Vector4[T] {
	result := Vector4[T]{}

	result.X = v.X*m.Raw[0] + v.Y*m.Raw[4] + v.Z*m.Raw[8] + v.W*m.Raw[12]
	result.Y = v.X*m.Raw[1] + v.Y*m.Raw[5] + v.Z*m.Raw[9] + v.W*m.Raw[13]
	result.Z = v.X*m.Raw[2] + v.Y*m.Raw[6] + v.Z*m.Raw[10] + v.W*m.Raw[14]
	result.W = v.X*m.Raw[3] + v.Y*m.Raw[7] + v.Z*m.Raw[11] + v.W*m.Raw[15]

	return result
}

func (v Vector4[T]) ToVector3XYZ() Vector3[T] {
	return Vector3[T]{v.X, v.Y, v.Z}
}
