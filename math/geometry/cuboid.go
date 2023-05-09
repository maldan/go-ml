package mmath_geom

import (
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"golang.org/x/exp/constraints"
	"math"
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

func (c Cuboid[T]) Add(v mmath_la.Vector3[T]) Cuboid[T] {
	c.MinX += v.X
	c.MaxX += v.X
	c.MinY += v.Y
	c.MaxY += v.Y
	c.MinZ += v.Z
	c.MaxZ += v.Z
	return c
}

func (c Cuboid[T]) Round() Cuboid[T] {
	c.MinX = T(math.Round(float64(c.MinX)))
	c.MaxX = T(math.Round(float64(c.MaxX)))
	c.MinY = T(math.Round(float64(c.MinY)))
	c.MaxY = T(math.Round(float64(c.MaxY)))
	c.MinZ = T(math.Round(float64(c.MinZ)))
	c.MaxZ = T(math.Round(float64(c.MaxZ)))
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
		MinX: mmath.Max(c.MinX, c2.MinX),
		MinY: mmath.Max(c.MinY, c2.MinY),
		MinZ: mmath.Max(c.MinZ, c2.MinZ),
		MaxX: mmath.Min(c.MaxX, c2.MaxX),
		MaxY: mmath.Min(c.MaxY, c2.MaxY),
		MaxZ: mmath.Min(c.MaxZ, c2.MaxZ),
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

func (c Cuboid[T]) SizeX() T {
	return mmath.Abs(c.MaxX - c.MinX)
}

func (c Cuboid[T]) SizeY() T {
	return mmath.Abs(c.MaxY - c.MinY)
}

func (c Cuboid[T]) SizeZ() T {
	return mmath.Abs(c.MaxZ - c.MinZ)
}

func (c Cuboid[T]) Size() mmath_la.Vector3[T] {
	return mmath_la.Vector3[T]{
		X: c.SizeX(),
		Y: c.SizeY(),
		Z: c.SizeZ(),
	}
}
