package masm_test

import (
	"fmt"
	masm "github.com/maldan/go-ml/asm"
	ml_console "github.com/maldan/go-ml/util/io/console"
	"testing"
)

func TestName(t *testing.T) {
	fn := masm.Function{Name: "Add3", StackOffset: 8}
	fn.AddArg("x", "*float32")
	fn.AddArg("y", "*float32")
	fn.AddArg("out", "*float32")

	fn.X86.MOV(masm.RAX, fn.ArgToOperand("x"))
	fn.X86.MOV(masm.RBX, fn.ArgToOperand("y"))
	fn.X86.MOV(masm.RCX, fn.ArgToOperand("out"))

	fn.X86.VMOVUPS(masm.YMM0, masm.RAX.SetAsPointer())
	fn.X86.VMOVUPS(masm.YMM1, masm.RBX.SetAsPointer())
	fn.X86.VADDPS(masm.YMM0, masm.YMM0, masm.YMM1)
	fn.X86.VMOVUPS(masm.RCX.SetAsPointer(), masm.YMM0)
	fn.X86.VZEROUPPER()
	fn.X86.RET()

	/*fn.X86.MOV(masm.RAX, fn.ArgToOperand("x"))
	fn.X86.MOV(masm.RBX, fn.ArgToOperand("y"))
	fn.X86.ADD(masm.RAX, masm.RBX)
	fn.X86.MOV(fn.ArgToOperand("out"), masm.RAX)
	fn.X86.RET()*/

	ml_console.PrintBytes(fn.ToBytes(), 32)
	fmt.Printf("%v\n", fn.ToGolang())
}
