package mfmt_test

import (
	"fmt"
	mcgi_color "github.com/maldan/go-ml/cgi/color"
	mcgi_image "github.com/maldan/go-ml/cgi/image"
	mfmt "github.com/maldan/go-ml/fmt"
	mmath "github.com/maldan/go-ml/math"
	"testing"
	"time"
)

func TestX(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("INT %v", 3))
	fmt.Printf("%v\n", mfmt.Sprintf("LEN %v", len([]int{1, 2, 3})))
	fmt.Printf("%v\n", mfmt.Sprintf("F32 %v", 0.13453))
	fmt.Printf("%v\n", mfmt.Sprintf("F32 %v", 1.13453))
	fmt.Printf("%v\n", mfmt.Sprintf("F32 %v", -0.13453))
	fmt.Printf("%v\n", mfmt.Sprintf("F32 %v", -1.13453))
	fmt.Printf("%v\n", mfmt.Sprintf("F32 %v", mmath.Pi))
}

func TestY(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", mcgi_color.RGBA8{R: 32}))
}

func TestZ(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", mcgi_image.ImageRGBA8{}.New(5, 5)))
}

func TestXX(t *testing.T) {
	fmt.Printf("%v\n", mfmt.Sprintf("XAX %v", nil))
}

func TestXX2(t *testing.T) {
	tt := time.Now()
	x := 0
	for i := 0; i < 1_000_000; i++ {
		x += 1
	}
	fmt.Printf("%v\n", time.Since(tt))
}
