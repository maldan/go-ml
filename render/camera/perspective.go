package mrender_camera

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

type PerspectiveCamera struct {
	Fov         float32
	AspectRatio float32
	Position    mmath_la.Vector3[float32]
	Rotation    mmath_la.Vector3[float32]
	Scale       mmath_la.Vector3[float32]
	Matrix      mmath_la.Matrix4x4[float32]
}

func (p *PerspectiveCamera) ApplyMatrix() {
	p.Matrix.Identity()

	p.Matrix.Perspective((p.Fov*3.141592653589793)/180, p.AspectRatio, 0.1, 1000.0)
	p.Matrix.Translate(p.Position)
	p.Matrix.Scale(p.Scale)
}
