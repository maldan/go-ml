package mrender_webgl

import (
	_ "embed"
	mr "github.com/maldan/go-ml/render"
	"syscall/js"
)

func Init(engine *mr.RenderEngine) {
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
