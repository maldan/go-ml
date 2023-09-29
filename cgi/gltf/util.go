package mcgi_gltf

import "encoding/base64"

func fromBase64(v string) ([]byte, error) {
	uDec, err := base64.StdEncoding.DecodeString(v)
	return uDec, err
}
