package ml_geom

import "golang.org/x/exp/constraints"

type Vector2[T constraints.Float] struct {
	X T
	Y T
}

type Rectangle[T constraints.Float] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

func (r *Rectangle[T]) Add(v Vector2[T]) {
	r.Left += v.X
	r.Right += v.X
	r.Top += v.Y
	r.Bottom += v.Y
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
