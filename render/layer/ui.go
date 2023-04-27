package mrender_layer

import (
	"fmt"
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"reflect"
	"unsafe"
)

type UI interface {
}

type UIElement struct {
	UvArea mmath_geom.Rectangle[float32]

	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Vector3[float32]
	Scale    mmath_la.Vector3[float32]
	UvOffset mmath_la.Vector2[float32]
	Color    ml_color.ColorRGBA[float32]

	Pivot mmath_la.Vector2[float32]

	IsVisible bool
	IsActive  bool
}

type UITextFont struct {
	Size       mmath_la.Vector2[float32]            `json:"size"`
	Symbol     map[string]mmath_la.Vector2[float32] `json:"symbol"`
	SymbolSize map[string]mmath_la.Vector2[float32] `json:"symbolSize"`
}

type UILayer struct {
	ElementList []UIElement
	FontMap     map[string]UITextFont

	VertexList []float32
	UvList     []float32
	NormalList []float32

	PositionList []float32
	RotationList []float32
	ScaleList    []float32
	ColorList    []float32

	IndexList []uint16

	VertexAmount int
	IndexAmount  int
	UvAmount     int
	ColorAmount  int

	//Camera mr_camera.OrthographicCamera

	InstanceId int

	state map[string]any
}

func (l *UILayer) Init() {
	l.VertexList = make([]float32, 65536*3)
	l.NormalList = make([]float32, 65536*3)
	l.UvList = make([]float32, 65536*2)
	l.PositionList = make([]float32, 65536*3)
	l.RotationList = make([]float32, 65536*3)
	l.ScaleList = make([]float32, 65536*3)
	l.ColorList = make([]float32, 65536*4)

	l.ElementList = make([]UIElement, 0, 1024)
	// l.MeshInstanceList = make([]mr_mesh.MeshInstance, 1024)
	l.IndexList = make([]uint16, 65536)
	l.FontMap = map[string]UITextFont{}

	fmt.Printf("Render allocated %v\n", cap(l.VertexList)*4*6)
}

func (l *UILayer) Render() {
	//l.Camera.ApplyMatrix()

	l.VertexAmount = 0
	l.IndexAmount = 0
	l.UvAmount = 0
	l.ColorAmount = 0

	vertexId := 0
	positionId := 0
	rotationId := 0
	scaleId := 0
	indexId := 0
	colorId := 0
	uvId := 0
	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.ElementList); i++ {
		element := l.ElementList[i]

		// Vertex
		l.VertexList[vertexId] = -0.5 + element.Pivot.X
		l.VertexList[vertexId+1] = -0.5 + element.Pivot.Y
		l.VertexList[vertexId+2] = 0

		l.VertexList[vertexId+3] = 0.5 + element.Pivot.X
		l.VertexList[vertexId+4] = -0.5 + element.Pivot.Y
		l.VertexList[vertexId+5] = 0

		l.VertexList[vertexId+6] = 0.5 + element.Pivot.X
		l.VertexList[vertexId+7] = 0.5 + element.Pivot.Y
		l.VertexList[vertexId+8] = 0

		l.VertexList[vertexId+9] = -0.5 + element.Pivot.X
		l.VertexList[vertexId+10] = 0.5 + element.Pivot.Y
		l.VertexList[vertexId+11] = 0

		vertexId += 3 * 4
		l.VertexAmount += 3 * 4

		// Position list
		l.PositionList[positionId] = element.Position.X
		l.PositionList[positionId+1] = element.Position.Y
		l.PositionList[positionId+2] = element.Position.Z

		l.PositionList[positionId+3] = element.Position.X
		l.PositionList[positionId+4] = element.Position.Y
		l.PositionList[positionId+5] = element.Position.Z

		l.PositionList[positionId+6] = element.Position.X
		l.PositionList[positionId+7] = element.Position.Y
		l.PositionList[positionId+8] = element.Position.Z

		l.PositionList[positionId+9] = element.Position.X
		l.PositionList[positionId+10] = element.Position.Y
		l.PositionList[positionId+11] = element.Position.Z

		positionId += 3 * 4

		// Rotation list
		l.RotationList[rotationId] = element.Rotation.X
		l.RotationList[rotationId+1] = element.Rotation.Y
		l.RotationList[rotationId+2] = element.Rotation.Z

		l.RotationList[rotationId+3] = element.Rotation.X
		l.RotationList[rotationId+4] = element.Rotation.Y
		l.RotationList[rotationId+5] = element.Rotation.Z

		l.RotationList[rotationId+6] = element.Rotation.X
		l.RotationList[rotationId+7] = element.Rotation.Y
		l.RotationList[rotationId+8] = element.Rotation.Z

		l.RotationList[rotationId+9] = element.Rotation.X
		l.RotationList[rotationId+10] = element.Rotation.Y
		l.RotationList[rotationId+11] = element.Rotation.Z

		rotationId += 3 * 4

		// Scale list
		l.ScaleList[scaleId] = element.Scale.X
		l.ScaleList[scaleId+1] = element.Scale.Y
		l.ScaleList[scaleId+2] = element.Scale.Z

		l.ScaleList[scaleId+3] = element.Scale.X
		l.ScaleList[scaleId+4] = element.Scale.Y
		l.ScaleList[scaleId+5] = element.Scale.Z

		l.ScaleList[scaleId+6] = element.Scale.X
		l.ScaleList[scaleId+7] = element.Scale.Y
		l.ScaleList[scaleId+8] = element.Scale.Z

		l.ScaleList[scaleId+9] = element.Scale.X
		l.ScaleList[scaleId+10] = element.Scale.Y
		l.ScaleList[scaleId+11] = element.Scale.Z

		scaleId += 3 * 4

		// Color list
		l.ColorList[colorId] = element.Color.R
		l.ColorList[colorId+1] = element.Color.G
		l.ColorList[colorId+2] = element.Color.B
		l.ColorList[colorId+3] = element.Color.A

		l.ColorList[colorId+4] = element.Color.R
		l.ColorList[colorId+5] = element.Color.G
		l.ColorList[colorId+6] = element.Color.B
		l.ColorList[colorId+7] = element.Color.A

		l.ColorList[colorId+8] = element.Color.R
		l.ColorList[colorId+9] = element.Color.G
		l.ColorList[colorId+10] = element.Color.B
		l.ColorList[colorId+11] = element.Color.A

		l.ColorList[colorId+12] = element.Color.R
		l.ColorList[colorId+13] = element.Color.G
		l.ColorList[colorId+14] = element.Color.B
		l.ColorList[colorId+15] = element.Color.A

		colorId += 4 * 4
		l.ColorAmount += 4 * 4

		// Uv
		l.UvList[uvId] = element.UvArea.Left
		l.UvList[uvId+1] = element.UvArea.Bottom

		l.UvList[uvId+2] = element.UvArea.Right
		l.UvList[uvId+3] = element.UvArea.Bottom

		l.UvList[uvId+4] = element.UvArea.Right
		l.UvList[uvId+5] = element.UvArea.Top

		l.UvList[uvId+6] = element.UvArea.Left
		l.UvList[uvId+7] = element.UvArea.Top

		uvId += 2 * 4
		l.UvAmount += 2 * 4

		// Indices
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

	if len(l.ElementList) > 0 {
		l.ElementList = l.ElementList[:0]
	}
}

func (l *UILayer) GetState() map[string]any {
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

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
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
