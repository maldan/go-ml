package ml_image

/*type ImageHighPrecision struct {
	Width  int
	Height int
	Pixels [][]ml_color.ColorRGBA[float32]
}

func (i *ImageHighPrecision) GetPixel(x int, y int) ml_color.ColorRGBA[float32] {
	if x < 0 {
		return ml_color.ColorRGBA[float32]{}
	}
	if y < 0 {
		return ml_color.ColorRGBA[float32]{}
	}
	if x > i.Width-1 {
		return ml_color.ColorRGBA[float32]{}
	}
	if y > i.Height-1 {
		return ml_color.ColorRGBA[float32]{}
	}
	return i.Pixels[x][y]
}

func (i *ImageHighPrecision) SetPixel(x int, y int, color ml_color.ColorRGBA[float32]) {
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

func (i *ImageHighPrecision) Clear() {
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			i.Pixels[x][y] = ml_color.ColorRGBA[float32]{
				math.SmallestNonzeroFloat32,
				math.SmallestNonzeroFloat32,
				math.SmallestNonzeroFloat32,
				math.SmallestNonzeroFloat32,
			}
		}
	}
}

func NewHighPrecision(width int, height int) ImageHighPrecision {
	img := ImageHighPrecision{}
	img.Width = width
	img.Height = height
	img.Pixels = make([][]ml_color.ColorRGBA[float32], width)
	for i := 0; i < width; i++ {
		img.Pixels[i] = make([]ml_color.ColorRGBA[float32], height)
	}

	img.Clear()

	return img
}
*/
