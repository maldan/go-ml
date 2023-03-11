package ml_fs

import (
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"io"
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

func Copy(from string, to string) error {
	source, err := os.Open(from)
	defer source.Close()
	if err != nil {
		return err
	}

	// Prepare dir
	os.MkdirAll(filepath.Dir(to), 0777)

	// Create destination file
	destination, err := os.Create(to)
	defer destination.Close()
	if err != nil {
		return err
	}

	// Copy
	_, err = io.Copy(destination, source)
	return err
}

func Rename(src string, dst string) error {
	return os.Rename(src, dst)
}

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

func ListAll(path string) ([]ml_file.File, error) {
	list := make([]ml_file.File, 0)

	curAbsPath, _ := filepath.Abs(path)
	curAbsPath = strings.ReplaceAll(curAbsPath, "\\", "/")

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip dir
			if info.IsDir() {
				return nil
			}

			absPath, _ := filepath.Abs(path)
			absPath = strings.ReplaceAll(absPath, "\\", "/")

			list = append(list, ml_file.File{
				Path: absPath,
				//FullPath:     absPath,
				//RelativePath: strings.Replace(absPath, curAbsPath, "", 1),
				//Dir:          strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
				//Name:         info.Name(),
			})

			return nil
		})
	if err != nil {
		return list, err
	}

	return list, nil
}

func DeleteFile(path string) error {
	return ml_file.New(path).Delete()
}
