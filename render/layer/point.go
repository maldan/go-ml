package mr_layer

import (
	mr_camera "github.com/maldan/go-ml/render/camera"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"reflect"
	"unsafe"
)

type PointLayer struct {
	PointList    []ml_geom.Vector3[float32]
	VertexList   []float32
	VertexAmount int
	Camera       mr_camera.PerspectiveCamera
}

func (l *PointLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.PointList = make([]ml_geom.Vector3[float32], 0, 1024)
}

func (l *PointLayer) Render() {
	l.Camera.ApplyMatrix()

	// Fill points
	vertexId := 0
	for i := 0; i < len(l.PointList); i++ {
		point := l.PointList[i]

		l.VertexList[vertexId] = point.X
		l.VertexList[vertexId+1] = point.Y
		l.VertexList[vertexId+2] = point.Z
		vertexId += 3
	}
	l.VertexAmount = vertexId

	if len(l.PointList) > 0 {
		l.PointList = l.PointList[:0]
	}
}

func (l *PointLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))

	return map[string]any{
		"vertexPointer": vertexHeader.Data,
		"vertexAmount":  l.VertexAmount,

		"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
	}
}
