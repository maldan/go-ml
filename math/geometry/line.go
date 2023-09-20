package mmath_geom

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
)

type Line2D[T constraints.Float] struct {
	From mmath_la.Vector2[float32]
	To   mmath_la.Vector2[float32]
}

func (l Line2D[T]) Intersect(l2 Line2D[T]) bool {
	p0_x := l.From.X
	p0_y := l.From.Y
	p1_x := l.To.X
	p1_y := l.To.Y
	p2_x := l2.From.X
	p2_y := l2.From.Y
	p3_x := l2.To.X
	p3_y := l2.To.Y

	s1_x := p1_x - p0_x
	s1_y := p1_y - p0_y
	s2_x := p3_x - p2_x
	s2_y := p3_y - p2_y

	s := (-s1_y*(p0_x-p2_x) + s1_x*(p0_y-p2_y)) / (-s2_x*s1_y + s1_x*s2_y)
	t := (s2_x*(p0_y-p2_y) - s2_y*(p0_x-p2_x)) / (-s2_x*s1_y + s1_x*s2_y)

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		return true
	}

	return false
}

type Line3D[T constraints.Float] struct {
	From mmath_la.Vector3[T]
	To   mmath_la.Vector3[T]
}
