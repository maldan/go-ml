package mr_layer

import (
	mr_camera "github.com/maldan/go-ml/render/camera"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"reflect"
	"unsafe"
)

type LineLayer struct {
	LineList     []ml_geom.Line[float32, ml_geom.Vector3[float32]]
	VertexList   []float32
	VertexAmount int
	Camera       mr_camera.PerspectiveCamera
}

func (l *LineLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.LineList = make([]ml_geom.Line[float32, ml_geom.Vector3[float32]], 0, 1024)
}

func (l *LineLayer) Render() {
	l.Camera.Fov = 45
	l.Camera.AspectRatio = 1
	l.Camera.Position = ml_geom.Vector3[float32]{0, 0, -1.5}
	l.Camera.Scale = ml_geom.Vector3[float32]{0.1, 0.1, 0.1}
	l.Camera.ApplyMatrix()

	// Fill points
	vertexId := 0
	for i := 0; i < len(l.LineList); i++ {
		line := l.LineList[i]

		l.VertexList[vertexId] = line.From.X
		l.VertexList[vertexId+1] = line.From.Y
		l.VertexList[vertexId+2] = line.From.Z
		l.VertexList[vertexId+3] = line.To.X
		l.VertexList[vertexId+4] = line.To.Y
		l.VertexList[vertexId+5] = line.To.Z
		vertexId += 6
	}
	l.VertexAmount = vertexId

	if len(l.LineList) > 0 {
		l.LineList = l.LineList[:0]
	}
}

func (l *LineLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))

	return map[string]any{
		"vertexPointer": vertexHeader.Data,
		"vertexAmount":  l.VertexAmount,

		"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
	}
}
