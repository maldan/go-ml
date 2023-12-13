package ml_vfs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/edsrzf/mmap-go"
	ml_number "github.com/maldan/go-ml/util/number"
	"os"
)

type IVFS interface {
	Open() int
	Read([]byte) (int, error)
	Write([]byte, int) (int error)
	Close(fd int)
}

type VFS struct {
	file      *os.File
	mmap      mmap.MMap
	pageSize  uint16
	filePath  string
	fileTable FileTable
}

func (f *VFS) Open() error {
	f.pageSize = 4096
	f.fileTable = FileTable{vfs: f}

	// Open or create file
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	f.file = file

	// Map file
	m, _ := mmap.Map(file, mmap.RDWR, 0)
	/*if err != nil {
		return err
	}*/
	f.mmap = m

	return nil
}

func (f *VFS) reopen() error {
	err := f.close()
	if err != nil {
		return err
	}

	// Open or create file
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	f.file = file

	// Map file
	m, _ := mmap.Map(file, mmap.RDWR, 0)
	f.mmap = m
	return nil
}

func (f *VFS) close() error {
	if f.mmap != nil {
		err := f.mmap.Unmap()
		f.mmap = nil
		if err != nil {
			return err
		}
	}
	if f.file != nil {
		err := f.file.Close()
		f.file = nil
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *VFS) numOfPages() int {
	return len(f.mmap) / int(f.pageSize)
}

func (f *VFS) getPage(index uint32) (*Page, error) {
	if int(index) > f.numOfPages() {
		return nil, errors.New(fmt.Sprintf("getPage: page %v out of range", index))
	}

	// Page
	page := Page{
		vfs:   f,
		Index: index,
	}

	return &page, nil
}

func (f *VFS) AllocatePage() (*Page, error) {
	// Get current file size
	info, err := f.file.Stat()
	if err != nil {
		return &Page{}, err
	}

	// Change file size
	index := int(info.Size()) / int(f.pageSize)
	err = f.file.Truncate(info.Size() + int64(f.pageSize))
	if err != nil {
		return &Page{}, err
	}

	// Reopen
	err = f.reopen()
	if err != nil {
		return &Page{}, err
	}

	// Page
	page := Page{
		vfs:   f,
		Index: uint32(index),
	}

	return &page, nil
}

func (f *VFS) getFreePage() (*Page, error) {
	total := f.numOfPages()
	for i := 0; i < total; i++ {
		p := Page{vfs: f, Index: uint32(i)}
		h := p.GetHeader()

		// Not used
		if !ml_number.CheckBitMask(h.Flags, FLAG_USED) {
			return &p, nil
		}
	}

	// Otherwise allocate new
	return f.AllocatePage()
}

func (f *VFS) GetFile(fullPath string) (File, error) {
	// Get first page
	page, err := f.getPage(0)
	if err != nil {
		return File{}, err
	}

	content := page.GetContent()

	// Get total number of records
	offset := 0
	totalRecords := int(binary.LittleEndian.Uint16(content[offset : offset+2]))
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
		if fullPath == path+"/"+name {
			return File{vfs: f, Name: name, Path: path, PageIndex: pageIndex}, nil
		}
	}

	return File{}, errors.New("not found")
}

func (f *VFS) CreateFile(file File) (File, error) {
	// Find is already exists
	file, err3 := f.GetFile(file.Name + "/" + file.Path)
	if err3 == nil {
		// It exists
		return File{}, errors.New("file already exists")
	}

	// Get free page
	page, err2 := f.getFreePage()
	if err2 != nil {
		return File{}, err2
	}

	// Set first page for file
	file.PageIndex = page.Index

	// Clear first page
	err := page.Clear()
	if err != nil {
		return File{}, err
	}

	// We allocate first page for file
	h := page.GetHeader()
	h.Flags |= FLAG_USED
	page.SetHeader(h)
	err = page.Save()
	if err != nil {
		return File{}, err
	}

	// Add file to file list

	return file, nil
}

/*
	Описание виртуальной файловой системы. Вся система построена по принципу страниц.
	Есть страница с информацией о файлах, она начинается всегда с 0 страницы.
	Все остальное страницы файлов. Типов страниц нет, так как это не имеет значения.

	У страниц есть флаги.
	- Занята ли страница
	- Используется ли сжатие
	- Используется ли шифрование

	У файла есть имя и дата создания, а так же путь где он лежит. Папок как таковых не существует в системе.
	Их нельзя создать или удалить. При создании файла указывается его путь. Так же у файла есть указание
	на какой странице он находится.

	В таблице индексов каждая страница будет содержать количество файлов в данной странице,
	и все файловые структуры будут неразрывны. То есть если информация о файле не влазит в страницу
	то переносить ее на следующую страницу. Это нужно для удобства парсинга.
*/
