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
	C string
}

type ETest struct {
	Expression string
	Bind       Test
	MustBe     any
}

func TestName(t *testing.T) {
	ee, _ := ml_expression.Parse("A == 5 && B == 4")
	ee.Bind(&Test{A: 5, B: 3})

	tt := time.Now()
	for i := 0; i < 100000; i++ {
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

	pp := govaluate.MapParameters(parameters)

	tt := time.Now()
	for i := 0; i < 100000; i++ {
		_, _ = expression.Evaluate(pp)
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestCompare(t *testing.T) {
	eList := []ETest{
		/*{"A == 5", Test{A: 5}, true},
		{"A == 4", Test{A: 5}, false},
		{"5 + 5", Test{A: 5}, 10},
		{"2 + 5", Test{A: 5}, 7},
		{"A + 1", Test{A: 5}, 6},
		{"C == 'hohol'", Test{C: "hohol"}, true},
		{"'ou' in C", Test{C: "hou houuu 123"}, true},*/
		{"A > 5", Test{A: 6}, true},
	}

	for i := 0; i < len(eList); i++ {
		ee, _ := ml_expression.Parse(eList[i].Expression)
		ee.Bind(&(eList[i].Bind))
		r := ee.Execute()
		if r != eList[i].MustBe {
			t.Fatalf("fuck %v - must be %v, got %v", eList[i].Expression, eList[i].MustBe, r)
		}
	}
}
