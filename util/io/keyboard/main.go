package ml_keyboard

var State = map[string]bool{}
var PressState = map[string]bool{}

const (
	KeyA      string = "KeyA"
	KeyD             = "KeyD"
	KeyW             = "KeyW"
	KeyS             = "KeyS"
	KeyE             = "KeyE"
	KeyP             = "KeyP"
	KeyQ             = "KeyQ"
	Space            = "Space"
	ArrowDown        = "ArrowDown"
	ArrowUp          = "ArrowUp"
)

func IsKeyDown(key string) bool {
	return State[key]
}

func IsKeyPressed(key string) bool {
	return PressState[key]
}

func ResetPressState() {
	for k, _ := range PressState {
		PressState[k] = false
	}
}
