package mrender_uv

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

func GetArea(x float32, y float32, sizeX float32, sizeY float32, textureWidth float32, textureHeight float32) []mmath_la.Vector2[float32] {
	x1 := x / textureWidth
	x2 := (x + sizeX) / textureWidth
	y1 := y / textureHeight
	y2 := (y + sizeY) / textureHeight

	return []mmath_la.Vector2[float32]{
		{x1, y1},
		{x2, y1},
		{x2, y2},
		{x1, y2},
	}
}

func GetOffset(x float32, y float32, textureWidth float32, textureHeight float32) mmath_la.Vector2[float32] {
	return mmath_la.Vector2[float32]{
		x / textureWidth,
		y / textureHeight,
	}
}
