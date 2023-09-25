package ml_color_test

import (
	"fmt"
	ml_color "github.com/maldan/go-ml/util/media/color"
	"testing"
)

func TestA(t *testing.T) {
	c := ml_color.RGB8{}.White()
	fmt.Printf("%v\n", c.MulScalar(-10.5))
}
