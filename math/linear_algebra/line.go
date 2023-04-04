package mmath_la

import "golang.org/x/exp/constraints"

type Line[T constraints.Float, C Vector2[T] | Vector3[T]] struct {
	From C
	To   C
}
