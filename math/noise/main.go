package mmath_noise

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_random "github.com/maldan/go-ml/util/random"
	"math"
)

/* Function to linearly interpolate between a0 and a1
 * Weight w should be in the range [0.0, 1.0]
 */
func interpolate(a0 float32, a1 float32, w float32) float32 {
	return (a1-a0)*w + a0

}

/* Create pseudorandom direction vector
 */
func randomGradient(ix int, iy int) mmath_la.Vector2[float32] {
	// No precomputed gradients mean this works for any number of grid coordinates
	w := 8 * 4
	s := w / 2 // rotation width
	a := ix
	b := iy
	a *= 3284157443
	b ^= a<<s | a>>w - s
	b *= 1911520717
	a ^= b<<s | b>>w - s
	a *= 2048419325
	random := ml_random.Range[float32](0, 3.1415*2) // in [0, 2*Pi]
	v := mmath_la.Vector2[float32]{}
	v.X = float32(math.Cos(float64(random)))
	v.Y = float32(math.Sin(float64(random)))
	return v
}

// Computes the dot product of the distance and gradient vectors.
func dotGridGradient(ix int, iy int, x float32, y float32) float32 {
	// Get gradient from integer coordinates
	gradient := randomGradient(ix, iy)

	// Compute the distance vector
	dx := x - float32(ix)
	dy := y - float32(iy)

	// Compute the dot-product
	return (dx*gradient.X + dy*gradient.Y)
}

// Compute Perlin noise at coordinates x, y
func Perlin(x float32, y float32) float32 {
	// Determine grid cell coordinates
	x0 := int(math.Floor(float64(x)))
	x1 := x0 + 1
	y0 := int(math.Floor(float64(y)))
	y1 := y0 + 1

	// Determine interpolation weights
	// Could also use higher order polynomial/s-curve here
	sx := x - float32(x0)
	sy := y - float32(y0)

	n0 := dotGridGradient(x0, y0, x, y)
	n1 := dotGridGradient(x1, y0, x, y)
	ix0 := interpolate(n0, n1, sx)

	n0 = dotGridGradient(x0, y1, x, y)
	n1 = dotGridGradient(x1, y1, x, y)
	ix1 := interpolate(n0, n1, sx)

	value := interpolate(ix0, ix1, sy)
	return value // Will return in range -1 to 1. To make it in range 0 to 1, multiply by 0.5 and add 0.5
}
