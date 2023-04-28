package mrender_layer

import (
	mr_mesh "github.com/maldan/go-ml/render/mesh"
	"reflect"
	"unsafe"
)

type LineLayer struct {
	VertexList []float32
	ColorList  []float32

	LineList []mr_mesh.Line

	state map[string]any
}

func (l *LineLayer) Init() {
	l.VertexList = make([]float32, 0, 1024)
	l.ColorList = make([]float32, 0, 1024)
	l.LineList = make([]mr_mesh.Line, 0, 1024)
}

func (l *LineLayer) Render() {
	// Clear before start
	l.VertexList = l.VertexList[:0]
	l.ColorList = l.ColorList[:0]

	// Fill points
	for i := 0; i < len(l.LineList); i++ {
		line := l.LineList[i]

		l.VertexList = append(l.VertexList,
			line.From.X, line.From.Y, line.From.Z,
			line.To.X, line.To.Y, line.To.Z,
		)

		l.ColorList = append(l.ColorList,
			line.Color.R, line.Color.G, line.Color.B, line.Color.A,
			line.Color.R, line.Color.G, line.Color.B, line.Color.A,
		)
	}

	// Clear lines
	if len(l.LineList) > 0 {
		l.LineList = l.LineList[:0]
	}
}

func (l *LineLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer": vertexHeader.Data,
			//"vertexAmount":  l.VertexAmount,
			"colorPointer": colorHeader.Data,
			//"colorAmount":   l.VertexAmount,
			// "projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		//l.state["vertexAmount"] = l.VertexAmount
		//l.state["colorAmount"] = l.VertexAmount
	}

	return l.state
}
