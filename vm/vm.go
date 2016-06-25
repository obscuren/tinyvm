// Copyright 2016 Jeffrey Wilcke
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode"

	"github.com/obscuren/tinyvm/asm"
)

const (
	Major = 0 // Major version
	Minor = 0 // Minor version
	Patch = 1 // Patch version
)

// VesionString represents the full version, including the name
// of the VM in string representation.
var VersionString = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

// VM is the Tiny Virtual Machine data structure. It contains all
// registers and data pointers.
type VM struct {
	registers [asm.MaxRegister]int64 // general purpose registers
	memory    map[uint64]int64       // memory addressed by pointer
	stack     []int64
}

// New returns a new initialised VM.
func New() *VM {
	return &VM{
		memory: make(map[uint64]int64),
	}
}

// Set sets the value to the receivers location. The receiver can be either
// register or memory.
func (vm *VM) Set(typ byte, loc byte, value int64) {
	switch typ {
	case asm.Reg:
		vm.registers[loc] = value
	case asm.Mem:
		vm.memory[uint64(loc)] = value
	case asm.Stack:
		vm.stack = append(vm.stack, value)
	}
}

// Get retrieves the value from the given storage type's location.
func (vm *VM) Get(typ byte, loc byte) int64 {
	switch typ {
	case asm.Reg:
		return vm.registers[loc]
	case asm.Mem:
		return vm.memory[uint64(loc)]
	case asm.Dec:
		return int64(loc)
	case asm.Stack:
		stackItem := vm.stack[len(vm.stack)-1]
		vm.stack = vm.stack[:len(vm.stack)-1]
		return stackItem
	}

	panic(fmt.Sprintf("vm.Get: invalid get type %d on %d", typ, loc))
}

// Exec executes the given byte code and returns the status of the
// program as well as the return value.
func (vm *VM) Exec(code []byte) error {
	var (
		cond      bool    // instruction condition
		callStack []int64 // call stack
	)
	// loop, read and execute each op code
	pc := vm.registers[0]
	for int(pc) < len(code) {
		switch op := asm.Op(code[pc]); op {
		case asm.Stop:
			return nil
		case asm.Mov:
			typl, loc, typv, locv := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			vm.Set(typl, loc, vm.Get(typv, locv))

			pc += 5
		case asm.Push:
			typv, locv := code[pc+1], code[pc+2]
			vm.Set(asm.Stack, 0, vm.Get(typv, locv))

			pc += 3
		case asm.Pop:
			vm.Get(asm.Stack, 0) // pop one item of stack and ignore it
			pc++
		case asm.Call:
			callStack = append(callStack, pc+3)
			pc = vm.Get(code[pc+1], code[pc+2])
		case asm.Ret:
			if len(callStack) == 0 {
				return nil
			}

			pc = callStack[len(callStack)-1]
			callStack = callStack[:len(callStack)-1]
		case asm.Add:
			typt, t := code[pc+1], code[pc+2]
			typa, a := code[pc+3], code[pc+4]
			typb, b := code[pc+5], code[pc+6]

			vm.Set(typt, t, vm.Get(typa, a)+vm.Get(typb, b))

			pc += 7
		case asm.Sub:
			typt, t := code[pc+1], code[pc+2]
			typa, a := code[pc+3], code[pc+4]
			typb, b := code[pc+5], code[pc+6]

			vm.Set(typt, t, vm.Get(typa, a)-vm.Get(typb, b))

			pc += 7
		case asm.Jmpi:
			typp, p := code[pc+1], code[pc+2]
			if cond {
				pc = vm.Get(typp, p)
			} else {
				pc += 3
			}
			cond = false // set the condition back to false after reading
		case asm.Jmpn:
			typp, p := code[pc+1], code[pc+2]
			if !cond {
				pc = vm.Get(typp, p)
			} else {
				pc += 3
			}
			cond = false // set the condition back to false after reading
		case asm.Jmp:
			pc = vm.Get(code[pc+1], code[pc+2])
		case asm.Lt:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) < vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Gt:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) > vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Lteq:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) <= vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Gteq:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) >= vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Eq:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) == vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Nq:
			typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			if vm.Get(typa, a) != vm.Get(typb, b) {
				cond = true
			}

			pc += 5
		case asm.Nop:
			pc++

		case asm.Dbg:
			value := vm.Get(code[pc+1], code[pc+2])

			fmt.Println("dbg:", value)

			pc += 3
		default:
			return fmt.Errorf("invalid opcode: %d", op)
		}
		vm.registers[0] = pc
	}

	return nil
}

// Stats prints the virtual machine internal statistics.
func (vm *VM) Stats() {
	fmt.Println("regs:")
	for register, value := range vm.registers {
		fmt.Println(asm.RegToString[asm.RegEntry(register)], ":", value)
	}

	fmt.Println()

	fmt.Println("mem:")
	for addr, value := range vm.memory {
		buff := new(bytes.Buffer)
		binary.Write(buff, binary.BigEndian, value)
		fmt.Printf("%04d: % x  ", addr, buff.Bytes())

		var str string
		for _, r := range buff.Bytes() {
			if r == 0 {
				str += "."
			} else if unicode.IsPrint(rune(r)) {
				str += fmt.Sprintf("%s", string(r))
			} else {
				str += "?"
			}
		}
		fmt.Println(str)
	}

	fmt.Println()
	fmt.Println("stack:")
	fmt.Printf("%x\n", vm.stack)

}
