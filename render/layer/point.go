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
	VertexList []float32
	ColorList  []float32

	PointList []Point
	state     map[string]any
}

func (l *PointLayer) Init() {
	l.VertexList = make([]float32, 0, 1024)
	l.ColorList = make([]float32, 0, 1024)
	l.PointList = make([]Point, 0, 1024)
}

func (l *PointLayer) Render() {
	// Clear before start
	l.VertexList = l.VertexList[:0]
	l.ColorList = l.ColorList[:0]

	for i := 0; i < len(l.PointList); i++ {
		point := l.PointList[i]

		l.VertexList = append(l.VertexList, point.Position.X, point.Position.Y, point.Position.Z, point.Size)
		l.ColorList = append(l.ColorList, point.Color.R, point.Color.G, point.Color.B, point.Color.A)
	}

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

			/*"vertexAmount": l.VertexAmount,
			"colorAmount":  l.VertexAmount,*/

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		/*l.state["vertexAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.VertexAmount*/
	}

	return l.state
}
