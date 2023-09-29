package mcgi_gltf

type PBR struct {
	BaseColorTexture struct {
		Index int `json:"index"`
	} `json:"baseColorTexture"`
	MetallicFactor  float32 `json:"metallicFactor"`
	RoughnessFactor float32 `json:"roughnessFactor"`
}

type Material struct {
	DoubleSided          bool   `json:"doubleSided"`
	Name                 string `json:"name"`
	PbrMetallicRoughness PBR    `json:"pbrMetallicRoughness"`
	gltf                 *GLTF
}

func (m Material) GetTexture() []byte {
	texture := m.gltf.Textures[m.PbrMetallicRoughness.BaseColorTexture.Index]
	return m.gltf.Images[texture.Source].content
}
