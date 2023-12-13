package ml_vfs

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type FileTable struct {
	vfs       *VFS
	PageIndex uint32

	// File table index per page
	// Amount of files - 2 byte
	// Free pointer - 2 byte
	// list of files struct...
}

func (ft *FileTable) GetMemoryOffset() uint64 {
	return uint64(ft.PageIndex) * uint64(ft.vfs.pageSize)
}

func (ft *FileTable) FirstPage() *Page {
	return &Page{vfs: ft.vfs, Index: 0}
}

/*func (ft *FileTable) GetNextTable() FileTable {
	return &FileTable{vfs: ft.vfs, PageIndex: 0}
}*/

func (ft *FileTable) CreateFile(filePath string, fileName string) (File, error) {
	// Get first page
	page := ft.FirstPage()
	content := page.GetContent()

	// Get total number of records
	offset := 0
	// Skip total records
	offset += 2

	// Get free pointer
	freePointer := binary.LittleEndian.Uint16(content[offset : offset+2])
	offset += 2

	// Write file info
	writeOffset := ft.GetMemoryOffset()
	writeOffset += uint64(freePointer)

	// Write name
	ft.vfs.mmap[writeOffset] = uint8(len(fileName))
	writeOffset += 1
	copy(ft.vfs.mmap[writeOffset:], fileName)
	writeOffset += uint64(len(fileName))

	// Write path
	ft.vfs.mmap[writeOffset] = uint8(len(filePath))
	writeOffset += 1
	copy(ft.vfs.mmap[writeOffset:], filePath)
	writeOffset += uint64(len(filePath))

	// Write page
	binary.LittleEndian.PutUint32(ft.vfs.mmap[writeOffset:], 0)
	writeOffset += 4

	// Update free pointer
	freePointer += uint16(1 + len(fileName) + 1 + len(filePath) + 4)
	startOfContent := page.GetMemoryOffsetToContent()
	binary.LittleEndian.PutUint16(ft.vfs.mmap[startOfContent+2:startOfContent+4], freePointer)

	return File{}, nil
}

func (ft *FileTable) FindFile(filePath string, fileName string) (File, error) {
	// Get first page
	page := ft.FirstPage()
	content := page.GetContent()

	// Get total number of records
	offset := 0
	totalRecords := int(binary.LittleEndian.Uint16(content[offset : offset+2]))
	offset += 2
	// skip free pointer
	offset += 2

	for i := 0; i < totalRecords; i++ {
		// Read name
		nl := int(content[offset])
		offset += 1
		name := string(content[offset : offset+nl])
		offset += nl
		fmt.Printf("Name: %v\n", name)

		// Read path
		pl := int(content[offset])
		offset += 1
		path := string(content[offset : offset+pl])
		offset += pl
		fmt.Printf("Path: %v\n", path)

		// Read page
		pageIndex := binary.LittleEndian.Uint32(content[offset : offset+4])
		offset += 4
		fmt.Printf("PageIndex: %v\n", pageIndex)

		// File is found
		if path == filePath && name == fileName {
			return File{vfs: ft.vfs, Name: name, Path: path, PageIndex: pageIndex}, nil
		}
	}

	return File{}, errors.New("not found")
}
