//go:build !wasm

package mcgi_image

import (
	"bytes"
	"errors"
	"fmt"
	mcgi_color "github.com/maldan/go-ml/cgi/color"
	mmath "github.com/maldan/go-ml/math"
	ml_fs "github.com/maldan/go-ml/util/io/fs"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_process "github.com/maldan/go-ml/util/process"
	"golang.org/x/image/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

// RGB8

func (i ImageRGB8) FromBytes(data []byte) (ImageRGB8, error) {
	image, err := readImageFromBytes(data)
	if err != nil {
		return ImageRGB8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]mcgi_color.RGB8, i.Width*i.Height)

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

			i.SetPixel(x, y, mcgi_color.RGB8{R: uint8(r), G: uint8(g), B: uint8(b)})
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
	i.Data = make([]mcgi_color.RGB8, i.Width*i.Height)

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

			i.SetPixel(x, y, mcgi_color.RGB8{R: uint8(r), G: uint8(g), B: uint8(b)})
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

// RGBA8

func (i ImageRGBA8) FromBytes(data []byte) (ImageRGBA8, error) {
	image, err := readImageFromBytes(data)
	if err != nil {
		return ImageRGBA8{}, err
	}

	i.Width = image.Bounds().Size().X
	i.Height = image.Bounds().Size().Y
	i.Data = make([]mcgi_color.RGBA8, i.Width*i.Height)

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

			(&i).SetPixel(x, y, mcgi_color.RGBA8{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
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
	i.Data = make([]mcgi_color.RGBA8, i.Width*i.Height)

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

			i.SetPixel(x, y, mcgi_color.RGBA8{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
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

type ImageMagickArgs struct {
	FromData []byte
	ToExt    string
	Quality  int
}

type ImageOptions struct {
	Format  string
	Quality int
	Mode    string
}

func ImageMagickConvertFromBytesTo(args ImageMagickArgs) ([]byte, error) {
	// Create temp file
	tempFile := ml_fs.GetTempFilePath()
	f1 := ml_file.New(tempFile)
	err := f1.Write(args.FromData)
	if err != nil {
		return nil, err
	}
	defer f1.Delete()

	outFilePath := ml_fs.GetTempFilePath() + "." + args.ToExt

	// Store photo
	ml_process.Exec("magick", tempFile, "-quality", fmt.Sprintf("%v", args.Quality), outFilePath)
	time.Sleep(time.Millisecond * 500)

	// Read
	f2 := ml_file.New(outFilePath)
	defer f2.Delete()
	dataOut, err := f2.ReadAll()
	return dataOut, err
}

func writeImageToFile(filePath string, width int, height int, pixels []byte, options ImageOptions) error {
	err := os.MkdirAll(path.Dir(filePath), 0777)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	defer f.Close()
	if err != nil {
		return err
	}

	switch options.Format {
	case "jpg", "jpeg":
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = pixels
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: options.Quality})
		return err
	case "png":
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = pixels
		err = png.Encode(f, img)
		return err
	case "gif":
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = pixels
		err = gif.Encode(f, img, &gif.Options{})
		return err
	default:
		return errors.New("unsupported format")
	}
}

func readImageFromFile(path string) (image.Image, error) {
	// Open
	imageFile, err := os.Open(path)
	defer imageFile.Close()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, imageFile)

	return readImageFromBytes(buf.Bytes())
}

func readImageFromBytes(data []byte) (image.Image, error) {
	imageFile := bytes.NewReader(data)

	// Read header
	header := make([]byte, 512)
	_, err1 := imageFile.Read(header)
	if err1 != nil {
		return nil, err1
	}
	_, err1 = imageFile.Seek(0, 0)
	if err1 != nil {
		return nil, err1
	}
	mimeType := http.DetectContentType(header)

	// Image
	var img image.Image

	// Mime
	switch mimeType {
	case "image/png":
		img2, err := png.Decode(imageFile)
		if err != nil {
			return nil, err
		}
		img = img2
		break
	case "image/jpeg":
		img2, err := jpeg.Decode(imageFile)
		if err != nil {
			return nil, err
		}
		img = img2
		break
	case "image/webp":
		img2, err := webp.Decode(imageFile)
		if err != nil {
			return nil, err
		}
		img = img2
		break
	case "image/gif":
		img2, err := gif.Decode(imageFile)
		if err != nil {
			return nil, err
		}
		img = img2
		break
	default:
		return nil, errors.New(fmt.Sprintf("unsupported '%v' file format", mimeType))
	}

	return img, nil
}
