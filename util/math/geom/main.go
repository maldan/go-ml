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
