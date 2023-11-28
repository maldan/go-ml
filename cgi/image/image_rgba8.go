package mcgi_image

import (
	mcgi_color "github.com/maldan/go-ml/cgi/color"
)

type ImageRGBA8 struct {
	Width  int
	Height int
	Data   []mcgi_color.RGBA8
}

func (i ImageRGBA8) New(width int, height int) ImageRGBA8 {
	i.Width = width
	i.Height = height
	i.Data = make([]mcgi_color.RGBA8, width*height)
	return i
}

func (i *ImageRGBA8) GetPixel(x int, y int) mcgi_color.RGBA8 {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return mcgi_color.RGBA8{}
	}
	return i.Data[index]
}

func (i *ImageRGBA8) SetPixel(x int, y int, color mcgi_color.RGBA8) {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return
	}
	i.Data[index] = color
}

func (i *ImageRGBA8) Clear() {
	for x := 0; x < len(i.Data); x++ {
		i.Data[x] = mcgi_color.RGBA8{}
	}
}

func (i ImageRGBA8) GetResolution() (int, int) {
	return i.Width, i.Height
}

func (i ImageRGBA8) Clone() ImageRGBA8 {
	d := make([]mcgi_color.RGBA8, 0)

	for x := 0; x < len(i.Data); x++ {
		d = append(d, i.Data[x])
	}

	return ImageRGBA8{
		Width:  i.Width,
		Height: i.Height,
		Data:   d,
	}
}
