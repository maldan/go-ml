package ml_number

import (
	mmath "github.com/maldan/go-ml/math"
	"golang.org/x/exp/constraints"
)

var NibbleLookup = []uint8{
	0, 1, 1, 2, 1, 2, 2, 3,
	1, 2, 2, 3, 2, 3, 3, 4,
}

func Remap[T constraints.Float | constraints.Integer](value T, low1 T, high1 T, low2 T, high2 T) T {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func CountSetBits(byte uint8) uint8 {
	return NibbleLookup[byte&0x0F] + NibbleLookup[byte>>4]
}

func CheckBitMask[T constraints.Integer](v T, mask T) bool {
	return v&mask == mask
}

/*func Clamp[T constraints.Integer | constraints.Float](v T, min T, max T) T {
	if v <= min {
		return min
	}
	if v >= max {
		return max
	}
	return v
}*/

func TowardsSmooth[T constraints.Float](from *T, to T, step T, delta T) {
	step = 1 / step
	*from += (to - *from) / (step / delta)
}

func MoveTowards[T constraints.Float](from *T, to T, step T) {
	if mmath.Abs(*from-to) <= step {
		*from = to
		return
	}

	if *from > to {
		*from -= step
	} else {
		*from += step
	}
}
