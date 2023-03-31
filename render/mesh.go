package ml_render

import ml_geom "github.com/maldan/go-ml/util/math/geom"

type Mesh struct {
	Vertices []ml_geom.Vector3[float32]
	Indices  []uint16
	UV       []ml_geom.Vector2[float32]

	Matrix ml_geom.Matrix4x4[float32]

	Position ml_geom.Vector3[float32]
	Rotation ml_geom.Vector3[float32]
	Scale    ml_geom.Vector3[float32]
}

func (m *Mesh) ApplyMatrix() {
	m.Matrix.Identity()
	m.Matrix.Translate(m.Position)
	m.Matrix.RotateX(m.Rotation.X)
	m.Matrix.RotateY(m.Rotation.Y)
	m.Matrix.RotateZ(m.Rotation.Z)
}
