package mcgi_image_test

import (
	"fmt"
	mcgi_color "github.com/maldan/go-ml/cgi/color"
	"github.com/maldan/go-ml/cgi/image"
	"testing"
)

/*func TestX(t *testing.T) {
	img, err := ml_image.FromFile("C:/Users/black/Desktop/CtOjmNXnY0y9AJNIHksu4xjg.webp")
	if err != nil {
		panic(err)
	}

	img.Map(func(x int, y int, color ml_image.Color) ml_image.Color {
		color.R /= 2
		color.G /= 2
		color.B /= 2
		return color
	})

	img = img.ResizeWidth(890 / 2)

	err = img.Save("C:/Users/black/Desktop/ttx.jpg", ml_image.ImageOptions{
		Format:  "jpg",
		Quality: 55,
	})
	fmt.Printf("%v\n", err)
}
*/

func TestA(t *testing.T) {
	i := mcgi_image.ImageRGBA8{}.New(2, 2)
	i.SetPixel(1, 1, mcgi_color.RGBA8{}.White())
	fmt.Printf("%v\n", i.Data)
}
