package ml_wasm

import (
	ml_keyboard "github.com/maldan/go-ml/util/io/keyboard"
	"syscall/js"
)

func BindKeyboard() {
	js.Global().Set("__golangWebglKeyboardBind_KeyDown", js.FuncOf(func(this js.Value, args []js.Value) any {
		ml_keyboard.State[args[0].Int()] = true
		return nil
	}))
	js.Global().Set("__golangWebglKeyboardBind_KeyUp", js.FuncOf(func(this js.Value, args []js.Value) any {
		ml_keyboard.State[args[0].Int()] = false
		return nil
	}))
}
