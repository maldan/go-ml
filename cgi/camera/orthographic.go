package mcgi_camera

import (
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

type OrthographicCamera struct {
	Area mmath_geom.Rectangle[float32]

	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Vector3[float32]
	Scale    mmath_la.Vector3[float32]
	Matrix   mmath_la.Matrix4x4[float32]
}

func (p *OrthographicCamera) ApplyMatrix() {
	p.Matrix.Identity()
	p.Matrix.Orthographic(p.Area.MinX, p.Area.MaxX, p.Area.MaxY, p.Area.MinY, -400, 400)
}
