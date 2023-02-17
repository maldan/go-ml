package ml_json

import (
	"encoding/json"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
)

func FromFile[T any](path string) (T, error) {
	outStruct := new(T)

	// Read content
	content, err := ml_file.New(path).ReadAll()
	if err != nil {
		return *outStruct, err
	}

	// Parse
	err = json.Unmarshal(content, outStruct)
	return *outStruct, err
}

func ToFile[T any](path string, v T) error {
	// Marshal
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// Write
	err = ml_file.New(path).Write(bytes)
	return err
}

func Update[T any](path string, f func(*T)) error {
	s, err := FromFile[T](path)
	if err != nil {
		return err
	}
	f(&s)
	err = ToFile(path, &s)
	return err
}
