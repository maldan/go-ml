package ml_geom

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Vector3[T constraints.Float] struct {
	X T
	Y T
	Z T
}

func (v *Vector3[T]) TransformMatrix4x4(mx Matrix4x4[T]) {
	x := v.X
	y := v.Y
	z := v.Z
	w := mx.Raw[3]*x + mx.Raw[7]*y + mx.Raw[11]*z + mx.Raw[15]
	if w == 0 {
		w = 1
	}
	v.X = (mx.Raw[0]*x + mx.Raw[4]*y + mx.Raw[8]*z + mx.Raw[12]) / w
	v.Y = (mx.Raw[1]*x + mx.Raw[5]*y + mx.Raw[9]*z + mx.Raw[13]) / w
	v.Z = (mx.Raw[2]*x + mx.Raw[6]*y + mx.Raw[10]*z + mx.Raw[14]) / w
}

func (v *Vector3[T]) Clone() Vector3[T] {
	return Vector3[T]{v.X, v.Y, v.Z}
}

func (v *Vector3[T]) ToVector2() Vector2[T] {
	return Vector2[T]{v.X, v.Y}
}

func (v *Vector3[T]) DistanceTo(to Vector3[T]) T {
	a := float64(v.X - to.X)
	b := float64(v.Y - to.Y)
	c := float64(v.Z - to.Z)

	return T(math.Sqrt(a*a + b*b + c*c))
}
