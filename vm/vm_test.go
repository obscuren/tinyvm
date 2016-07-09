package vm

import (
	"testing"

	"github.com/obscuren/tinyvm/asm"
)

func TestExecution(t *testing.T) {
	for i, test := range []struct {
		code   string
		result uint32
	}{
		{"mov r0 #10", 10},
		{"add r0 r0 #1", 1},
		{"mov r0 #2\nsub r0 r0 #1", 1},
		{"mov r0 #1\nrsb r0 r0 #2", 1},
		{"mov r0 #2\nmul r0 r0 #2", 4},
		{"mov r0 #2\ndiv r0 r0 #2", 1},
		{"mov r0 #2\ndiv r0 r0 #2", 1},
		{"mov r0 #1\nand r0 r0 #2", 0},
		{"mov r0 #2\nxor r0 r0 #1", 3},
		{"mov r0 #1\norr r0 r0 #2", 3},
		{"mov r0 #1\nlsl r0 r0 #1", 2},
		{"mov r0 #2\nlsr r0 r0 #1", 1},
	} {
		code, err := asm.Assemble(test.code)
		if err != nil {
			t.Errorf("%d failed: %v", i, err)
			continue
		}
		vm := New(false)
		err = vm.Exec(code)
		if err != nil {
			t.Errorf("%d failed: %v", i, err)
			continue
		}
		if r0 := vm.Get(asm.Reg, asm.R0); r0 != test.result {
			t.Errorf("%d failed: expected %d got %d", i, test.result, r0)
		}
	}
}

func TestStack(t *testing.T) {
	for i, test := range []struct {
		code string
		r0   uint32
		r1   uint32
		sp   uint32
	}{
		{"mov r0 #1\npush r0\npush r0", 1, 0, StackSize - 3},    // double push r0
		{"mov r0 #1\npush r0\npop r1", 1, 1, StackSize - 1},     // push r0 pop in to r1
		{"mov r0 #1\npush r0\nldm r1 r13", 1, 1, StackSize - 2}, // push r0 manual store sp pos in r1
	} {
		code, err := asm.Assemble(test.code)
		if err != nil {
			t.Errorf("%d failed: %v", i, err)
			continue
		}
		vm := New(false)
		err = vm.Exec(code)
		if err != nil {
			t.Errorf("%d failed: %v", i, err)
			continue
		}
		if r0 := vm.Get(asm.Reg, asm.R0); r0 != test.r0 {
			t.Errorf("%d failed: expected %d got %d", i, test.r0, r0)
		}
		if r1 := vm.Get(asm.Reg, asm.R1); r1 != test.r1 {
			t.Errorf("%d failed: expected %d got %d", i, test.r1, r1)
		}
		if sp := vm.Get(asm.Reg, asm.SP); sp != test.sp {
			t.Errorf("%d failed: expected %d got %d", i, test.sp, sp)
		}
	}
}
