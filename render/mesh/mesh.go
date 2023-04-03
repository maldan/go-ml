package mr_mesh

import (
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	ml_color "github.com/maldan/go-ml/util/media/color"
	ml_number "github.com/maldan/go-ml/util/number"
)

type Mesh struct {
	Id       int
	Vertices []ml_geom.Vector3[float32]
	Indices  []uint16
	UV       []ml_geom.Vector2[float32]
	Normal   []ml_geom.Vector3[float32]

	/*Matrix ml_geom.Matrix4x4[float32]

	Position ml_geom.Vector3[float32]
	Rotation ml_geom.Vector3[float32]
	Scale    ml_geom.Vector3[float32]

	*/
}

type MeshInstance struct {
	Id int

	// Matrix ml_geom.Matrix4x4[float32]

	Position ml_geom.Vector3[float32]
	Rotation ml_geom.Vector3[float32]
	Scale    ml_geom.Vector3[float32]
	UvOffset ml_geom.Vector2[float32]
	Color    ml_color.ColorRGBA[float32]

	IsVisible bool
	IsActive  bool
}

func New() *Mesh {
	return &Mesh{
		Id: -1,
	}
}

/*func NewInstance() MeshInstance {
	return MeshInstance{
		Id:        -1,
		Scale:     ml_geom.Vector3[float32]{1, 1, 1},
		IsVisible: true,
	}
}

func NewInstanceOf(mesh *Mesh) MeshInstance {
	return MeshInstance{
		Id:        mesh.Id,
		Scale:     ml_geom.Vector3[float32]{1, 1, 1},
		IsVisible: true,
	}
}*/

/*func (m *Mesh) ApplyMatrix() {
	m.Matrix.Identity()
	m.Matrix.Translate(m.Position)
	m.Matrix.RotateX(m.Rotation.X)
	m.Matrix.RotateY(m.Rotation.Y)
	m.Matrix.RotateZ(m.Rotation.Z)
}*/

func (m *Mesh) ScaleUV(size ml_geom.Vector2[float32]) {
	for i := 0; i < len(m.UV); i++ {
		m.UV[i].X *= size.X
		m.UV[i].Y *= size.Y
	}
}

// MakeCube 0b11_11_11_00 [Front, Back, Top, Bottom, Right, Left]
func (m *Mesh) MakeCube(size ml_geom.Vector3[float32], side uint8) {
	sideAmount := ml_number.CountSetBits(side)

	m.Vertices = make([]ml_geom.Vector3[float32], 0, 4*sideAmount)
	m.Normal = make([]ml_geom.Vector3[float32], 0, 4*sideAmount)

	// Front
	if side&0b1000_0000 == 0b1000_0000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{0.0, 0.0, 1.0})
		}
	}

	// Back
	if side&0b0100_0000 == 0b0100_0000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{0.0, 0.0, -1.0})
		}
	}

	// Top
	if side&0b0010_0000 == 0b0010_0000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{0.0, 1.0, 0.0})
		}
	}

	// Bottom
	if side&0b0001_0000 == 0b0001_0000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{0.0, -1.0, 0.0})
		}
	}

	// Right
	if side&0b0000_1000 == 0b0000_1000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{1.0, -1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{1.0, 0.0, 0.0})
		}
	}

	// Left
	if side&0b0000_1000 == 0b0000_1000 {
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, ml_geom.Vector3[float32]{-1.0, 1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, ml_geom.Vector3[float32]{-1.0, 0.0, 0.0})
		}
	}

	// UV
	m.UV = make([]ml_geom.Vector2[float32], 0, 4*sideAmount)
	for i := 0; i < int(sideAmount); i++ {
		m.UV = append(m.UV, ml_geom.Vector2[float32]{0.0, 0.0})
		m.UV = append(m.UV, ml_geom.Vector2[float32]{1.0, 0.0})
		m.UV = append(m.UV, ml_geom.Vector2[float32]{1.0, 1.0})
		m.UV = append(m.UV, ml_geom.Vector2[float32]{0.0, 1.0})
	}

	for i := 0; i < len(m.Vertices); i++ {
		m.Vertices[i].X *= size.X
		m.Vertices[i].Y *= size.Y
		m.Vertices[i].Z *= size.Z
	}

	m.Indices = make([]uint16, 0, 6*sideAmount)
	for i := 0; i < int(sideAmount); i++ {
		next := uint16(i * 4)
		m.Indices = append(m.Indices, next, 1+next, 2+next, next, 2+next, 3+next)
	}
}
