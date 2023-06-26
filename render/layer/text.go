package mrender_layer

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
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
	NewLineChars     string
	Color            ml_color.ColorRGBA[float32]
	IsUI             bool
	MaxWidth         int
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
	NewLineChars     string
	PivotOffset      mmath_la.Vector2[float32]
	Color            ml_color.ColorRGBA[float32]
	IsUI             bool
	MaxWidth         int
}

func (l *TextLayer) Draw(args TextDrawArgs) {
	l.TextList = append(l.TextList, Text{
		Font:             args.FontName,
		Content:          args.Text,
		Size:             args.Size,
		Position:         args.Position,
		Rotation:         args.Rotation,
		NewLineDirection: args.NewLineDirection,
		NewLineChars:     args.NewLineChars,
		Color:            args.Color,
		IsUI:             args.IsUI,
		MaxWidth:         args.MaxWidth,
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

		textOffset := mmath_la.Vector3[float32]{}
		for _, ch := range text.Content {
			letter := ch
			size := text.Size
			isUI := float32(0)
			if text.IsUI {
				isUI = 1
			}

			letterVertexLength := 0
			switch letter {

			case '0':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case '1':
				letterVertexLength = 7
				l.VertexList = append(l.VertexList, 0*size, -1*size, isUI, 0.25*size, -1*size, isUI, 0.25*size, 0.75*size, isUI, 0*size, 0.75*size, isUI, -0.25*size, 0.5*size, isUI, 0*size, 0.5*size, isUI, -0.25*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 7
				break

			case '2':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.75*size, 0*size, isUI, 0.5*size, 0*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 6+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 11+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex)
				lastMaxIndex += 16
				break

			case '3':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.25*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.25*size, 0*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case '4':
				letterVertexLength = 11
				l.VertexList = append(l.VertexList, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.75*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0.75*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex)
				lastMaxIndex += 11
				break

			case '5':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex)
				lastMaxIndex += 16
				break

			case '6':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 3+lastMaxIndex, 17+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 18
				break

			case '7':
				letterVertexLength = 7
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, -0.5*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 7
				break

			case '8':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 12+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex)
				lastMaxIndex += 16
				break

			case '9':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.75*size, 0*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0.5*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 11+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex)
				lastMaxIndex += 16
				break

			case '!':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, -0.75*size, isUI, -0.125*size, -0.75*size, isUI, -0.125*size, -0.5*size, isUI, 0.125*size, -0.5*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 8
				break

			case '(':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.25*size, -1*size, isUI, 0*size, -1*size, isUI, -0.5*size, -0.75*size, isUI, -0.25*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, -0.25*size, 0.5*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 3+lastMaxIndex, 0+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 8
				break

			case ')':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.25*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.25*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 3+lastMaxIndex, 0+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 8
				break

			case '+':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.5*size, -0.125*size, isUI, 0.5*size, -0.125*size, isUI, 0.5*size, 0.125*size, isUI, -0.5*size, 0.125*size, isUI, -0.125*size, -0.5*size, isUI, 0.125*size, -0.5*size, isUI, 0.125*size, 0.5*size, isUI, -0.125*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 8
				break

			case ',':
				letterVertexLength = 5
				l.VertexList = append(l.VertexList, -0.125*size, -0.75*size, isUI, 0.125*size, -0.75*size, isUI, 0.125*size, -0.5*size, isUI, -0.125*size, -0.5*size, isUI, -0.25*size, -1*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex)
				lastMaxIndex += 5
				break

			case '-':
				letterVertexLength = 4
				l.VertexList = append(l.VertexList, -0.5*size, -0.125*size, isUI, 0.5*size, -0.125*size, isUI, 0.5*size, 0.125*size, isUI, -0.5*size, 0.125*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex)
				lastMaxIndex += 4
				break

			case '.':
				letterVertexLength = 4
				l.VertexList = append(l.VertexList, -0.125*size, -0.75*size, isUI, 0.125*size, -0.75*size, isUI, 0.125*size, -0.5*size, isUI, -0.125*size, -0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex)
				lastMaxIndex += 4
				break

			case '/':
				letterVertexLength = 4
				l.VertexList = append(l.VertexList, -0.5*size, -1*size, isUI, -0.25*size, -1*size, isUI, 0.5*size, 0.75*size, isUI, 0.25*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex)
				lastMaxIndex += 4
				break

			case '=':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.5*size, 0.125*size, isUI, 0.5*size, 0.125*size, isUI, 0.5*size, 0.375*size, isUI, -0.5*size, 0.375*size, isUI, -0.5*size, -0.375*size, isUI, 0.5*size, -0.375*size, isUI, 0.5*size, -0.125*size, isUI, -0.5*size, -0.125*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 8
				break

			case '?':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, -0.75*size, isUI, -0.125*size, -0.75*size, isUI, -0.125*size, -0.375*size, isUI, 0.125*size, -0.5*size, isUI, 0.75*size, 0*size, isUI, 0.5*size, 0.125*size, isUI, 0.5*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.25*size, 0.5*size, isUI, -0.5*size, 0*size, isUI, -0.25*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 10+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'A':
				letterVertexLength = 15
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.5*size, isUI, 0.5*size, -0.5*size, isUI, 0.5*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex)
				lastMaxIndex += 15
				break

			case 'B':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 5+lastMaxIndex, 14+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex, 16+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex, 16+lastMaxIndex, 14+lastMaxIndex, 17+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'C':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'D':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 7+lastMaxIndex, 10+lastMaxIndex, 7+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'E':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'F':
				letterVertexLength = 11
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex)
				lastMaxIndex += 11
				break

			case 'G':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, -0.75*size, isUI, 0.75*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, 0.25*size, 0*size, isUI, 0.25*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 11+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 12+lastMaxIndex, 14+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 16
				break

			case 'H':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'I':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, -0.75*size, isUI, -0.125*size, -0.75*size, isUI, -0.125*size, 0.5*size, isUI, 0.125*size, 0.5*size, isUI, -0.125*size, 0.75*size, isUI, 0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 8
				break

			case 'J':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, 0.25*size, -1*size, isUI, 0.5*size, -1*size, isUI, 0.5*size, 0.75*size, isUI, 0.25*size, 0.75*size, isUI, -0.75*size, -1*size, isUI, -0.75*size, -0.75*size, isUI, 0.25*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, -0.75*size, -0.5*size, isUI, -0.5*size, -0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 5+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex, 8+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'K':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.375*size, isUI, -0.5*size, 0*size, isUI, 0.125*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.125*size, -1*size, isUI, 0.5*size, -1*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'L':
				letterVertexLength = 7
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 1+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 1+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 7
				break

			case 'M':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, 0.375*size, isUI, 0*size, 0.25*size, isUI, 0*size, -0.125*size, isUI, 0.5*size, 0.375*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 2+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'N':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, -0.5*size, 0.25*size, isUI, 0.5*size, -0.5*size, isUI, 0.5*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'O':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'P':
				letterVertexLength = 13
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, 0.5*size, 0*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex)
				lastMaxIndex += 13
				break

			case 'Q':
				letterVertexLength = 17
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.25*size, -0.625*size, isUI, 0.375*size, -0.5*size, isUI, 0.75*size, -1*size, isUI, 0.625*size, -0.875*size, isUI, 0.875*size, -1*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 16+lastMaxIndex, 15+lastMaxIndex)
				lastMaxIndex += 17
				break

			case 'R':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, 0.5*size, 0*size, isUI, 0.5*size, 0.5*size, isUI, -0.5*size, -0.5*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 16
				break

			case 'S':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex)
				lastMaxIndex += 16
				break

			case 'T':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, 0.5*size, isUI, -0.125*size, 0.5*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 8
				break

			case 'U':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 1+lastMaxIndex, 7+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'V':
				letterVertexLength = 6
				l.VertexList = append(l.VertexList, 0*size, -0.375*size, isUI, 0*size, -1*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 3+lastMaxIndex, 0+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 4+lastMaxIndex, 0+lastMaxIndex, 5+lastMaxIndex)
				lastMaxIndex += 6
				break

			case 'W':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -0.75*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0*size, -0.75*size, isUI, 0*size, -0.5*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 1+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 4+lastMaxIndex, 11+lastMaxIndex, 8+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'X':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, -0.25*size, -0.25*size, isUI, -0.25*size, 0*size, isUI, 0*size, 0*size, isUI, 0*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 4+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 4+lastMaxIndex, 11+lastMaxIndex, 8+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 12+lastMaxIndex, 6+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 13+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 8+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 8+lastMaxIndex, 13+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Y':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, 0*size, 0.125*size, isUI, 0*size, -0.25*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, -0.125*size, -0.125*size, isUI, 0.125*size, -0.125*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 3+lastMaxIndex, 0+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 4+lastMaxIndex, 0+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 8+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'Z':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.375*size, -0.75*size, isUI, 0.375*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 3+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'А':
				letterVertexLength = 15
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.5*size, isUI, 0.5*size, -0.5*size, isUI, 0.5*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex)
				lastMaxIndex += 15
				break

			case 'Б':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 3+lastMaxIndex, 17+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'В':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 5+lastMaxIndex, 14+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex, 16+lastMaxIndex, 11+lastMaxIndex, 14+lastMaxIndex, 16+lastMaxIndex, 14+lastMaxIndex, 17+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'Г':
				letterVertexLength = 7
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 3+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 1+lastMaxIndex)
				lastMaxIndex += 7
				break

			case 'Д':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.75*size, -0.5*size, isUI, -0.5*size, -0.5*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.5*size, -0.5*size, isUI, 0.75*size, -0.5*size, isUI, -0.5*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.25*size, 0.5*size, isUI, 0.25*size, 0.5*size, isUI, -0.75*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI, -0.25*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 3+lastMaxIndex, 0+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 7+lastMaxIndex, 13+lastMaxIndex, 2+lastMaxIndex, 13+lastMaxIndex, 12+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 10+lastMaxIndex, 14+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 9+lastMaxIndex, 16+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'Е':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Ж':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, -0.25*size, -0.25*size, isUI, -0.25*size, 0*size, isUI, 0*size, 0*size, isUI, 0*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, -0.125*size, 0.75*size, isUI, 0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 4+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 4+lastMaxIndex, 11+lastMaxIndex, 8+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 12+lastMaxIndex, 6+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 13+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 8+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 8+lastMaxIndex, 13+lastMaxIndex, 9+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 17+lastMaxIndex, 14+lastMaxIndex, 17+lastMaxIndex, 16+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'З':
				letterVertexLength = 20
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, 0*size, 0*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0.5*size, isUI, -0.75*size, -0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.5*size, isUI, -0.75*size, 0.25*size, isUI, -0.5*size, 0.25*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex, 3+lastMaxIndex, 15+lastMaxIndex, 16+lastMaxIndex, 3+lastMaxIndex, 16+lastMaxIndex, 14+lastMaxIndex, 17+lastMaxIndex, 18+lastMaxIndex, 19+lastMaxIndex, 17+lastMaxIndex, 19+lastMaxIndex, 4+lastMaxIndex)
				lastMaxIndex += 20
				break

			case 'И':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, 0.125*size, isUI, 0.5*size, 0.5*size, isUI, -0.5*size, -0.375*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'Й':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, 0.125*size, isUI, 0.5*size, 0.5*size, isUI, -0.5*size, -0.375*size, isUI, -0.25*size, 0.75*size, isUI, 0.25*size, 0.75*size, isUI, -0.25*size, 1*size, isUI, 0.25*size, 1*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 12+lastMaxIndex, 15+lastMaxIndex, 14+lastMaxIndex)
				lastMaxIndex += 16
				break

			case 'К':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0.125*size, isUI, 0.125*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.125*size, -1*size, isUI, 0.5*size, -1*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'Л':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'М':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, 0.375*size, isUI, 0*size, 0.25*size, isUI, 0*size, -0.125*size, isUI, 0.5*size, 0.375*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 2+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 7+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'Н':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'О':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'П':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'Р':
				letterVertexLength = 13
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.25*size, isUI, 0.75*size, -0.25*size, isUI, 0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, 0.5*size, 0*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex)
				lastMaxIndex += 13
				break

			case 'С':
				letterVertexLength = 10
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 5+lastMaxIndex)
				lastMaxIndex += 10
				break

			case 'Т':
				letterVertexLength = 8
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, 0.5*size, isUI, -0.125*size, 0.5*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 8
				break

			case 'У':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.75*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0.75*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.75*size, -1*size, isUI, 0.5*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 0+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Ф':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, -0.125*size, -1*size, isUI, 0.125*size, -1*size, isUI, 0.125*size, 0.875*size, isUI, -0.125*size, 0.875*size, isUI, -0.75*size, 0.375*size, isUI, 0.75*size, 0.375*size, isUI, -0.75*size, 0.625*size, isUI, 0.75*size, 0.625*size, isUI, -0.75*size, -0.5*size, isUI, 0.75*size, -0.5*size, isUI, 0.75*size, -0.25*size, isUI, -0.75*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0.375*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0.375*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 4+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 11+lastMaxIndex, 13+lastMaxIndex, 4+lastMaxIndex, 14+lastMaxIndex, 10+lastMaxIndex, 5+lastMaxIndex, 14+lastMaxIndex, 5+lastMaxIndex, 15+lastMaxIndex)
				lastMaxIndex += 16
				break

			case 'Х':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, -0.25*size, -0.25*size, isUI, -0.25*size, 0*size, isUI, 0*size, 0*size, isUI, 0*size, -0.25*size, isUI, 0.25*size, -0.25*size, isUI, 0.25*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 4+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 4+lastMaxIndex, 11+lastMaxIndex, 8+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 12+lastMaxIndex, 6+lastMaxIndex, 12+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 13+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 0+lastMaxIndex, 8+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 8+lastMaxIndex, 13+lastMaxIndex, 9+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Ц':
				letterVertexLength = 12
				l.VertexList = append(l.VertexList, -0.75*size, -0.75*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.75*size, -0.75*size, isUI, 0.75*size, -0.5*size, isUI, -0.75*size, -0.5*size, isUI, -0.5*size, -0.5*size, isUI, 0.5*size, -0.5*size, isUI, 1*size, -1*size, isUI, 1*size, -0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 1+lastMaxIndex, 7+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 5+lastMaxIndex, 11+lastMaxIndex, 6+lastMaxIndex)
				lastMaxIndex += 12
				break

			case 'Ч':
				letterVertexLength = 11
				l.VertexList = append(l.VertexList, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.75*size, -0.25*size, isUI, -0.5*size, -0.25*size, isUI, -0.5*size, 0.75*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 5+lastMaxIndex, 7+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex)
				lastMaxIndex += 11
				break

			case 'Ш':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI, -0.125*size, -0.75*size, isUI, 0.125*size, -0.75*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 1+lastMaxIndex, 7+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Щ':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.75*size, isUI, 1*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.5*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI, -0.125*size, -0.75*size, isUI, 0.125*size, -0.75*size, isUI, 0.125*size, 0.75*size, isUI, -0.125*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 0+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 1+lastMaxIndex, 7+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 9+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 9+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 10+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Ы':
				letterVertexLength = 22
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.25*size, -1*size, isUI, 0.25*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0*size, -0.75*size, isUI, 0*size, 0*size, isUI, 0.25*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 3+lastMaxIndex, 17+lastMaxIndex, 11+lastMaxIndex, 18+lastMaxIndex, 19+lastMaxIndex, 20+lastMaxIndex, 18+lastMaxIndex, 20+lastMaxIndex, 21+lastMaxIndex)
				lastMaxIndex += 22
				break

			case 'Ь':
				letterVertexLength = 18
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0*size, isUI, 0.75*size, 0*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, -0.75*size, 0*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI, -0.5*size, -0.75*size, isUI, -0.5*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 2+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 11+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 13+lastMaxIndex, 15+lastMaxIndex, 4+lastMaxIndex, 3+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 3+lastMaxIndex, 17+lastMaxIndex, 11+lastMaxIndex)
				lastMaxIndex += 18
				break

			case 'Э':
				letterVertexLength = 14
				l.VertexList = append(l.VertexList, -0.75*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, -0.75*size, -0.75*size, isUI, -0.75*size, 0.5*size, isUI, 0.75*size, 0.5*size, isUI, 0.75*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.5*size, 0*size, isUI, 0.5*size, -0.75*size, isUI, 0.5*size, 0.5*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 10+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 2+lastMaxIndex, 5+lastMaxIndex, 12+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex)
				lastMaxIndex += 14
				break

			case 'Ю':
				letterVertexLength = 20
				l.VertexList = append(l.VertexList, -0.25*size, -0.75*size, isUI, 0*size, -1*size, isUI, 0*size, 0.75*size, isUI, -0.25*size, 0.5*size, isUI, 0*size, 0.5*size, isUI, 0.5*size, 0.5*size, isUI, 0.5*size, 0.75*size, isUI, 0.75*size, 0.5*size, isUI, 0.5*size, -1*size, isUI, 0.75*size, -0.75*size, isUI, 0.5*size, -0.75*size, isUI, 0*size, -0.75*size, isUI, -0.75*size, -1*size, isUI, -0.5*size, -1*size, isUI, -0.5*size, 0.75*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, -0.25*size, isUI, -0.25*size, -0.25*size, isUI, -0.25*size, 0*size, isUI, -0.5*size, 0*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 8+lastMaxIndex, 9+lastMaxIndex, 7+lastMaxIndex, 8+lastMaxIndex, 7+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 2+lastMaxIndex, 1+lastMaxIndex, 8+lastMaxIndex, 10+lastMaxIndex, 1+lastMaxIndex, 10+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 12+lastMaxIndex, 14+lastMaxIndex, 15+lastMaxIndex, 16+lastMaxIndex, 17+lastMaxIndex, 18+lastMaxIndex, 16+lastMaxIndex, 18+lastMaxIndex, 19+lastMaxIndex)
				lastMaxIndex += 20
				break

			case 'Я':
				letterVertexLength = 16
				l.VertexList = append(l.VertexList, 0.5*size, -1*size, isUI, 0.75*size, -1*size, isUI, 0.75*size, 0.75*size, isUI, 0.5*size, 0.75*size, isUI, -0.75*size, -0.25*size, isUI, 0.5*size, -0.25*size, isUI, 0.5*size, 0*size, isUI, -0.75*size, 0*size, isUI, 0.5*size, 0.5*size, isUI, -0.75*size, 0.5*size, isUI, -0.75*size, 0.75*size, isUI, -0.5*size, 0*size, isUI, -0.5*size, 0.5*size, isUI, -0.75*size, -1*size, isUI, -0.375*size, -1*size, isUI, 0.125*size, -0.25*size, isUI)
				l.IndexList = append(l.IndexList, 0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex, 0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex, 4+lastMaxIndex, 5+lastMaxIndex, 6+lastMaxIndex, 4+lastMaxIndex, 6+lastMaxIndex, 7+lastMaxIndex, 9+lastMaxIndex, 8+lastMaxIndex, 3+lastMaxIndex, 9+lastMaxIndex, 3+lastMaxIndex, 10+lastMaxIndex, 7+lastMaxIndex, 11+lastMaxIndex, 12+lastMaxIndex, 7+lastMaxIndex, 12+lastMaxIndex, 9+lastMaxIndex, 13+lastMaxIndex, 14+lastMaxIndex, 5+lastMaxIndex, 13+lastMaxIndex, 5+lastMaxIndex, 15+lastMaxIndex)
				lastMaxIndex += 16
				break
			}

			for j := 0; j < letterVertexLength; j++ {
				// Position list
				p := text.Position.Add(textOffset)
				l.PositionList = append(l.PositionList, p.X, p.Y, p.Z)

				// Rotation list
				r := text.Rotation
				l.RotationList = append(l.RotationList, r.X, r.Y, r.Z)

				// Color list
				c := text.Color
				l.ColorList = append(l.ColorList, c.R, c.G, c.B, c.A)
			}

			/*char := font.Symbol[fmt.Sprintf("%c", ch)]

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

			// Indices
			l.IndexList = append(l.IndexList,
				0+lastMaxIndex, 1+lastMaxIndex, 2+lastMaxIndex,
				0+lastMaxIndex, 2+lastMaxIndex, 3+lastMaxIndex,
			)

			lastMaxIndex += 4*/

			textOffset.X += (text.Size * 2)

			// New line
			if text.MaxWidth > 0 && int(textOffset.X) > text.MaxWidth {
				textOffset.X = 0
				s := text.NewLineDirection.Scale(text.Size * 2)
				textOffset = textOffset.Add(s)
				continue
			}

			// New line
			for _, newLineChar := range text.NewLineChars {
				if ch == newLineChar {
					textOffset.X = 0
					s := text.NewLineDirection.Scale(text.Size * 2)
					textOffset = textOffset.Add(s)
					break
				}
			}
		}
	}

	if len(l.TextList) > 0 {
		l.TextList = l.TextList[:0]
	}
}

/*func (l *TextLayer) Render() {
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
}*/

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
