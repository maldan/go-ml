package ml_color

import (
	ml_number "github.com/maldan/go-ml/util/number"
	"golang.org/x/exp/constraints"
)

type ColorRGB[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
}

type ColorRGBA[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
	A T
}

func (c *ColorRGBA[T]) SetRGB(r T, g T, b T) {
	c.R = r
	c.G = g
	c.B = b
}

func (c ColorRGBA[T]) AddRGB(r T, g T, b T) ColorRGBA[T] {
	c.R += r
	c.G += g
	c.B += b
	return c
}

func (c ColorRGBA[T]) Mix(c2 ColorRGBA[T], t float32) ColorRGBA[T] {
	c.R = T(ml_number.Clamp(float32(c.R)*(1.0-t)+float32(c2.R)*t, 0, 255))
	c.G = T(ml_number.Clamp(float32(c.G)*(1.0-t)+float32(c2.G)*t, 0, 255))
	c.B = T(ml_number.Clamp(float32(c.B)*(1.0-t)+float32(c2.B)*t, 0, 255))
	return c
}

func (c ColorRGBA[T]) MulF32(v float32) ColorRGBA[T] {
	c.R = T(ml_number.Clamp(float32(c.R)*v, 0, 255))
	c.G = T(ml_number.Clamp(float32(c.G)*v, 0, 255))
	c.B = T(ml_number.Clamp(float32(c.B)*v, 0, 255))
	return c
}

func (c ColorRGBA[T]) To01() ColorRGBA[float32] {
	c2 := ColorRGBA[float32]{}
	c2.R = float32(c.R) / 255.0
	c2.G = float32(c.G) / 255.0
	c2.B = float32(c.B) / 255.0
	return c2
}

func (c ColorRGBA[T]) To255() ColorRGBA[uint8] {
	c2 := ColorRGBA[uint8]{}
	c2.R = uint8(ml_number.Clamp(float32(c.R)*255.0, 0, 255))
	c2.G = uint8(ml_number.Clamp(float32(c.G)*255.0, 0, 255))
	c2.B = uint8(ml_number.Clamp(float32(c.B)*255.0, 0, 255))
	return c2
}

func Lerp[T constraints.Integer | constraints.Float](from ColorRGBA[T], to ColorRGBA[T], t T) ColorRGBA[T] {
	return ColorRGBA[T]{
		A: ml_number.Lerp(from.A, to.A, t),
		R: ml_number.Lerp(from.R, to.R, t),
		G: ml_number.Lerp(from.G, to.G, t),
		B: ml_number.Lerp(from.B, to.B, t),
	}
}
