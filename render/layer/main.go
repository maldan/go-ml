package mr_layer

import (
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	"reflect"
	"unsafe"
)

type Layer interface {
}

type MainLayer struct {
	MeshList   []*mr_mesh.Mesh
	VertexList []float32
	UvList     []float32
	NormalList []float32

	PositionList []float32
	RotationList []float32
	ScaleList    []float32

	IndexList []uint16

	VertexAmount int
	IndexAmount  int
	UvAmount     int

	Camera mr_camera.PerspectiveCamera
}

func (l *MainLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.NormalList = make([]float32, 65536*3)
	l.UvList = make([]float32, 65536*2)
	l.PositionList = make([]float32, 65536*3)
	l.RotationList = make([]float32, 65536*3)
	l.ScaleList = make([]float32, 65536*3)

	l.MeshList = make([]*mr_mesh.Mesh, 0, 1024)
	l.IndexList = make([]uint16, 65536)
}

func (l *MainLayer) Render() {
	l.Camera.ApplyMatrix()

	l.VertexAmount = 0
	l.IndexAmount = 0
	l.UvAmount = 0

	vertexId := 0
	indexId := 0
	uvIndex := 0
	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.MeshList); i++ {
		mesh := l.MeshList[i]

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

			p := mesh.Position
			l.PositionList[vertexId] = p.X
			l.PositionList[vertexId+1] = p.Y
			l.PositionList[vertexId+2] = p.Z

			p = mesh.Rotation
			l.RotationList[vertexId] = p.X
			l.RotationList[vertexId+1] = p.Y
			l.RotationList[vertexId+2] = p.Z

			p = mesh.Scale
			l.ScaleList[vertexId] = p.X
			l.ScaleList[vertexId+1] = p.Y
			l.ScaleList[vertexId+2] = p.Z

			vertexId += 3
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

			uvIndex += 2
		}
		l.UvAmount += len(mesh.UV) * 2
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

	return map[string]any{
		"vertexPointer":   vertexHeader.Data,
		"normalPointer":   normalHeader.Data,
		"uvPointer":       uvHeader.Data,
		"indexPointer":    indexHeader.Data,
		"positionPointer": positionHeader.Data,
		"rotationPointer": rotationHeader.Data,
		"scalePointer":    scaleHeader.Data,

		"vertexAmount":   l.VertexAmount,
		"normalAmount":   l.VertexAmount,
		"uvAmount":       l.UvAmount,
		"indexAmount":    l.IndexAmount,
		"positionAmount": l.VertexAmount,
		"rotationAmount": l.VertexAmount,
		"scaleAmount":    l.VertexAmount,

		"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
	}
}
