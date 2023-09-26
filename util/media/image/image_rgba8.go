package ml_image

import (
	mmath "github.com/maldan/go-ml/math"
	ml_color "github.com/maldan/go-ml/util/media/color"
)

type ImageRGBA8 struct {
	Width  int
	Height int
	Data   []ml_color.RGBA8
}

func (i ImageRGBA8) New(width int, height int) ImageRGBA8 {
	i.Width = width
	i.Height = height
	i.Data = make([]ml_color.RGBA8, width*height)
	return i
}

func (i *ImageRGBA8) GetPixel(x int, y int) ml_color.RGBA8 {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return ml_color.RGBA8{}
	}
	return i.Data[index]
}

func (i *ImageRGBA8) SetPixel(x int, y int, color ml_color.RGBA8) {
	index := y*i.Width + x
	if index < 0 || index >= len(i.Data) {
		return
	}
	i.Data[index] = color
}

func (i *ImageRGBA8) Clear() {
	for x := 0; x < len(i.Data); x++ {
		i.Data[x] = ml_color.RGBA8{}
	}
}

func (i ImageRGBA8) FromBytes(data []byte) (ImageRGBA8, error) {
	image, err := readImageFromBytes(data)
	if err != nil {
		return ImageRGBA8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]ml_color.RGBA8, i.Width*i.Height)

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			r, g, b, a := image.At(x, y).RGBA()

			if r > 255 {
				r = r >> 8
			}
			if g > 255 {
				g = g >> 8
			}
			if b > 255 {
				b = b >> 8
			}
			if a > 255 {
				a = a >> 8
			}

			r = mmath.Clamp(r, 0, 255)
			g = mmath.Clamp(g, 0, 255)
			b = mmath.Clamp(b, 0, 255)
			a = mmath.Clamp(a, 0, 255)

			(&i).SetPixel(x, y, ml_color.RGBA8{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	return i, nil
}

func (i ImageRGBA8) FromFile(filePath string) (ImageRGBA8, error) {
	image, err := readImageFromFile(filePath)
	if err != nil {
		return ImageRGBA8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]ml_color.RGBA8, i.Width*i.Height)

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			r, g, b, a := image.At(x, y).RGBA()

			if r > 255 {
				r = r >> 8
			}
			if g > 255 {
				g = g >> 8
			}
			if b > 255 {
				b = b >> 8
			}
			if a > 255 {
				a = a >> 8
			}

			r = mmath.Clamp(r, 0, 255)
			g = mmath.Clamp(g, 0, 255)
			b = mmath.Clamp(b, 0, 255)
			a = mmath.Clamp(a, 0, 255)

			i.SetPixel(x, y, ml_color.RGBA8{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	return i, nil
}

func (i *ImageRGBA8) ToFile(filePath string, options ImageOptions) error {
	// Prepare pixels
	pixels := make([]byte, i.Width*i.Height*4)
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			p := i.GetPixel(x, y)
			pixels = append(pixels, p.R, p.G, p.B, p.A)
		}
	}

	return writeImageToFile(filePath, i.Width, i.Height, pixels, options)
}
