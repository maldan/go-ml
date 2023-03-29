package ml_geom

import (
	"golang.org/x/exp/constraints"
	"math"
)

type Matrix4x4[T constraints.Float] struct {
	Raw [16]T
}

func (m *Matrix4x4[T]) Identity() {
	m.Raw[0] = 1
	m.Raw[1] = 0
	m.Raw[2] = 0
	m.Raw[3] = 0
	m.Raw[4] = 0
	m.Raw[5] = 1
	m.Raw[6] = 0
	m.Raw[7] = 0
	m.Raw[8] = 0
	m.Raw[9] = 0
	m.Raw[10] = 1
	m.Raw[11] = 0
	m.Raw[12] = 0
	m.Raw[13] = 0
	m.Raw[14] = 0
	m.Raw[15] = 1
}

func (m *Matrix4x4[T]) Translate(v Vector3[T]) {
	a00 := m.Raw[0]
	a01 := m.Raw[1]
	a02 := m.Raw[2]
	a03 := m.Raw[3]
	a10 := m.Raw[4]
	a11 := m.Raw[5]
	a12 := m.Raw[6]
	a13 := m.Raw[7]
	a20 := m.Raw[8]
	a21 := m.Raw[9]
	a22 := m.Raw[10]
	a23 := m.Raw[11]

	m.Raw[0] = a00
	m.Raw[1] = a01
	m.Raw[2] = a02
	m.Raw[3] = a03
	m.Raw[4] = a10
	m.Raw[5] = a11
	m.Raw[6] = a12
	m.Raw[7] = a13
	m.Raw[8] = a20
	m.Raw[9] = a21
	m.Raw[10] = a22
	m.Raw[11] = a23

	m.Raw[12] = a00*v.X + a10*v.Y + a20*v.Z + m.Raw[12]
	m.Raw[13] = a01*v.X + a11*v.Y + a21*v.Z + m.Raw[13]
	m.Raw[14] = a02*v.X + a12*v.Y + a22*v.Z + m.Raw[14]
	m.Raw[15] = a03*v.X + a13*v.Y + a23*v.Z + m.Raw[15]
}

func (m *Matrix4x4[T]) Scale(v Vector3[T]) {
	m.Raw[0] *= v.X
	m.Raw[1] *= v.X
	m.Raw[2] *= v.X
	m.Raw[3] *= v.X
	m.Raw[4] *= v.Y
	m.Raw[5] *= v.Y
	m.Raw[6] *= v.Y
	m.Raw[7] *= v.Y
	m.Raw[8] *= v.Z
	m.Raw[9] *= v.Z
	m.Raw[10] *= v.Z
	m.Raw[11] *= v.Z
}

func (m *Matrix4x4[T]) RotateX(rad T) {
	s := T(math.Sin(float64(rad)))
	c := T(math.Cos(float64(rad)))

	a10 := m.Raw[4]
	a11 := m.Raw[5]
	a12 := m.Raw[6]
	a13 := m.Raw[7]
	a20 := m.Raw[8]
	a21 := m.Raw[9]
	a22 := m.Raw[10]
	a23 := m.Raw[11]

	// Perform axis-specific matrix multiplication
	m.Raw[4] = a10*c + a20*s
	m.Raw[5] = a11*c + a21*s
	m.Raw[6] = a12*c + a22*s
	m.Raw[7] = a13*c + a23*s
	m.Raw[8] = a20*c - a10*s
	m.Raw[9] = a21*c - a11*s
	m.Raw[10] = a22*c - a12*s
	m.Raw[11] = a23*c - a13*s
}

func (m *Matrix4x4[T]) RotateY(rad T) {
	s := T(math.Sin(float64(rad)))
	c := T(math.Cos(float64(rad)))

	a00 := m.Raw[0]
	a01 := m.Raw[1]
	a02 := m.Raw[2]
	a03 := m.Raw[3]
	a20 := m.Raw[8]
	a21 := m.Raw[9]
	a22 := m.Raw[10]
	a23 := m.Raw[11]

	// Perform axis-specific matrix multiplication
	m.Raw[0] = a00*c - a20*s
	m.Raw[1] = a01*c - a21*s
	m.Raw[2] = a02*c - a22*s
	m.Raw[3] = a03*c - a23*s
	m.Raw[8] = a00*s + a20*c
	m.Raw[9] = a01*s + a21*c
	m.Raw[10] = a02*s + a22*c
	m.Raw[11] = a03*s + a23*c
}

func (m *Matrix4x4[T]) RotateZ(rad T) {
	s := T(math.Sin(float64(rad)))
	c := T(math.Cos(float64(rad)))

	a00 := m.Raw[0]
	a01 := m.Raw[1]
	a02 := m.Raw[2]
	a03 := m.Raw[3]
	a10 := m.Raw[4]
	a11 := m.Raw[5]
	a12 := m.Raw[6]
	a13 := m.Raw[7]

	// Perform axis-specific matrix multiplication
	m.Raw[0] = a00*c + a10*s
	m.Raw[1] = a01*c + a11*s
	m.Raw[2] = a02*c + a12*s
	m.Raw[3] = a03*c + a13*s
	m.Raw[4] = a10*c - a00*s
	m.Raw[5] = a11*c - a01*s
	m.Raw[6] = a12*c - a02*s
	m.Raw[7] = a13*c - a03*s
}

func (m *Matrix4x4[T]) Perspective(fov T, aspect T, near T, far T) {
	f := 1.0 / T(math.Tan(float64(fov)/2.0))
	m.Raw[0] = f / aspect
	m.Raw[1] = 0
	m.Raw[2] = 0
	m.Raw[3] = 0
	m.Raw[4] = 0
	m.Raw[5] = f
	m.Raw[6] = 0
	m.Raw[7] = 0
	m.Raw[8] = 0
	m.Raw[9] = 0
	m.Raw[11] = -1
	m.Raw[12] = 0
	m.Raw[13] = 0
	m.Raw[15] = 0

	nf := 1 / (near - far)
	m.Raw[10] = (far + near) * nf
	m.Raw[14] = 2 * far * near * nf
}

func (m *Matrix4x4[T]) Multiply(b Matrix4x4[T]) {
	a00 := m.Raw[0]
	a01 := m.Raw[1]
	a02 := m.Raw[2]
	a03 := m.Raw[3]
	a10 := m.Raw[4]
	a11 := m.Raw[5]
	a12 := m.Raw[6]
	a13 := m.Raw[7]
	a20 := m.Raw[8]
	a21 := m.Raw[9]
	a22 := m.Raw[10]
	a23 := m.Raw[11]
	a30 := m.Raw[12]
	a31 := m.Raw[13]
	a32 := m.Raw[14]
	a33 := m.Raw[15]

	// Cache only the current line of the second matrix
	b0 := b.Raw[0]
	b1 := b.Raw[1]
	b2 := b.Raw[2]
	b3 := b.Raw[3]
	m.Raw[0] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	m.Raw[1] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	m.Raw[2] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	m.Raw[3] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0 = b.Raw[4]
	b1 = b.Raw[5]
	b2 = b.Raw[6]
	b3 = b.Raw[7]
	m.Raw[4] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	m.Raw[5] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	m.Raw[6] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	m.Raw[7] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0 = b.Raw[8]
	b1 = b.Raw[9]
	b2 = b.Raw[10]
	b3 = b.Raw[11]
	m.Raw[8] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	m.Raw[9] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	m.Raw[10] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	m.Raw[11] = b0*a03 + b1*a13 + b2*a23 + b3*a33

	b0 = b.Raw[12]
	b1 = b.Raw[13]
	b2 = b.Raw[14]
	b3 = b.Raw[15]
	m.Raw[12] = b0*a00 + b1*a10 + b2*a20 + b3*a30
	m.Raw[13] = b0*a01 + b1*a11 + b2*a21 + b3*a31
	m.Raw[14] = b0*a02 + b1*a12 + b2*a22 + b3*a32
	m.Raw[15] = b0*a03 + b1*a13 + b2*a23 + b3*a33
}

func (m *Matrix4x4[T]) Rotate(rad T, axis Vector3[T]) {
	ln := T(math.Sqrt(float64(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)))
	if ln < 0.000001 {
		return
	}

	ln = 1 / ln
	axis.X *= ln
	axis.Y *= ln
	axis.Z *= ln

	s := T(math.Sin(float64(rad)))
	c := T(math.Cos(float64(rad)))
	t := 1 - c

	a00 := m.Raw[0]
	a01 := m.Raw[1]
	a02 := m.Raw[2]
	a03 := m.Raw[3]
	a10 := m.Raw[4]
	a11 := m.Raw[5]
	a12 := m.Raw[6]
	a13 := m.Raw[7]
	a20 := m.Raw[8]
	a21 := m.Raw[9]
	a22 := m.Raw[10]
	a23 := m.Raw[11]

	// Construct the elements of the rotation matrix
	b00 := axis.X*axis.X*t + c
	b01 := axis.Y*axis.X*t + axis.Z*s
	b02 := axis.Z*axis.X*t - axis.Y*s
	b10 := axis.X*axis.Y*t - axis.Z*s
	b11 := axis.Y*axis.Y*t + c
	b12 := axis.Z*axis.Y*t + axis.X*s
	b20 := axis.X*axis.Z*t + axis.Y*s
	b21 := axis.Y*axis.Z*t - axis.X*s
	b22 := axis.Z*axis.Z*t + c

	// Perform rotation-specific matrix multiplication
	m.Raw[0] = a00*b00 + a10*b01 + a20*b02
	m.Raw[1] = a01*b00 + a11*b01 + a21*b02
	m.Raw[2] = a02*b00 + a12*b01 + a22*b02
	m.Raw[3] = a03*b00 + a13*b01 + a23*b02
	m.Raw[4] = a00*b10 + a10*b11 + a20*b12
	m.Raw[5] = a01*b10 + a11*b11 + a21*b12
	m.Raw[6] = a02*b10 + a12*b11 + a22*b12
	m.Raw[7] = a03*b10 + a13*b11 + a23*b12
	m.Raw[8] = a00*b20 + a10*b21 + a20*b22
	m.Raw[9] = a01*b20 + a11*b21 + a21*b22
	m.Raw[10] = a02*b20 + a12*b21 + a22*b22
	m.Raw[11] = a03*b20 + a13*b21 + a23*b22
}
