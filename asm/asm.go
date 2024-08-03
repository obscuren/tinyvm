package asm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	Reg   byte = iota // register
	Mem               // memory
	Dec               // decimal
	Stack             // stack
)

// assembler contains the necessary fields to compile a
// successful tinyvm program.
type assembler struct {
	labels     map[string]int
	setLabels  map[int]string
	pc         int
}

// Assemble takes code as input and returns the compiled binary code
// or an error if it failed.
func Assemble(code string) ([]byte, error) {
	assembler := &assembler{
		labels:    make(map[string]int),
		setLabels: make(map[int]string),
	}
	return assembler.assemble(code)
}

// assemble take code as input and assembles the instructions and returns
// an error if it failed.
func (p assembler) assemble(code string) ([]byte, error) {
	var instructions []Instruction
	for _, line := range strings.Split(code, "\n") {
		// trim comments
		if idx := strings.Index(line, comment); idx > 0 {
			line = line[:idx]
		}

		// trim all whitespace
		line = strings.TrimSpace(strings.Replace(line, "\t", " ", -1))
		if len(line) == 0 {
			continue
		}

		switch {
		case isLabel(line):
			line = strings.TrimSuffix(line, labelType)
			p.labels[line] = p.pc
		default:
			var splitStr []string
			for _, str := range strings.Split(line, " ") {
				if len(str) > 0 {
					splitStr = append(splitStr, str)
				}
			}
			splitStr[0] = strings.TrimSpace(splitStr[0])

			instrs, err := p.parseInstrs(splitStr)
			if err != nil {
				return nil, err
			}

			instructions = append(instructions, instrs...)
			// increment program count by the amount of instructions
			p.pc += len(instrs)
		}
	}
	// link the instructions
	p.link(instructions)

	// encode to binary
	writer := new(bytes.Buffer)
	for _, instr := range instructions {
		encoded, err := EncodeInstruction(instr)
		if err != nil {
			return nil, fmt.Errorf("%s: unexpected error: %v", instr.Op, err)
		}
		binary.Write(writer, binary.BigEndian, encoded)
	}

	return writer.Bytes(), nil
}

// parseInstrs attemps to parse the given args in a set of instructions
func (a assembler) parseInstrs(args []string) ([]Instruction, error) {
	var (
		instructions []Instruction
		op, cond, s  = a.parseOp(args[0])
	)
	args = args[1:]

	// If the instruction is a pseudo op code take special care.
	// Usually these instruction involve returning multiple parse
	// instructions.
	if isPseudoInstr(op) {
		switch op {
		case Push:
			if len(args) != 1 {
				return nil, opArgError(op, 1, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}
			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}

			instructions = []Instruction{
				{Op: Sub, Dst: SP, Ops1: SP, Immediate: true, Value: 1},
				{Op: Stm, Mode: DataTransfer, Dst: RegEntry(dst), Ops1: SP},
			}
		case Pop:
			if len(args) != 1 {
				return nil, opArgError(op, 1, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}
			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}

			instructions = []Instruction{
				{Op: Ldm, Mode: DataTransfer, Dst: RegEntry(dst), Ops1: SP},
				{Op: Add, Dst: SP, Ops1: SP, Immediate: true, Value: 1},
			}
		}
	} else {
		instr := Instruction{
			Cond: cond,
			Op:   op,
			S:    s,
		}
		switch op {
		case Cmp:
			if len(args) != 2 {
				return nil, opArgError(op, 2, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}
			if !isRegister(args[1]) {
				return nil, fmt.Errorf("%s: ops1 must be register: %s", op, args[0])
			}

			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}
			ops1, err := strconv.Atoi(args[1][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}
			instr.Dst = RegEntry(dst)
			instr.Ops1 = RegEntry(ops1)
		case Mov:
			if len(args) != 2 {
				return nil, opArgError(op, 2, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}

			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}

			instr.Dst = RegEntry(dst)
			if isImmediate(args[1]) {
				ops, err := strconv.Atoi(args[1][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
				}
				instr.Immediate = true
				instr.Value = uint32(ops)
			} else {
				ops, err := strconv.Atoi(args[1][1:])
				if err != nil {
					// Expect a string. TODO fix this
					a.setLabels[a.pc] = args[1]
				}
				instr.Ops1 = RegEntry(ops)
			}
		case Add, Sub, Mul, Div, Rsb, And, Xor, Orr, Lsl, Lsr:
			if len(args) != 3 {
				return nil, opArgError(op, 3, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}
			if !isRegister(args[1]) {
				return nil, fmt.Errorf("%s: ops1 must be register: %s", op, args[0])
			}

			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}
			ops1, err := strconv.Atoi(args[1][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}

			instr.Dst = RegEntry(dst)
			instr.Ops1 = RegEntry(ops1)

			ops2, err := strconv.Atoi(args[2][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
			}
			if isImmediate(args[2]) {
				instr.Immediate = true
				instr.Value = uint32(ops2)
			} else {
				instr.Ops2 = RegEntry(ops2)
			}
		case Call:
			if len(args) != 1 {
				return nil, opArgError(op, 1, len(args))
			}
			ops, err := strconv.Atoi(args[0][1:])
			if err != nil {
				// Expect a string. TODO fix this
				a.setLabels[a.pc] = args[0]
			}
			instr.Dst = RegEntry(ops)
			instr.Mode = Branching
		case Ret:
			instr.Mode = Branching
		case Ldm, Stm:
			if len(args) != 2 {
				return nil, opArgError(op, 2, len(args))
			}
			if !isRegister(args[0]) {
				return nil, fmt.Errorf("%s: dst must be register: %s", op, args[0])
			}

			dst, err := strconv.Atoi(args[0][1:])
			if err != nil {
				return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
			}

			instr.Dst = RegEntry(dst)
			if isImmediate(args[1]) {
				ops, err := strconv.Atoi(args[1][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
				}
				instr.Immediate = true
				instr.Value = uint32(ops)
			} else {
				ops, err := strconv.Atoi(args[1][1:])
				if err != nil {
					// Expect a string. TODO fix this
					a.setLabels[a.pc] = args[1]
				}
				instr.Ops1 = RegEntry(ops)
			}
			instr.Mode = DataTransfer
		}
		instructions = []Instruction{instr}
	}
	return instructions, nil
}

// parseOp parses the given op string and returns the opcode
// conditional value and the S flag.
func (a assembler) parseOp(strOp string) (Op, Cond, bool) {
	var (
		op  Op   // operation
		con Cond // condition
	)
	if len(strOp) > 4 {
		switch strOp[len(strOp)-4:] {
		case "gte":
			con = Gte
		case "lte":
			con = Lte
		}
		if con != NoCond {
			op = OpString[strOp[:len(strOp)-4]]
		}
	}
	if len(strOp) > 2 {
		switch strOp[len(strOp)-2:] {
		case "gt":
			con = Gt
		case "lt":
			con = Lt
		case "eq":
			con = Eq
		case "ne":
			con = Ne
		}
		if con != NoCond {
			op = OpString[strOp[:len(strOp)-2]]
		}
	}

	// yuck clean me up please
	var sSet bool
	if con == NoCond {
		if strOp[len(strOp)-1] == 's' {
			op = OpString[strOp[:len(strOp)-1]]
			sSet = true
		} else {
			op = OpString[strOp]
		}
	}
	return op, con, sSet
}

// link links the labels and instructions together.
func (a assembler) link(instructions []Instruction) {
	for pc, label := range a.setLabels {
		instructions[pc].Immediate = true
		instructions[pc].Value = uint32(a.labels[label])
	}
}

// opArgsError is a helper function for to report argument errors.
func opArgError(op Op, must, count int) error {
	return fmt.Errorf("[ %s ] requires %d argumenst but got %d", op, must, count)
}
