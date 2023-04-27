package mrender_layer

import (
	"fmt"
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	"reflect"
	"unsafe"
)

type Layer interface {
}

type MainLayer struct {
	VertexList   []float32
	UvList       []float32
	NormalList   []float32
	PositionList []float32
	RotationList []float32
	ScaleList    []float32
	ColorList    []float32
	IndexList    []uint16

	VertexAmount int
	UvAmount     int
	ColorAmount  int
	IndexAmount  int

	AllocatedMesh    []mr_mesh.Mesh
	MeshInstanceList []mr_mesh.MeshInstance

	Camera mr_camera.PerspectiveCamera

	InstanceId int

	state map[string]any
}

func (l *MainLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.NormalList = make([]float32, 65536*3)
	l.UvList = make([]float32, 65536*2)
	l.PositionList = make([]float32, 65536*3)
	l.RotationList = make([]float32, 65536*3)
	l.ScaleList = make([]float32, 65536*3)
	l.ColorList = make([]float32, 65536*4)

	l.AllocatedMesh = make([]mr_mesh.Mesh, 0, 8192)
	l.MeshInstanceList = make([]mr_mesh.MeshInstance, 1024)
	l.IndexList = make([]uint16, 65536)

	fmt.Printf("Render allocated %v\n", cap(l.VertexList)*4*6)
}

func (l *MainLayer) Render() {
	l.Camera.ApplyMatrix()

	l.VertexAmount = 0
	l.IndexAmount = 0
	l.UvAmount = 0
	l.ColorAmount = 0

	vertexId := 0
	vertexId2 := 0
	indexId := 0
	uvIndex := 0
	colorId := 0
	lastMaxIndex := uint16(0)

	for i := 0; i < len(l.MeshInstanceList); i++ {
		instance := l.MeshInstanceList[i]

		if !instance.IsVisible {
			continue
		}
		if !instance.IsActive {
			continue
		}
		if instance.Id < 0 {
			continue
		}

		mesh := l.AllocatedMesh[instance.Id]

		// Copy vertex
		for j := 0; j < len(mesh.Vertices); j++ {
			v := mesh.Vertices[j]
			l.VertexList[vertexId] = v.X
			l.VertexList[vertexId+1] = v.Y
			l.VertexList[vertexId+2] = v.Z

			// Copy normal
			n := mesh.Normal[j]
			l.NormalList[vertexId] = n.X
			l.NormalList[vertexId+1] = n.Y
			l.NormalList[vertexId+2] = n.Z

			vertexId += 3
		}

		for j := 0; j < len(mesh.Vertices); j++ {
			p := instance.Position
			l.PositionList[vertexId2] = p.X
			l.PositionList[vertexId2+1] = p.Y
			l.PositionList[vertexId2+2] = p.Z

			p = instance.Rotation
			l.RotationList[vertexId2] = p.X
			l.RotationList[vertexId2+1] = p.Y
			l.RotationList[vertexId2+2] = p.Z

			p = instance.Scale
			l.ScaleList[vertexId2] = p.X
			l.ScaleList[vertexId2+1] = p.Y
			l.ScaleList[vertexId2+2] = p.Z

			vertexId2 += 3
		}
		l.VertexAmount += len(mesh.Vertices) * 3

		// Copy index
		maxIndex := lastMaxIndex
		for j := 0; j < len(mesh.Indices); j++ {
			l.IndexList[indexId] = mesh.Indices[j] + maxIndex
			if l.IndexList[indexId] > lastMaxIndex {
				lastMaxIndex = l.IndexList[indexId]
			}
			indexId += 1
		}
		lastMaxIndex += 1
		l.IndexAmount += len(mesh.Indices)

		// Copy uv
		for j := 0; j < len(mesh.UV); j++ {
			v := mesh.UV[j]

			l.UvList[uvIndex] = v.X
			l.UvList[uvIndex+1] = v.Y

			if instance.UvOffset.X != 0 || instance.UvOffset.Y != 0 {
				l.UvList[uvIndex] += instance.UvOffset.X
				l.UvList[uvIndex+1] -= instance.UvOffset.Y
			}

			uvIndex += 2
		}
		l.UvAmount += len(mesh.UV) * 2

		// Copy color
		for j := 0; j < len(mesh.Vertices); j++ {
			c := instance.Color
			l.ColorList[colorId] = c.R
			l.ColorList[colorId+1] = c.G
			l.ColorList[colorId+2] = c.B
			l.ColorList[colorId+3] = c.A
			colorId += 4
		}
		l.ColorAmount += len(mesh.Vertices) * 4
	}
}

func (l *MainLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	normalHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.NormalList))
	uvHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.UvList))
	positionHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.PositionList))
	rotationHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.RotationList))
	scaleHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ScaleList))
	indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.IndexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer":   vertexHeader.Data,
			"normalPointer":   normalHeader.Data,
			"uvPointer":       uvHeader.Data,
			"indexPointer":    indexHeader.Data,
			"positionPointer": positionHeader.Data,
			"rotationPointer": rotationHeader.Data,
			"scalePointer":    scaleHeader.Data,
			"colorPointer":    colorHeader.Data,

			"vertexAmount":   l.VertexAmount,
			"normalAmount":   l.VertexAmount,
			"uvAmount":       l.UvAmount,
			"indexAmount":    l.IndexAmount,
			"positionAmount": l.VertexAmount,
			"rotationAmount": l.VertexAmount,
			"scaleAmount":    l.VertexAmount,
			"colorAmount":    l.ColorAmount,

			"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		l.state["vertexAmount"] = l.VertexAmount
		l.state["normalAmount"] = l.VertexAmount
		l.state["uvAmount"] = l.UvAmount
		l.state["indexAmount"] = l.IndexAmount
		l.state["positionAmount"] = l.VertexAmount
		l.state["rotationAmount"] = l.VertexAmount
		l.state["scaleAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.ColorAmount
	}

	return l.state
}
