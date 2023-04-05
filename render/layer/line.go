package mrender_layer

import (
	mr_camera "github.com/maldan/go-ml/render/camera"
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	"reflect"
	"unsafe"
)

type LineLayer struct {
	LineList     []mr_mesh.Line
	VertexList   []float32
	VertexAmount int
	ColorList    []float32
	Camera       mr_camera.PerspectiveCamera

	state map[string]any
}

func (l *LineLayer) Init() {
	l.ColorList = make([]float32, 65536*3)
	l.VertexList = make([]float32, 65536*3)
	l.LineList = make([]mr_mesh.Line, 0, 1024)
}

func (l *LineLayer) Render() {
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

		l.ColorList[vertexId] = line.Color.R
		l.ColorList[vertexId+1] = line.Color.G
		l.ColorList[vertexId+2] = line.Color.B
		l.ColorList[vertexId+3] = line.Color.R
		l.ColorList[vertexId+4] = line.Color.G
		l.ColorList[vertexId+5] = line.Color.B

		vertexId += 6
	}
	l.VertexAmount = vertexId

	if len(l.LineList) > 0 {
		l.LineList = l.LineList[:0]
	}
}

func (l *LineLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer":           vertexHeader.Data,
			"vertexAmount":            l.VertexAmount,
			"colorPointer":            colorHeader.Data,
			"colorAmount":             l.VertexAmount,
			"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		l.state["vertexAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.VertexAmount
	}

	return l.state
}
