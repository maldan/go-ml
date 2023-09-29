package mcgi_color_test

import (
	"fmt"
	"github.com/maldan/go-ml/cgi/color"
	"testing"
)

func TestA(t *testing.T) {
	c := mcgi_color.RGB8{}.White()
	fmt.Printf("%v\n", c.MulScalar(-10.5))
}
