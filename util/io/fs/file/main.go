package ml_file

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	ml_convert "github.com/maldan/go-ml/util/convert"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
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

func (f File) ToBytes() ([]byte, error) {
	out := make([]byte, 0)

	// Write name
	out = append(out, uint8(len(f.Name())))
	out = append(out, f.Name()...)

	// Write size
	size := f.Size()
	out = append(out, uint8(size))
	out = append(out, uint8(size>>8))
	out = append(out, uint8(size>>16))
	out = append(out, uint8(size>>24))

	// Write content
	if f.isVirtual {
		out = append(out, f.content...)
	} else {
		err := f.Load()
		if err != nil {
			return nil, err
		}
		out = append(out, f.content...)
	}

	return out, nil
}

func (f File) ToDataUrl() ([]byte, error) {
	// Get content
	content, err := f.ReadAll()
	if err != nil {
		return []byte{}, err
	}

	// Prepare
	out := make([]byte, 0, len(content)+128)
	// out = append(out, '"')
	out = append(out, []byte("data:"+f.mime+";base64,")...)

	// Add
	out = append(out, base64.StdEncoding.EncodeToString(content)...)
	// out = append(out, '"')
	return out, nil
}

func (f *File) Load() error {
	content, err := f.ReadAll()
	f.content = content
	return err
}

func (f *File) ReadAll() ([]byte, error) {
	if f.isVirtual {
		return f.content, nil
	}
	content, err := os.ReadFile(f.Path)
	f.mime = http.DetectContentType(content)
	return content, err
}

func (f *File) ImageDimension() (int, int, error) {
	reader, err := os.Open(f.Path)
	defer reader.Close()
	if err != nil {
		return 0, 0, err
	}
	im, _, err := image.DecodeConfig(reader)
	return im.Width, im.Height, err
}

func (f *File) Save() error {
	return f.Write(f.content)
}

func (f *File) Write(content []byte) error {
	if f.Path == "" {
		return errors.New("path not defined")
	}

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
	if f.isVirtual {
		return int64(len(f.content))
	}
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

/*func (f *File) RelativePath() string {
	return filepath.Rel(f.Path)
}*/

func (f *File) Created() time.Time {
	s, err := os.Stat(f.Path)
	if err != nil {
		return time.Time{}
	}
	return s.ModTime()
}

func (f *File) Delete() error {
	if f.isVirtual {
		return errors.New("file is virtual")
	}

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

func (f *File) Sha256() (string, error) {
	ff, err := os.Open(f.Path)
	defer ff.Close()
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	_, err = io.Copy(hasher, ff)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
