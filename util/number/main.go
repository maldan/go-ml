package ml_number

import (
	"golang.org/x/exp/constraints"
)

func Remap[T constraints.Float | constraints.Integer](value T, low1 T, high1 T, low2 T, high2 T) T {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

var NibbleLookup = []uint8{
	0, 1, 1, 2, 1, 2, 2, 3,
	1, 2, 2, 3, 2, 3, 3, 4,
}

func CountSetBits(byte uint8) uint8 {
	return NibbleLookup[byte&0x0F] + NibbleLookup[byte>>4]
}

func CheckBitMask[T constraints.Integer](v T, mask T) bool {
	return v&mask == mask
}

func Lerp[T constraints.Integer | constraints.Float](start T, end T, t T) T {
	return (1-t)*start + t*end
}
