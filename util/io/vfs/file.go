package ml_vfs

import "time"

type File struct {
	vfs     *VFS
	Name    string
	Path    string
	Page    uint32
	Created time.Time
}

func (f File) Read(buf []byte) (int, error) {
	return 0, nil
}

func (f File) Write(buf []byte) {

}

func (f File) Delete() {

}
