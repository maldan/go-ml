package mcgi_color

import mmath "github.com/maldan/go-ml/math"

type RGBA32F struct {
	R float32
	G float32
	B float32
	A float32
}

func (c RGBA32F) ToRGB8() RGB8 {
	return RGB8{
		R: uint8(mmath.Clamp(c.R*255.0, 0, 255)),
		G: uint8(mmath.Clamp(c.G*255.0, 0, 255)),
		B: uint8(mmath.Clamp(c.B*255.0, 0, 255)),
	}
}

func (c RGBA32F) ToRGBA8() RGBA8 {
	return RGBA8{
		R: uint8(mmath.Clamp(c.R*255.0, 0, 255)),
		G: uint8(mmath.Clamp(c.G*255.0, 0, 255)),
		B: uint8(mmath.Clamp(c.B*255.0, 0, 255)),
		A: uint8(mmath.Clamp(c.A*255.0, 0, 255)),
	}
}
