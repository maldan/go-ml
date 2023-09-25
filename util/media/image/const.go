package ml_image

type ImageOptions struct {
	Format  string
	Quality int
	Mode    string
}

/*type PixelMatrix struct {
	Data   []ml_color.RGBA8
	Width  int
	Height int
}

func (m PixelMatrix) Init(width int, height int) PixelMatrix {
	m.Width = width
	m.Height = height
	m.Data = make([]ml_color.RGBA8, width*height)
	return m
}

func (m *PixelMatrix) SetPixel(x int, y int, color ml_color.RGBA8) {
	index := y*m.Width + x

	if index < 0 || index >= len(m.Data) {
		return
	}

	m.Data[index] = color
}

func (m *PixelMatrix) GetPixel(x int, y int) ml_color.RGBA8 {
	index := y*m.Width + x

	if index < 0 || index >= len(m.Data) {
		return ml_color.RGBA8{}
	}

	return m.Data[index]
}
*/
