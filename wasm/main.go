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
	/*mousemove := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Get("view").Get("innerWidth").Float()
		height := args[0].Get("view").Get("innerHeight").Float()

		px := (args[0].Get("pageX").Float()/args[0].Get("view").Get("innerWidth").Float())*2 - 1
		py := (args[0].Get("pageY").Float()/args[0].Get("view").Get("innerHeight").Float())*2 - 1
		ml_mouse.Position.X = float32(px)
		ml_mouse.Position.Y = -float32(py)

		// If there is canvas
		canvas := js.Global().Get("window").Get("go").Get("canvas")
		if !canvas.IsUndefined() {
			fx := canvas.Call("getBoundingClientRect").Get("width").Float() / width
			fy := canvas.Call("getBoundingClientRect").Get("height").Float() / height

			ml_mouse.Position.X /= float32(fx)
			ml_mouse.Position.Y /= float32(fy)
		}

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
	js.Global().Get("document").Call("addEventListener", "mouseup", mouseup)*/

	ExportFunction("setMousePosition", func(args []js.Value) any {
		ml_mouse.Position.X = float32(args[0].Float())
		ml_mouse.Position.Y = float32(args[1].Float())
		return nil
	})

	ExportFunction("setMouseDown", func(args []js.Value) any {
		ml_mouse.State[args[0].Int()] = args[1].Bool()
		return nil
	})
	ExportFunction("setMouseClick", func(args []js.Value) any {
		ml_mouse.ClickState[args[0].Int()] = args[1].Bool()
		return nil
	})
}

func ExportFunction(name string, fn func(a []js.Value) any) {
	if js.Global().Get("window").Get("go").IsUndefined() {
		js.Global().Get("window").Set("go", map[string]any{})
	}

	js.Global().Get("window").Get("go").Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		return fn(args)
	}))
}

func ExportFunction2(name string, fn func(this js.Value, a []js.Value) any) {
	if js.Global().Get("window").Get("go").IsUndefined() {
		js.Global().Get("window").Set("go", map[string]any{})
	}

	js.Global().Get("window").Get("go").Set(name, js.FuncOf(fn))
}

func InitRender(engine *mrender.RenderEngine) {
	state := map[string]any{
		"dynamicMeshLayer": engine.Main.GetState(),
		"staticMeshLayer":  engine.StaticMesh.GetState(),
		"pointLayer":       engine.Point.GetState(),
		"lineLayer":        engine.Line.GetState(),
		"textLayer":        engine.Text.GetState(),
		"uiLayer":          engine.UI.GetState(),
	}
	ExportFunction2("renderState", func(this js.Value, args []js.Value) any {
		state["dynamicMeshLayer"] = engine.Main.GetState()
		state["staticMeshLayer"] = engine.StaticMesh.GetState()
		state["pointLayer"] = engine.Point.GetState()
		state["lineLayer"] = engine.Line.GetState()
		state["textLayer"] = engine.Text.GetState()
		state["uiLayer"] = engine.UI.GetState()
		return state
	})

	ExportFunction2("renderFrame", func(this js.Value, args []js.Value) any {
		engine.Render()
		return nil
	})

	ExportFunction2("renderResize", func(this js.Value, args []js.Value) any {
		engine.ScreenSize.X = float32(args[0].Float())
		engine.ScreenSize.Y = float32(args[1].Float())

		engine.GlobalCamera.AspectRatio = float32(args[0].Float() / args[1].Float())
		engine.UI.Camera.Area.Right = float32(args[0].Float())
		engine.UI.Camera.Area.Bottom = float32(args[1].Float())
		return nil
	})
}

func InitSound() {
	ExportFunction("soundState", func(args []js.Value) any {
		return maudio.GetState()
	})
	ExportFunction("soundTick", func(args []js.Value) any {
		maudio.Tick(args[0].Int())
		return nil
	})
}
