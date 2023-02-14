package ml_file

import (
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path string
	mime string
}

func New(path string) *File {
	return &File{Path: path}
}

func NewWithMime(path string, mime string) *File {
	return &File{Path: path, mime: mime}
}

func (f *File) ReadAll() ([]byte, error) {
	return os.ReadFile(f.Path)
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
