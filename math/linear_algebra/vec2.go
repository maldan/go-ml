package mmath_la

import (
	"encoding/binary"
	mmath "github.com/maldan/go-ml/math"
	"golang.org/x/exp/constraints"
	"math"
)

type Vector2[T constraints.Float] struct {
	X T `json:"x"`
	Y T `json:"y"`
}

func (v Vector2[T]) Add(v2 Vector2[T]) Vector2[T] {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (v Vector2[T]) Sub(v2 Vector2[T]) Vector2[T] {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

func (v Vector2[T]) SubXY(x T, y T) Vector2[T] {
	v.X -= x
	v.Y -= y
	return v
}

func (v *Vector2[T]) Length() T {
	ax := float64(v.X)
	ay := float64(v.Y)
	return T(math.Sqrt((ax * ax) + (ay * ay)))
}

func (v Vector2[T]) Normalize() Vector2[T] {
	l := v.Length()
	if l == 0 {
		return Vector2[T]{}
	}
	v.X /= l
	v.Y /= l
	return v
}

func (v Vector2[T]) Floor() Vector2[T] {
	v.X = T(int(v.X))
	v.Y = T(int(v.Y))
	//v.X = T(math.Floor(float64(v.X)))
	//v.Y = T(math.Floor(float64(v.Y)))
	return v
}

func (v Vector2[T]) Round() Vector2[T] {
	v.X = T(math.Round(float64(v.X)))
	v.Y = T(math.Round(float64(v.Y)))
	return v
}

func (v Vector2[T]) Ceil() Vector2[T] {
	v.X = mmath.Ceil(v.X)
	v.Y = mmath.Ceil(v.Y)
	return v
}

func (v Vector2[T]) Scale(v2 T) Vector2[T] {
	v.X *= v2
	v.Y *= v2
	return v
}

func (v Vector2[T]) AddXY(x T, y T) Vector2[T] {
	v.X += x
	v.Y += y
	return v
}

func (v Vector2[T]) MulXY(x T, y T) Vector2[T] {
	v.X *= x
	v.Y *= y
	return v
}

func (v *Vector2[T]) Clone() Vector2[T] {
	return Vector2[T]{v.X, v.Y}
}

func (v Vector2[T]) DistanceTo(to Vector2[T]) T {
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

func (v Vector2[T]) ToVector3XY() Vector3[T] {
	return Vector3[T]{v.X, v.Y, 0}
}

func (v *Vector2[T]) ToVector3XZ() Vector3[T] {
	return Vector3[T]{v.X, 0, v.Y}
}

func (v Vector2[T]) FromBytes(data []byte) Vector2[T] {
	v.X = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0:4])))
	v.Y = T(math.Float32frombits(binary.LittleEndian.Uint32(data[0+4 : 4+4])))
	return v
}
