package ml_number_test

import (
	ml_console "github.com/maldan/go-ml/util/io/console"
	ml_number "github.com/maldan/go-ml/util/number"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"testing"
)

func TestName(t *testing.T) {
	f := make([]float64, 0)
	for i := 0; i < 100; i++ {
		f = append(f, ml_number.RandFloat(-1, 1))
	}
	f = ml_slice.SortAZ(f)
	ml_console.PrettyPrint(f)
}

func TestName2(t *testing.T) {
	f := make([]int, 0)
	for i := 0; i < 100; i++ {
		f = append(f, ml_number.RandInt(0, 5))
	}
	f = ml_slice.SortAZ(f)
	ml_console.PrettyPrint(f)
}
