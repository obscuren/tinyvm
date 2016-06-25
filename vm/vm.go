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
}

// New returns a new initialised VM.
func New() *VM {
	return &VM{
		memory: make(map[uint64]int64),
	}
}

func (vm *VM) Set(typ byte, loc byte, value int64) {
	switch typ {
	case 0:
		vm.registers[loc] = value
	case 1:
		vm.memory[uint64(loc)] = value
	}
}

func (vm *VM) Get(typ byte, loc byte) int64 {
	switch typ {
	case 0:
		return vm.registers[loc]
	case 1:
		return vm.memory[uint64(loc)]
	case 10:
		return int64(loc)
	}

	panic(fmt.Sprintf("vm.Get: invalid get type %d on %d", typ, loc))
}

// Exec executes the given byte code and returns the status of the
// program as well as the return value.
func (vm *VM) Exec(code []byte) ([]byte, error) {
	var cond bool // instruction condition
	// loop, read and execute each op code
	pc := vm.registers[0]
	for int(pc) < len(code) {
		switch op := asm.Op(code[pc]); op {
		case asm.Mov:
			typl, loc, typv, locv := code[pc+1], code[pc+2], code[pc+3], code[pc+4]
			vm.Set(typl, loc, vm.Get(typv, locv))

			pc += 5
		case asm.Ret:
			typ, loc := code[pc+1], code[pc+2]
			buffer := new(bytes.Buffer)
			binary.Write(buffer, binary.BigEndian, vm.Get(typ, loc))

			return buffer.Bytes(), nil
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
			return nil, fmt.Errorf("invalid opcode: %d", op)
		}
		vm.registers[0] = pc
	}

	return nil, nil
}

func (vm *VM) Stats() {
	fmt.Println("regs:")
	for register, value := range vm.registers {
		fmt.Println(asm.RegToString[asm.Reg(register)], ":", value)
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
}
