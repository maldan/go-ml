package ml_keyboard

var State = map[string]bool{}

const (
	KeyA  string = "KeyA"
	KeyD         = "KeyD"
	KeyW         = "KeyW"
	KeyS         = "KeyS"
	KeyE         = "KeyE"
	Space        = "Space"
)

func IsKeyDown(key string) bool {
	return State[key]
}
