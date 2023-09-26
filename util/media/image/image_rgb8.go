package ml_image

import (
	mmath "github.com/maldan/go-ml/math"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type ImageRGB8 struct {
	Width  int
	Height int
	Data   []ml_color.RGB8
}

func (i ImageRGB8) New(width int, height int) ImageRGB8 {
	i.Width = width
	i.Height = height
	i.Data = make([]ml_color.RGB8, width*height)
	return i
}

func (i *ImageRGB8) GetPixel(x int, y int) ml_color.RGB8 {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return ml_color.RGB8{}
	}
	return i.Data[index]
}

func (i *ImageRGB8) SetPixel(x int, y int, color ml_color.RGB8) {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return
	}
	i.Data[index] = color
}

func (i *ImageRGB8) Clear() {
	for x := 0; x < len(i.Data); x++ {
		i.Data[x] = ml_color.RGB8{}
	}
}

func (i ImageRGB8) FromBytes(data []byte) (ImageRGB8, error) {
	image, err := readImageFromBytes(data)
	if err != nil {
		return ImageRGB8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]ml_color.RGB8, i.Width*i.Height)

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			r, g, b, _ := image.At(x, y).RGBA()

			if r > 255 {
				r = r >> 8
			}
			if g > 255 {
				g = g >> 8
			}
			if b > 255 {
				b = b >> 8
			}

			r = mmath.Clamp(r, 0, 255)
			g = mmath.Clamp(g, 0, 255)
			b = mmath.Clamp(b, 0, 255)

			i.SetPixel(x, y, ml_color.RGB8{R: uint8(r), G: uint8(g), B: uint8(b)})
		}
	}

	return i, nil
}

func (i ImageRGB8) FromFile(filePath string) (ImageRGB8, error) {
	image, err := readImageFromFile(filePath)
	if err != nil {
		return ImageRGB8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]ml_color.RGB8, i.Width*i.Height)

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			r, g, b, _ := image.At(x, y).RGBA()

			if r > 255 {
				r = r >> 8
			}
			if g > 255 {
				g = g >> 8
			}
			if b > 255 {
				b = b >> 8
			}

			r = mmath.Clamp(r, 0, 255)
			g = mmath.Clamp(g, 0, 255)
			b = mmath.Clamp(b, 0, 255)

			i.SetPixel(x, y, ml_color.RGB8{R: uint8(r), G: uint8(g), B: uint8(b)})
		}
	}

	return i, nil
}

func (i *ImageRGB8) ToFile(filePath string, options ImageOptions) error {
	// Prepare pixels
	pixels := make([]byte, i.Width*i.Height*4)
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			p := i.GetPixel(x, y)
			pixels = append(pixels, p.R, p.G, p.B, 255)
		}
	}

	return writeImageToFile(filePath, i.Width, i.Height, pixels, options)
}
