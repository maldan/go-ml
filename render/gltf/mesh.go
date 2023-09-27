package mrender_gltf

import (
	"encoding/binary"
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
)

type Mesh struct {
	Name       string      `json:"name"`
	Primitives []Primitive `json:"primitives"`
	gltf       *GLTF
}

type Attribute struct {
	POSITION   *int `json:"POSITION"`
	NORMAL     *int `json:"NORMAL"`
	TEXCOORD_0 *int `json:"TEXCOORD_0"`
	JOINTS_0   *int `json:"JOINTS_0"`
	WEIGHTS_0  *int `json:"WEIGHTS_0"`
}

type Target struct {
	POSITION *int `json:"POSITION"`
	NORMAL   *int `json:"NORMAL"`
}

type Primitive struct {
	Attributes Attribute `json:"attributes"`
	Targets    []Target  `json:"targets"`
	Indices    int       `json:"indices"`
	Material   *int      `json:"material"`
}

func (m Mesh) GetTargetVertices(id int) []mmath_la.Vector3[float32] {
	p := m.Primitives[0].Targets[id].POSITION
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector3[float32], 0)
	if accessor.Type == "VEC3" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector3[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetTargetNormals(id int) []mmath_la.Vector3[float32] {
	p := m.Primitives[0].Targets[id].NORMAL
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector3[float32], 0)
	if accessor.Type == "VEC3" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector3[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetJoints() []mmath_la.Vector4[float32] {
	p := m.Primitives[0].Attributes.JOINTS_0
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector4[float32], 0)
	if accessor.Type == "VEC4" {
		for i := 0; i < accessor.Count; i++ {
			out = append(out, mmath_la.Vector4[float32]{
				X: float32(buf[offset : offset+4][0]),
				Y: float32(buf[offset : offset+4][1]),
				Z: float32(buf[offset : offset+4][2]),
				W: float32(buf[offset : offset+4][3]),
			})
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetWeight() []mmath_la.Vector4[float32] {
	p := m.Primitives[0].Attributes.WEIGHTS_0
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector4[float32], 0)
	if accessor.Type == "VEC4" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector4[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetNormals() []mmath_la.Vector3[float32] {
	p := m.Primitives[0].Attributes.NORMAL
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector3[float32], 0)
	if accessor.Type == "VEC3" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector3[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetVertices() []mmath_la.Vector3[float32] {
	p := m.Primitives[0].Attributes.POSITION
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector3[float32], 0)
	if accessor.Type == "VEC3" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector3[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetUV() []mmath_la.Vector2[float32] {
	p := m.Primitives[0].Attributes.TEXCOORD_0
	if p == nil {
		return nil
	}

	accessor := m.gltf.Accessors[*p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]mmath_la.Vector2[float32], 0)
	if accessor.Type == "VEC2" {
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector2[float32]{}.FromBytes(buf[offset:])
			out = append(out, v)
			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) GetIndices() []uint32 {
	p := m.Primitives[0].Indices
	accessor := m.gltf.Accessors[p]
	finalView := m.gltf.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := m.gltf.Buffers[finalView.Buffer].content

	out := make([]uint32, 0)
	if accessor.Type == "SCALAR" {
		for i := 0; i < accessor.Count; i++ {
			if byteSize == 2 {
				num := binary.LittleEndian.Uint16(buf[offset : offset+2])
				out = append(out, uint32(num))
			}
			if byteSize == 4 {
				num := binary.LittleEndian.Uint32(buf[offset : offset+4])
				out = append(out, num)
			}

			offset += byteSize * componentAmount
		}
	}

	return out
}

func (m Mesh) HasMaterial() bool {
	return m.Primitives[0].Material != nil
}

func (m Mesh) GetMaterial() Material {
	return m.gltf.Materials[*m.Primitives[0].Material]
}
