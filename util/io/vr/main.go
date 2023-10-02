package ml_vr

import mmath_la "github.com/maldan/go-ml/math/linear_algebra"

var LeftAxis = mmath_la.Vector2[float32]{}
var RightAxis = mmath_la.Vector2[float32]{}
var LeftTrigger = false
var RightTrigger = false

var HeadTransform = mmath_la.Matrix4x4[float32]{}
var LeftControllerTransform = mmath_la.Matrix4x4[float32]{}
var RightControllerTransform = mmath_la.Matrix4x4[float32]{}

var LeftEyeTransform = mmath_la.Matrix4x4[float32]{}
var RightEyeTransform = mmath_la.Matrix4x4[float32]{}

var LeftProjectionMatrix = mmath_la.Matrix4x4[float32]{}
var RightProjectionMatrix = mmath_la.Matrix4x4[float32]{}

var positionOffset = mmath_la.Vector3[float32]{}
var rotateOffset = mmath_la.Quaternion[float32]{}.Identity()

func OffsetPosition(dir mmath_la.Vector3[float32]) {
	hh := rotateOffset.Invert().Mul(HeadTransform.GetRotation())

	head := mmath_la.Matrix4x4[float32]{}.Identity().RotateQuaternion(hh)

	dirNew := dir.ToVector4(1.0).MulMatrix(head)
	// dirNew = dirNew.MulMatrix(rotateOffset.ToMatrix4x4())

	positionOffset = positionOffset.Add(dirNew.ToVector3XYZ())
}

func OffsetRotation(dir mmath_la.Vector3[float32]) {
	rotateOffset = rotateOffset.Mul(mmath_la.Quaternion[float32]{}.FromEuler(dir))
}

func Calculate() {
	// Offset dir
	mx := mmath_la.Matrix4x4[float32]{}.Identity()
	mx = mx.RotateQuaternion(HeadTransform.GetRotation())

	offsetTransform := mmath_la.Matrix4x4[float32]{}.
		Identity().
		RotateQuaternion(rotateOffset).
		Translate(positionOffset.Invert())

	LeftEyeTransform = LeftEyeTransform.Multiply(offsetTransform)
	RightEyeTransform = RightEyeTransform.Multiply(offsetTransform)
}
