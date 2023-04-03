package ml_mouse

import mgeom "github.com/maldan/go-ml/math/geom"

var State = map[int]bool{}
var Position = mgeom.Vector2[float32]{}

const (
	LeftButton   int = 0
	RightButton  int = 1
	MiddleButton int = 2
)

func IsMouseDown(key int) bool {
	return State[key]
}

func GetPosition() mgeom.Vector2[float32] {
	return Position
}
