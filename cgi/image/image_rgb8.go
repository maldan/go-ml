package mcgi_image

import (
	mcgi_color "github.com/maldan/go-ml/cgi/color"
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
