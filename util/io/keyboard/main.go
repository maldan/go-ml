package ml_keyboard

var State = map[int]bool{}

const (
	KeyA int = 65
	KeyD     = 68
	KeyW     = 87
	KeyS     = 83
)

func IsKeyDown(key int) bool {
	return State[key]
}
