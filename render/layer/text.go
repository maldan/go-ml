package mrender_layer

import (
	"fmt"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	mrender_uv "github.com/maldan/go-ml/render/uv"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"reflect"
	"unsafe"
)

type Text struct {
	Font             string
	Content          string
	Size             float32
	Position         mmath_la.Vector3[float32]
	Rotation         mmath_la.Vector3[float32]
	NewLineDirection mmath_la.Vector3[float32]
	Color            ml_color.ColorRGBA[float32]
}

/*type TextFont struct {
	Symbol map[uint8]mmath_geom.Rectangle[float32]
}*/

type TextLayer struct {
	VertexList   []float32
	UvList       []float32
	PositionList []float32
	RotationList []float32
	ColorList    []float32
	IndexList    []uint16

	TextList []Text
	FontMap  map[string]TextFont

	state map[string]any
}

type TextDrawArgs struct {
	FontName         string
	Text             string
	Position         mmath_la.Vector3[float32]
	Rotation         mmath_la.Vector3[float32]
	Size             float32
	NewLineDirection mmath_la.Vector3[float32]
	PivotOffset      mmath_la.Vector2[float32]
	Color            ml_color.ColorRGBA[float32]
}

func (l *TextLayer) Draw(args TextDrawArgs) {
	l.TextList = append(l.TextList, Text{
		Font:             args.FontName,
		Content:          args.Text,
		Size:             args.Size,
		Position:         args.Position,
		Rotation:         args.Rotation,
		NewLineDirection: args.NewLineDirection,
		Color:            args.Color,
	})
}

func (l *TextLayer) Init() {
	l.VertexList = make([]float32, 0, 1024)
	l.UvList = make([]float32, 0, 1024)
	l.PositionList = make([]float32, 0, 1024)
	l.RotationList = make([]float32, 0, 1024)
	l.ColorList = make([]float32, 0, 1024)
	l.IndexList = make([]uint16, 0, 1024)

	l.TextList = make([]Text, 1024)
	l.FontMap = map[string]TextFont{}
}

/*func (l *TextLayer) Render() {
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

			l.UvList[uvId] = rect.MinX
			l.UvList[uvId+1] = rect.MaxY

			l.UvList[uvId+2] = rect.MaxX
			l.UvList[uvId+3] = rect.MaxY

			l.UvList[uvId+4] = rect.MaxX
			l.UvList[uvId+5] = rect.MinY

			l.UvList[uvId+6] = rect.MinX
			l.UvList[uvId+7] = rect.MinY

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
}*/

func (l *TextLayer) Render() {
	// Clear before start
	l.VertexList = l.VertexList[:0]
	l.UvList = l.UvList[:0]
	l.PositionList = l.PositionList[:0]
	l.RotationList = l.RotationList[:0]
	l.ColorList = l.ColorList[:0]
	l.IndexList = l.IndexList[:0]

	lastMaxIndex := uint16(0)
	for i := 0; i < len(l.TextList); i++ {
		text := l.TextList[i]
		font := l.FontMap[text.Font]

		textOffset := mmath_la.Vector3[float32]{}
		for _, ch := range text.Content {
			if ch == '\n' {
				textOffset.X = 0
				textOffset = textOffset.Add(text.NewLineDirection)
				continue
			}
			char := font.Symbol[fmt.Sprintf("%c", ch)]
			rect := mrender_uv.GetArea(char.X, char.Y, font.Size.X, font.Size.Y, 1024, 1024)

			l.VertexList = append(l.VertexList,
				-0.5*text.Size, -0.5*text.Size, 0,
				0.5*text.Size, -0.5*text.Size, 0,
				0.5*text.Size, 0.5*text.Size, 0,
				-0.5*text.Size, 0.5*text.Size, 0,
			)

			// Position list
			p := text.Position.Add(textOffset)
			l.PositionList = append(l.PositionList, p.X, p.Y, p.Z, p.X, p.Y, p.Z, p.X, p.Y, p.Z, p.X, p.Y, p.Z)

			// Rotation list
			r := text.Rotation
			l.RotationList = append(l.RotationList, r.X, r.Y, r.Z, r.X, r.Y, r.Z, r.X, r.Y, r.Z, r.X, r.Y, r.Z)

			// Color list
			c := text.Color
			l.ColorList = append(l.ColorList, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A, c.R, c.G, c.B, c.A)

			// Uv
			l.UvList = append(l.UvList,
				rect.MinX, rect.MaxY,
				rect.MaxX, rect.MaxY,
				rect.MaxX, rect.MinY,
				rect.MinX, rect.MinY,
			)

			// Indices
			l.IndexList = append(l.IndexList,
				0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex,
				0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex,
			)

			lastMaxIndex += 4

			textOffset.X += text.Size
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
			/*"vertexAmount":    l.VertexAmount,
			"uvAmount":        l.UvAmount,
			"indexAmount":     l.IndexAmount,
			"positionAmount":  l.VertexAmount,
			"colorAmount":     l.VertexAmount,*/

			//"projectionMatrixPointer": uintptr(unsafe.Pointer(&l.Camera.Matrix.Raw)),
		}
	} else {
		/*l.state["vertexAmount"] = l.VertexAmount
		l.state["positionAmount"] = l.VertexAmount
		l.state["colorAmount"] = l.VertexAmount

		l.state["uvAmount"] = l.UvAmount
		l.state["indexAmount"] = l.IndexAmount*/
	}

	return l.state
}
