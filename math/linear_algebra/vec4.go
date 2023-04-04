package mmath_la

import "golang.org/x/exp/constraints"

type Vector4[T constraints.Float] struct {
	X T
	Y T
	Z T
	W T
}

func (v *Vector4[T]) Clone() Vector4[T] {
	return Vector4[T]{v.X, v.Y, v.Z, v.W}
}
