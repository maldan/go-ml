package expression_test

import (
	"github.com/maldan/go-ml/util/expression"
	ml_console "github.com/maldan/go-ml/util/io/console"
	"testing"
)

type Test struct {
	A int
	B int
}

func TestName(t *testing.T) {
	ee, _ := expression.Parse("A == 5 && B == 4")
	ee.Bind(&Test{A: 5, B: 3})
	res := ee.Execute()
	ml_console.PrettyPrint(res)
}
