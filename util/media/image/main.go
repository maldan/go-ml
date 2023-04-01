package ml_image

import (
	"errors"
	"fmt"
	"golang.org/x/exp/constraints"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

type ColorRGB[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
}

type ColorRGBA[T constraints.Integer | constraints.Float] struct {
	R T
	G T
	B T
	A T
}

type Image struct {
	Width  int
	Height int
	Pixels [][]ColorRGBA[uint8]
}

type ImageOptions struct {
	Format  string
	Quality int
	Mode    string
}

func (i *Image) ResizeWidth(width int) Image {
	aspectX := float64(width) / float64(i.Width)
	aspectY := (float64(i.Height) * aspectX) / float64(i.Height)

	// Prepare image
	outImage := Image{
		Width:  width,
		Height: int(float64(i.Height) * aspectX),
	}

	outImage.Pixels = make([][]ColorRGBA[uint8], width)
	for ii := 0; ii < outImage.Width; ii++ {
		outImage.Pixels[ii] = make([]ColorRGBA[uint8], outImage.Height)
	}

	for y := 0; y < outImage.Height; y++ {
		for x := 0; x < width; x++ {
			outImage.SetPixel(x, y, i.GetPixel(int(float64(x)/aspectX), int(float64(y)/aspectY)))
		}
	}

	return outImage
}

func (i *Image) Resize(width int, height int) Image {
	aspectX := float64(width) / float64(i.Width)
	aspectY := float64(height) / float64(i.Height)

	// Prepare image
	outImage := Image{
		Width:  width,
		Height: height,
	}
	outImage.Pixels = make([][]ColorRGBA[uint8], width)
	for ii := 0; ii < outImage.Width; ii++ {
		outImage.Pixels[ii] = make([]ColorRGBA[uint8], height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			outImage.SetPixel(x, y, i.GetPixel(int(float64(x)/aspectX), int(float64(y)/aspectY)))
		}
	}

	return outImage
}

func (i *Image) GetPixel(x int, y int) ColorRGBA[uint8] {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return i.Pixels[x][y]
}

func (i *Image) SetPixel(x int, y int, color ColorRGBA[uint8]) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	i.Pixels[x][y] = color
}

func (i *Image) ForEach(fn func(x int, y int, color ColorRGBA[uint8])) {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			fn(x, y, i.Pixels[x][y])
		}
	}
}

func (i *Image) Map(fn func(x int, y int, color ColorRGBA[uint8]) ColorRGBA[uint8]) {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			i.Pixels[x][y] = fn(x, y, i.Pixels[x][y])
		}
	}
}

func (i *Image) Save(path string, options ImageOptions) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	switch options.Format {
	case "jpg", "jpeg":
		img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
		img.Pix = make([]byte, 0, i.Width*i.Height*4)
		for y := 0; y < i.Height; y++ {
			for x := 0; x < i.Width; x++ {
				color := i.Pixels[x][y]
				img.Pix = append(img.Pix, color.R, color.G, color.B, color.A)
			}
		}
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: options.Quality})
		return err
	case "png":
		img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
		img.Pix = make([]byte, 0, i.Width*i.Height*4)
		for y := 0; y < i.Height; y++ {
			for x := 0; x < i.Width; x++ {
				color := i.Pixels[x][y]
				img.Pix = append(img.Pix, color.R, color.G, color.B, color.A)
			}
		}
		err = png.Encode(f, img)
		return err
	default:
		return errors.New("unsupported format")
	}
}

func FromFile(path string) (Image, error) {
	// Open
	imageFile, err := os.Open(path)
	defer imageFile.Close()
	if err != nil {
		return Image{}, err
	}

	// Read header
	header := make([]byte, 512)
	imageFile.Read(header)
	imageFile.Seek(0, 0)
	mimeType := http.DetectContentType(header)

	// Image
	var img image.Image

	// Mime
	switch mimeType {
	case "image/png":
		img2, err := png.Decode(imageFile)
		if err != nil {
			return Image{}, err
		}
		img = img2
		break
	case "image/jpeg":
		img2, err := jpeg.Decode(imageFile)
		if err != nil {
			return Image{}, err
		}
		img = img2
		break
	case "image/webp":
		img2, err := webp.Decode(imageFile)
		if err != nil {
			return Image{}, err
		}
		img = img2
		break
	default:
		return Image{}, errors.New(fmt.Sprintf("unsupported '%v' file format", mimeType))
	}

	// Read image
	imageOut := Image{}
	imageOut.Width = img.Bounds().Size().X
	imageOut.Height = img.Bounds().Size().Y
	imageOut.Pixels = make([][]ColorRGBA[uint8], imageOut.Width)

	for i := 0; i < imageOut.Width; i++ {
		imageOut.Pixels[i] = make([]ColorRGBA[uint8], imageOut.Height)
		for j := 0; j < imageOut.Height; j++ {
			r, g, b, a := img.At(i, j).RGBA()
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

			imageOut.Pixels[i][j] = ColorRGBA[uint8]{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	}

	return imageOut, nil
}
