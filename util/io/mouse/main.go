package ml_mouse

import ml_geom "github.com/maldan/go-ml/util/math/geom"

var State = map[int]bool{}
var Position = ml_geom.Vector2[float32]{}

const (
	LeftButton   int = 0
	RightButton  int = 1
	MiddleButton int = 2
)

func IsMouseDown(key int) bool {
	return State[key]
}

func GetPosition() ml_geom.Vector2[float32] {
	return Position
}
