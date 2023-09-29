package mcgi_rasterizer

import (
	"github.com/maldan/go-ml/cgi/color"
	mmath_geom "github.com/maldan/go-ml/math/geometry"
)

type screen struct {
	Buffer  []uint8
	ZBuffer []float32
	Width   int
	Height  int
}

func (s *screen) Init(width int, height int) {
	s.Width = width
	s.Height = height
	s.Buffer = make([]uint8, width*height*4)
	s.ZBuffer = make([]float32, width*height)
}

func (s *screen) IsVisible(triangle mmath_geom.Triangle3D[float32]) bool {
	if triangle.A.X <= 0 && triangle.B.X <= 0 && triangle.C.X <= 0 {
		return false
	}
	if triangle.A.Y <= 0 && triangle.B.Y <= 0 && triangle.C.Y <= 0 {
		return false
	}
	if triangle.A.X > float32(s.Width) && triangle.B.X > float32(s.Width) && triangle.C.X > float32(s.Width) {
		return false
	}
	if triangle.A.Y > float32(s.Height) && triangle.B.Y > float32(s.Height) && triangle.C.Y > float32(s.Height) {
		return false
	}
	return true
}

func (s *screen) SetPixel(x int, y int, color ml_color.RGBA8) {
	index := (y*s.Width + x) * 4

	if index < 0 || index >= len(s.ZBuffer) {
		return
	}

	s.Buffer[index] = color.R
	s.Buffer[index+1] = color.G
	s.Buffer[index+2] = color.B
}

func (s *screen) Clear() {
	for i := 0; i < len(s.ZBuffer); i++ {
		s.ZBuffer[i] = 0
	}
	for i := 0; i < len(s.Buffer); i++ {
		s.Buffer[i] = 0
	}
}

func (s *screen) SetZPixel(x int, y int, v float32) {
	index := y*s.Width + x

	if index < 0 || index >= len(s.ZBuffer) {
		return
	}

	s.ZBuffer[index] = v
}

func (s *screen) GetZPixel(x int, y int) float32 {
	index := y*s.Width + x

	if index < 0 || index >= len(s.ZBuffer) {
		return 0
	}

	return s.ZBuffer[index]
}
