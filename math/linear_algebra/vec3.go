package mmath_la

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Vector3[T constraints.Float] struct {
	X T
	Y T
	Z T
}

func (v Vector3[T]) TransformMatrix4x4(mx Matrix4x4[T]) Vector3[T] {
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

	return v
}

func (v *Vector3[T]) Clone() Vector3[T] {
	return Vector3[T]{v.X, v.Y, v.Z}
}

func (v *Vector3[T]) Length() T {
	ax := float64(v.X)
	ay := float64(v.Y)
	az := float64(v.Z)
	return T(math.Sqrt((ax * ax) + (ay * ay) + (az * az)))
}

func (v Vector3[T]) Normalize() Vector3[T] {
	l := v.Length()
	v.X /= l
	v.Y /= l
	v.Z /= l
	return v
}

func (v Vector3[T]) Dot(v2 Vector3[T]) T {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v Vector3[T]) Cross(v2 Vector3[T]) Vector3[T] {
	out := Vector3[T]{}
	out.X = v.Y*v2.Z - v.Z*v2.Y
	out.Y = v.Z*v2.X - v.X*v2.Z
	out.Z = v.X*v2.Y - v.Y*v2.X
	return out
}

/*
func (v Vector3[T]) Cross(v2 Vector3[T]) T {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}
*/

func (v Vector3[T]) Divide(v2 T) Vector3[T] {
	v.X /= v2
	v.Y /= v2
	v.Z /= v2
	return v
}

func (v Vector3[T]) Scale(v2 T) Vector3[T] {
	v.X *= v2
	v.Y *= v2
	v.Z *= v2
	return v
}

func (v Vector3[T]) Sin(v2 T) Vector3[T] {
	v.X = T(math.Sin(float64(v2)))
	v.Y = T(math.Sin(float64(v2)))
	v.Z = T(math.Sin(float64(v2)))
	return v
}

func (v Vector3[T]) Abs() Vector3[T] {
	if v.X < 0 {
		v.X *= -1
	}
	if v.Y < 0 {
		v.Y *= -1
	}
	if v.Z < 0 {
		v.Z *= -1
	}
	return v
}

func (v Vector3[T]) Add(v2 Vector3[T]) Vector3[T] {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

func (v Vector3[T]) AddXYZ(x T, y T, z T) Vector3[T] {
	v.X += x
	v.Y += y
	v.Z += z
	return v
}

func (v Vector3[T]) Mul(v2 Vector3[T]) Vector3[T] {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z
	return v
}

func (v Vector3[T]) Sub(v2 Vector3[T]) Vector3[T] {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

func (v Vector3[T]) SubXYZ(x T, y T, z T) Vector3[T] {
	v.X -= x
	v.Y -= y
	v.Z -= z
	return v
}

func (v Vector3[T]) ToVector2() Vector2[T] {
	return Vector2[T]{v.X, v.Y}
}

func (v *Vector3[T]) DistanceTo(to Vector3[T]) T {
	a := float64(v.X - to.X)
	b := float64(v.Y - to.Y)
	c := float64(v.Z - to.Z)

	return T(math.Sqrt(a*a + b*b + c*c))
}

func (v Vector3[T]) DirectionXZToAngle() T {
	return T(math.Atan2(float64(v.X), float64(-v.Z)))
}

/*func AngleBetweenPoints[T constraints.Float](p1 Vector2[T], p2 T) T {
	rad := math.Atan2(p2.Y-p1.Y, p2.X-p1.X)
	return T(rad)
}*/
