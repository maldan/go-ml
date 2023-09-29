package mcgi_gltf

type Image struct {
	BufferView int    `json:"bufferView"`
	MimeType   string `json:"mimeType"`
	Name       string `json:"name"`
	content    []byte
	gltf       *GLTF
}

func (i *Image) Load() {
	view := i.gltf.BufferViews[i.BufferView]
	i.content = i.gltf.Buffers[view.Buffer].content[view.ByteOffset : view.ByteOffset+view.ByteLength]
}
