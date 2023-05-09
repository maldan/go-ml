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

func Lerp[T constraints.Integer | constraints.Float](from ColorRGBA[T], to ColorRGBA[T], t T) ColorRGBA[T] {
	return ColorRGBA[T]{
		A: ml_number.Lerp(from.A, to.A, t),
		R: ml_number.Lerp(from.R, to.R, t),
		G: ml_number.Lerp(from.G, to.G, t),
		B: ml_number.Lerp(from.B, to.B, t),
	}
}
