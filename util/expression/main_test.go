package ml_expression_test

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/maldan/go-ml/util/expression"
	"testing"
	"time"
)

type Test struct {
	A int
	B int
}

func TestName(t *testing.T) {
	ee, _ := ml_expression.Parse("A == 5 && B == 4")
	ee.Bind(&Test{A: 5, B: 3})

	tt := time.Now()
	for i := 0; i < 10000; i++ {
		_ = ee.Execute()
	}
	fmt.Printf("%v\n", time.Since(tt))
	// ml_console.PrettyPrint(res)
}

func TestName2(t *testing.T) {
	expression, _ := govaluate.NewEvaluableExpression("A == 5 && B == 4")

	parameters := make(map[string]interface{}, 8)
	parameters["A"] = 5
	parameters["B"] = 3

	tt := time.Now()
	for i := 0; i < 10000; i++ {
		_, _ = expression.Evaluate(parameters)
	}
	fmt.Printf("%v\n", time.Since(tt))
}
