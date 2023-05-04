package mrender_mesh

import (
	"encoding/binary"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_color "github.com/maldan/go-ml/util/media/color"
	ml_number "github.com/maldan/go-ml/util/number"
	"math"
)

type Mesh struct {
	Id       int
	Vertices []mmath_la.Vector3[float32]
	Indices  []uint16
	UV       []mmath_la.Vector2[float32]
	Normal   []mmath_la.Vector3[float32]
	Color    []ml_color.ColorRGBA[float32]
}

type MeshInstance struct {
	Id int

	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Vector3[float32]
	Scale    mmath_la.Vector3[float32]
	UvOffset mmath_la.Vector2[float32]
	Color    ml_color.ColorRGBA[float32]

	IsVisible bool
	IsActive  bool
}

func New() Mesh {
	return Mesh{
		Id: -1,
	}
}

func BytesToMesh(bytes []byte) Mesh {
	offset := 0
	// fmt.Printf("%v\n", len(bytes))
	out := New()

	// Read vertex
	amount := binary.LittleEndian.Uint16(bytes[0:2])
	offset += 2
	for i := 0; i < int(amount); i++ {
		x := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4
		y := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4
		z := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4

		out.Vertices = append(out.Vertices, mmath_la.Vector3[float32]{x, y, z})
	}

	// Read normal
	amount = binary.LittleEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	for i := 0; i < int(amount); i++ {
		x := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4
		y := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4
		z := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4

		out.Normal = append(out.Normal, mmath_la.Vector3[float32]{x, y, z})
	}

	// Read index
	amount = binary.LittleEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	for i := 0; i < int(amount); i++ {
		index := binary.LittleEndian.Uint16(bytes[offset : offset+2])
		offset += 2
		out.Indices = append(out.Indices, index)
	}

	// Read uv
	amount = binary.LittleEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	for i := 0; i < int(amount); i++ {
		x := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4
		y := math.Float32frombits(binary.LittleEndian.Uint32(bytes[offset : offset+4]))
		offset += 4

		out.UV = append(out.UV, mmath_la.Vector2[float32]{x, y})
	}

	return out
}

/*func (m *Mesh) ApplyMatrix() {
	m.Matrix.Identity()
	m.Matrix.Translate(m.Position)
	m.Matrix.RotateX(m.Rotation.X)
	m.Matrix.RotateY(m.Rotation.Y)
	m.Matrix.RotateZ(m.Rotation.Z)
}*/

func (m *Mesh) SetPosition(position mmath_la.Vector3[float32]) {
	mx := mmath_la.Matrix4x4[float32]{}
	mx.Identity()
	mx.Translate(position)

	for i := 0; i < len(m.Vertices); i++ {
		m.Vertices[i] = m.Vertices[i].TransformMatrix4x4(mx)
	}
}

func (m *Mesh) RotateY(rad float32) {
	mx := mmath_la.Matrix4x4[float32]{}
	mx.Identity()
	mx.RotateY(rad)

	for i := 0; i < len(m.Vertices); i++ {
		m.Vertices[i] = m.Vertices[i].TransformMatrix4x4(mx)
		m.Normal[i] = m.Normal[i].TransformMatrix4x4(mx)
	}
}

func (m *Mesh) ScaleUV(size mmath_la.Vector2[float32]) {
	for i := 0; i < len(m.UV); i++ {
		m.UV[i].X *= size.X
		m.UV[i].Y *= size.Y
	}
}

func (m *Mesh) OffsetUv(offset mmath_la.Vector2[float32]) {
	for i := 0; i < len(m.UV); i++ {
		m.UV[i].X += offset.X
		m.UV[i].Y += offset.Y
	}
}

// MakeCube 0b11_11_11_00 [Front, Back, Top, Bottom, Left, Right]
func (m *Mesh) MakeCube(size mmath_la.Vector3[float32], side uint8) {
	sideAmount := ml_number.CountSetBits(side)

	m.Vertices = make([]mmath_la.Vector3[float32], 0, 4*sideAmount)
	m.Normal = make([]mmath_la.Vector3[float32], 0, 4*sideAmount)

	// Front
	if side&0b1000_0000 == 0b1000_0000 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{0.0, 0.0, 1.0})
		}
	}

	// Back
	if side&0b0100_0000 == 0b0100_0000 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{0.0, 0.0, -1.0})
		}
	}

	// Top
	if side&0b0010_0000 == 0b0010_0000 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{0.0, 1.0, 0.0})
		}
	}

	// Bottom
	if side&0b0001_0000 == 0b0001_0000 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{0.0, -1.0, 0.0})
		}
	}

	// Left
	if side&0b0000_1000 == 0b0000_1000 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, -1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{-1.0, 1.0, -1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{-1.0, 0.0, 0.0})
		}
	}

	// Right
	if side&0b0000_0100 == 0b0000_0100 {
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, -1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, 1.0, 1.0})
		m.Vertices = append(m.Vertices, mmath_la.Vector3[float32]{1.0, -1.0, 1.0})

		for i := 0; i < 4; i++ {
			m.Normal = append(m.Normal, mmath_la.Vector3[float32]{1.0, 0.0, 0.0})
		}
	}

	// UV
	m.UV = make([]mmath_la.Vector2[float32], 0, 4*sideAmount)
	for i := 0; i < int(sideAmount); i++ {
		m.UV = append(m.UV, mmath_la.Vector2[float32]{0.0, 0.0})
		m.UV = append(m.UV, mmath_la.Vector2[float32]{1.0, 0.0})
		m.UV = append(m.UV, mmath_la.Vector2[float32]{1.0, 1.0})
		m.UV = append(m.UV, mmath_la.Vector2[float32]{0.0, 1.0})
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
