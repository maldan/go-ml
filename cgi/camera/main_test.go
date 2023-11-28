package mcgi_camera_test

import (
	"fmt"
	mcgi_camera "github.com/maldan/go-ml/cgi/camera"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	"testing"
)

func TestX(t *testing.T) {
	c := mcgi_camera.PerspectiveCamera{
		Fov:                45,
		AspectRatio:        1,
		Near:               0.01,
		Far:                32,
		Position:           mmath_la.Vec3f(0.0, 0.7, 3.5),
		Scale:              mmath_la.Vector3[float32]{1, 1, 1},
		IsInversePositionX: true,
		IsInversePositionY: true,
		IsInversePositionZ: true,
	}
	c.Calculate()
	fmt.Printf("%v\n", c.ViewMatrix)
}
