package mmath_la

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Vector2[T constraints.Float] struct {
	X T
	Y T
}

func (v *Vector2[T]) Clone() Vector2[T] {
	return Vector2[T]{v.X, v.Y}
}

func (v *Vector2[T]) DistanceTo(to Vector2[T]) T {
	a := float64(v.X - to.X)
	b := float64(v.Y - to.Y)

	return T(math.Sqrt(a*a + b*b))
}

func (v *Vector2[T]) AngleBetween(v2 Vector2[T]) T {
	return T(math.Atan2(float64(v2.Y-v.Y), float64(v2.X-v.X)))
}

func (v *Vector2[T]) DirectionFromAngle(rad T) {
	v.X = T(math.Cos(float64(rad)))
	v.Y = T(math.Sin(float64(rad)))
}

func (v *Vector2[T]) ToVector3() Vector3[T] {
	return Vector3[T]{v.X, v.Y, 0}
}
