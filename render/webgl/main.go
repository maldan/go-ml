package ml_webgl

import (
	_ "embed"
	ml_render "github.com/maldan/go-ml/render"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"reflect"
	"syscall/js"
	"unsafe"
)

type RenderEngine struct {
	MeshList     []*ml_render.Mesh
	VertexList   []float32
	PositionList []float32
	RotationList []float32
	ScaleList    []float32
	IndexList    []uint16

	PointList       []ml_geom.Vector3[float32]
	PointVertexList []float32

	PointAmount  int
	VertexAmount int
	IndexAmount  int

	ProjectionMatrix ml_geom.Matrix4x4[float32]
}

var State RenderEngine = RenderEngine{}

func Init() {
	js.Global().Set("goWasmRenderState", js.FuncOf(func(this js.Value, args []js.Value) any {
		vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.VertexList))
		positionHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.PositionList))
		rotationHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.RotationList))
		indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.IndexList))

		pointHeader := (*reflect.SliceHeader)(unsafe.Pointer(&State.PointVertexList))

		return map[string]any{
			/*"vertexArrayLength":    State.VertexAmount,
			"vertexArrayPointer":   vertexHeader.Data,
			"positionArrayPointer": positionHeader.Data,
			"rotationArrayPointer": rotationHeader.Data,

			"indexArrayLength":  State.IndexAmount,
			"indexArrayPointer": indexHeader.Data,

			"projectionMatrixPointer": uintptr(unsafe.Pointer(&State.ProjectionMatrix.Raw)),

			"pointArrayLength":  State.PointAmount * 3,
			"pointArrayPointer": pointHeader.Data,*/

			"mainLayer": map[string]any{
				"vertexPointer": vertexHeader.Data,
				"vertexAmount":  State.VertexAmount,
				"indexPointer":  indexHeader.Data,
				"indexAmount":   State.IndexAmount,

				"positionPointer": positionHeader.Data,
				"rotationPointer": rotationHeader.Data,
				"scalePointer":    0,

				"projectionMatrixPointer": uintptr(unsafe.Pointer(&State.ProjectionMatrix.Raw)),
			},
			"pointLayer": map[string]any{
				"vertexPointer": pointHeader.Data,
				"vertexAmount":  State.PointAmount * 3,
			},
		}
	}))

	js.Global().Set("goWasmRenderFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		Render()
		return nil
	}))

	State.VertexList = make([]float32, 65536*3)
	State.PositionList = make([]float32, 65536*3)
	State.RotationList = make([]float32, 65536*3)
	State.ScaleList = make([]float32, 65536*3)

	State.MeshList = make([]*ml_render.Mesh, 0, 1024)
	State.IndexList = make([]uint16, 65536)

	State.PointList = make([]ml_geom.Vector3[float32], 0, 1024)
	State.PointVertexList = make([]float32, 1024*4)
}

func AddMesh(mesh *ml_render.Mesh) {
	State.MeshList = append(State.MeshList, mesh)
}

type Color32 struct {
	R float32
	G float32
	B float32
	A float32
}

func DebugLine(from ml_geom.Vector3[float32], to ml_geom.Vector3[float32], color Color32) {
	// State.PointList = append(State.PointList, )
}

func DebugPoint(to ml_geom.Vector3[float32]) {
	State.PointList = append(State.PointList, to)
}

func Render() {
	State.ProjectionMatrix.Identity()
	State.ProjectionMatrix.Perspective((45*3.14)/180, 1, 0.1, 100.0)
	State.ProjectionMatrix.Translate(ml_geom.Vector3[float32]{0, 0, -1.5})
	State.ProjectionMatrix.Scale(ml_geom.Vector3[float32]{0.1, 0.1, 0.1})

	State.VertexAmount = 0
	State.IndexAmount = 0

	vertexId := 0
	indexId := 0
	lastMaxIndex := uint16(0)
	for i := 0; i < len(State.MeshList); i++ {
		mesh := State.MeshList[i]

		// Copy vertex
		for j := 0; j < len(mesh.Vertices); j++ {
			v := mesh.Vertices[j]
			State.VertexList[vertexId] = v.X
			State.VertexList[vertexId+1] = v.Y
			State.VertexList[vertexId+2] = v.Z

			p := mesh.Position
			State.PositionList[vertexId] = p.X
			State.PositionList[vertexId+1] = p.Y
			State.PositionList[vertexId+2] = p.Z

			p = mesh.Rotation
			State.RotationList[vertexId] = p.X
			State.RotationList[vertexId+1] = p.Y
			State.RotationList[vertexId+2] = p.Z

			p = mesh.Scale
			State.ScaleList[vertexId] = p.X
			State.ScaleList[vertexId+1] = p.Y
			State.ScaleList[vertexId+2] = p.Z

			vertexId += 3
		}
		State.VertexAmount += len(mesh.Vertices) * 3

		// Copy index
		maxIndex := lastMaxIndex
		for j := 0; j < len(mesh.Indices); j++ {
			State.IndexList[indexId] = mesh.Indices[j] + maxIndex
			if State.IndexList[indexId] > lastMaxIndex {
				lastMaxIndex = State.IndexList[indexId]
			}
			indexId += 1
		}
		lastMaxIndex += 1
		State.IndexAmount += len(mesh.Indices)
	}

	// Fill points
	vertexId = 0
	State.PointAmount = len(State.PointList)
	for i := 0; i < len(State.PointList); i++ {
		point := State.PointList[i]

		State.PointVertexList[vertexId] = point.X
		State.PointVertexList[vertexId+1] = point.Y
		State.PointVertexList[vertexId+2] = point.Z
		vertexId += 3
	}

	if len(State.PointList) > 0 {
		State.PointList = State.PointList[:0]
	}
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
