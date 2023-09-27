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

func TestQuaternion(t *testing.T) {
	q1 := mmath_la.Quaternion[float32]{}.FromEuler(mmath_la.Vec3f(45, 0, 0).ToRad())
	//fmt.Printf("%v\n", q1.ToEuler())
	//fmt.Printf("%v\n", q1.ToEuler().ToDeg())

	q2 := mmath_la.Quaternion[float32]{}.FromEuler(mmath_la.Vec3f(90, 0, 0).ToRad())
	fmt.Printf("%v\n", q1.Lerp(q2, 0).ToEuler().ToDeg())

	//let q1 = Quaternion::from_euler(Vector3::new(45.0f32, 0.0, 0.0).to_radians());
	//let q2 = Quaternion::from_euler(Vector3::new(90.0f32, 0.0, 0.0).to_radians());
}
