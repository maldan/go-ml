package mmath_geom

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
	"math"
)

type Triangle3D[T constraints.Float] struct {
	A mmath_la.Vector3[T]
	B mmath_la.Vector3[T]
	C mmath_la.Vector3[T]
}

type Triangle2D[T constraints.Float] struct {
	A mmath_la.Vector2[T]
	B mmath_la.Vector2[T]
	C mmath_la.Vector2[T]
}

func (t Triangle3D[T]) TransformMatrix4x4(mx mmath_la.Matrix4x4[T]) Triangle3D[T] {
	t.A = t.A.TransformMatrix4x4(mx)
	t.B = t.B.TransformMatrix4x4(mx)
	t.C = t.C.TransformMatrix4x4(mx)
	return t
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

func (t Triangle2D[T]) OnEachPixel(fn func(v mmath_la.Vector2[T], a T, b T, g T), fId int) {
	w := int(math.Round(float64(t.Width() + 2)))
	h := int(math.Round(float64(t.Height() + 2)))
	minX := T(math.Round(float64(t.MinX() - 1)))
	minY := T(math.Round(float64(t.MinY() - 1)))

	if w > 320*2 {
		w = 320 * 2
	}
	if h > 240*2 {
		h = 240 * 2
	}

	square := t.Square()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			p := mmath_la.Vector2[T]{
				X: minX + T(x),
				Y: minY + T(y),
			}

			if p.X < 0 || p.Y < 0 || p.X > 320*2 || p.Y > 240*2 {
				continue
			}

			/*if int(p.X)%4 == fId || int(p.Y)%4 == fId {
				continue
			}*/

			a := t.A
			b := t.B
			c := t.C

			// Interpolation coefficient
			s1 := 0.5 * ((b.X-p.X)*(c.Y-p.Y) - (c.X-p.X)*(b.Y-p.Y))
			s2 := 0.5 * ((p.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(p.Y-a.Y))
			s3 := 0.5 * ((b.X-a.X)*(p.Y-a.Y) - (p.X-a.X)*(b.Y-a.Y))

			alpha := s1 / square
			beta := s2 / square
			gamma := s3 / square

			// Point inside triangle
			if t.IntersectPoint(p) {
				fn(p, alpha, beta, gamma)
			}
		}
	}
}

func (t Triangle2D[T]) MinX() T {
	return mmath.Min(t.A.X, mmath.Min(t.B.X, t.C.X))
}

func (t Triangle2D[T]) MaxX() T {
	return mmath.Max(t.A.X, mmath.Max(t.B.X, t.C.X))
}

func (t Triangle2D[T]) Width() T {
	return mmath.Abs(t.MaxX() - t.MinX())
}

func (t Triangle2D[T]) MinY() T {
	return mmath.Min(t.A.Y, mmath.Min(t.B.Y, t.C.Y))
}

func (t Triangle2D[T]) MaxY() T {
	return mmath.Max(t.A.Y, mmath.Max(t.B.Y, t.C.Y))
}

func (t Triangle2D[T]) Height() T {
	return mmath.Abs(t.MaxY() - t.MinY())
}

func (t Triangle2D[T]) IntersectPoint(p mmath_la.Vector2[T]) bool {
	// First edge
	v1 := t.A
	v2 := t.B
	check := ((v2.X-v1.X)*(p.Y-v1.Y) - (v2.Y-v1.Y)*(p.X-v1.X)) > 0.0
	if !check {
		return false
	}

	// Second edge
	v1 = t.B
	v2 = t.C
	check = ((v2.X-v1.X)*(p.Y-v1.Y) - (v2.Y-v1.Y)*(p.X-v1.X)) > 0
	if !check {
		return false
	}

	// Third edge
	v1 = t.C
	v2 = t.A
	check = ((v2.X-v1.X)*(p.Y-v1.Y) - (v2.Y-v1.Y)*(p.X-v1.X)) > 0
	if !check {
		return false
	}

	return true
}

func (t Triangle2D[T]) Square() T {
	a := t.A
	b := t.B
	c := t.C

	return 0.5 * ((b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y))
}
