package mmath_geom

import "golang.org/x/exp/constraints"

type Rectangle[T constraints.Float] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

func (r *Rectangle[T]) Offset(x T, y T) {
	r.Left += x
	r.Right += x
	r.Top += y
	r.Bottom += y
}

func (r Rectangle[T]) Intersect(r2 Rectangle[T]) bool {
	return !(r2.Left > r.Right ||
		r2.Right < r.Left ||
		r2.Top > r.Bottom ||
		r2.Bottom < r.Top)
}

func IsRectanglesIntersect[T constraints.Float](r1 Rectangle[T], r2 Rectangle[T]) bool {
	return !(r2.Left > r1.Right || r2.Right < r1.Left || r2.Top > r1.Bottom || r2.Bottom < r1.Top)
}
