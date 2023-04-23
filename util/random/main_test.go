package ml_random_test

import (
	ml_console "github.com/maldan/go-ml/util/io/console"
	ml_random "github.com/maldan/go-ml/util/random"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"testing"
)

func TestName(t *testing.T) {
	f := make([]float64, 0)
	for i := 0; i < 100; i++ {
		f = append(f, ml_random.Range[float64](-5, 10))
	}
	f = ml_slice.SortAZ(f)
	ml_console.PrettyPrint(f)
}

func TestName2(t *testing.T) {
	f := make([]int, 0)
	for i := 0; i < 100; i++ {
		f = append(f, ml_random.RangeInt(-5, 5))
	}
	f = ml_slice.SortAZ(f)
	ml_console.PrettyPrint(f)
}
