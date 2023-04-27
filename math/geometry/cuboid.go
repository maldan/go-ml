package mmath_geom

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
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
	return c.MinX <= c2.MaxX &&
		c.MaxX >= c2.MinX &&
		c.MinY <= c2.MaxY &&
		c.MaxY >= c2.MinY &&
		c.MinZ <= c2.MaxZ &&
		c.MaxZ >= c2.MinZ
}
