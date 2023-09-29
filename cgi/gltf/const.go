package mcgi_gltf

import (
	mmath_la "github.com/maldan/go-ml/math/linear_algebra"
	ml_convert "github.com/maldan/go-ml/util/convert"
	"strings"
)

type BufferView struct {
	Buffer     int `json:"buffer"`
	ByteLength int `json:"byteLength"`
	ByteOffset int `json:"byteOffset"`
	Target     int `json:"target"`
}

type Buffer struct {
	ByteLength int    `json:"byteLength"`
	Uri        string `json:"uri"`
	content    []byte
}

func (b *Buffer) Parse() {
	t := strings.Split(b.Uri, ",")[1]
	b.content, _ = ml_convert.FromBase64(t)
}

type Accessor struct {
	BufferView    int    `json:"bufferView"`
	ComponentType int    `json:"componentType"`
	Count         int    `json:"count"`
	Type          string `json:"type"`
}

type Node struct {
	Mesh        int        `json:"mesh"`
	Name        string     `json:"name"`
	Rotation    *[]float32 `json:"rotation"`
	Scale       *[]float32 `json:"scale"`
	Translation *[]float32 `json:"translation"`
}

type Texture struct {
	Sampler int `json:"sampler"`
	Source  int `json:"source"`
}

type GLTF struct {
	Asset       any          `json:"asset"`
	Scene       int          `json:"scene"`
	Scenes      any          `json:"scenes"`
	Nodes       []Node       `json:"nodes"`
	Materials   []Material   `json:"materials"`
	Meshes      []Mesh       `json:"meshes"`
	Images      []Image      `json:"images"`
	Textures    []Texture    `json:"textures"`
	Accessors   []Accessor   `json:"accessors"`
	BufferViews []BufferView `json:"bufferViews"`
	Buffers     []Buffer     `json:"buffers"`
	Animations  []Animation  `json:"animations"`
	Skins       []Skin       `json:"skins"`
}

func (g GLTF) ParseAccessor(accessor Accessor) {
	finalView := g.BufferViews[accessor.BufferView]

	componentAmount := numberOfComponents(accessor.Type)
	byteSize := byteLength(accessor.ComponentType)
	offset := finalView.ByteOffset
	buf := g.Buffers[finalView.Buffer].content

	if accessor.Type == "VEC3" {
		bb := make([]mmath_la.Vector3[float32], 0)
		for i := 0; i < accessor.Count; i++ {
			v := mmath_la.Vector3[float32]{}.FromBytes(buf[offset:])
			bb = append(bb, v)
			offset += byteSize * componentAmount
		}
	}
}

func numberOfComponents(t string) int {
	switch t {
	case "SCALAR":
		return 1
	case "VEC2":
		return 2
	case "VEC3":
		return 3
	case "VEC4":
		return 4
	case "MAT4":
		return 16
	default:
		return 0
	}
}

func componentTypes(n int) string {
	switch n {
	case 5120:
		return "byte"
	case 5121:
		return "unsigned_byte"
	case 5122:
		return "short"
	case 5123:
		return "unsigned_short"
	case 5125:
		return "unsigned_int"
	case 5126:
		return "float"
	}
	return ""
}

func byteLength(n int) int {
	switch n {
	case 5120:
		return 1
	case 5121:
		return 1
	case 5122:
		return 2
	case 5123:
		return 2
	case 5125: // unsigned int
		return 4
	case 5126:
		return 4
	}
	return 0
}
