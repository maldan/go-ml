package mrender_layer

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"reflect"
	"unsafe"
)

type Point struct {
	Position mmath_la.Vector3[float32]
	Size     float32
	Color    ml_color.ColorRGBA[float32]
}

type PointLayer struct {
	PointList    []Point
	VertexList   []float32
	ColorList    []float32
	VertexAmount int
	// Camera       mr_camera.PerspectiveCamera

	state map[string]any
}

func (l *PointLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.ColorList = make([]float32, 65536*3)
	l.PointList = make([]Point, 0, 1024)
}

func (l *PointLayer) Render() {
	//l.Camera.ApplyMatrix()

	// Fill points
	vertexId := 0
	for i := 0; i < len(l.PointList); i++ {
		point := l.PointList[i]

		l.VertexList[vertexId] = point.Position.X
		l.VertexList[vertexId+1] = point.Position.Y
		l.VertexList[vertexId+2] = point.Position.Z
		l.VertexList[vertexId+3] = point.Size

		l.ColorList[vertexId] = point.Color.R
		l.ColorList[vertexId+1] = point.Color.G
		l.ColorList[vertexId+2] = point.Color.B
		l.ColorList[vertexId+3] = point.Color.A

		vertexId += 4
	}
	l.VertexAmount = vertexId

	if len(l.PointList) > 0 {
		l.PointList = l.PointList[:0]
	}
}

func (l *PointLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer": vertexHeader.Data,
			"colorPointer":  colorHeader.Data,

			"vertexAmount": l.VertexAmount,
			"colorAmount":  l.VertexAmount,

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		l.state["vertexAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.VertexAmount
	}

	return l.state
}
