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

func FloorInt[T constraints.Float](a T) int {
	return int(a)
}

func RoundInt(a float32) int {
	return int(Round(a))
}

func Ceil[T constraints.Integer | constraints.Float](a T) T {
	return T(math.Ceil(float64(a)))
}

func Floor[T constraints.Float](a T) T {
	return T(int(a))
	// return T(math.Floor(float64(a)))
}

func Round(x float32) float32 {
	// Round is a faster implementation of:
	//
	// func Round(x float64) float64 {
	//   t := Trunc(x)
	//   if Abs(x-t) >= 0.5 {
	//     return t + Copysign(1, x)
	//   }
	//   return t
	// }

	bits := Float32bits(x)
	e := uint(bits>>shift) & mask
	if e < bias {
		// Round abs(x) < 1 including denormals.
		bits &= signMask // +-0
		if e == bias-1 {
			bits |= uvone // +-1
		}
	} else if e < bias+shift {
		// Round any abs(x) >= 1 containing a fractional component [0,1).
		//
		// Numbers with larger exponents are returned unchanged since they
		// must be either an integer, infinity, or NaN.
		const half = 1 << (shift - 1)
		e -= bias
		bits += half >> e
		bits &^= fracMask >> e
	}
	return Float32frombits(bits)
	//return T(math.Round(float64(a)))
}

func Mod[T constraints.Float](x, y T) T {
	result := x - y*T(int(x/y))
	if result < 0 {
		result += y
	}
	return result
}

func Clamp[T constraints.Integer | constraints.Float](v T, min T, max T) T {
	if v <= min {
		return min
	}
	if v >= max {
		return max
	}
	return v
}

func Clamp01[T constraints.Integer | constraints.Float](v T) T {
	return Clamp(v, 0, 1)
}

func Lerp[T constraints.Integer | constraints.Float](start T, end T, t T) T {
	return (1-t)*start + t*end
}

func Remap[T constraints.Float | constraints.Integer](value T, low1 T, high1 T, low2 T, high2 T) T {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func DegToRad[T constraints.Float](v T) T {
	return v * (T(Pi) / 180.0)
}
