package mr_mesh

import (
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	ml_image "github.com/maldan/go-ml/util/media/image"
)

type Line struct {
	From  ml_geom.Vector3[float32]
	To    ml_geom.Vector3[float32]
	Color ml_image.ColorRGB[float32]
}
