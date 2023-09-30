package ml_vr

import mmath_la "github.com/maldan/go-ml/math/linear_algebra"

var LeftAxis = mmath_la.Vector2[float32]{}
var RightAxis = mmath_la.Vector2[float32]{}
var LeftTrigger = false
var RightTrigger = false

var HeadTransform = mmath_la.Matrix4x4[float32]{}
var LeftControllerTransform = mmath_la.Matrix4x4[float32]{}
var RightControllerTransform = mmath_la.Matrix4x4[float32]{}
