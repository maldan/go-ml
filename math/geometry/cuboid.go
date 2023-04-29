package mmath_geom

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_number "github.com/maldan/go-ml/util/number"
	"golang.org/x/exp/constraints"
)

type Cuboid[T constraints.Float] struct {
	/*Left  T
	Top   T
	Front T

	Right  T
	Bottom T
	Back   T*/

	MinX T
	MaxX T
	MinY T
	MaxY T
	MinZ T
	MaxZ T
}

func (c Cuboid[T]) Translate(v mmath_la.Vector3[T]) Cuboid[T] {
	c.MinX += v.X
	c.MaxX += v.X
	c.MinY += v.Y
	c.MaxY += v.Y
	c.MinZ += v.Z
	c.MaxZ += v.Z
	return c
}

func (c Cuboid[T]) Intersect(c2 Cuboid[T]) bool {
	return c.MinX < c2.MaxX &&
		c.MaxX > c2.MinX &&
		c.MinY < c2.MaxY &&
		c.MaxY > c2.MinY &&
		c.MinZ < c2.MaxZ &&
		c.MaxZ > c2.MinZ
}

func (c Cuboid[T]) Overlap(c2 Cuboid[T]) Cuboid[T] {
	return Cuboid[T]{
		MinX: ml_number.Max(c.MinX, c2.MinX),
		MinY: ml_number.Max(c.MinY, c2.MinY),
		MinZ: ml_number.Max(c.MinZ, c2.MinZ),
		MaxX: ml_number.Min(c.MaxX, c2.MaxX),
		MaxY: ml_number.Min(c.MaxY, c2.MaxY),
		MaxZ: ml_number.Min(c.MaxZ, c2.MaxZ),
	}
}

func (c Cuboid[T]) ToRelative() Cuboid[T] {
	cx := c.MaxX - c.MinX
	cy := c.MaxY - c.MinY
	cz := c.MaxZ - c.MinZ

	return Cuboid[T]{
		MaxX: cx,
		MaxY: cy,
		MaxZ: cz,
	}
}
