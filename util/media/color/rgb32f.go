package ml_color

import mmath "github.com/maldan/go-ml/math"

type RGB32F struct {
	R float32
	G float32
	B float32
}

func (c RGB32F) ToRGB8() RGB8 {
	return RGB8{
		R: uint8(mmath.Clamp(c.R*255.0, 0, 255)),
		G: uint8(mmath.Clamp(c.G*255.0, 0, 255)),
		B: uint8(mmath.Clamp(c.B*255.0, 0, 255)),
	}
}
