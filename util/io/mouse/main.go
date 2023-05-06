package ml_mouse

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

var ClickState = [...]bool{false, false, false, false, false, false, false, false}
var State = [...]bool{false, false, false, false, false, false, false, false}
var Position = mmath_la.Vector2[float32]{}

const (
	LeftButton   int = 0
	RightButton  int = 1
	MiddleButton int = 2
)

func IsMouseClick(key int) bool {
	if key > len(ClickState)-1 {
		return false
	}
	return ClickState[key]
}

func IsMouseDown(key int) bool {
	if key > len(State)-1 {
		return false
	}
	return State[key]
}

func ResetClickState() {
	for i := 0; i < len(ClickState); i++ {
		ClickState[i] = false
	}
}

func GetPosition() mmath_la.Vector2[float32] {
	return Position
}
