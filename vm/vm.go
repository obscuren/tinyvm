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

	StackSize = 1024
)

// VesionString represents the full version, including the name
// of the VM in string representation.
var VersionString = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)

// VM is the Tiny Virtual Machine data structure. It contains all
// registers and data pointers.
type VM struct {
	registers [asm.MaxRegister]uint32 // general purpose registers
	memory    []uint32                // memory

	debug bool
}

// New returns a new initialised VM.
func New(debug bool) *VM {
	vm := &VM{
		memory: make([]uint32, StackSize),
		debug:  debug,
	}
	vm.Set(asm.Reg, asm.R13, StackSize-1)
	return vm
}

// Set sets the value to the receivers location. The receiver can be either
// register or memory.
func (vm *VM) Set(typ byte, loc uint32, value uint32) {
	switch typ {
	case asm.Reg:
		vm.registers[loc] = value
	case asm.Mem:
		vm.memory[loc] = value
	}
}

// Get retrieves the value from the given storage type's location.
func (vm *VM) Get(typ byte, loc uint32) uint32 {
	switch typ {
	case asm.Reg:
		return vm.registers[byte(loc)]
	case asm.Mem:
		return vm.memory[loc]
	case asm.Dec:
		return loc
	}

	panic(fmt.Sprintf("vm.Get: invalid get type %d on %d", typ, loc))
}

func getOps2(vm *VM, instr asm.Instruction) uint32 {
	var ops2 uint32
	if instr.Immediate {
		ops2 = instr.Value
	} else {
		ops2 = vm.Get(asm.Reg, uint32(instr.Ops2))
	}
	return ops2
}

func getOps1(vm *VM, instr asm.Instruction) uint32 {
	var ops2 uint32
	if instr.Immediate {
		ops2 = instr.Value
	} else {
		ops2 = vm.Get(asm.Reg, uint32(instr.Ops1))
	}
	return ops2
}

// Exec executes the given byte code and returns the status of the
// program as well as the return value.
func (vm *VM) Exec(code []byte) error {
	var (
		callStack        []uint32               // call stack
		instrPos         = vm.registers[15] * 4 // instruction to read
		conditionalValue int32                  // condition value used by conditional instructions
	)

	// iterate over the instructions
	for int(instrPos) < len(code) {
		// loop, read and execute each op code
		pc := vm.registers[15]
		branch := pc // for branch tracking

		instr := asm.DecodeInstruction(binary.BigEndian.Uint32(code[instrPos : instrPos+4]))
		if vm.debug {
			fmt.Printf("instruction: %032b\n", instr.Raw)
			fmt.Printf("state: cv=%d\n", conditionalValue)
			fmt.Printf("cond= %s m=%v op=%s (pc=%d) dst=r%v ops1=r%d ops2=r%d I=%v S=%v value=%v\n", instr.Cond, instr.Mode, instr.Op, pc, instr.Dst, instr.Ops1, instr.Ops2, instr.Immediate, instr.S, instr.Value)
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
		case asm.Lte:
			if conditionalValue < 0 {
				skipInstr = true
			}
		case asm.Gte:
			if conditionalValue > 0 {
				skipInstr = true
			}
		}
		// reset the conditional value once it has been read
		conditionalValue = 0

		// instructions are skipped based on the conditional value
		// and the instruction condition.
		if !skipInstr {
			switch instr.Mode {
			case asm.DataProcessing:
				switch instr.Op {
				case asm.Mov:
					vm.Set(asm.Reg, uint32(instr.Dst), getOps1(vm, instr))
					pc++
				case asm.Add:
					ops2 := getOps2(vm, instr)

					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))+ops2)
					pc++
				case asm.Sub, asm.Rsb:
					ops2 := getOps2(vm, instr)

					var a, b uint32
					if instr.Op == asm.Sub {
						a, b = vm.Get(asm.Reg, uint32(instr.Ops1)), ops2
					} else {
						a, b = ops2, vm.Get(asm.Reg, uint32(instr.Ops1))
					}
					fmt.Println(a, b)

					vm.Set(asm.Reg, uint32(instr.Dst), a-b)
					pc++
				case asm.Mul:
					ops2 := getOps2(vm, instr)

					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))*ops2)
					pc++
				case asm.Div:
					ops2 := getOps2(vm, instr)

					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))/ops2)
					pc++
				case asm.And:
					ops2 := getOps2(vm, instr)
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))&ops2)

					pc++
				case asm.Xor:
					ops2 := getOps2(vm, instr)
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))^ops2)

					pc++
				case asm.Orr:
					ops2 := getOps2(vm, instr)
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))|ops2)

					pc++
				case asm.Lsl:
					ops2 := getOps2(vm, instr)
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))<<ops2)

					pc++
				case asm.Lsr:
					ops2 := getOps2(vm, instr)
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Reg, uint32(instr.Ops1))>>ops2)

					pc++
				case asm.Cmp:
					conditionalValue = int32(vm.Get(asm.Reg, uint32(instr.Dst)) - vm.Get(asm.Reg, uint32(instr.Ops1)))
					pc++
				default:
					return fmt.Errorf("invalid opcode: %d", instr.Op)
				}
			case asm.DataTransfer:
				switch instr.Op {
				case asm.Ldm:
					vm.Set(asm.Reg, uint32(instr.Dst), vm.Get(asm.Mem, getOps1(vm, instr)))

					pc++
				case asm.Stm:
					vm.Set(asm.Mem, getOps1(vm, instr), vm.Get(asm.Reg, uint32(instr.Dst)))
					pc++
				}
			case asm.Branching:
				switch instr.Op {
				case asm.Call:
					callStack = append(callStack, pc+1)
					vm.Set(asm.Reg, uint32(asm.R15), instr.Value)
				case asm.Ret:
					if len(callStack) == 0 {
						return nil
					}
					pc = callStack[len(callStack)-1]
					callStack = callStack[:len(callStack)-1]
				}
			}
			// set conditional value if S is set
			if instr.S {
				conditionalValue = int32(vm.Get(asm.Reg, uint32(instr.Dst)))
			}
		} else {
			// increment the program counter
			pc++
		}

		// Track branch. If modified don't increment
		if branch == vm.registers[15] {
			vm.registers[15] = pc
		}
		instrPos = vm.registers[15] * 4
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
				str += string(r)
			} else {
				str += "?"
			}
		}
		fmt.Println(str)
	}

	fmt.Println()
}
