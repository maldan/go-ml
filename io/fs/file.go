package ml_fs

import (
	ml_file "github.com/maldan/go-ml/io/fs/file"
	"os"
	"path/filepath"
	"strings"
)

/*type FileInfo struct {
	RelativePath string
	FullPath     string
	Name         string
	Ext          string
	Dir          string
	IsDir        bool
}*/

func Mkdir(path string) error {
	// Create path for file
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}
	return nil
}

/*func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	return data, nil
}*/

/*func WriteFile(path string, content []byte) error {
	// Create path for file
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	// Write content
	err = os.WriteFile(path, content, 0777)
	return err
}*/
/*
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}*/

func List(path string) ([]ml_file.File, error) {
	// Get list of files and dirs
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	out := make([]ml_file.File, 0)
	for _, f := range list {
		absPath, _ := filepath.Abs(path + "/" + f.Name())
		absPath = strings.ReplaceAll(absPath, "\\", "/")

		// ext := strings.Split(f.Name(), ".")

		out = append(out, ml_file.File{
			Path: absPath,
			/*RelativePath: f.Name(),
			FullPath:     absPath,
			Name:         f.Name(),
			Ext:          ext[len(ext)-1],
			Dir:          strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
			IsDir:        f.IsDir(),*/
		})
	}
	return out, nil
}
