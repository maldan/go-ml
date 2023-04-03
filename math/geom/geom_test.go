package mgeom_test

import (
	"fmt"
	mgeom "github.com/maldan/go-ml/math/geom"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	tt := time.Now()

	vx := mgeom.Vector3[float32]{}
	mx := mgeom.Matrix4x4[float32]{}

	for i := 0; i < 1_000_000; i++ {
		vx.TransformMatrix4x4(mx)
	}

	fmt.Printf("%v\n", time.Since(tt))
}
