package ml_geom_test

import (
	"fmt"
	ml_geom "github.com/maldan/go-ml/util/math/geom"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	tt := time.Now()

	vx := ml_geom.Vector3[float32]{}
	mx := ml_geom.Matrix4x4[float32]{}

	for i := 0; i < 1_000_000; i++ {
		vx.TransformMatrix4x4(mx)
	}

	fmt.Printf("%v\n", time.Since(tt))
}
