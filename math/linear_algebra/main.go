package mmath_la

import "golang.org/x/exp/constraints"

func Vec2f[T constraints.Integer | constraints.Float](x T, y T) Vector2[float32] {
	return Vector2[float32]{X: float32(x), Y: float32(y)}
}

func Vec2fZero() Vector2[float32] {
	return Vector2[float32]{}
}

func Vec2d(x float64, y float64) Vector2[float64] {
	return Vector2[float64]{X: x, Y: y}
}

func Vec3f[T constraints.Integer | constraints.Float](x T, y T, z T) Vector3[float32] {
	return Vector3[float32]{X: float32(x), Y: float32(y), Z: float32(z)}
}
