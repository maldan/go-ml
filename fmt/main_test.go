package mfmt_test

import (
	"fmt"
	mcgi_color "github.com/maldan/go-ml/cgi/color"
	mcgi_image "github.com/maldan/go-ml/cgi/image"
	mfmt "github.com/maldan/go-ml/fmt"
	"testing"
)

func TestX(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", 3))
}

func TestY(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", mcgi_color.RGBA8{R: 32}))
}

func TestZ(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", mcgi_image.ImageRGBA8{}.New(5, 5)))
}
