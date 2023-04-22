package mmath_geom

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
)

type Rectangle[T constraints.Float] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

func (r Rectangle[T]) Add(x T, y T) Rectangle[T] {
	r.Left += x
	r.Right += x
	r.Top += y
	r.Bottom += y
	return r
}

func (r Rectangle[T]) Intersect(r2 Rectangle[T]) bool {
	return !(r2.Left > r.Right ||
		r2.Right < r.Left ||
		r2.Top > r.Bottom ||
		r2.Bottom < r.Top)
}

func (r Rectangle[T]) IntersectPoint(p mmath_la.Vector2[T]) bool {
	if p.X > r.Left && p.X < r.Right && p.Y > r.Top && p.Y < r.Bottom {
		return true
	}
	return false
}

func IsRectanglesIntersect[T constraints.Float](r1 Rectangle[T], r2 Rectangle[T]) bool {
	return !(r2.Left > r1.Right || r2.Right < r1.Left || r2.Top > r1.Bottom || r2.Bottom < r1.Top)
}
