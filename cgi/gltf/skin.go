package mcgi_gltf

import mmath_la "github.com/maldan/go-ml/math/linear_algebra"

type Skin struct {
	InverseBindMatrices int    `json:"inverseBindMatrices"`
	Joints              []int  `json:"joints"`
	Name                string `json:"name"`
	gltf                *GLTF
}

type SkinBone struct {
	Id       int
	Name     string
	Position mmath_la.Vector3[float32]
	Rotation mmath_la.Quaternion[float32]
	Scale    mmath_la.Vector3[float32]
}

func (s Skin) GetInverseBindMatrixList() []mmath_la.Matrix4x4[float32] {
	accessor := s.gltf.Accessors[s.InverseBindMatrices]
	finalView := s.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := s.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Matrix4x4[float32], 0)
	if accessor.Type == "MAT4" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Matrix4x4[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (s Skin) GetBones() []SkinBone {
	out := make([]SkinBone, 0)

	for i := 0; i < len(s.Joints); i++ {
		joint := s.Joints[i]
		node := s.gltf.Nodes[joint]
		bone := SkinBone{
			Id:   joint,
			Name: node.Name,
		}
		if node.Translation != nil {
			bone.Position = mmath_la.Vector3[float32]{
				X: (*node.Translation)[0],
				Y: (*node.Translation)[1],
				Z: (*node.Translation)[2],
			}
		}
		if node.Rotation != nil {
			bone.Rotation = mmath_la.Quaternion[float32]{
				X: (*node.Rotation)[0],
				Y: (*node.Rotation)[1],
				Z: (*node.Rotation)[2],
				W: (*node.Rotation)[3],
			}
		}
		if node.Scale != nil {
			bone.Scale = mmath_la.Vector3[float32]{
				X: (*node.Scale)[0],
				Y: (*node.Scale)[1],
				Z: (*node.Scale)[2],
			}
		}
		out = append(out, bone)
	}

	return out
}
