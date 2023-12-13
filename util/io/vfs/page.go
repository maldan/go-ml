package ml_vfs

import "encoding/binary"

type PageHeader struct {
	Flags         uint8
	NextPageIndex uint32
	ContentLength uint16
}

type Page struct {
	vfs   *VFS
	Index uint32

	// flags - 1 byte
	// nextPageIndex - 4 byte
	// contentLength - 2 byte
	// content - PAGE_SIZE - 1 - 4 - 2 bytes
}

const offsetToPageContent = 1 + 4 + 2

const FLAG_USED = 1
const FLAG_COMPRESION = 2
const FLAG_ENCRYPTION = 4

func (p *Page) GetMemoryOffset() uint64 {
	return uint64(p.Index) * uint64(p.vfs.pageSize)
}

func (p *Page) GetMemoryOffsetToContent() uint64 {
	return (uint64(p.Index) * uint64(p.vfs.pageSize)) + offsetToPageContent
}

func (p *Page) GetHeader() PageHeader {
	offset := p.GetMemoryOffset()

	return PageHeader{
		Flags:         p.vfs.mmap[offset],
		NextPageIndex: binary.LittleEndian.Uint32(p.vfs.mmap[offset+1 : offset+1+4]),
		ContentLength: binary.LittleEndian.Uint16(p.vfs.mmap[offset+1+4 : offset+1+4+2]),
	}
}

func (p *Page) SetHeader(header PageHeader) {
	offset := p.GetMemoryOffset()

	p.vfs.mmap[offset] = header.Flags
	binary.LittleEndian.PutUint32(p.vfs.mmap[offset+1:offset+1+4], header.NextPageIndex)
	binary.LittleEndian.PutUint16(p.vfs.mmap[offset+1+4:offset+1+4+2], header.ContentLength)
}

func (p *Page) GetNextPage() *Page {
	h := p.GetHeader()
	if h.NextPageIndex == 0 {
		return nil
	}
	return &Page{vfs: p.vfs, Index: h.NextPageIndex}
}

func (p *Page) GetContent() []byte {
	h := p.GetHeader()
	offset := p.GetMemoryOffset()

	if h.ContentLength == 0 {
		return nil
	}

	return p.vfs.mmap[offset+offsetToPageContent : offset+offsetToPageContent+uint64(h.ContentLength)]
}

func (p *Page) WriteContent(data []byte) ([]byte, int) {
	// Copy content to page
	offset := p.GetMemoryOffset()
	pp := 0
	for i := offsetToPageContent; i < int(p.vfs.pageSize); i++ {
		p.vfs.mmap[int(offset)+i] = data[pp]
		pp += 1
		if pp >= len(data) {
			break
		}
	}

	// There is reminder
	if len(data) > pp {
		return data[pp:], pp
	}

	return nil, pp
}

func (p *Page) Clear() error {
	p.SetHeader(PageHeader{})
	return p.Save()
}

func (p *Page) Save() error {
	return p.vfs.mmap.Flush()
}
