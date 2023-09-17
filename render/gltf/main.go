package mrender_gltf

import (
	"encoding/json"
	"os"
)

func FromFile(path string) (GLTF, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return GLTF{}, err
	}
	gltf := GLTF{}
	err2 := json.Unmarshal(data, &gltf)
	if err2 != nil {
		return GLTF{}, err2
	}

	for i := 0; i < len(gltf.Buffers); i++ {
		gltf.Buffers[i].Parse()
	}

	for i := 0; i < len(gltf.Meshes); i++ {
		gltf.Meshes[i].gltf = &gltf
	}

	return gltf, nil
}
