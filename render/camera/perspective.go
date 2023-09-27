package mrender_camera

import (
	mmath "github.com/maldan/go-ml/math"
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

func (p *PerspectiveCamera) MoveToWhereLooking() {
	proj := mmath_la.Matrix4x4[float32]{}.
		Perspective(mmath.DegToRad(p.Fov), p.AspectRatio, p.Near, p.Far)

	position := p.Position
	position.X *= -1
	position.Y *= -1
	position.Z *= -1

	// Position
	p.Matrix = p.Matrix.
		Identity().
		Translate(position).
		RotateX(p.Rotation.X).
		RotateY(p.Rotation.Y).
		RotateZ(p.Rotation.Z).
		Scale(p.Scale)

	// fmt.Printf("MATRIX: %v\n", proj.Raw)

	p.Matrix = proj.Multiply(p.Matrix)
}

func (p *PerspectiveCamera) ApplyMatrix() {
	proj := mmath_la.Matrix4x4[float32]{}.
		Perspective(mmath.DegToRad(p.Fov), p.AspectRatio, p.Near, p.Far)

	position := p.Position
	position.X *= -1
	position.Y *= -1
	position.Z *= -1

	// Position
	p.Matrix = p.Matrix.
		Identity().
		RotateX(p.Rotation.X).
		RotateY(p.Rotation.Y).
		Translate(position).
		Scale(p.Scale)

	// fmt.Printf("MATRIX: %v\n", proj.Raw)

	p.Matrix = proj.Multiply(p.Matrix)
}
