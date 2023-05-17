package mmath_geom

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
	"math"
)

/*type Rectangle[T constraints.Float] struct {
	Left   T
	Top    T
	Right  T
	Bottom T
}

func (r Rectangle[T]) CenterXY() (T, T) {
	l := r.Left + r.Width()/2
	t := r.Top + r.Height()/2
	return l, t
}

func (r Rectangle[T]) Center() mmath_la.Vector2[T] {
	l, t := r.CenterXY()
	return mmath_la.Vector2[T]{l, t}
}

func (r Rectangle[T]) Width() T {
	return mmath.Abs(r.Right - r.Left)
}

func (r Rectangle[T]) Height() T {
	return mmath.Abs(r.Top - r.Bottom)
}

func (r Rectangle[T]) AddXY(x T, y T) Rectangle[T] {
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
*/

type Rectangle[T constraints.Float] struct {
	MinX T
	MaxX T
	MinY T
	MaxY T
}

func (r Rectangle[T]) Shrink(v T) Rectangle[T] {
	r.MinX += v
	r.MaxX -= v
	r.MinY += v
	r.MaxY -= v
	return r
}

func (r Rectangle[T]) Expand(v T) Rectangle[T] {
	r.MinX -= v
	r.MaxX += v
	r.MinY -= v
	r.MaxY += v
	return r
}

func (r Rectangle[T]) CenterXY() (T, T) {
	l := r.MinX + r.Width()/2
	t := r.MinY + r.Height()/2
	return l, t
}

func (r Rectangle[T]) Min() mmath_la.Vector2[T] {
	return mmath_la.Vector2[T]{r.MinX, r.MinY}
}

func (r Rectangle[T]) Max() mmath_la.Vector2[T] {
	return mmath_la.Vector2[T]{r.MaxX, r.MaxY}
}

func (r Rectangle[T]) Center() mmath_la.Vector2[T] {
	l, t := r.CenterXY()
	return mmath_la.Vector2[T]{l, t}
}

func (r Rectangle[T]) Width() T {
	return mmath.Abs(r.MaxX - r.MinX)
}

func (r Rectangle[T]) Height() T {
	return mmath.Abs(r.MaxY - r.MinY)
}

func (r Rectangle[T]) Add(v mmath_la.Vector2[T]) Rectangle[T] {
	return r.AddXY(v.X, v.Y)
}

func (r Rectangle[T]) Sub(v mmath_la.Vector2[T]) Rectangle[T] {
	return r.SubXY(v.X, v.Y)
}

func (r Rectangle[T]) AddXY(x T, y T) Rectangle[T] {
	r.MinX += x
	r.MaxX += x
	r.MinY += y
	r.MaxY += y
	return r
}

func (r Rectangle[T]) Floor() Rectangle[T] {
	r.MinX = T(math.Floor(float64(r.MinX)))
	r.MaxX = T(math.Floor(float64(r.MaxX)))
	r.MinY = T(math.Floor(float64(r.MinY)))
	r.MaxY = T(math.Floor(float64(r.MaxY)))
	return r
}

func (r Rectangle[T]) Ceil() Rectangle[T] {
	r.MinX = T(math.Ceil(float64(r.MinX)))
	r.MaxX = T(math.Ceil(float64(r.MaxX)))
	r.MinY = T(math.Ceil(float64(r.MinY)))
	r.MaxY = T(math.Ceil(float64(r.MaxY)))
	return r
}

func (r Rectangle[T]) SubXY(x T, y T) Rectangle[T] {
	r.MinX -= x
	r.MaxX -= x
	r.MinY -= y
	r.MaxY -= y
	return r
}

func (r Rectangle[T]) Intersect(r2 Rectangle[T]) bool {
	return r.MinX < r2.MaxX &&
		r.MaxX > r2.MinX &&
		r.MinY < r2.MaxY &&
		r.MaxY > r2.MinY
	/*return !(r2.Left > r.Right ||
	r2.Right < r.Left ||
	r2.Top > r.Bottom ||
	r2.Bottom < r.Top)*/
}

func (r Rectangle[T]) Inside(r2 Rectangle[T]) bool {
	if r.MinX < r2.MinX {
		return false
	}
	if r.MaxX > r2.MaxX {
		return false
	}
	if r.MinY < r2.MinY {
		return false
	}
	if r.MaxY > r2.MaxY {
		return false
	}

	return true
}

func (r Rectangle[T]) IntersectPoint(p mmath_la.Vector2[T]) bool {
	if p.X > r.MinX && p.X < r.MaxX && p.Y > r.MinY && p.Y < r.MaxY {
		return true
	}
	return false
}
