package mr_uv

import ml_geom "github.com/maldan/go-ml/util/math/geom"

func GetArea(x float32, y float32, sizeX float32, sizeY float32, textureWidth float32, textureHeight float32) []ml_geom.Vector2[float32] {
	x1 := x / textureWidth
	x2 := (x + sizeX) / textureWidth
	y1 := y / textureHeight
	y2 := (y + sizeY) / textureHeight

	return []ml_geom.Vector2[float32]{
		{x1, y1},
		{x2, y1},
		{x2, y2},
		{x1, y2},
	}
}
