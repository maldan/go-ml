package mmath_geom

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
	"math"
)

type Triangle2D[T constraints.Float] struct {
	A mmath_la.Vector2[T]
	B mmath_la.Vector2[T]
	C mmath_la.Vector2[T]
}

func (t Triangle2D[T]) OnEachPixel(fn func(v mmath_la.Vector2[T], a T, b T, g T)) int {
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

	square := t.Area()

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

	return w * h
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

func (t Triangle2D[T]) Area() T {
	x1 := t.A.X
	x2 := t.B.X
	x3 := t.C.X

	y1 := t.A.Y
	y2 := t.B.Y
	y3 := t.C.Y

	//return 0.5 * ((b.X-a.X)*(c.Y-a.Y) - (c.X-a.X)*(b.Y-a.Y))
	return 0.5 * mmath.Abs(x1*(y2-y3)+x2*(y3-y1)+x3*(y1-y2))
}

func (t Triangle2D[T]) InterpolateABG(alpha T, beta T, gamma T) mmath_la.Vector2[T] {
	uva := t.A.MulScalar(alpha)
	uvb := t.B.MulScalar(beta)
	uvc := t.C.MulScalar(gamma)
	return (uva.Add(uvb)).Add(uvc)
}

func (t Triangle2D[T]) TopMiddleBottom() (mmath_la.Vector2[T], mmath_la.Vector2[T], mmath_la.Vector2[T]) {
	var top, middle, bottom mmath_la.Vector2[T]

	x1 := t.A.X
	x2 := t.B.X
	x3 := t.C.X
	y1 := t.A.Y
	y2 := t.B.Y
	y3 := t.C.Y

	if y1 <= y2 && y1 <= y3 {
		top = mmath_la.Vector2[T]{x1, y1}
		if y2 <= y3 {
			middle = mmath_la.Vector2[T]{x2, y2}
			bottom = mmath_la.Vector2[T]{x3, y3}
		} else {
			middle = mmath_la.Vector2[T]{x3, y3}
			bottom = mmath_la.Vector2[T]{x2, y2}
		}
		return top, middle, bottom
	}

	if y2 <= y1 && y2 <= y3 {
		top = mmath_la.Vector2[T]{x2, y2}
		if y1 <= y3 {
			middle = mmath_la.Vector2[T]{x1, y1}
			bottom = mmath_la.Vector2[T]{x3, y3}
		} else {
			middle = mmath_la.Vector2[T]{x3, y3}
			bottom = mmath_la.Vector2[T]{x1, y1}
		}
		return top, middle, bottom
	}

	top = mmath_la.Vector2[T]{x3, y3}
	if y1 <= y2 {
		middle = mmath_la.Vector2[T]{x1, y1}
		bottom = mmath_la.Vector2[T]{x2, y2}
	} else {
		middle = mmath_la.Vector2[T]{x2, y2}
		bottom = mmath_la.Vector2[T]{x1, y1}
	}
	return top, middle, bottom
}

/*func (t Triangle2D[T]) Top() mmath_la.Vector2[T] {
	A := t.A
	B := t.B
	C := t.C

	if A.Y < B.Y && A.Y < C.Y {
		return A
	}

	if B.Y < A.Y && B.Y < C.Y {
		return B
	}

	return C
}

func (t Triangle2D[T]) Middle() mmath_la.Vector2[T] {
	A := t.A
	B := t.B
	C := t.C

	if A.Y < B.Y && A.Y > C.Y {
		return A
	}

	if B.Y > A.Y && B.Y < C.Y {
		return B
	}

	return C
}

func (t Triangle2D[T]) Bottom() mmath_la.Vector2[T] {
	A := t.A
	B := t.B
	C := t.C

	if A.Y > B.Y && A.Y > C.Y {
		return A
	}

	if B.Y > A.Y && B.Y > C.Y {
		return B
	}

	return C
}*/

func (t Triangle2D[T]) SplitY() (top Triangle2D[T], bottom Triangle2D[T]) {
	A := t.A
	B := t.B
	C := t.C

	_top := mmath_la.Vector2[T]{}
	_middle := mmath_la.Vector2[T]{}
	_bottom := mmath_la.Vector2[T]{}

	// If A top
	if A.Y < B.Y && A.Y < C.Y {
		_top = A

		if B.Y < C.Y {
			_middle = B
			_bottom = C
		} else {
			_middle = C
			_bottom = B
		}

		return Triangle2D[T]{A: _top, B: _middle, C: _bottom}, Triangle2D[T]{A: _middle, B: _bottom}
	}

	// If B top
	if B.Y < A.Y && B.Y < C.Y {
		_top = B
	}

	// If C top
	if C.Y < A.Y && C.Y < B.Y {
		_top = C
	}

	return Triangle2D[T]{A: _top, B: _middle}, Triangle2D[T]{A: _middle, B: _bottom}
}
