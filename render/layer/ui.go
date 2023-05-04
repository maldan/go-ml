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
	VertexList   []float32
	UvList       []float32
	NormalList   []float32
	PositionList []float32
	RotationList []float32
	ScaleList    []float32
	ColorList    []float32
	IndexList    []uint16

	ElementList []UIElement
	FontMap     map[string]UITextFont

	InstanceId int

	state map[string]any
}

func (l *UILayer) Init() {
	l.VertexList = make([]float32, 0, 1024)
	l.NormalList = make([]float32, 0, 1024)
	l.UvList = make([]float32, 0, 1024)
	l.PositionList = make([]float32, 0, 1024)
	l.RotationList = make([]float32, 0, 1024)
	l.ScaleList = make([]float32, 0, 1024)
	l.ColorList = make([]float32, 0, 1024)

	l.ElementList = make([]UIElement, 0, 1024)
	// l.MeshInstanceList = make([]mr_mesh.MeshInstance, 1024)
	l.IndexList = make([]uint16, 0, 1024)
	l.FontMap = map[string]UITextFont{}

	fmt.Printf("Render allocated %v\n", cap(l.VertexList)*4*6)
}

func (l *UILayer) Render() {
	// Clear before start
	l.VertexList = l.VertexList[:0]
	l.UvList = l.UvList[:0]
	l.NormalList = l.NormalList[:0]
	l.PositionList = l.PositionList[:0]
	l.RotationList = l.RotationList[:0]
	l.ScaleList = l.ScaleList[:0]
	l.ColorList = l.ColorList[:0]
	l.IndexList = l.IndexList[:0]

	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.ElementList); i++ {
		element := l.ElementList[i]

		// Vertex
		pv := element.Pivot
		l.VertexList = append(l.VertexList,
			-0.5+pv.X, -0.5+pv.Y, 0,
			0.5+pv.X, -0.5+pv.Y, 0,
			0.5+pv.X, 0.5+pv.Y, 0,
			-0.5+pv.X, 0.5+pv.Y, 0,
		)

		// Position list
		p := element.Position
		l.PositionList = append(l.PositionList, p.X, p.Y, p.Z, p.X, p.Y, p.Z, p.X, p.Y, p.Z, p.X, p.Y, p.Z)

		// Rotation list
		r := element.Rotation
		l.RotationList = append(l.RotationList, r.X, r.Y, r.Z, r.X, r.Y, r.Z, r.X, r.Y, r.Z, r.X, r.Y, r.Z)

		// Scale list
		s := element.Scale
		l.ScaleList = append(l.ScaleList, s.X, s.Y, s.Z, s.X, s.Y, s.Z, s.X, s.Y, s.Z, s.X, s.Y, s.Z)

		// Color list
		c := element.Color
		l.ColorList = append(l.ColorList, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A)

		// Uv
		l.UvList = append(l.UvList,
			element.UvArea.MinX, element.UvArea.MaxY,
			element.UvArea.MaxX, element.UvArea.MaxY,
			element.UvArea.MaxX, element.UvArea.MinY,
			element.UvArea.MinX, element.UvArea.MinY,
		)

		// Indices
		l.IndexList = append(l.IndexList,
			0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex,
			0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex,
		)

		lastMaxIndex += 4
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

			/*"vertexAmount":   l.VertexAmount,
			"normalAmount":   l.VertexAmount,
			"uvAmount":       l.UvAmount,
			"indexAmount":    l.IndexAmount,
			"positionAmount": l.VertexAmount,
			"rotationAmount": l.VertexAmount,
			"scaleAmount":    l.VertexAmount,
			"colorAmount":    l.ColorAmount,*/

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		/*l.state["vertexAmount"] = l.VertexAmount
		l.state["normalAmount"] = l.VertexAmount
		l.state["uvAmount"] = l.UvAmount
		l.state["indexAmount"] = l.IndexAmount
		l.state["positionAmount"] = l.VertexAmount
		l.state["rotationAmount"] = l.VertexAmount
		l.state["scaleAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.ColorAmount*/
	}

	return l.state
}
