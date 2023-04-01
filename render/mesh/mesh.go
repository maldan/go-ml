package mr_mesh

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

func (m *Mesh) MakeCube(position ml_geom.Vector3[float32], size float32) {
	m.Vertices = []ml_geom.Vector3[float32]{
		// Front
		{-1.0, -1.0, 1.0},
		{1.0, -1.0, 1.0},
		{1.0, 1.0, 1.0},
		{-1.0, 1.0, 1.0},

		// Back
		{-1.0, -1.0, -1.0},
		{-1.0, 1.0, -1.0},
		{1.0, 1.0, -1.0},
		{1.0, -1.0, -1.0},

		// Top
		{-1.0, 1.0, -1.0},
		{-1.0, 1.0, 1.0},
		{1.0, 1.0, 1.0},
		{1.0, 1.0, -1.0},

		// Bottom
		{-1.0, -1.0, -1.0},
		{1.0, -1.0, -1.0},
		{1.0, -1.0, 1.0},
		{-1.0, -1.0, 1.0},

		// Right
		{1.0, -1.0, -1.0},
		{1.0, 1.0, -1.0},
		{1.0, 1.0, 1.0},
		{1.0, -1.0, 1.0},

		// Left
		{-1.0, -1.0, -1.0},
		{-1.0, -1.0, 1.0},
		{-1.0, 1.0, 1.0},
		{-1.0, 1.0, -1.0},
	}

	/*for i := 0; i < len(m.Vertices); i++ {
		m.Vertices[i].X *= size
		m.Vertices[i].Y *= size
		m.Vertices[i].Z *= size
	}*/

	m.Indices = []uint16{
		0,
		1,
		2,
		0,
		2,
		3, // front
		4,
		5,
		6,
		4,
		6,
		7, // back
		8,
		9,
		10,
		8,
		10,
		11, // top
		12,
		13,
		14,
		12,
		14,
		15, // bottom
		16,
		17,
		18,
		16,
		18,
		19, // right
		20,
		21,
		22,
		20,
		22,
		23, // left
	}
}
