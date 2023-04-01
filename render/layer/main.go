package mr_layer

import (
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"reflect"
	"unsafe"
)

type Layer interface {
}

type MainLayer struct {
	MeshList     []*mr_mesh.Mesh
	VertexList   []float32
	PositionList []float32
	RotationList []float32
	ScaleList    []float32
	IndexList    []uint16
	VertexAmount int
	IndexAmount  int
	Camera       mr_camera.PerspectiveCamera
}

func (l *MainLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.PositionList = make([]float32, 65536*3)
	l.RotationList = make([]float32, 65536*3)
	l.ScaleList = make([]float32, 65536*3)

	l.MeshList = make([]*mr_mesh.Mesh, 0, 1024)
	l.IndexList = make([]uint16, 65536)
}

func (l *MainLayer) Render() {
	l.Camera.Fov = 45
	l.Camera.AspectRatio = 1
	l.Camera.Position = ml_geom.Vector3[float32]{0, 0, -1.5}
	l.Camera.Scale = ml_geom.Vector3[float32]{0.1, 0.1, 0.1}
	l.Camera.ApplyMatrix()

	l.VertexAmount = 0
	l.IndexAmount = 0

	vertexId := 0
	indexId := 0
	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.MeshList); i++ {
		mesh := l.MeshList[i]

		// Copy vertex
		for j := 0; j < len(mesh.Vertices); j++ {
			v := mesh.Vertices[j]

			l.VertexList[vertexId] = v.X
			l.VertexList[vertexId+1] = v.Y
			l.VertexList[vertexId+2] = v.Z

			p := mesh.Position
			l.PositionList[vertexId] = p.X
			l.PositionList[vertexId+1] = p.Y
			l.PositionList[vertexId+2] = p.Z

			r := mesh.Rotation
			l.RotationList[vertexId] = r.X
			l.RotationList[vertexId+1] = r.Y
			l.RotationList[vertexId+2] = r.Z

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
	}
}

func (l *MainLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	positionHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.PositionList))
	rotationHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.RotationList))
	indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.IndexList))

	return map[string]any{
		"vertexPointer":   vertexHeader.Data,
		"indexPointer":    indexHeader.Data,
		"positionPointer": positionHeader.Data,
		"rotationPointer": rotationHeader.Data,
		"scalePointer":    0,

		"vertexAmount":   l.VertexAmount,
		"indexAmount":    l.IndexAmount,
		"positionAmount": l.VertexAmount,
		"rotationAmount": l.VertexAmount,
		"scaleAmount":    l.VertexAmount,

		"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
	}
}
