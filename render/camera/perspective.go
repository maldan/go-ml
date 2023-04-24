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
	Target      mmath_la.Vector3[float32]
	Matrix      mmath_la.Matrix4x4[float32]
}

func (p *PerspectiveCamera) ApplyMatrix() {
	proj := mmath_la.Matrix4x4[float32]{}
	proj.Perspective((p.Fov*3.141592653589793)/180, p.AspectRatio, 0.1, 1000.0)

	position := p.Position
	position.X *= -1
	position.Y *= -1
	position.Z *= -1

	// Position
	p.Matrix.Identity()

	p.Matrix.RotateX(p.Rotation.X)
	p.Matrix.Translate(position)
	//p.Matrix.RotateY(p.Rotation.Y)
	//p.Matrix.RotateZ(p.Rotation.Z)
	p.Matrix.Scale(p.Scale)

	//targetMx := mmath_la.Matrix4x4[float32]{}
	//targetMx.Identity()
	//targetMx.TargetTo(p.Position, target, mmath_la.Vector3[float32]{0, 1, 0})

	//p.Matrix.Multiply(targetMx)

	proj.Multiply(p.Matrix)
	p.Matrix = proj
}
