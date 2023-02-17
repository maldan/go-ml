package ml_file

import (
	"encoding/base64"
	"errors"
	ml_convert "github.com/maldan/go-ml/util/convert"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path      string
	isVirtual bool
	mime      string
	content   []byte
}

func New(path string) *File {
	return &File{Path: path, mime: "application/octet-stream"}
}

func NewWithMime(path string, mime string) *File {
	return &File{Path: path, mime: mime}
}

func (f *File) UnmarshalJSON(b []byte) error {
	removeQuotes := b[1 : len(b)-1]
	data, mime, err := ml_convert.DataUrlToBytes(string(removeQuotes))
	if err != nil {
		return err
	}

	f.isVirtual = true
	f.mime = mime
	f.content = data

	return nil
}

func (f File) MarshalJSON() ([]byte, error) {
	mime := f.mime
	if mime == "" {
		mime = "application/octet-stream"
	}

	// Get content
	content, err := f.ReadAll()
	if err != nil {
		return []byte{}, err
	}

	// Prepare
	out := make([]byte, 0, len(content)+128)
	out = append(out, '"')
	out = append(out, []byte("data:"+mime+";base64,")...)

	// Add
	out = append(out, base64.StdEncoding.EncodeToString(content)...)
	out = append(out, '"')

	return out, nil
}

func (f *File) ReadAll() ([]byte, error) {
	if f.isVirtual {
		return f.content, nil
	}
	content, err := os.ReadFile(f.Path)
	f.mime = http.DetectContentType(content)
	return content, err
}

func (f *File) Save() error {
	return f.Write(f.content)
}

func (f *File) Write(content []byte) error {
	// Create path for file
	err := os.MkdirAll(filepath.Dir(f.Path), 0777)
	if err != nil {
		return err
	}

	// Write content
	err = os.WriteFile(f.Path, content, 0777)
	return err
}

func (f *File) Exists() bool {
	_, err := os.Stat(f.Path)
	if err == nil {
		return true
	}
	return false
}

func (f *File) Size() int64 {
	s, err := os.Stat(f.Path)
	if err != nil {
		return 0
	}
	return s.Size()
}

func (f *File) Mime() string {
	return f.mime
}

func (f *File) Name() string {
	return filepath.Base(f.Path)
}

func (f *File) Created() time.Time {
	s, err := os.Stat(f.Path)
	if err != nil {
		return time.Time{}
	}
	return s.ModTime()
}

func (f *File) Delete() error {
	stat, err := os.Stat(f.Path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		err = os.Remove(f.Path)
		if err != nil {
			return err
		}
	} else {
		return errors.New("path is directory")
	}

	return nil
}
