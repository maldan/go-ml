package ml_webgl

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

	/*State.VertexList = make([]float32, 65536*3)
	State.PositionList = make([]float32, 65536*3)
	State.RotationList = make([]float32, 65536*3)
	State.ScaleList = make([]float32, 65536*3)

	State.MeshList = make([]*ml_render.Mesh, 0, 1024)
	State.IndexList = make([]uint16, 65536)

	State.PointList = make([]ml_geom.Vector3[float32], 0, 1024)
	State.PointVertexList = make([]float32, 1024*4)*/
}

/*func Start() {
	meshList := make([]ml_render.Mesh, 0)
	meshList = append(meshList, ml_render.Mesh{
		Vertices: []ml_geom.Vector3[float32]{
			{0.5, 0.5, 0},
			{-0.5, 0.5, 0},
			{0.5, -0.5, 0},
			{-0.5, -0.5, 0},
		},
	})
	meshList = append(meshList, ml_render.Mesh{
		Vertices: []ml_geom.Vector3[float32]{
			{0.5, 0.5, 0},
			{-0.5, 0.5, 0},
			{0.5, -0.5, 0},
			{-0.5, -0.5, 0},
		},
	})

	projectionMatrix := ml_geom.Matrix4x4[float32]{}
	projectionMatrix.Identity()
	projectionMatrix.Perspective((45*3.14)/180, 1, 0.1, 100.0)
	projectionMatrix.Translate(ml_geom.Vector3[float32]{0, 0, -1.5})
	projectionMatrix.Scale(ml_geom.Vector3[float32]{0.5, 0.5, 0.5})

	file, _ := os.OpenFile("./index.html", os.O_RDWR, 0777)
	info, _ := file.Stat()
	fmt.Printf("FileSize: %v\n", info.Size())
	sas := make([]byte, info.Size())
	file.Read(sas)
	fmt.Printf("%v\n", string(sas))

	for {
		time.Sleep(time.Millisecond * 16)
	}
}
*/
