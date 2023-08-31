package ml_image

import (
	"bytes"
	"errors"
	"fmt"
	ml_fs "github.com/maldan/go-ml/util/io/fs"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	ml_color "github.com/maldan/go-ml/util/media/color"
	ml_process "github.com/maldan/go-ml/util/process"
	"golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type Image struct {
	Width  int
	Height int
	Pixels [][]ml_color.ColorRGBA[uint8]
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

	outImage.Pixels = make([][]ml_color.ColorRGBA[uint8], width)
	for ii := 0; ii < outImage.Width; ii++ {
		outImage.Pixels[ii] = make([]ml_color.ColorRGBA[uint8], outImage.Height)
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
	outImage.Pixels = make([][]ml_color.ColorRGBA[uint8], width)
	for ii := 0; ii < outImage.Width; ii++ {
		outImage.Pixels[ii] = make([]ml_color.ColorRGBA[uint8], height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			outImage.SetPixel(x, y, i.GetPixel(int(float64(x)/aspectX), int(float64(y)/aspectY)))
		}
	}

	return outImage
}

func (i *Image) GetPixel(x int, y int) ml_color.ColorRGBA[uint8] {
	if x < 0 {
		return ml_color.ColorRGBA[uint8]{}
	}
	if y < 0 {
		return ml_color.ColorRGBA[uint8]{}
	}
	if x > i.Width-1 {
		return ml_color.ColorRGBA[uint8]{}
	}
	if y > i.Height-1 {
		return ml_color.ColorRGBA[uint8]{}
	}
	return i.Pixels[x][y]
}

func (i *Image) SetPixel(x int, y int, color ml_color.ColorRGBA[uint8]) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x > i.Width-1 {
		return
	}
	if y > i.Height-1 {
		return
	}
	i.Pixels[x][y] = color
}

func (i *Image) ForEach(fn func(x int, y int, color ml_color.ColorRGBA[uint8])) {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			fn(x, y, i.Pixels[x][y])
		}
	}
}

func (i *Image) Map(fn func(x int, y int, color ml_color.ColorRGBA[uint8]) ml_color.ColorRGBA[uint8]) {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			i.Pixels[x][y] = fn(x, y, i.Pixels[x][y])
		}
	}
}

func (i *Image) Save(filePath string, options ImageOptions) error {
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
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, imageFile)

	return FromBytes(buf.Bytes())
}

func FromBytes(data []byte) (Image, error) {
	imageFile := bytes.NewReader(data)

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
	imageOut.Pixels = make([][]ml_color.ColorRGBA[uint8], imageOut.Width)

	/*for y := 0; y < imageOut.Height; y++ {
		imageOut.Pixels[y] = make([]ml_color.ColorRGBA[uint8], imageOut.Width)

		for x := 0; x < imageOut.Height; x++ {
			r, g, b, a := img.At(x, y).RGBA()
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

			imageOut.Pixels[y][x] = ml_color.ColorRGBA[uint8]{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	}*/
	for i := 0; i < imageOut.Width; i++ {
		imageOut.Pixels[i] = make([]ml_color.ColorRGBA[uint8], imageOut.Height)
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

			imageOut.Pixels[i][j] = ml_color.ColorRGBA[uint8]{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
		}
	}

	return imageOut, nil
}

type ImageMagickArgs struct {
	FromData []byte
	ToExt    string
	Quality  int
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
