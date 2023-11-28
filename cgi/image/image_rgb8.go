package mcgi_image

import (
	mcgi_color "github.com/maldan/go-ml/cgi/color"
	"math"
)

type ImageRGB8 struct {
	Width  int
	Height int
	Data   []mcgi_color.RGB8
}

func (i ImageRGB8) New(width int, height int) ImageRGB8 {
	i.Width = width
	i.Height = height
	i.Data = make([]mcgi_color.RGB8, width*height)
	return i
}

func (i *ImageRGB8) GetPixel(x int, y int) mcgi_color.RGB8 {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return mcgi_color.RGB8{}
	}
	return i.Data[index]
}

func (i *ImageRGB8) SetPixel(x int, y int, color mcgi_color.RGB8) {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return
	}
	i.Data[index] = color
}

func (i *ImageRGB8) Clear() {
	for x := 0; x < len(i.Data); x++ {
		i.Data[x] = mcgi_color.RGB8{}
	}
}

func (i ImageRGB8) GetResolution() (int, int) {
	return i.Width, i.Height
}

func (i ImageRGB8) Clone() ImageRGB8 {
	d := make([]mcgi_color.RGB8, 0)

	for x := 0; x < len(i.Data); x++ {
		d = append(d, i.Data[x])
	}

	return ImageRGB8{
		Width:  i.Width,
		Height: i.Height,
		Data:   d,
	}
}

func (i *ImageRGB8) EdgeDetection() *ImageRGB8 {
	width := i.Width
	height := i.Height
	lineArt := i.Clone()

	img := i

	// Примените фильтр Собеля для обнаружения границ
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			px00 := img.GetPixel(x-1, y-1)
			px01 := img.GetPixel(x, y-1)
			px02 := img.GetPixel(x+1, y-1)
			px10 := img.GetPixel(x-1, y)
			// px11 := img.GetPixel(x, y)
			px12 := img.GetPixel(x+1, y)
			px20 := img.GetPixel(x-1, y+1)
			px21 := img.GetPixel(x, y+1)
			px22 := img.GetPixel(x+1, y+1)

			gx := -int(px00.R) - 2*int(px01.R) - int(px02.R) +
				int(px20.R) + 2*int(px21.R) + int(px22.R)
			gy := -int(px00.R) - 2*int(px10.R) - int(px20.R) +
				int(px02.R) + 2*int(px12.R) + int(px22.R)
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))

			// Примените пороговый фильтр, чтобы определить границы
			if magnitude > 255 {
				lineArt.SetPixel(x, y, mcgi_color.RGB8{}.White())
			} else {
				lineArt.SetPixel(x, y, mcgi_color.RGB8{}.Black())
			}
		}
	}

	return &lineArt
}

func (i *ImageRGB8) Multiply(color mcgi_color.RGB8) *ImageRGB8 {
	for x := 0; x < len(i.Data); x++ {
		i.Data[x] = i.Data[x].Mul(color)
	}

	return i
}
