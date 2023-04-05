package mwasm

import (
	maudio "github.com/maldan/go-ml/audio"
	mrender "github.com/maldan/go-ml/render"
	ml_keyboard "github.com/maldan/go-ml/util/io/keyboard"
	ml_mouse "github.com/maldan/go-ml/util/io/mouse"
	"syscall/js"
)

func BindKeyboard() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ml_keyboard.State[args[0].Get("code").String()] = true
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "keydown", cb)

	cb2 := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ml_keyboard.State[args[0].Get("code").String()] = false
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "keyup", cb2)
}

func BindMouse() {
	mousemove := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		px := (args[0].Get("pageX").Float()/args[0].Get("view").Get("innerWidth").Float())*2 - 1
		py := (args[0].Get("pageY").Float()/args[0].Get("view").Get("innerHeight").Float())*2 - 1
		ml_mouse.Position.X = float32(px)
		ml_mouse.Position.Y = -float32(py)
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "mousemove", mousemove)

	mousedown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ml_mouse.State[args[0].Get("button").Int()] = true
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "mousedown", mousedown)

	mouseup := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ml_mouse.State[args[0].Get("button").Int()] = false
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "mouseup", mouseup)
}

func ExportFunction(name string, fn func(args []js.Value)) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		fn(args)
		return nil
	}))
}

func InitRender(engine *mrender.RenderEngine) {
	js.Global().Set("goWasmRenderState", js.FuncOf(func(this js.Value, args []js.Value) any {
		return map[string]any{
			"mainLayer":  engine.Main.GetState(),
			"pointLayer": engine.Point.GetState(),
			"lineLayer":  engine.Line.GetState(),
		}
	}))

	js.Global().Set("goWasmRenderFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		engine.Render()
		return nil
	}))
}

func InitSound() {
	js.Global().Set("goWasmSoundState", js.FuncOf(func(this js.Value, args []js.Value) any {
		return maudio.GetState()
	}))

	js.Global().Set("goWasmSoundTick", js.FuncOf(func(this js.Value, args []js.Value) any {
		maudio.Tick(args[0].Int())
		return nil
	}))
}
