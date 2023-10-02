package mmath_la

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
)

type Matrix4x4[T constraints.Float] struct {
	Raw [16]T
}

func (m Matrix4x4[T]) Identity() Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) Transpose() Matrix4x4[T] {
	result := [16]T{}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			result[i*4+j] = m.Raw[j*4+i]
		}
	}

	m.Raw = result

	return m
}

func (m *Matrix4x4[T]) Clone() Matrix4x4[T] {
	return Matrix4x4[T]{Raw: m.Raw}
}

func (m Matrix4x4[T]) Invert() Matrix4x4[T] {
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

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	// Calculate the determinant
	det :=
		b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06

	if det == 0 {
		return m
	}
	det = 1.0 / det

	m.Raw[0] = (a11*b11 - a12*b10 + a13*b09) * det
	m.Raw[1] = (a02*b10 - a01*b11 - a03*b09) * det
	m.Raw[2] = (a31*b05 - a32*b04 + a33*b03) * det
	m.Raw[3] = (a22*b04 - a21*b05 - a23*b03) * det
	m.Raw[4] = (a12*b08 - a10*b11 - a13*b07) * det
	m.Raw[5] = (a00*b11 - a02*b08 + a03*b07) * det
	m.Raw[6] = (a32*b02 - a30*b05 - a33*b01) * det
	m.Raw[7] = (a20*b05 - a22*b02 + a23*b01) * det
	m.Raw[8] = (a10*b10 - a11*b08 + a13*b06) * det
	m.Raw[9] = (a01*b08 - a00*b10 - a03*b06) * det
	m.Raw[10] = (a30*b04 - a31*b02 + a33*b00) * det
	m.Raw[11] = (a21*b02 - a20*b04 - a23*b00) * det
	m.Raw[12] = (a11*b07 - a10*b09 - a12*b06) * det
	m.Raw[13] = (a00*b09 - a01*b07 + a02*b06) * det
	m.Raw[14] = (a31*b01 - a30*b03 - a32*b00) * det
	m.Raw[15] = (a20*b03 - a21*b01 + a22*b00) * det

	return m
}

func (m Matrix4x4[T]) Translate(v Vector3[T]) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) Scale(v Vector3[T]) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) RotateX(rad T) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) RotateY(rad T) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) RotateZ(rad T) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) RotateQuaternion(q Quaternion[T]) Matrix4x4[T] {
	mm := q.ToMatrix4x4()
	return m.Multiply(mm)
}

func (m Matrix4x4[T]) Perspective(fov T, aspect T, near T, far T) Matrix4x4[T] {
	f := T(1.0 / (math.Tan(float64(fov) / 2.0)))

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
	m.Raw[11] = -1.0
	m.Raw[12] = 0
	m.Raw[13] = 0
	m.Raw[15] = 0

	// nf := 1 / (near - far)
	m.Raw[10] = (far + near) / (near - far)
	m.Raw[14] = (2.0 * far * near) / (near - far)

	return m
}

func (m Matrix4x4[T]) Orthographic(left T, right T, bottom T, top T, near T, far T) Matrix4x4[T] {
	m.Raw[0] = 2 / (right - left)
	m.Raw[1] = 0
	m.Raw[2] = 0
	m.Raw[3] = 0

	m.Raw[4] = 0
	m.Raw[5] = 2 / (top - bottom)
	m.Raw[6] = 0
	m.Raw[7] = 0

	m.Raw[8] = 0
	m.Raw[9] = 0
	m.Raw[10] = 2 / (near - far)
	m.Raw[11] = 0

	m.Raw[12] = (left + right) / (left - right)
	m.Raw[13] = (bottom + top) / (bottom - top)
	m.Raw[14] = (near + far) / (near - far)
	m.Raw[15] = 1

	/*lr := 1 / (left - right)
	bt := 1 / (bottom - top)
	nf := 1 / (near - far)

	m.Raw[0] = -2 * lr
	m.Raw[1] = 0
	m.Raw[2] = 0
	m.Raw[3] = 0
	m.Raw[4] = 0
	m.Raw[5] = -2 * bt
	m.Raw[6] = 0
	m.Raw[7] = 0
	m.Raw[8] = 0
	m.Raw[9] = 0
	m.Raw[10] = 2 * nf
	m.Raw[11] = 0
	m.Raw[12] = (left + right) * lr
	m.Raw[13] = (top + bottom) * bt
	m.Raw[14] = (far + near) * nf
	m.Raw[15] = 1*/

	return m
}

func (m Matrix4x4[T]) Multiply(b Matrix4x4[T]) Matrix4x4[T] {
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

	return m
}

func (m Matrix4x4[T]) MultiplyVector4(v Vector4[T]) Vector4[T] {
	result := Vector4[T]{}

	result.X = v.X*m.Raw[0] + v.Y*m.Raw[1] + v.Z*m.Raw[2] + v.W*m.Raw[3]
	result.Y = v.X*m.Raw[4] + v.Y*m.Raw[5] + v.Z*m.Raw[6] + v.W*m.Raw[7]
	result.Z = v.X*m.Raw[8] + v.Y*m.Raw[9] + v.Z*m.Raw[10] + v.W*m.Raw[11]
	result.W = v.X*m.Raw[12] + v.Y*m.Raw[13] + v.Z*m.Raw[14] + v.W*m.Raw[15]

	return result
}

func (m Matrix4x4[T]) Rotate(rad T, axis Vector3[T]) Matrix4x4[T] {
	ln := T(math.Sqrt(float64(axis.X*axis.X + axis.Y*axis.Y + axis.Z*axis.Z)))
	if ln < T(0.000001) {
		return Matrix4x4[T]{}
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

	return m
}

func (m *Matrix4x4[T]) TargetTo(from Vector3[T], to Vector3[T], up Vector3[T]) {
	z0 := from.X - to.X
	z1 := from.Y - to.Y
	z2 := from.Z - to.Z

	ln := z0*z0 + z1*z1 + z2*z2
	if ln > 0 {
		ln = 1 / T(math.Sqrt(float64(ln)))
		z0 *= ln
		z1 *= ln
		z2 *= ln
	}

	x0 := up.Y*z2 - up.Z*z1
	x1 := up.Z*z0 - up.X*z2
	x2 := up.X*z1 - up.Y*z0

	ln = x0*x0 + x1*x1 + x2*x2
	if ln > 0 {
		ln = 1 / T(math.Sqrt(float64(ln)))
		x0 *= ln
		x1 *= ln
		x2 *= ln
	}

	m.Raw[0] = x0
	m.Raw[1] = x1
	m.Raw[2] = x2
	m.Raw[3] = 0
	m.Raw[4] = z1*x2 - z2*x1
	m.Raw[5] = z2*x0 - z0*x2
	m.Raw[6] = z0*x1 - z1*x0
	m.Raw[7] = 0
	m.Raw[8] = z0
	m.Raw[9] = z1
	m.Raw[10] = z2
	m.Raw[11] = 0
	m.Raw[12] = from.X
	m.Raw[13] = from.Y
	m.Raw[14] = from.Z
	m.Raw[15] = 1
}

func (m Matrix4x4[T]) ToString() string {
	return fmt.Sprintf(
		"Matrix4x4 \n(%.2f %.2f %.2f %.2f\n%.2f %.2f %.2f %.2f\n%.2f %.2f %.2f %.2f\n%.2f %.2f %.2f %.2f)",
		m.Raw[0], m.Raw[1], m.Raw[2], m.Raw[3],
		m.Raw[4], m.Raw[5], m.Raw[6], m.Raw[7],
		m.Raw[8], m.Raw[9], m.Raw[10], m.Raw[11],
		m.Raw[12], m.Raw[13], m.Raw[14], m.Raw[15],
	)
}

func (m Matrix4x4[T]) FromBytes(data []byte) Matrix4x4[T] {
	for i := 0; i < 16; i += 4 {
		v1 := Vector4[T]{}.FromBytes(data[i*4:])
		m.Raw[i] = v1.X
		m.Raw[i+1] = v1.Y
		m.Raw[i+2] = v1.Z
		m.Raw[i+3] = v1.W
	}

	return m
}

func (m Matrix4x4[T]) GetPosition() Vector3[T] {
	return Vector3[T]{X: m.Raw[12], Y: m.Raw[13], Z: m.Raw[14]}
}

func (m Matrix4x4[T]) GetScale() Vector3[T] {
	m11 := m.Raw[0]
	m12 := m.Raw[1]
	m13 := m.Raw[2]
	m21 := m.Raw[4]
	m22 := m.Raw[5]
	m23 := m.Raw[6]
	m31 := m.Raw[8]
	m32 := m.Raw[9]
	m33 := m.Raw[10]

	return Vector3[T]{
		X: T(math.Sqrt(float64(m11*m11 + m12*m12 + m13*m13))),
		Y: T(math.Sqrt(float64(m21*m21 + m22*m22 + m23*m23))),
		Z: T(math.Sqrt(float64(m31*m31 + m32*m32 + m33*m33))),
	}
}

func (m Matrix4x4[T]) GetRotation() Quaternion[T] {
	scaling := m.GetScale()

	is1 := 1.0 / scaling.X
	is2 := 1.0 / scaling.Y
	is3 := 1.0 / scaling.Z

	sm11 := m.Raw[0] * is1
	sm12 := m.Raw[1] * is2
	sm13 := m.Raw[2] * is3
	sm21 := m.Raw[4] * is1
	sm22 := m.Raw[5] * is2
	sm23 := m.Raw[6] * is3
	sm31 := m.Raw[8] * is1
	sm32 := m.Raw[9] * is2
	sm33 := m.Raw[10] * is3

	trace := sm11 + sm22 + sm33
	s := T(0.0)
	out := Quaternion[T]{}.Identity()

	if trace > 0.0 {
		s = T(math.Sqrt(float64(trace+1.0)) * 2.0)
		out.X = (sm23 - sm32) / s
		out.Y = (sm31 - sm13) / s
		out.Z = (sm12 - sm21) / s
		out.W = 0.25 * s
	} else if sm11 > sm22 && sm11 > sm33 {
		s = T(math.Sqrt(float64(1.0+sm11-sm22-sm33))) * 2.0
		out.X = 0.25 * s
		out.Y = (sm12 + sm21) / s
		out.Z = (sm31 + sm13) / s
		out.W = (sm23 - sm32) / s
	} else if sm22 > sm33 {
		s = T(math.Sqrt(float64(1.0+sm22-sm11-sm33))) * 2.0
		out.X = (sm12 + sm21) / s
		out.Y = 0.25 * s
		out.Z = (sm23 + sm32) / s
		out.W = (sm31 - sm13) / s
	} else {
		s = T(math.Sqrt(float64(1.0+sm33-sm11-sm22))) * 2.0
		out.X = (sm31 + sm13) / s
		out.Y = (sm23 + sm32) / s
		out.Z = 0.25 * s
		out.W = (sm12 - sm21) / s
	}

	return out
}
