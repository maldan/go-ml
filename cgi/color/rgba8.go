package mcgi_color

import mmath "github.com/maldan/go-ml/math"

type RGBA8 struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func (c RGBA8) Black() RGBA8 {
	return RGBA8{R: 0, G: 0, B: 0, A: 0}
}

func (c RGBA8) White() RGBA8 {
	return RGBA8{R: 255, G: 255, B: 255, A: 255}
}

func (c RGBA8) ToRGB32F() RGB32F {
	return RGB32F{
		R: mmath.Clamp01(float32(c.R) / 255.0),
		G: mmath.Clamp01(float32(c.G) / 255.0),
		B: mmath.Clamp01(float32(c.B) / 255.0),
	}
}

func (c RGBA8) Lerp(to RGBA8, t float32) RGBA8 {
	t = mmath.Clamp01(t)
	return RGBA8{
		A: uint8(mmath.Lerp(float32(c.A), float32(to.A), t)),
		R: uint8(mmath.Lerp(float32(c.R), float32(to.R), t)),
		G: uint8(mmath.Lerp(float32(c.G), float32(to.G), t)),
		B: uint8(mmath.Lerp(float32(c.B), float32(to.B), t)),
	}
}

func (c RGBA8) MulScalar(s float32) RGBA8 {
	c.R = uint8(mmath.Clamp(float32(c.R)*s, 0, 255))
	c.G = uint8(mmath.Clamp(float32(c.G)*s, 0, 255))
	c.B = uint8(mmath.Clamp(float32(c.B)*s, 0, 255))
	c.A = uint8(mmath.Clamp(float32(c.A)*s, 0, 255))
	return c
}
