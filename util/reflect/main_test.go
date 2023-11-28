package ml_reflect_test

import (
	"fmt"
	ml_reflect "github.com/maldan/go-ml/util/reflect"
	"testing"
)

func TestA(t *testing.T) {
	type X struct {
		Name string `json:"name"`
	}

	x := X{}
	ml_reflect.SetStructField(&x, "Name", "gas")
	fmt.Printf("%v\n", x)
}
