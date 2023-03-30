package ml_webgl

import (
	_ "embed"
	"fmt"
	ml_render "github.com/maldan/go-ml/render"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"os"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"
)

func Init() {
	js.Global().Set("__webglEngineState", js.FuncOf(func(this js.Value, args []js.Value) any {
		vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&state.VertexList))
		indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&state.IndexList))

		return map[string]any{
			"vertexArrayLength":  state.VertexAmount,
			"vertexArrayPointer": vertexHeader.Data,

			"indexArrayLength":  state.IndexAmount,
			"indexArrayPointer": indexHeader.Data,
		}
	}))
	js.Global().Set("__webglEngineRender", js.FuncOf(func(this js.Value, args []js.Value) any {
		Render()
		return nil
	}))

	state.VertexList = make([]float32, 65536*3)
	state.MeshList = make([]*ml_render.Mesh, 0, 32)
	state.IndexList = make([]uint16, 65536)
}

type RenderEngine struct {
	MeshList     []*ml_render.Mesh
	VertexList   []float32
	IndexList    []uint16
	VertexAmount int
	IndexAmount  int
}

var state RenderEngine = RenderEngine{}

func AddMesh(mesh *ml_render.Mesh) {
	state.MeshList = append(state.MeshList, mesh)
}

func Render() {
	projectionMatrix := ml_geom.Matrix4x4[float32]{}
	projectionMatrix.Identity()
	projectionMatrix.Perspective((45*3.14)/180, 1, 0.1, 100.0)
	projectionMatrix.Translate(ml_geom.Vector3[float32]{0, 0, -1.5})
	projectionMatrix.Scale(ml_geom.Vector3[float32]{0.1, 0.1, 0.1})

	state.VertexAmount = 0
	state.IndexAmount = 0

	vertexId := 0
	indexId := 0
	lastMaxIndex := uint16(0)
	for i := 0; i < len(state.MeshList); i++ {
		finalMx := ml_geom.Matrix4x4[float32]{}
		finalMx.Identity()

		state.MeshList[i].ApplyMatrix()

		finalMx.Multiply(projectionMatrix)
		finalMx.Multiply(state.MeshList[i].Matrix)

		// Copy vertex
		for j := 0; j < len(state.MeshList[i].Vertices); j++ {
			newVertices := state.MeshList[i].Vertices[j].Clone()
			newVertices.TransformMatrix4x4(finalMx)

			state.VertexList[vertexId] = newVertices.X
			state.VertexList[vertexId+1] = newVertices.Y
			state.VertexList[vertexId+2] = newVertices.Z
			/*newVertices := state.MeshList[i].Vertices[j].Clone()

			finalMx := ml_geom.Matrix4x4[float32]{}
			finalMx.Identity()

			mx := ml_geom.Matrix4x4[float32]{}
			mx.Identity()
			mx.Translate(state.MeshList[i].Position)
			mx.RotateX(state.MeshList[i].Rotation.X)
			mx.RotateY(state.MeshList[i].Rotation.Y)
			mx.RotateZ(state.MeshList[i].Rotation.Z)

			finalMx.Multiply(projectionMatrix)
			finalMx.Multiply(mx)
			newVertices.TransformMatrix4x4(finalMx)

			state.VertexList[vertexId] = newVertices.X
			state.VertexList[vertexId+1] = newVertices.Y
			state.VertexList[vertexId+2] = newVertices.Z*/
			vertexId += 3
		}
		state.VertexAmount += len(state.MeshList[i].Vertices) * 3

		// Copy index
		maxIndex := lastMaxIndex
		for j := 0; j < len(state.MeshList[i].Indices); j++ {
			state.IndexList[indexId] = state.MeshList[i].Indices[j] + maxIndex
			if state.IndexList[indexId] > lastMaxIndex {
				lastMaxIndex = state.IndexList[indexId]
			}
			indexId += 1
		}
		lastMaxIndex += 1
		state.IndexAmount += len(state.MeshList[i].Indices)
	}
}

func Start() {
	//vertexList := make([]float32, 128)
	//indexList := make([]int16, 128)
	//exportSlice("getVertexListPointer", vertexList)
	//exportSlice("getIndexListPointer", indexList)

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

	// fmt.Printf("%v\n", engineJs)
	// js.Global().Set("engineJs", engineJs)

	projectionMatrix := ml_geom.Matrix4x4[float32]{}
	projectionMatrix.Identity()
	projectionMatrix.Perspective((45*3.14)/180, 1, 0.1, 100.0)
	projectionMatrix.Translate(ml_geom.Vector3[float32]{0, 0, -1.5})
	projectionMatrix.Scale(ml_geom.Vector3[float32]{0.5, 0.5, 0.5})

	xxx := float32(0)

	/*r := ml_net.Get("./index.html", nil)
	b, _ := r.Bytes()
	fmt.Printf("%v\n", string(b))*/

	file, _ := os.OpenFile("./index.html", os.O_RDWR, 0777)
	info, _ := file.Stat()
	fmt.Printf("FileSize: %v\n", info.Size())
	sas := make([]byte, info.Size())
	file.Read(sas)
	fmt.Printf("%v\n", string(sas))

	/*sas := make([]byte, 12)
	ll, err2 := file.Read(sas)
	fmt.Printf("Read length: %v\n", ll)
	fmt.Printf("Read err: %v\n", err2)
	fmt.Printf("%v\n", file)
	fmt.Printf("%v\n", err)
	fmt.Printf("%v\n", sas)*/

	/*f := ml_file.New("./index.html")
	d, err := f.ReadAll()
	fmt.Printf("DATA: %v\n", d)
	fmt.Printf("ER:: %v\n", err)*/

	for {
		// mx := ml_geom.Matrix4x4[float32]{}
		// mx.Identity()
		// mx.Scale(ml_geom.Vector3[float32]{xxx, xxx, xxx})
		// mx.RotateZ(xxx)
		/*for i := 0; i < len(meshList[0].Vertices); i++ {
			meshList[0].Vertices[i].TransformMatrix4x4(mx)
		}*/
		// fmt.Printf("%v\n", meshList[0].Vertices)
		// meshList[0].Matrix.Identity()
		// fmt.Printf("A %v\n", meshList[0].Matrix)
		// meshList[0].Matrix.Rotate(xxx, ml_geom.Vector3[float32]{0, 0, 1})
		// fmt.Printf("B %v\n", meshList[0].Matrix)
		// meshList[0].Matrix.RotateZ(xxx)
		// meshList[0].Matrix.Multiply(projectionMatrix)
		// meshList[0].RotationX = xxx
		// meshList[0].Position.X = float32(math.Sin(float64(xxx)))
		// meshList[0].Position.Y = float32(math.Sin(float64(xxx)))
		// meshList[0].Position.Z = float32(math.Sin(float64(xxx)))
		meshList[0].Rotation.Z = xxx
		xxx += 0.01

		// meshToVertex(meshList, vertexList, projectionMatrix)
		/*for i := 0; i < len(vertexList); i++ {
			vertexList[i] = ml_number.RandFloat32(-1, 1)
		}*/

		time.Sleep(time.Millisecond * 16)
	}
}
