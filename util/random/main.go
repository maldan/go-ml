package ml_random

import (
	"golang.org/x/exp/constraints"
	"math/rand"
)

func SetSeed(seed int64) {
	rand.Seed(seed)
}

func RangeInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func Range[T constraints.Float](min T, max T) T {
	return min + T(rand.Float64())*(max-min)
}

/*func RandFloat32(min float32, max float32) float32 {
	return remap(rand.Float32(), 0, 1, min, max)
}

func remap[T constraints.Float | constraints.Integer](value T, low1 T, high1 T, low2 T, high2 T) T {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}*/
