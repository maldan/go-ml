package ml_mouse

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

var ClickState = [4]bool{false, false, false, false}
var State = [4]bool{false, false, false, false}
var Position = mmath_la.Vector2[float32]{}

const (
	LeftButton   int = 0
	RightButton  int = 1
	MiddleButton int = 2
)

func IsMouseClick(key int) bool {
	return ClickState[key]
}

func IsMouseDown(key int) bool {
	return State[key]
}

func GetPosition() mmath_la.Vector2[float32] {
	return Position
}
