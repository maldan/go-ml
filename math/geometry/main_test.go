package mmath_geom

import (
	"fmt"
	mmath "github.com/maldan/go-ml/math"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"testing"
)

func TestZ(t *testing.T) {
	tr := Triangle2D[float32]{
		A: mmath_la.Vec2f(0, 0),
		B: mmath_la.Vec2f(10, 0),
		C: mmath_la.Vec2f(12, -10),
	}
	top, middle, bottom := tr.TopMiddleBottom()

	// Natural flat top
	if middle.Y == bottom.Y {
		if bottom.X < middle.X {
			middle, bottom = bottom, middle
		}
		top, middle, bottom = middle, bottom, top
		// bottom, middle = middle, bottom

		fmt.Printf("%v %v %v\n", top, middle, bottom)

		// First slope
		dx := top.X - bottom.X
		dy := bottom.Y - top.Y
		slope1 := dx / dy
		if dy == 0 {
			slope1 = 0
		}
		fmt.Printf("Slope 1: %v\n", slope1)

		// Second slope
		dx = bottom.X - middle.X
		dy = middle.Y - bottom.Y
		slope2 := dx / dy
		if dy == 0 {
			slope2 = 0
		}
		fmt.Printf("Slope 2: %v\n", slope2)

		startY, endY := int(top.Y), int(bottom.Y)
		height := startY - endY
		for y := 0; y < height; y++ {
			index := float32(y)
			xStart := mmath.RoundInt((top.X) + index*slope1)
			xEnd := mmath.RoundInt((middle.X) + index*slope2)

			fmt.Printf("%v - %v - [%v]\n", xStart, xEnd, index)
		}
	}
}

func TestY(t *testing.T) {
	tr := Triangle2D[float32]{
		A: mmath_la.Vec2f(0, 0),
		B: mmath_la.Vec2f(10, 0),
		C: mmath_la.Vec2f(10, 10),
	}
	top, middle, bottom := tr.TopMiddleBottom()

	// Natural flat top
	if top.Y == middle.Y {
		if top.X > middle.X {
			top, middle = middle, top
		}

		fmt.Printf("%v %v %v\n", top, middle, bottom)

		// First slope
		dx := bottom.X - top.X
		dy := bottom.Y - top.Y
		slope1 := dx / dy
		if dy == 0 {
			slope1 = 0
		}
		fmt.Printf("Slope 1: %v\n", slope1)

		// Second slope
		dx = bottom.X - middle.X
		dy = bottom.Y - middle.Y
		slope2 := dx / dy
		if dy == 0 {
			slope2 = 0
		}
		fmt.Printf("Slope 2: %v\n", slope2)

		startY, endY := int(top.Y), int(bottom.Y)
		height := endY - startY
		for y := 0; y < height; y++ {
			index := float32(y)
			xStart := mmath.RoundInt((top.X) + index*slope1)
			xEnd := mmath.RoundInt((middle.X) + index*slope2)

			fmt.Printf("%v - %v - [%v]\n", xStart, xEnd, index)
		}
	}
}

func TestX(t *testing.T) {
	// {329.3846 423.8154} {332.49133 420.11914} {328.30826 425.6057}}
	// {308.60156 402.1911} {308.6596 403.19974} {308.09467 403.33582}
	tr := Triangle2D[float32]{
		A: mmath_la.Vec2f(308.60156, 402.1911),
		B: mmath_la.Vec2f(308.6596, 403.19974),
		C: mmath_la.Vec2f(308.09467, 403.33582),
	}
	top, middle, bottom := tr.TopMiddleBottom()
	top = top.Floor()
	middle = middle.Floor()
	bottom = bottom.Floor()

	fmt.Printf("%v %v %v\n", top, middle, bottom)
	slope1 := (middle.X - top.X) / (middle.Y - top.Y)
	slope2 := (bottom.X - top.X) / (bottom.Y - top.Y)
	fmt.Printf("Slope1: %v\n", slope1)
	fmt.Printf("Slope2: %v\n", slope2)

	//ll := (bottom.Y - top.Y) / (middle.Y - top.Y)
	//fmt.Printf("%v\n", ll)

	for y := int(top.Y); y <= int(middle.Y); y += 1 {
		yy := float32(y)
		xStart := top.X + (yy-top.Y)*slope1
		xEnd := (top.X + (yy-top.Y)*(slope2))
		fmt.Printf("<X: %v X>: %v | Y: %v\n", xStart, xEnd, y)
	}
	fmt.Printf("\n")
	slope1 = (bottom.X - middle.X) / (bottom.Y - middle.Y)
	fmt.Printf("Slope1: %v\n", slope1)
	for y := int(middle.Y); y <= int(bottom.Y); y++ {
		yy := float32(y)
		xStart := ((middle.X) + (yy-middle.Y)*slope1)
		xEnd := ((top.X) + (yy-top.Y)*slope2)
		fmt.Printf("<X: %v X>: %v | Y: %v\n", xStart, xEnd, y)
		fmt.Printf("%v - %v - %v\n", middle.X, yy-middle.Y, slope1)
		fmt.Printf("W: %v\n", mmath.Abs(int(xEnd)-int(xStart)))
	}
}
