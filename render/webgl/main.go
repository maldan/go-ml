package ml_webgl

import (
	_ "embed"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"
)

//go:embed js/engine.js
var engineJs string

type Mesh struct {
	Vertices []ml_geom.Vector3[float32]
	Indices  []uint16
	UV       []ml_geom.Vector2[float32]
	// Matrix   ml_geom.Matrix4x4[float32]

	Position ml_geom.Vector3[float32]
	Rotation ml_geom.Vector3[float32]
}

func exportSlice[T any](name string, list []T) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		header := (*reflect.SliceHeader)(unsafe.Pointer(&list))
		return header.Data
	}))
}

func meshToVertex(meshList []Mesh, vertexList []float32, project ml_geom.Matrix4x4[float32]) {
	id := 0
	for i := 0; i < len(meshList); i++ {
		for j := 0; j < len(meshList[i].Vertices); j++ {
			newVertices := meshList[i].Vertices[j].Clone()

			finalMx := ml_geom.Matrix4x4[float32]{}
			finalMx.Identity()

			mx := ml_geom.Matrix4x4[float32]{}
			mx.Identity()
			mx.RotateX(meshList[i].Rotation.X)
			mx.RotateY(meshList[i].Rotation.Y)
			mx.RotateZ(meshList[i].Rotation.Z)
			mx.Translate(meshList[i].Position)

			finalMx.Multiply(project)
			finalMx.Multiply(mx)
			newVertices.TransformMatrix4x4(finalMx)

			vertexList[id] = newVertices.X
			vertexList[id+1] = newVertices.Y
			vertexList[id+2] = newVertices.Z
			id += 3
		}
	}
}

func Start() {
	vertexList := make([]float32, 128)
	indexList := make([]int16, 128)
	exportSlice("getVertexListPointer", vertexList)
	exportSlice("getIndexListPointer", indexList)

	meshList := make([]Mesh, 0)
	meshList = append(meshList, Mesh{
		Vertices: []ml_geom.Vector3[float32]{
			{0.5, 0.5, 0},
			{-0.5, 0.5, 0},
			{0.5, -0.5, 0},
			{-0.5, -0.5, 0},
		},
	})
	meshList = append(meshList, Mesh{
		Vertices: []ml_geom.Vector3[float32]{
			{0.5, 0.5, 0},
			{-0.5, 0.5, 0},
			{0.5, -0.5, 0},
			{-0.5, -0.5, 0},
		},
	})

	// fmt.Printf("%v\n", engineJs)
	js.Global().Set("engineJs", engineJs)

	projectionMatrix := ml_geom.Matrix4x4[float32]{}
	projectionMatrix.Identity()
	projectionMatrix.Perspective((45*3.14)/180, 1, 0.1, 100.0)
	projectionMatrix.Translate(ml_geom.Vector3[float32]{0, 0, -1.5})
	projectionMatrix.Scale(ml_geom.Vector3[float32]{0.5, 0.5, 0.5})

	xxx := float32(0)

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
		//meshList[0].Position.X = float32(math.Sin(float64(xxx)))
		//meshList[0].Position.Y = float32(math.Sin(float64(xxx)))
		// meshList[0].Position.Z = float32(math.Sin(float64(xxx)))
		meshList[0].Rotation.Z = xxx
		xxx += 0.01

		meshToVertex(meshList, vertexList, projectionMatrix)
		/*for i := 0; i < len(vertexList); i++ {
			vertexList[i] = ml_number.RandFloat32(-1, 1)
		}*/

		time.Sleep(time.Millisecond * 16)
	}
}
