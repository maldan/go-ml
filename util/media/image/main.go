package ml_image

/*type Image struct {
	Width   int
	Height  int
	Pixels  [][]ml_color.ColorRGBA[uint8]
	Pixels2 []ml_color.ColorRGBA[uint8]
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
}*/
/*
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
*/
/*
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

func (i *Image) Clear() {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			i.Pixels[x][y] = ml_color.ColorRGBA[uint8]{0, 0, 0, 255}
		}
	}
}*/

/*func (i *Image) Save(filePath string, options ImageOptions) error {
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
}*/

/*func New(width int, height int) Image {
	img := Image{}
	img.Width = width
	img.Height = height
	img.Pixels = make([][]ml_color.ColorRGBA[uint8], width)
	for i := 0; i < width; i++ {
		img.Pixels[i] = make([]ml_color.ColorRGBA[uint8], height)
	}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Pixels[i][j].A = 255
		}
	}
	return img
}

func New2(width int, height int) Image {
	img := Image{}
	img.Width = width
	img.Height = height
	img.Pixels = make([][]ml_color.ColorRGBA[uint8], height)
	for i := 0; i < height; i++ {
		img.Pixels[i] = make([]ml_color.ColorRGBA[uint8], width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Pixels[y][x].A = 255
		}
	}
	return img
}

func New3(width int, height int) Image {
	img := Image{}
	img.Width = width
	img.Height = height
	img.Pixels2 = make([]ml_color.ColorRGBA[uint8], width*height)
	return img
}
*/
