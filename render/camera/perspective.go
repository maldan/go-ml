package mr_camera

import ml_geom "github.com/maldan/go-ml/util/math/geom"

type PerspectiveCamera struct {
	Fov         float32
	AspectRatio float32
	Position    ml_geom.Vector3[float32]
	Rotation    ml_geom.Vector3[float32]
	Scale       ml_geom.Vector3[float32]
	Matrix      ml_geom.Matrix4x4[float32]
}

func (p *PerspectiveCamera) ApplyMatrix() {
	p.Matrix.Identity()
	p.Matrix.Perspective((p.Fov*3.141592653589793)/180, p.AspectRatio, 0.1, 1000.0)
	p.Matrix.Translate(p.Position)
	p.Matrix.Scale(p.Scale)
}
