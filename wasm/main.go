package mwasm

import (
	maudio "github.com/maldan/go-ml/audio"
	mrender "github.com/maldan/go-ml/render"
	ml_keyboard "github.com/maldan/go-ml/util/io/keyboard"
	ml_mouse "github.com/maldan/go-ml/util/io/mouse"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"
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

	/*ExportFunction("setMousePosition", func(args []js.Value) any {
		ml_mouse.Position.X = float32(args[0].Float())
		ml_mouse.Position.Y = float32(args[1].Float())
		return nil
	})

	ExportFunction("setMouseDown", func(args []js.Value) any {
		ml_mouse.State[args[0].Int()] = args[1].Bool()
		return nil
	})*/

	ExportFunction("setMouseClick", func(args []js.Value) any {
		ml_mouse.ClickState[args[0].Int()] = args[1].Bool()
		return nil
	})

	ExportPointer("mousePosition", unsafe.Pointer(&ml_mouse.Position))
	ExportPointer("mouseDown", unsafe.Pointer(&ml_mouse.State))
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

func ExportPointer(name string, pointer unsafe.Pointer) {
	if js.Global().Get("window").Get("go").IsUndefined() {
		js.Global().Get("window").Set("go", map[string]any{})
	}
	if js.Global().Get("window").Get("go").Get("pointer").IsUndefined() {
		js.Global().Get("window").Get("go").Set("pointer", map[string]any{})
	}
	js.Global().Get("window").Get("go").Get("pointer").Set(name, uintptr(pointer))
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

		engine.Camera.AspectRatio = float32(args[0].Float() / args[1].Float())
		engine.UICamera.Area.MaxX = float32(args[0].Float())
		engine.UICamera.Area.MaxY = float32(args[1].Float())
		return nil
	})

	// Render state
	ExportPointer("renderState", unsafe.Pointer(&engine.State))

	// Export light
	ExportPointer("renderLight", unsafe.Pointer(&engine.Light))

	// Export camera
	ExportPointer("renderCamera_matrix", unsafe.Pointer(&engine.Camera.Matrix.Raw))
	ExportPointer("renderUICamera_matrix", unsafe.Pointer(&engine.UICamera.Matrix.Raw))

	// Export dynamic mesh layer
	ExportPointer("renderDynamicMeshLayer_vertex", unsafe.Pointer(&engine.Main.VertexList))
	ExportPointer("renderDynamicMeshLayer_uv", unsafe.Pointer(&engine.Main.UvList))
	ExportPointer("renderDynamicMeshLayer_normal", unsafe.Pointer(&engine.Main.NormalList))
	ExportPointer("renderDynamicMeshLayer_position", unsafe.Pointer(&engine.Main.PositionList))
	ExportPointer("renderDynamicMeshLayer_rotation", unsafe.Pointer(&engine.Main.RotationList))
	ExportPointer("renderDynamicMeshLayer_scale", unsafe.Pointer(&engine.Main.ScaleList))
	ExportPointer("renderDynamicMeshLayer_color", unsafe.Pointer(&engine.Main.ColorList))
	ExportPointer("renderDynamicMeshLayer_index", unsafe.Pointer(&engine.Main.IndexList))

	// Export static mesh layer
	ExportPointer("renderStaticMeshLayer_vertex", unsafe.Pointer(&engine.StaticMesh.VertexList))
	ExportPointer("renderStaticMeshLayer_uv", unsafe.Pointer(&engine.StaticMesh.UvList))
	ExportPointer("renderStaticMeshLayer_normal", unsafe.Pointer(&engine.StaticMesh.NormalList))
	ExportPointer("renderStaticMeshLayer_color", unsafe.Pointer(&engine.StaticMesh.ColorList))
	ExportPointer("renderStaticMeshLayer_index", unsafe.Pointer(&engine.StaticMesh.IndexList))

	// Export line layer
	ExportPointer("renderLineLayer_vertex", unsafe.Pointer(&engine.Line.VertexList))
	ExportPointer("renderLineLayer_color", unsafe.Pointer(&engine.Line.ColorList))

	// Export point layer
	ExportPointer("renderPointLayer_vertex", unsafe.Pointer(&engine.Point.VertexList))
	ExportPointer("renderPointLayer_color", unsafe.Pointer(&engine.Point.ColorList))

	// Export ui layer
	ExportPointer("renderUILayer_vertex", unsafe.Pointer(&engine.UI.VertexList))
	ExportPointer("renderUILayer_uv", unsafe.Pointer(&engine.UI.UvList))
	ExportPointer("renderUILayer_position", unsafe.Pointer(&engine.UI.PositionList))
	ExportPointer("renderUILayer_rotation", unsafe.Pointer(&engine.UI.RotationList))
	ExportPointer("renderUILayer_scale", unsafe.Pointer(&engine.UI.ScaleList))
	ExportPointer("renderUILayer_color", unsafe.Pointer(&engine.UI.ColorList))
	ExportPointer("renderUILayer_index", unsafe.Pointer(&engine.UI.IndexList))
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

var fileMap = map[string][]byte{}

func LoadFile(path string) ([]byte, error) {
	// Load from cache
	_, ok := fileMap[path]
	if ok {
		return fileMap[path], nil
	}

	size := uint32(0)
	js.Global().Get("window").Get("go").Get("fs").Call("openFile", path, uintptr(unsafe.Pointer(&size)))

	// Wait until it's ready
	for {
		if size > 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}

	// Allocate size
	fileMap[path] = make([]byte, size)
	arr, _ := fileMap[path]
	arr2 := (*reflect.SliceHeader)(unsafe.Pointer(&arr))

	// Read to
	js.Global().Get("window").Get("go").Get("fs").Call("readFile", path, arr2.Data)

	// Done
	return fileMap[path], nil
}
