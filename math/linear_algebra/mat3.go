package mmath_la

import "golang.org/x/exp/constraints"

type Matrix3x3[T constraints.Float] struct {
	Raw [9]T
}

func (m Matrix3x3[T]) Identity() Matrix3x3[T] {
	m.Raw[0] = 1
	m.Raw[1] = 0
	m.Raw[2] = 0

	m.Raw[3] = 0
	m.Raw[4] = 1
	m.Raw[5] = 0

	m.Raw[6] = 0
	m.Raw[7] = 0
	m.Raw[8] = 1

	return m
}

func (m Matrix3x3[T]) Transpose() Matrix3x3[T] {
	result := [9]T{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result[i*3+j] = m.Raw[j*3+i]
		}
	}

	m.Raw = result

	return m
}

func (m Matrix3x3[T]) Determinant() T {
	return m.Raw[0]*(m.Raw[4]*m.Raw[8]-m.Raw[5]*m.Raw[7]) -
		m.Raw[1]*(m.Raw[3]*m.Raw[8]-m.Raw[5]*m.Raw[6]) +
		m.Raw[2]*(m.Raw[3]*m.Raw[7]-m.Raw[4]*m.Raw[6])
}

func (m Matrix3x3[T]) Invert() Matrix3x3[T] {
	det := m.Determinant()

	if det == 0 {
		panic("Определитель матрицы равен 0, матрица необратима")
	}

	inverted := Matrix3x3[T]{}

	inverted.Raw[0] = (m.Raw[4]*m.Raw[8] - m.Raw[5]*m.Raw[7]) / det
	inverted.Raw[1] = (m.Raw[2]*m.Raw[7] - m.Raw[1]*m.Raw[8]) / det
	inverted.Raw[2] = (m.Raw[1]*m.Raw[5] - m.Raw[2]*m.Raw[4]) / det
	inverted.Raw[3] = (m.Raw[5]*m.Raw[6] - m.Raw[3]*m.Raw[8]) / det
	inverted.Raw[4] = (m.Raw[0]*m.Raw[8] - m.Raw[2]*m.Raw[6]) / det
	inverted.Raw[5] = (m.Raw[2]*m.Raw[3] - m.Raw[0]*m.Raw[5]) / det
	inverted.Raw[6] = (m.Raw[3]*m.Raw[7] - m.Raw[4]*m.Raw[6]) / det
	inverted.Raw[7] = (m.Raw[1]*m.Raw[6] - m.Raw[0]*m.Raw[7]) / det
	inverted.Raw[8] = (m.Raw[0]*m.Raw[4] - m.Raw[1]*m.Raw[3]) / det

	return inverted
}

func (m Matrix3x3[T]) FromVectorRows(r1 Vector3[T], r2 Vector3[T], r3 Vector3[T]) Matrix3x3[T] {
	m.Raw[0] = r1.X
	m.Raw[1] = r1.Y
	m.Raw[2] = r1.Z

	m.Raw[3] = r2.X
	m.Raw[4] = r2.Y
	m.Raw[5] = r2.Z

	m.Raw[6] = r3.X
	m.Raw[7] = r3.Y
	m.Raw[8] = r3.Z

	return m
}

func (m Matrix3x3[T]) FromVectorColumns(r1 Vector3[T], r2 Vector3[T], r3 Vector3[T]) Matrix3x3[T] {
	m.Raw[0] = r1.X
	m.Raw[1] = r2.X
	m.Raw[2] = r3.X

	m.Raw[3] = r1.Y
	m.Raw[4] = r2.Y
	m.Raw[5] = r3.Y

	m.Raw[6] = r1.Z
	m.Raw[7] = r2.Z
	m.Raw[8] = r3.Z

	return m
}
