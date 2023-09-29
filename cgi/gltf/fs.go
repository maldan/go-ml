//go:build !wasm

package mcgi_gltf

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
	for i := 0; i < len(gltf.Images); i++ {
		gltf.Images[i].gltf = &gltf
		gltf.Images[i].Load()
	}
	for i := 0; i < len(gltf.Materials); i++ {
		gltf.Materials[i].gltf = &gltf
	}
	for i := 0; i < len(gltf.Skins); i++ {
		gltf.Skins[i].gltf = &gltf
	}

	return gltf, nil
}
