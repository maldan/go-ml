package masm

import (
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"log"
)

type Arg struct {
	Name   string
	Type   string
	Size   int
	Offset int
}

type Function struct {
	Name        string
	StackOffset int
	Args        []Arg
	X86         X86
}

func (f *Function) AddArg(name string, kind string) {
	arg := Arg{Name: name, Type: kind}
	if kind[0] == '*' {
		arg.Size = 8
	}
	offset := 0
	for i := 0; i < len(f.Args); i++ {
		offset += f.Args[i].Size
	}
	arg.Offset = offset
	f.Args = append(f.Args, arg)
	fmt.Printf("%+v\n", f.Args)
}

func (f *Function) ArgToOperand(name string) Operand {
	offset := 0
	for i := 0; i < len(f.Args); i++ {
		if f.Args[i].Name == name {
			offset = f.Args[i].Offset
			break
		}
	}

	return Operand{
		IsRegisterAsPointer: true,
		Offset:              offset + f.StackOffset,
		Register: Register{
			Id: rsp,
		},
	}
}

func (f Function) ToBytes() []byte {
	out := make([]byte, 0)
	for i := 0; i < len(f.X86.Code); i++ {
		for j := 0; j < len(f.X86.Code[i].Data); j++ {
			out = append(out, f.X86.Code[i].Data[j])
		}
	}
	return out
}

func (f Function) ToGolang() string {
	out := fmt.Sprintf("TEXT ·%v(SB),NOSPLIT,$0-24\n", f.Name)
	for i := 0; i < len(f.X86.Code); i++ {
		line := f.X86.Code[i]
		/*for j := 0; j < len(line.Data); j++ {
			out += fmt.Sprintf("BYTE $0x%X; ", line.Data[j])
		}*/
		outLine := "    "
		for {
			block := ml_slice.PullFirst(&line.Data, 8)
			/*if len(block) == 8 {
				out += fmt.Sprintf("QWORD $0x%X; ", binary.LittleEndian.Uint64(block))
				continue
			}*/
			if len(block) >= 4 {
				subBlock := ml_slice.PullFirst(&block, 4)
				outLine += fmt.Sprintf("LONG $0x%X; ", binary.LittleEndian.Uint32(subBlock))
				line.Data = ml_slice.Combine(block, line.Data)
				continue
			}
			if len(block) >= 2 {
				subBlock := ml_slice.PullFirst(&block, 2)
				outLine += fmt.Sprintf("WORD $0x%X; ", binary.LittleEndian.Uint16(subBlock))
				line.Data = ml_slice.Combine(block, line.Data)
				continue
			}

			if len(block) == 0 {
				break
			}
			for j := 0; j < len(block); j++ {
				outLine += fmt.Sprintf("BYTE $0x%X; ", block[j])
			}
		}

		out += outLine
		for j := 0; j < (48 - len(outLine)); j++ {
			out += " "
		}
		out += " // " + line.Comment
		out += "\n"
		// out += fmt.Sprintf("BYTE $0x%X; ", f.X86.Code[i])
	}
	return out
}

func MapRegion() {
	m, err := mmap.MapRegion(nil, 100, mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		log.Fatal(err)
	}
	// m acts as a writable slice of bytes that is not managed by the Go runtime.
	fmt.Println(len(m))

	// Because the region is not managed by the Go runtime, the Unmap method should
	// be called when finished with it to avoid leaking memory.
	if err := m.Unmap(); err != nil {
		log.Fatal(err)
	}

	// Output: 100
}

/*func T() uintptr {
	fn := Function{}
	fn.X86.MOV(RAX, RSP.Offset(8))
	//fn.X86.MOV(RBX, RBP.Offset(8))
	//fn.X86.ADD(RAX, RBX)
	fn.X86.MOV(RSP.Offset(16+8), RAX)
	fn.X86.RET()
	ml_console.PrintBytes(fn.X86.code, 32)

	// Allocate memory
	mmapFunc, err := mmap.MapRegion(nil, len(fn.X86.code), mmap.EXEC|mmap.RDWR, mmap.ANON, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Copy code to memory
	copy(mmapFunc, fn.X86.code)

	//type execFunc func(int64, int64) int64
	//unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	//f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	// fmt.Printf("SUK: %v\n", f(10, 20))
	return (uintptr)(unsafe.Pointer(&mmapFunc))
}
*/
