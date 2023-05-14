package mrender_layer

import (
	"fmt"
	mmath_geom "github.com/maldan/go-ml/math/geometry"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mrender_uv "github.com/maldan/go-ml/render/uv"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"reflect"
	"unsafe"
)

type UI interface {
}

type UIElement struct {
	UvArea mmath_geom.Rectangle[float32]

	Position mmath_la.Vector2[float32]
	Rotation float32
	Scale    mmath_la.Vector2[float32]
	UvOffset mmath_la.Vector2[float32]
	Color    ml_color.ColorRGBA[float32]

	Pivot mmath_la.Vector2[float32]
}

type TextFont struct {
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
	FontMap     map[string]TextFont

	InstanceId int

	state map[string]any
}

type UIDrawArgs struct {
	UV          mmath_geom.Rectangle[float32]
	Position    mmath_la.Vector2[float32]
	Size        mmath_la.Vector2[float32]
	Rotation    float32
	PivotOffset mmath_la.Vector2[float32]
	Color       ml_color.ColorRGBA[float32]
}

type UIDrawTextArgs struct {
	FontName    string
	Text        string
	Position    mmath_la.Vector2[float32]
	Size        float32
	PivotOffset mmath_la.Vector2[float32]
	Color       ml_color.ColorRGBA[float32]
}

func (l *UILayer) Draw(args UIDrawArgs) {
	l.ElementList = append(l.ElementList, UIElement{
		UvArea:   args.UV,
		Position: args.Position,
		Rotation: args.Rotation,
		Scale:    args.Size,
		Color:    args.Color,
		Pivot:    args.PivotOffset,
	})
}

func (l *UILayer) DrawText(args UIDrawTextArgs) {
	font, ok := l.FontMap[args.FontName]
	if !ok {
		fmt.Printf("Font %v not found\n", args.FontName)
	}

	offsetX := float32(0)
	offsetY := float32(0)
	for _, ch := range args.Text {
		if ch == '\n' {
			offsetX = 0
			offsetY += font.Size.Y
			continue
		}
		ch2 := fmt.Sprintf("%c", ch)
		char := font.Symbol[ch2]
		rect := mrender_uv.GetArea(char.X, char.Y, font.Size.X, font.Size.Y, 1024, 1024)

		l.Draw(
			UIDrawArgs{
				UV:          rect,
				Position:    args.Position.AddXY(offsetX, offsetY),
				Size:        font.Size.Scale(args.Size),
				Color:       args.Color,
				PivotOffset: args.PivotOffset,
			},
		)

		if font.SymbolSize[ch2].X == 0 {
			offsetX += font.Size.X * args.Size
		} else {
			offsetX += font.SymbolSize[ch2].X * args.Size
		}
	}
	/*for i := 0; i < len(text); i++ {
		font[text[i]]
	}*/
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
	l.FontMap = map[string]TextFont{}

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
			-0.5+pv.X, -0.5+pv.Y,
			0.5+pv.X, -0.5+pv.Y,
			0.5+pv.X, 0.5+pv.Y,
			-0.5+pv.X, 0.5+pv.Y,
		)

		// Position list
		p := element.Position
		l.PositionList = append(l.PositionList, p.X, p.Y, p.X, p.Y, p.X, p.Y, p.X, p.Y)

		// Rotation list
		r := element.Rotation
		l.RotationList = append(l.RotationList, r, r, r, r)

		// Scale list
		s := element.Scale
		l.ScaleList = append(l.ScaleList, s.X, s.Y, s.X, s.Y, s.X, s.Y, s.X, s.Y)

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
