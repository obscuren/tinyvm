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
	registers [asm.MaxRegister]uint32 // general purpose registers
	memory    map[uint32]uint32       // memory addressed by pointer
	stack     []uint32

	debug bool
}

// New returns a new initialised VM.
func New(debug bool) *VM {
	return &VM{
		memory: make(map[uint32]uint32),
		debug:  debug,
	}
}

// Set sets the value to the receivers location. The receiver can be either
// register or memory.
func (vm *VM) Set(typ byte, loc byte, value uint32) {
	switch typ {
	case asm.Reg:
		vm.registers[loc] = value
	case asm.Mem:
		vm.memory[uint32(loc)] = value
	case asm.Stack:
		vm.stack = append(vm.stack, value)
	}
}

// Get retrieves the value from the given storage type's location.
func (vm *VM) Get(typ byte, loc byte) uint32 {
	switch typ {
	case asm.Reg:
		return vm.registers[loc]
	case asm.Mem:
		return vm.memory[uint32(loc)]
	case asm.Dec:
		return uint32(loc)
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
		//callStack []uint32 // call stack
		codePos          = vm.registers[15] * 4
		conditionalValue int32
	)

	for int(codePos) < len(code) {
		// loop, read and execute each op code
		pc := vm.registers[15]
		branch := pc // for branch tracking

		instr := asm.DecodeInstruction(binary.BigEndian.Uint32(code[codePos : codePos+4]))
		if vm.debug {
			fmt.Printf("instruction: %032b\n", instr.Raw)
			fmt.Printf("state: cv=%d\n", conditionalValue)
			fmt.Printf("cond= %s op=%s (pc=%d) dst=r%v ops1=r%d ops2=r%d I=%v S=%v value=%v\n", instr.Cond, instr.Op, pc, instr.Dst, instr.Ops1, instr.Ops2, instr.Immediate, instr.S, instr.Value)
		}

		// boolean determining whether we should skip the instruction
		// based on the instructions conditional value.
		var skipInstr bool
		switch instr.Cond {
		case asm.Eq:
			if conditionalValue != 0 {
				skipInstr = true
			}
		case asm.Ne:
			if conditionalValue == 0 {
				skipInstr = true
			}
		case asm.Lt:
			if conditionalValue >= 0 {
				skipInstr = true
			}
		case asm.Gt:
			if conditionalValue <= 0 {
				skipInstr = true
			}
		case asm.Lteq:
			if conditionalValue < 0 {
				skipInstr = true
			}
		case asm.Gteq:
			if conditionalValue > 0 {
				skipInstr = true
			}
		}
		conditionalValue = 0

		if !skipInstr {
			switch instr.Op {
			case asm.Mov:
				if instr.Immediate {
					vm.Set(asm.Reg, byte(instr.Dst), instr.Value)
				} else {
					vm.Set(asm.Reg, byte(instr.Dst), vm.Get(asm.Reg, byte(instr.Ops1)))
				}
				pc++
			case asm.Add:
				var ops2 uint32
				if instr.Immediate {
					ops2 = instr.Value
				} else {
					ops2 = vm.Get(asm.Reg, byte(instr.Ops2))
				}
				vm.Set(asm.Reg, byte(instr.Dst), vm.Get(asm.Reg, byte(instr.Ops1))+ops2)
				pc++
			case asm.Sub:
				var ops2 uint32
				if instr.Immediate {
					ops2 = instr.Value
				} else {
					ops2 = vm.Get(asm.Reg, byte(instr.Ops2))
				}
				vm.Set(asm.Reg, byte(instr.Dst), vm.Get(asm.Reg, byte(instr.Ops1))-ops2)
				pc++
			case asm.Cmp:
				conditionalValue = int32(vm.Get(asm.Reg, byte(instr.Dst)) - vm.Get(asm.Reg, byte(instr.Ops1)))
				pc++
			}
			// set conditional value if S is set
			if instr.S {
				conditionalValue = int32(vm.Get(asm.Reg, byte(instr.Dst)))
			}
		} else {
			pc++
		}

		/*
			op := asm.Op(code[pc])
			switch op {
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
			case asm.Jmpt:
				typc, c, typp, p := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
				if vm.Get(typc, c) > 0 {
					pc = vm.Get(typp, p)
				} else {
					pc += 5
				}
			case asm.Jmpf:
				typc, c, typp, p := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
				if vm.Get(typc, c) <= 0 {
					pc = vm.Get(typp, p)
				} else {
					pc += 5
				}
			case asm.Lt:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) < vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Gt:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) > vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Lteq:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) <= vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Gteq:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) >= vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Eq:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) == vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Nq:
				typr, r, typa, a, typb, b := code[pc+1], code[pc+2], code[pc+3], code[pc+4], code[pc+5], code[pc+6]
				var v int32
				if vm.Get(typa, a) != vm.Get(typb, b) {
					v = 1
				}
				vm.Set(typr, r, v)

				pc += 7
			case asm.Nop:
				pc++

			case asm.Dbg:
				value := vm.Get(code[pc+1], code[pc+2])

				fmt.Println("dbg:", value)

				pc += 3
			default:
				return fmt.Errorf("invalid opcode: %d", op)
			}
		*/
		// Track branch. If modified don't increment
		if branch == vm.registers[15] {
			vm.registers[15] = pc
		}
		codePos = vm.registers[15] * 4
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
