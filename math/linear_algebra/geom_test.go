package mmath_la_test

import (
	"fmt"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"

	"testing"
	"time"
)

func TestName(t *testing.T) {
	tt := time.Now()

	vx := mmath_la.Vector3[float32]{}
	mx := mmath_la.Matrix4x4[float32]{}

	for i := 0; i < 1_000_000; i++ {
		vx.TransformMatrix4x4(mx)
	}

	fmt.Printf("%v\n", time.Since(tt))
}
