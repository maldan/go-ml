package masm

import "fmt"

type CodeLine struct {
	Data    []byte
	Comment string
}

type X86 struct {
	Code []CodeLine
}

const (
	rax = iota
	rcx
	rdx
	rbx
	rsp
	rbp
	rsi
	rdi
)

const (
	ymm0 = iota + 20
	ymm1
	ymm2
	ymm3
	ymm4
	ymm5
	ymm6
	ymm7
	ymm8
	ymm9
	ymm10
)

type Register struct {
	Id uint8
}

type Operand struct {
	IsRegister          bool
	IsRegisterAsPointer bool
	IsConstant          bool
	Register            Register
	Value               int
	Offset              int
}

func (o Operand) SetAsPointer() Operand {
	o.IsRegister = false
	o.IsRegisterAsPointer = true
	return o
}

func (o Operand) SetOffset(offset int) Operand {
	o.Offset = offset
	o.IsRegister = false
	o.IsRegisterAsPointer = true
	return o
}

func (o Operand) ToString() string {
	if o.IsRegister {
		return o.Register.ToString()
	}
	if o.IsRegisterAsPointer {
		return fmt.Sprintf("[%v+%v]", o.Register.ToString(), o.Offset)
	}
	return "?"
}

func (r Register) Size() int {
	if r.Id == rax || r.Id == rbx || r.Id == rsp {
		return 8
	}
	return 0
}

func (r Register) ToString() string {
	switch r.Id {
	case rax:
		return "RAX"
	case rbx:
		return "RBX"
	case rcx:
		return "RCX"
	case rbp:
		return "RBP"
	case rsp:
		return "RSP"
	case ymm0:
		return "YMM0"
	case ymm1:
		return "YMM1"
	case ymm2:
		return "YMM2"
	}
	return "?"
}

var RAX = Operand{IsRegister: true, Register: Register{Id: rax}}
var RBX = Operand{IsRegister: true, Register: Register{Id: rbx}}
var RCX = Operand{IsRegister: true, Register: Register{Id: rcx}}
var RSP = Operand{IsRegister: true, Register: Register{Id: rsp}}
var RBP = Operand{IsRegister: true, Register: Register{Id: rbp}}

var YMM0 = Operand{IsRegister: true, Register: Register{Id: ymm0}}
var YMM1 = Operand{IsRegister: true, Register: Register{Id: ymm1}}
var YMM2 = Operand{IsRegister: true, Register: Register{Id: ymm2}}

func (x *X86) MOV(dst Operand, src Operand) {
	if dst.IsRegister && src.IsRegisterAsPointer {
		if dst.Register.Id == rax && src.Register.Id == rsp {
			if src.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8b, 0x04, 0x24},
					Comment: "mov rax, [rsp]",
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8b, 0x44, 0x24, byte(src.Offset)},
					Comment: fmt.Sprintf("mov rax, [rsp+%v]", src.Offset),
				})
			}
			return
		}

		if dst.Register.Id == rbx && src.Register.Id == rsp {
			if src.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8b, 0x1C, 0x24},
					Comment: "mov rbx, [rsp]",
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8b, 0x5C, 0x24, byte(src.Offset)},
					Comment: fmt.Sprintf("mov rbx, [rsp+%v]", src.Offset),
				})
			}
			return
		}

		if dst.Register.Id == rcx && src.Register.Id == rsp {
			if src.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8B, 0x0C, 0x24},
					Comment: "mov rcx, [rsp]",
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x8B, 0x4C, 0x24, byte(src.Offset)},
					Comment: fmt.Sprintf("mov rcx, [rsp+%v]", src.Offset),
				})
			}
			return
		}

		if dst.Register.Id == rax && src.Register.Id == rbp {
			x.Code = append(x.Code, CodeLine{
				Data:    []byte{0x48, 0x8B, 0x45, byte(src.Offset)},
				Comment: fmt.Sprintf("mov rax, [rbp+%v]", src.Offset),
			})
			return
		}
	}

	if src.IsRegister && dst.IsRegisterAsPointer {
		if src.Register.Id == rax && dst.Register.Id == rsp {
			if dst.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x89, 0x04, 0x24},
					Comment: fmt.Sprintf("mov [rsp], rax"),
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0x48, 0x89, 0x44, 0x24, byte(dst.Offset)},
					Comment: fmt.Sprintf("mov [rsp+%v], rax", dst.Offset),
				})
			}
			return
		}

		/*if src.Register.Id == rax && dst.Register.Id == rbp {
			// mov [rbp+n], rax
			x.Code = append(x.Code, 0x48, 0x89, 0x45, byte(dst.Offset))
			return
		}*/
	}

	panic(fmt.Sprintf("Unsupported scenario MOV %+v, %+v", dst.ToString(), src.ToString()))
}

func (x *X86) ADD(dst Operand, src Operand) {
	if dst.IsRegister && src.IsRegister {
		if dst.Register.Id == rax && src.Register.Id == rbx {
			x.Code = append(x.Code, CodeLine{
				Data:    []byte{0x48, 0x01, 0xD8},
				Comment: fmt.Sprintf("add rax, rbx"),
			})
			return
		}
	}

	panic(fmt.Sprintf("Unsupported scenario ADD %+v, %+v", dst.ToString(), src.ToString()))
}

// VMOVDQU - (Move Double Quadword Unaligned)
// Не требует строгого выравнивания данных в памяти.
// Это позволяет перемещать данные, которые могут быть не выровнены по требованиям инструкции,
// но она менее эффективна с точки зрения производительности.
// Она может использоваться в случаях, когда необходимо перемещать данные, которые не могут быть выровнены.
func (x *X86) VMOVDQU(dst Operand, src Operand) {
	if dst.IsRegister && src.IsRegisterAsPointer {
		if dst.Register.Id == ymm0 {
			if src.Register.Id == rsp {
				if src.Offset == 0 {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFE, 0x6F, 0x04, 0x24},
						Comment: fmt.Sprintf("vmovdqu ymm0, [rsp]"),
					})
				} else {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFE, 0x6F, 0x44, 0x24, byte(src.Offset)},
						Comment: fmt.Sprintf("vmovdqu ymm0, [rsp+%v]", byte(src.Offset)),
					})
				}
				return
			}
		}
	}

	panic(fmt.Sprintf("Unsupported scenario VMOVDQU %+v, %+v", dst.ToString(), src.ToString()))
}

func (x *X86) VMOVUPS(dst Operand, src Operand) {
	if dst.IsRegister && src.IsRegisterAsPointer {
		if dst.Register.Id == ymm0 {
			if src.Register.Id == rax {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x10, 0x00},
					Comment: fmt.Sprintf("vmovups ymm0, [rax]"),
				})
				return
			}
			if src.Register.Id == rbx {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x10, 0x03},
					Comment: fmt.Sprintf("vmovups ymm0, [rbx]"),
				})
				return
			}
			if src.Register.Id == rsp {
				if src.Offset == 0 {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x10, 0x04, 0x24},
						Comment: fmt.Sprintf("vmovups ymm0, [rsp]"),
					})
				} else {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x10, 0x44, 0x24, byte(src.Offset)},
						Comment: fmt.Sprintf("vmovups ymm0, [rsp+%v]", byte(src.Offset)),
					})
				}
				return
			}
		}
		if dst.Register.Id == ymm1 {
			if src.Register.Id == rax {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x10, 0x08},
					Comment: fmt.Sprintf("vmovups ymm1, [rax]"),
				})
				return
			}
			if src.Register.Id == rbx {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x10, 0x0B},
					Comment: fmt.Sprintf("vmovups ymm1, [rbx]"),
				})
				return
			}
			if src.Register.Id == rsp {
				if src.Offset == 0 {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x10, 0x0C, 0x24},
						Comment: fmt.Sprintf("vmovups ymm1, [rsp]"),
					})
				} else {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x10, 0x4C, 0x24, byte(src.Offset)},
						Comment: fmt.Sprintf("vmovups ymm1, [rsp+%v]", byte(src.Offset)),
					})
				}
				return
			}
		}
	}

	if dst.IsRegisterAsPointer && src.IsRegister {
		if src.Register.Id == ymm0 && dst.Register.Id == rsp {
			if dst.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x11, 0x04, 0x24},
					Comment: fmt.Sprintf("vmovups [rsp], ymm0"),
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x11, 0x44, 0x24, byte(dst.Offset)},
					Comment: fmt.Sprintf("vmovups [rsp+%v], ymm0", dst.Offset),
				})
			}

			return
		}
		if src.Register.Id == ymm0 && dst.Register.Id == rcx {
			if dst.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x11, 0x01},
					Comment: fmt.Sprintf("vmovups [rcx], ymm0"),
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x11, 0x41, byte(dst.Offset)},
					Comment: fmt.Sprintf("vmovups [rcx+%v], ymm0", dst.Offset),
				})
			}

			return
		}
	}

	panic(fmt.Sprintf("Unsupported scenario VMOVUPS %+v, %+v", dst.ToString(), src.ToString()))
}

func (x *X86) VMOVAPS(dst Operand, src Operand) {
	if dst.IsRegister && src.IsRegisterAsPointer {
		if dst.Register.Id == ymm0 {
			if src.Register.Id == rax {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x28, 0x00},
					Comment: fmt.Sprintf("vmovaps ymm0, [rax]"),
				})
				return
			}
			if src.Register.Id == rbx {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x28, 0x03},
					Comment: fmt.Sprintf("vmovaps ymm0, [rbx]"),
				})
				return
			}
			if src.Register.Id == rsp {
				if src.Offset == 0 {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x28, 0x04, 0x24},
						Comment: fmt.Sprintf("vmovaps ymm0, [rsp]"),
					})
				} else {
					x.Code = append(x.Code, CodeLine{
						Data:    []byte{0xC5, 0xFC, 0x28, 0x44, 0x24, byte(src.Offset)},
						Comment: fmt.Sprintf("vmovaps ymm0, [rsp+%v]", byte(src.Offset)),
					})
				}
				return
			}
		}
		if dst.Register.Id == ymm1 {
			if src.Register.Id == rax {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x28, 0x08},
					Comment: fmt.Sprintf("vmovaps ymm1, [rax]"),
				})
				return
			}
			if src.Register.Id == rbx {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x28, 0x0B},
					Comment: fmt.Sprintf("vmovaps ymm1, [rbx]"),
				})
				return
			}
		}
	}

	if dst.IsRegisterAsPointer && src.IsRegister {
		if src.Register.Id == ymm0 && dst.Register.Id == rsp {
			if dst.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x29, 0x04, 0x24},
					Comment: fmt.Sprintf("vmovaps [rsp], ymm0"),
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x29, 0x44, 0x24, byte(dst.Offset)},
					Comment: fmt.Sprintf("vmovaps [rsp+%v], ymm0", dst.Offset),
				})
			}

			return
		}
		if src.Register.Id == ymm0 && dst.Register.Id == rcx {
			if dst.Offset == 0 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x29, 0x01},
					Comment: fmt.Sprintf("vmovaps [rcx], ymm0"),
				})
			} else {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x29, 0x41, byte(dst.Offset)},
					Comment: fmt.Sprintf("vmovaps [rcx+%v], ymm0", dst.Offset),
				})
			}

			return
		}
	}

	panic(fmt.Sprintf("Unsupported scenario VMOVAPS %+v, %+v", dst.ToString(), src.ToString()))
}

func (x *X86) VADDPS(dst Operand, src1 Operand, src2 Operand) {
	if dst.IsRegister && src1.IsRegister && src2.IsRegister {
		if dst.Register.Id == ymm0 {
			if src1.Register.Id == ymm0 && src2.Register.Id == ymm1 {
				x.Code = append(x.Code, CodeLine{
					Data:    []byte{0xC5, 0xFC, 0x58, 0xC1},
					Comment: fmt.Sprintf("vaddps ymm0, ymm0, ymm1"),
				})
				return
			}
		}
	}

	panic(fmt.Sprintf("Unsupported scenario VADDPS %+v, %+v, %+v", dst.ToString(), src1.ToString(), src2.ToString()))
}

func (x *X86) VZEROUPPER() {
	x.Code = append(x.Code, CodeLine{
		Data:    []byte{0xC5, 0xF8, 0x77},
		Comment: fmt.Sprintf("vzeroupper"),
	})
}

func (x *X86) RET() {
	x.Code = append(x.Code, CodeLine{
		Data:    []byte{0xC3},
		Comment: "ret",
	})
}
