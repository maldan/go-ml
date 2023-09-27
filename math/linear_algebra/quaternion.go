package mmath_la

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Quaternion[T constraints.Float] struct {
	X T
	Y T
	Z T
	W T
}

func (q Quaternion[T]) Identity() Quaternion[T] {
	return Quaternion[T]{0, 0, 0, 1}
}

func (q Quaternion[T]) Mul(q2 Quaternion[T]) Quaternion[T] {
	x1 := q.X
	y1 := q.Y
	z1 := q.Z
	w1 := q.W

	x2 := q2.X
	y2 := q2.Y
	z2 := q2.Z
	w2 := q2.W

	return Quaternion[T]{
		X: w1*x2 + x1*w2 + y1*z2 - z1*y2,
		Y: w1*y2 + y1*w2 + z1*x2 - x1*z2,
		Z: w1*z2 + z1*w2 + x1*y2 - y1*x2,
		W: w1*w2 - x1*x2 - y1*y2 - z1*z2,
	}
}

func (q Quaternion[T]) FromEuler(v Vector3[T]) Quaternion[T] {
	_x := v.X * 0.5
	_y := v.Y * 0.5
	_z := v.Z * 0.5

	c_x := T(math.Cos(float64(_x)))
	c_y := T(math.Cos(float64(_y)))
	c_z := T(math.Cos(float64(_z)))

	s_x := T(math.Sin(float64(_x)))
	s_y := T(math.Sin(float64(_y)))
	s_z := T(math.Sin(float64(_z)))

	return Quaternion[T]{
		W: c_x*c_y*c_z - s_x*s_y*s_z,
		X: c_y*c_z*s_x + c_x*s_y*s_z,
		Y: c_x*c_z*s_y - c_y*s_x*s_z,
		Z: c_x*c_y*s_z + c_z*s_x*s_y,
	}
}

func (q Quaternion[T]) Lerp(b Quaternion[T], t T) Quaternion[T] {
	result := Quaternion[T]{}
	tInv := 1.0 - t

	// Linear interpolation for the quaternion components
	result.W = q.W*tInv + b.W*t
	result.X = q.X*tInv + b.X*t
	result.Y = q.Y*tInv + b.Y*t
	result.Z = q.Z*tInv + b.Z*t

	// Normalize the resulting quaternion
	norm := T(math.Sqrt(float64(result.W*result.W + result.X*result.X + result.Y*result.Y + result.Z*result.Z)))
	result.W /= norm
	result.X /= norm
	result.Y /= norm
	result.Z /= norm

	return result
}

func (q Quaternion[T]) Inverse() Quaternion[T] {
	w := q.W
	x := q.X
	y := q.Y
	z := q.Z

	normSq := w*w + x*x + y*y + z*z
	if normSq == 0.0 {
		return Quaternion[T]{}
	}

	normSq = 1.0 / normSq
	return Quaternion[T]{
		X: -x * normSq,
		Y: -y * normSq,
		Z: -z * normSq,
		W: w * normSq,
	}
}

func (q Quaternion[T]) ToMatrix4x4() Matrix4x4[T] {
	w := q.W
	x := q.X
	y := q.Y
	z := q.Z

	x2 := x + x
	y2 := y + y
	z2 := z + z
	xx := x * x2
	xy := x * y2
	xz := x * z2
	yy := y * y2
	yz := y * z2
	zz := z * z2
	wx := w * x2
	wy := w * y2
	wz := w * z2

	mx := Matrix4x4[T]{}

	mx.Raw[0] = 1.0 - (yy + zz)
	mx.Raw[4] = xy - wz
	mx.Raw[8] = xz + wy

	mx.Raw[1] = xy + wz
	mx.Raw[5] = 1.0 - (xx + zz)
	mx.Raw[9] = yz - wx

	mx.Raw[2] = xz - wy
	mx.Raw[6] = yz + wx
	mx.Raw[10] = 1.0 - (xx + yy)

	// last column
	mx.Raw[3] = 0.0
	mx.Raw[7] = 0.0
	mx.Raw[11] = 0.0

	// bottom row
	mx.Raw[12] = 0.0
	mx.Raw[13] = 0.0
	mx.Raw[14] = 0.0
	mx.Raw[15] = 1.0

	return mx
}

/*func (q Quaternion[T]) ToEulerDeg() Vector3[T] {
	qq := q.ToEulerRad()
	return Vector3[T]{
		X: mmath.RadToDeg(qq.X),
		Y: mmath.RadToDeg(qq.Y),
		Z: mmath.RadToDeg(qq.Z),
	}
}*/

func (q Quaternion[T]) ToEuler() Vector3[T] {
	t := 2.0 * (q.W*q.Y - q.Z*q.X)
	v := Vector3[T]{}

	// Set X
	a := 2.0 * (q.W*q.X + q.Y*q.Z)
	v.X = T(math.Atan2(float64(a), float64(1.0-2.0*(q.X*q.X+q.Y*q.Y))))

	// Set Y
	if t >= 1.0 {
		v.Y = T(math.Pi / 2.0)
	} else {
		if t <= -1.0 {
			v.Y = T(-math.Pi / 2.0)
		} else {
			v.Y = T(math.Asin(float64(t)))
		}
	}

	// Set Z
	a = 2.0 * (q.W*q.Z + q.X*q.Y)
	v.Z = T(math.Atan2(float64(a), float64(1.0-2.0*(q.Y*q.Y+q.Z*q.Z))))

	return v
}

func (q Quaternion[T]) ToVector() Vector4[T] {
	return Vector4[T]{X: q.X, Y: q.Y, Z: q.Z, W: q.W}
}

func (q Quaternion[T]) Normalize() Quaternion[T] {
	magnitude := q.Magnitude()
	return Quaternion[T]{
		X: q.X / magnitude,
		Y: q.Y / magnitude,
		Z: q.Z / magnitude,
		W: q.W / magnitude,
	}
}

func (q Quaternion[T]) Magnitude() T {
	a := math.Pow(float64(q.X), 2) + math.Pow(float64(q.Y), 2) + math.Pow(float64(q.Z), 2) + math.Pow(float64(q.W), 2)
	return T(math.Sqrt(a))
}

func (q Quaternion[T]) Difference(q2 Quaternion[T]) Quaternion[T] {
	return q.Mul(q2.Inverse())
}
