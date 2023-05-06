package mmath

import (
	"golang.org/x/exp/constraints"
	"math"
)

func Min[T constraints.Integer | constraints.Float](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Integer | constraints.Float](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T constraints.Integer | constraints.Float](a T) T {
	if a < 0 {
		return -a
	}

	return a
}

func CeilInt[T constraints.Float](a T) int {
	return int(math.Ceil(float64(a)))
}

func Ceil[T constraints.Float](a T) T {
	return T(math.Ceil(float64(a)))
}
