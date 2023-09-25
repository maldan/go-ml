package ml_color

import (
	mmath "github.com/maldan/go-ml/math"
	ml_random "github.com/maldan/go-ml/util/random"
)

type RGB8 struct {
	R uint8
	G uint8
	B uint8
}

func (c RGB8) White() RGB8 {
	return RGB8{R: 255, G: 255, B: 255}
}

func (c RGB8) Red() RGB8 {
	return RGB8{R: 255, G: 0, B: 0}
}

func (c RGB8) Green() RGB8 {
	return RGB8{R: 0, G: 255, B: 0}
}

func (c RGB8) Blue() RGB8 {
	return RGB8{R: 0, G: 0, B: 255}
}

func (c RGB8) Random() RGB8 {
	return RGB8{
		R: uint8(ml_random.RangeInt(0, 255)),
		G: uint8(ml_random.RangeInt(0, 255)),
		B: uint8(ml_random.RangeInt(0, 255)),
	}
}

func (c RGB8) ToRGB32F() RGB32F {
	return RGB32F{
		R: mmath.Clamp(float32(c.R)/255.0, 0, 1),
		G: mmath.Clamp(float32(c.G)/255.0, 0, 1),
		B: mmath.Clamp(float32(c.B)/255.0, 0, 1),
	}
}

func (c RGB8) Add(to RGB8) RGB8 {
	c.R = uint8(mmath.Clamp(int(c.R)+int(to.R), 0, 255))
	c.G = uint8(mmath.Clamp(int(c.G)+int(to.G), 0, 255))
	c.B = uint8(mmath.Clamp(int(c.B)+int(to.B), 0, 255))
	return c
}

func (c RGB8) Sub(to RGB8) RGB8 {
	c.R = uint8(mmath.Clamp(int(c.R)-int(to.R), 0, 255))
	c.G = uint8(mmath.Clamp(int(c.G)-int(to.G), 0, 255))
	c.B = uint8(mmath.Clamp(int(c.B)-int(to.B), 0, 255))
	return c
}

func (c RGB8) MulScalar(s float32) RGB8 {
	c.R = uint8(mmath.Clamp(float32(c.R)*s, 0, 255))
	c.G = uint8(mmath.Clamp(float32(c.G)*s, 0, 255))
	c.B = uint8(mmath.Clamp(float32(c.B)*s, 0, 255))
	return c
}

func (c RGB8) Mix(c2 RGB8, t float32) RGB8 {
	t = mmath.Clamp(t, 0, 1)
	c.R = uint8(mmath.Clamp(float32(c.R)*(1.0-t)+float32(c2.R)*t, 0, 255))
	c.G = uint8(mmath.Clamp(float32(c.G)*(1.0-t)+float32(c2.G)*t, 0, 255))
	c.B = uint8(mmath.Clamp(float32(c.B)*(1.0-t)+float32(c2.B)*t, 0, 255))
	return c
}
