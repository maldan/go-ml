package ml_gamepad

var KeyState0 = [...]bool{
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}
var KeyPressState0 = [...]bool{
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false,
}
var AxisState0 = [...]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

const (
	KeyA     uint8 = 0
	KeyB           = 1
	KeyX           = 2
	KeyY           = 3
	KeyLB          = 4
	KeyRB          = 5
	KeyLT          = 6
	KeyRT          = 7
	KeyPause       = 8
	KeyMenu        = 9
	KeyL           = 10
	KeyR           = 11
	KeyTop         = 12
	KeyDown        = 13
	KeyLeft        = 14
	KeyRight       = 15
)

func Axis(id uint8, axis uint8) float32 {
	return AxisState0[axis]
}

func IsKeyDown(id uint8, key uint8) bool {
	return KeyState0[key]
}

func IsKeyPressed(key uint8) bool {
	return KeyPressState0[key]
}

func ResetPressState() {
	for id := range KeyPressState0 {
		KeyPressState0[id] = false
	}
}
