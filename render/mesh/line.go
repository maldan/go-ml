package mrender_mesh

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type Line struct {
	From  mmath_la.Vector3[float32]
	To    mmath_la.Vector3[float32]
	Color ml_color.ColorRGBA[float32]
	IsUi  bool
}
