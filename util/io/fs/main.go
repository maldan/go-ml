package ml_fs

import (
	ml_crypto "github.com/maldan/go-ml/util/crypto"
	ml_file "github.com/maldan/go-ml/util/io/fs/file"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*type FileInfo struct {
	RelativePath string
	FullPath     string
	Name         string
	Ext          string
	Dir          string
	IsDir        bool
}*/

type FileInfo struct {
	RelativePath string    `json:"relativePath"`
	FullPath     string    `json:"fullPath"`
	Name         string    `json:"name"`
	Ext          string    `json:"ext"`
	Dir          string    `json:"dir"`
	IsDir        bool      `json:"isDir"`
	Size         int64     `json:"size"`
	Created      time.Time `json:"created"`
}

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

func List(path string) ([]FileInfo, error) {
	// Get list of files and dirs
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	out := make([]FileInfo, 0)
	for _, f := range list {
		absPath, _ := filepath.Abs(path + "/" + f.Name())
		absPath = strings.ReplaceAll(absPath, "\\", "/")
		ext := strings.Split(f.Name(), ".")

		info, _ := f.Info()

		out = append(out, FileInfo{
			RelativePath: f.Name(),
			FullPath:     absPath,
			Name:         f.Name(),
			Ext:          ext[len(ext)-1],
			Dir:          strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
			IsDir:        f.IsDir(),
			Size:         info.Size(),
			Created:      info.ModTime(),
		})
	}
	return out, nil
}

func ListAll(path string) ([]FileInfo, error) {
	list := make([]FileInfo, 0)

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

			ext := strings.Split(info.Name(), ".")

			list = append(list, FileInfo{
				FullPath:     absPath,
				RelativePath: strings.Replace(absPath, curAbsPath, "", 1),
				Name:         info.Name(),
				Ext:          ext[len(ext)-1],
				Dir:          strings.ReplaceAll(filepath.Dir(absPath), "\\", "/"),
				IsDir:        info.IsDir(),
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

func GetTempFilePath() string {
	return os.TempDir() + "/" + ml_crypto.UID(16)
}
