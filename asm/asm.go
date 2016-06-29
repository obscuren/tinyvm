package asm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	labelType       = ":"
	comment         = ";"
	registerPrefix  = "r"
	immediatePrefix = "#"
)

// isLabel returns whether s is of type label
func isLabel(s string) bool {
	return strings.HasSuffix(s, labelType)
}

// isRegister returns whether s is of type register
func isRegister(s string) bool {
	return strings.HasPrefix(s, registerPrefix)
}

// isImmediate returns whether s is of type immediate
func isImmediate(s string) bool {
	return strings.HasPrefix(s, immediatePrefix)
}

type parser struct {
	parsedCode []byte
	labels     map[string]int
	toFill     map[int]string
	pc         int
}

func Parse(code string) ([]byte, error) {
	parser := &parser{
		labels: make(map[string]int),
		toFill: make(map[int]string),
	}
	return parser.parse(code)
}

func opArgError(op Op, must, count int) error {
	return fmt.Errorf("%s requires %d argumenst but got %d", op, must, count)
}

func (p parser) parse(code string) ([]byte, error) {
	var instructions []Instruction
	for _, line := range strings.Split(code, "\n") {
		// trim comments
		if idx := strings.Index(line, comment); idx > 0 {
			line = line[:idx]
		}

		// trim all whitespace
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		line = strings.Replace(line, "\t", " ", -1)

		switch {
		case isLabel(line):
			line = strings.TrimSuffix(line, labelType)
			p.labels[line] = p.pc
			p.parsedCode = append(p.parsedCode, byte(Nop))
			// decrement program counter as a measure of "ignore"
			// this "instruction".
			p.pc--
		default:
			var splitStr []string
			for _, str := range strings.Split(line, " ") {
				if len(str) > 0 {
					splitStr = append(splitStr, str)
				}
			}

			var (
				op  Op   // op code
				con Cond // condition
			)

			splitStr[0] = strings.TrimSpace(splitStr[0])
			if len(splitStr[0]) > 4 {
				switch splitStr[0][len(splitStr[0])-4:] {
				case "gteq":
					con = Gteq
				case "lteq":
					con = Lteq
				}
				if con != NoCond {
					op = OpString[splitStr[0][:len(splitStr[0])-4]]
				}
			}
			if len(splitStr[0]) > 2 {
				switch splitStr[0][len(splitStr[0])-2:] {
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
					op = OpString[splitStr[0][:len(splitStr[0])-2]]
				}
			}

			// yuck clean me up please
			var sSet bool
			if con == NoCond {
				if splitStr[0][len(splitStr[0])-1] == 's' {
					op = OpString[splitStr[0][:len(splitStr[0])-1]]
					sSet = true
				} else {
					op = OpString[splitStr[0]]
				}
			}

			var instr Instruction
			instr = Instruction{
				Cond: con,
				Op:   op,
				S:    sSet,
			}

			switch op {
			case Cmp:
				if len(splitStr) != 3 {
					return nil, opArgError(op, 2, len(splitStr))
				}
				if !isRegister(splitStr[1]) {
					return nil, fmt.Errorf("%s: dst must be register: %s", op, splitStr[1])
				}
				if !isRegister(splitStr[2]) {
					return nil, fmt.Errorf("%s: ops1 must be register: %s", op, splitStr[1])
				}

				dst, err := strconv.Atoi(splitStr[1][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}
				ops1, err := strconv.Atoi(splitStr[2][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}
				instr.Dst = RegEntry(dst)
				instr.Ops1 = RegEntry(ops1)
			case Mov:
				if len(splitStr) != 3 {
					return nil, opArgError(op, 2, len(splitStr))
				}
				if !isRegister(splitStr[1]) {
					return nil, fmt.Errorf("%s: dst must be register: %s", op, splitStr[1])
				}

				dst, err := strconv.Atoi(splitStr[1][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}

				instr.Dst = RegEntry(dst)
				if isImmediate(splitStr[2]) {
					ops, err := strconv.Atoi(splitStr[2][1:])
					if err != nil {
						return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
					}
					instr.Immediate = true
					instr.Value = uint32(ops)
				} else {
					ops, err := strconv.Atoi(splitStr[2][1:])
					if err != nil {
						// Expect a string. TODO fix this
						p.toFill[p.pc] = splitStr[2]
					}
					instr.Ops1 = RegEntry(ops)
				}
			case Add, Sub:
				if len(splitStr) != 4 {
					return nil, opArgError(op, 3, len(splitStr))
				}
				if !isRegister(splitStr[1]) {
					return nil, fmt.Errorf("%s: dst must be register: %s", op, splitStr[1])
				}
				if !isRegister(splitStr[2]) {
					return nil, fmt.Errorf("%s: ops1 must be register: %s", op, splitStr[1])
				}

				dst, err := strconv.Atoi(splitStr[1][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}
				ops1, err := strconv.Atoi(splitStr[2][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}

				instr.Dst = RegEntry(dst)
				instr.Ops1 = RegEntry(ops1)

				ops2, err := strconv.Atoi(splitStr[3][1:])
				if err != nil {
					return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
				}
				if isImmediate(splitStr[3]) {
					instr.Immediate = true
					instr.Value = uint32(ops2)
				} else {
					instr.Ops2 = RegEntry(ops2)
				}
			case Call:
				if len(splitStr) != 2 {
					return nil, opArgError(op, 1, len(splitStr))
				}
				ops, err := strconv.Atoi(splitStr[1][1:])
				if err != nil {
					// Expect a string. TODO fix this
					p.toFill[p.pc] = splitStr[1]
				}
				instr.Dst = RegEntry(ops)
			}
			instructions = append(instructions, instr)
		}
		p.pc++
	}
	for pc, label := range p.toFill {
		instructions[pc].Immediate = true
		instructions[pc].Value = uint32(p.labels[label])
	}

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

const (
	Reg byte = iota
	Mem
	Dec
	Stack
)

func (p parser) parseLoc(pos int, s string) []byte {
	switch {
	case s == "pop":
		return []byte{Stack, 0}

	case strings.HasPrefix(s, "r"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{Reg, byte(R0 + n)}

	case strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{Mem, byte(n)}

	case strings.HasPrefix(s, "#"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{Dec, byte(n)}

	default:
		p.toFill[pos+1] = s
		return []byte{Dec, 0}
	}
}
