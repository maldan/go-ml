package ml_geom

import "golang.org/x/exp/constraints"

type Vector4[T constraints.Float] struct {
	X T
	Y T
	Z T
	W T
}

/*func (v *Vector4[T]) TransformMatrix4x4(mx Matrix4x4[T]) {
	v.X = (mx.Raw[0]*v.X + mx.Raw[4]*v.Y + mx.Raw[8]*v.Z + mx.Raw[12]) / w
	v.Y = (mx.Raw[1]*v.X + mx.Raw[5]*v.Y + mx.Raw[9]*v.Z + mx.Raw[13]) / w
	v.Z = (mx.Raw[2]*v.X + mx.Raw[6]*v.Y + mx.Raw[10]*v.Z + mx.Raw[14]) / w
}*/

func (v *Vector4[T]) Clone() Vector4[T] {
	return Vector4[T]{v.X, v.Y, v.Z, v.W}
}
