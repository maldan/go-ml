package mrender_layer

import (
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mr_camera "github.com/maldan/go-ml/render/camera"
	"reflect"
	"unsafe"
)

type Text struct {
	Font     string
	Content  string
	Size     float32
	Position mmath_la.Vector3[float32]
}

type TextFont struct {
	Symbol map[uint8]mmath_geom.Rectangle[float32]
}

type TextLayer struct {
	TextList []Text
	FontMap  map[string]TextFont

	VertexList []float32
	UvList     []float32

	PositionList []float32
	ColorList    []float32

	IndexList []uint16

	VertexAmount int
	IndexAmount  int
	UvAmount     int
	ColorAmount  int

	InstanceId int

	Camera mr_camera.PerspectiveCamera

	state map[string]any
}

func (l *TextLayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.UvList = make([]float32, 65536*2)
	l.PositionList = make([]float32, 65536*3)
	l.ColorList = make([]float32, 65536*4)
	l.IndexList = make([]uint16, 65536)

	l.TextList = make([]Text, 1024)
	l.FontMap = map[string]TextFont{}
}

func (l *TextLayer) Render() {
	l.VertexAmount = 0
	l.IndexAmount = 0
	l.UvAmount = 0
	l.ColorAmount = 0

	vertexId := 0
	positionId := 0
	uvId := 0
	colorId := 0
	indexId := uint16(0)
	lastMaxIndex := uint16(0)

	for i := 0; i < len(l.TextList); i++ {
		text := l.TextList[i]

		// Vertex
		for j := 0; j < len(text.Content); j++ {
			l.VertexList[vertexId] = -0.5 * text.Size
			l.VertexList[vertexId+1] = -0.5 * text.Size
			l.VertexList[vertexId+2] = 0

			l.VertexList[vertexId+3] = 0.5 * text.Size
			l.VertexList[vertexId+4] = -0.5 * text.Size
			l.VertexList[vertexId+5] = 0

			l.VertexList[vertexId+6] = 0.5 * text.Size
			l.VertexList[vertexId+7] = 0.5 * text.Size
			l.VertexList[vertexId+8] = 0

			l.VertexList[vertexId+9] = -0.5 * text.Size
			l.VertexList[vertexId+10] = 0.5 * text.Size
			l.VertexList[vertexId+11] = 0

			vertexId += 3 * 4
			l.VertexAmount += 3 * 4
		}

		// Position list
		for j := 0; j < len(text.Content); j++ {
			l.PositionList[positionId] = text.Position.X + float32(j)*text.Size
			l.PositionList[positionId+1] = text.Position.Y
			l.PositionList[positionId+2] = text.Position.Z

			l.PositionList[positionId+3] = text.Position.X + float32(j)*text.Size
			l.PositionList[positionId+4] = text.Position.Y
			l.PositionList[positionId+5] = text.Position.Z

			l.PositionList[positionId+6] = text.Position.X + float32(j)*text.Size
			l.PositionList[positionId+7] = text.Position.Y
			l.PositionList[positionId+8] = text.Position.Z

			l.PositionList[positionId+9] = text.Position.X + float32(j)*text.Size
			l.PositionList[positionId+10] = text.Position.Y
			l.PositionList[positionId+11] = text.Position.Z

			positionId += 3 * 4
		}

		// Color list
		for j := 0; j < len(text.Content); j++ {
			l.ColorList[colorId] = 1.0
			l.ColorList[colorId+1] = 1.0
			l.ColorList[colorId+2] = 1.0

			l.ColorList[colorId+3] = 1.0
			l.ColorList[colorId+4] = 1.0
			l.ColorList[colorId+5] = 1.0

			l.ColorList[colorId+6] = 1.0
			l.ColorList[colorId+7] = 1.0
			l.ColorList[colorId+8] = 1.0

			l.ColorList[colorId+9] = 1.0
			l.ColorList[colorId+10] = 1.0
			l.ColorList[colorId+11] = 1.0

			colorId += 3 * 4
		}

		// Uv
		for j := 0; j < len(text.Content); j++ {
			rect := l.FontMap[text.Font].Symbol[text.Content[j]]

			l.UvList[uvId] = rect.Left
			l.UvList[uvId+1] = rect.Bottom

			l.UvList[uvId+2] = rect.Right
			l.UvList[uvId+3] = rect.Bottom

			l.UvList[uvId+4] = rect.Right
			l.UvList[uvId+5] = rect.Top

			l.UvList[uvId+6] = rect.Left
			l.UvList[uvId+7] = rect.Top

			uvId += 2 * 4
			l.UvAmount += 2 * 4
		}

		// Indices
		for j := 0; j < len(text.Content); j++ {
			l.IndexList[indexId] = 0 + lastMaxIndex
			l.IndexList[indexId+1] = 1 + lastMaxIndex
			l.IndexList[indexId+2] = 2 + lastMaxIndex
			l.IndexList[indexId+3] = 0 + lastMaxIndex
			l.IndexList[indexId+4] = 2 + lastMaxIndex
			l.IndexList[indexId+5] = 3 + lastMaxIndex

			indexId += 6
			lastMaxIndex += 4
			l.IndexAmount += 6
		}
	}

	if len(l.TextList) > 0 {
		l.TextList = l.TextList[:0]
	}
}

func (l *TextLayer) GetState() map[string]any {
	vertexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.VertexList))
	uvHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.UvList))
	positionHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.PositionList))
	indexHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.IndexList))
	colorHeader := (*reflect.SliceHeader)(unsafe.Pointer(&l.ColorList))

	if l.state == nil {
		l.state = map[string]any{
			"vertexPointer":   vertexHeader.Data,
			"uvPointer":       uvHeader.Data,
			"indexPointer":    indexHeader.Data,
			"positionPointer": positionHeader.Data,
			"colorPointer":    colorHeader.Data,
			"vertexAmount":    l.VertexAmount,
			"uvAmount":        l.UvAmount,
			"indexAmount":     l.IndexAmount,
			"positionAmount":  l.VertexAmount,
			"colorAmount":     l.VertexAmount,

			"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		l.state["vertexAmount"] = l.VertexAmount
		l.state["positionAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.VertexAmount

		l.state["uvAmount"] = l.UvAmount
		l.state["indexAmount"] = l.IndexAmount
	}

	return l.state
}