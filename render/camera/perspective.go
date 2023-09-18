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
	Near        float32
	Far         float32
}

func (p *PerspectiveCamera) ApplyMatrix() {
	proj := mmath_la.Matrix4x4[float32]{}.
		Perspective((p.Fov*3.141592653589793)/180, p.AspectRatio, p.Near, p.Far)

	position := p.Position
	position.X *= -1
	position.Y *= -1
	position.Z *= -1

	// Position
	p.Matrix = p.Matrix.
		Identity().
		RotateX(p.Rotation.X).
		Translate(position).
		Scale(p.Scale)

	//targetMx := mmath_la.Matrix4x4[float32]{}
	//targetMx.Identity()
	//targetMx.TargetTo(p.Position, target, mmath_la.Vector3[float32]{0, 1, 0})

	//p.Matrix.Multiply(targetMx)

	p.Matrix = proj.Multiply(p.Matrix)
}
