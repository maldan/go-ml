package ml_vfs

import "time"

type File struct {
	vfs       *VFS
	Name      string
	Path      string
	PageIndex uint32
	Created   time.Time

	// name len - 1 byte
	// name - ? byte
	// path len - 1 byte
	// path - ? byte
	// page - 4 byte
}

func (f File) Read(buf []byte) (int, error) {
	return 0, nil
}

func (f File) Write(buf []byte) {

}

func (f File) Delete() {

}
