package mcgi_mesh

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

type MeshTarget struct {
	TriangleList []mmath_geom.Triangle3D[float32]
	NormalList   []mmath_geom.Triangle3D[float32]
	IsChanged    []bool
}

type MeshTexture struct {
	Diffuse   ml_image.ImageRGBA8
	AO        ml_image.ImageRGBA8
	Normal    ml_image.ImageRGBA8
	Roughness ml_image.ImageRGBA8
}

type Mesh struct {
	Vertices []mmath_la.Vector3[float32]
	Normal   []mmath_la.Vector3[float32]
	Indices  []uint32
	UV       []mmath_la.Vector2[float32]
	Texture  MeshTexture

	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Quaternion[float32]
	Scale    mmath_la.Vector3[float32]

	NormalMatrix mmath_la.Matrix4x4[float32]
	WorldMatrix  mmath_la.Matrix4x4[float32]
	RenderMatrix mmath_la.Matrix4x4[float32]

	OriginalTriangleList       []mmath_geom.Triangle3D[float32]
	OriginalNormalTriangleList []mmath_geom.Triangle3D[float32]

	TriangleList       []mmath_geom.Triangle3D[float32]
	NormalTriangleList []mmath_geom.Triangle3D[float32]
	UVTriangleList     []mmath_geom.Triangle2D[float32]

	TargetList   []MeshTarget
	TargetWeight []float32
}

func (m *Mesh) CalculateWorldMatrix() {
	m.WorldMatrix = mmath_la.Matrix4x4[float32]{}.
		Identity().
		Translate(m.Position).
		RotateX(m.Rotation.X * (mmath.Pi / 180.0)).
		RotateY(m.Rotation.Y * (mmath.Pi / 180.0)).
		RotateZ(m.Rotation.Z * (mmath.Pi / 180.0)).
		Scale(m.Scale)
}

func (m *Mesh) CalculateNormalMatrix() {
	m.NormalMatrix = m.WorldMatrix.Invert().Transpose()
}

func (m *Mesh) Prepare() {
	for j := 0; j < len(m.Indices); j += 3 {
		// Triangle
		m.OriginalTriangleList = append(m.OriginalTriangleList, mmath_geom.Triangle3D[float32]{
			A: m.Vertices[m.Indices[j]],
			B: m.Vertices[m.Indices[j+1]],
			C: m.Vertices[m.Indices[j+2]],
		})
		m.TriangleList = append(m.TriangleList, mmath_geom.Triangle3D[float32]{
			A: m.Vertices[m.Indices[j]],
			B: m.Vertices[m.Indices[j+1]],
			C: m.Vertices[m.Indices[j+2]],
		})

		// UV
		m.UVTriangleList = append(m.UVTriangleList, mmath_geom.Triangle2D[float32]{
			A: m.UV[m.Indices[j]],
			B: m.UV[m.Indices[j+1]],
			C: m.UV[m.Indices[j+2]],
		})

		// Normal
		m.NormalTriangleList = append(m.NormalTriangleList, mmath_geom.Triangle3D[float32]{
			A: m.Normal[m.Indices[j]],
			B: m.Normal[m.Indices[j+1]],
			C: m.Normal[m.Indices[j+2]],
		})
		m.OriginalNormalTriangleList = append(m.OriginalNormalTriangleList, mmath_geom.Triangle3D[float32]{
			A: m.Normal[m.Indices[j]],
			B: m.Normal[m.Indices[j+1]],
			C: m.Normal[m.Indices[j+2]],
		})
	}
}
