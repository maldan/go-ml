package mmath_geom

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
)

type Triangle4D[T constraints.Float] struct {
	A mmath_la.Vector4[T]
	B mmath_la.Vector4[T]
	C mmath_la.Vector4[T]
}

func (t Triangle4D[T]) MultiplyMatrix4x4(mx mmath_la.Matrix4x4[T]) Triangle4D[T] {
	t.A = t.A.MulMatrix(mx)
	t.B = t.B.MulMatrix(mx)
	t.C = t.C.MulMatrix(mx)
	return t
}

func (t Triangle4D[T]) ToTriangle3D() Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.ToVector3XYZ(),
		B: t.B.ToVector3XYZ(),
		C: t.C.ToVector3XYZ(),
	}
}

func (t Triangle4D[T]) TopMiddleBottom() (mmath_la.Vector4[T], mmath_la.Vector4[T], mmath_la.Vector4[T]) {
	var top, middle, bottom mmath_la.Vector4[T]

	x1 := t.A.X
	x2 := t.B.X
	x3 := t.C.X
	y1 := t.A.Y
	y2 := t.B.Y
	y3 := t.C.Y
	z1 := t.A.Z
	z2 := t.B.Z
	z3 := t.C.Z
	w1 := t.A.W
	w2 := t.B.W
	w3 := t.C.W

	if y1 <= y2 && y1 <= y3 {
		top = mmath_la.Vector4[T]{x1, y1, z1, w1}
		if y2 <= y3 {
			middle = mmath_la.Vector4[T]{x2, y2, z2, w2}
			bottom = mmath_la.Vector4[T]{x3, y3, z3, w3}
		} else {
			middle = mmath_la.Vector4[T]{x3, y3, z3, w3}
			bottom = mmath_la.Vector4[T]{x2, y2, z2, w2}
		}
		return top, middle, bottom
	}

	if y2 <= y1 && y2 <= y3 {
		top = mmath_la.Vector4[T]{x2, y2, z2, w2}
		if y1 <= y3 {
			middle = mmath_la.Vector4[T]{x1, y1, z1, w1}
			bottom = mmath_la.Vector4[T]{x3, y3, z3, w3}
		} else {
			middle = mmath_la.Vector4[T]{x3, y3, z3, w3}
			bottom = mmath_la.Vector4[T]{x1, y1, z1, w1}
		}
		return top, middle, bottom
	}

	top = mmath_la.Vector4[T]{x3, y3, z3, w3}
	if y1 <= y2 {
		middle = mmath_la.Vector4[T]{x1, y1, z1, w1}
		bottom = mmath_la.Vector4[T]{x2, y2, z2, w2}
	} else {
		middle = mmath_la.Vector4[T]{x2, y2, z2, w2}
		bottom = mmath_la.Vector4[T]{x1, y1, z1, w1}
	}
	return top, middle, bottom
}

func (t Triangle4D[T]) MinZ() T {
	return mmath.Min(mmath.Min(t.A.Z, t.B.Z), t.C.Z)
}

func (t Triangle4D[T]) MinW() T {
	return mmath.Min(mmath.Min(t.A.W, t.B.W), t.C.W)
}

func (t Triangle4D[T]) MaxZ() T {
	return mmath.Max(mmath.Max(t.A.Z, t.B.Z), t.C.Z)
}

func (t Triangle4D[T]) MaxW() T {
	return mmath.Max(mmath.Max(t.A.W, t.B.W), t.C.W)
}

/*func (t Triangle3D[T]) InterpolateABG(alpha T, beta T, gamma T) mmath_la.Vector3[T] {
	uva := t.A.MulScalar(alpha)
	uvb := t.B.MulScalar(beta)
	uvc := t.C.MulScalar(gamma)
	return (uva.Add(uvb)).Add(uvc)
}

func (t Triangle3D[T]) Lerp(to Triangle3D[T], weight T) Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.Lerp(to.A, weight),
		B: t.B.Lerp(to.B, weight),
		C: t.C.Lerp(to.C, weight),
	}
}

func (t Triangle3D[T]) Add(to Triangle3D[T]) Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.Add(to.A),
		B: t.B.Add(to.B),
		C: t.C.Add(to.C),
	}
}

func (t Triangle3D[T]) Sub(to Triangle3D[T]) Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.Sub(to.A),
		B: t.B.Sub(to.B),
		C: t.C.Sub(to.C),
	}
}

func (t Triangle3D[T]) Mul(to Triangle3D[T]) Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.Mul(to.A),
		B: t.B.Mul(to.B),
		C: t.C.Mul(to.C),
	}
}

func (t Triangle3D[T]) MulScalar(s T) Triangle3D[T] {
	return Triangle3D[T]{
		A: t.A.MulScalar(s),
		B: t.B.MulScalar(s),
		C: t.C.MulScalar(s),
	}
}

func (t Triangle3D[T]) IsZero() bool {
	return t.A.IsZero() && t.B.IsZero() && t.C.IsZero()
}

func (t Triangle3D[T]) GetOrientation(view mmath_la.Vector3[T]) int {
	v1 := t.B.Sub(t.A)
	v2 := t.C.Sub(t.A)

	cross := v1.Cross(v2)
	observerVector := view.Sub(t.A)

	dotProduct := cross.Dot(observerVector)

	if dotProduct > 0 {
		return 1
	} else if dotProduct < 0 {
		return -1
	}
	return 0
}

func (t Triangle3D[T]) TopMiddleBottom() (mmath_la.Vector3[T], mmath_la.Vector3[T], mmath_la.Vector3[T]) {
	var top, middle, bottom mmath_la.Vector3[T]

	x1 := t.A.X
	x2 := t.B.X
	x3 := t.C.X
	y1 := t.A.Y
	y2 := t.B.Y
	y3 := t.C.Y
	z1 := t.A.Z
	z2 := t.B.Z
	z3 := t.C.Z

	if y1 <= y2 && y1 <= y3 {
		top = mmath_la.Vector3[T]{x1, y1, z1}
		if y2 <= y3 {
			middle = mmath_la.Vector3[T]{x2, y2, z2}
			bottom = mmath_la.Vector3[T]{x3, y3, z3}
		} else {
			middle = mmath_la.Vector3[T]{x3, y3, z3}
			bottom = mmath_la.Vector3[T]{x2, y2, z2}
		}
		return top, middle, bottom
	}

	if y2 <= y1 && y2 <= y3 {
		top = mmath_la.Vector3[T]{x2, y2, z2}
		if y1 <= y3 {
			middle = mmath_la.Vector3[T]{x1, y1, z1}
			bottom = mmath_la.Vector3[T]{x3, y3, z3}
		} else {
			middle = mmath_la.Vector3[T]{x3, y3, z3}
			bottom = mmath_la.Vector3[T]{x1, y1, z1}
		}
		return top, middle, bottom
	}

	top = mmath_la.Vector3[T]{x3, y3, z3}
	if y1 <= y2 {
		middle = mmath_la.Vector3[T]{x1, y1, z1}
		bottom = mmath_la.Vector3[T]{x2, y2, z2}
	} else {
		middle = mmath_la.Vector3[T]{x2, y2, z2}
		bottom = mmath_la.Vector3[T]{x1, y1, z1}
	}
	return top, middle, bottom
}

func (t Triangle3D[T]) Area() T {
	// Задайте вершины треугольника в 3D координатах
	vertex1 := t.A
	vertex2 := t.B
	vertex3 := t.C

	// Вычислите векторы между вершинами
	edge1 := mmath_la.Vector3[T]{
		X: vertex2.X - vertex1.X,
		Y: vertex2.Y - vertex1.Y,
		Z: vertex2.Z - vertex1.Z,
	}

	edge2 := mmath_la.Vector3[T]{
		X: vertex3.X - vertex1.X,
		Y: vertex3.Y - vertex1.Y,
		Z: vertex3.Z - vertex1.Z,
	}

	// Вычислите векторное произведение (кросс-произведение) между edge1 и edge2
	crossProduct := mmath_la.Vector3[T]{
		X: edge1.Y*edge2.Z - edge1.Z*edge2.Y,
		Y: edge1.Z*edge2.X - edge1.X*edge2.Z,
		Z: edge1.X*edge2.Y - edge1.Y*edge2.X,
	}

	// Вычислите площадь треугольника как половину длины вектора crossProduct
	area := 0.5 * T(math.Sqrt(float64(crossProduct.X*crossProduct.X+crossProduct.Y*crossProduct.Y+crossProduct.Z*crossProduct.Z)))
	return area
}
*/
