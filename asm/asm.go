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

func isLabel(s string) bool {
	return strings.HasSuffix(s, labelType)
}

func isRegister(s string) bool {
	return strings.HasPrefix(s, registerPrefix)
}

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
	writer := new(bytes.Buffer)

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
			p.pc++
		default:
			var splitStr []string
			for _, str := range strings.Split(line, " ") {
				if len(str) > 0 {
					splitStr = append(splitStr, str)
				}
			}
			p.pc++

			var (
				op    = OpString[strings.TrimSpace(splitStr[0])]
				instr Instruction
			)

			switch op {
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

				instr = Instruction{
					Op:  op,
					Dst: RegEntry(dst),
				}
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
						return nil, fmt.Errorf("%s: unexepected error: %v", op, err)
					}
					instr.Ops1 = RegEntry(ops)
				}
				encoded, err := EncodeInstruction(instr)
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}
				binary.Write(writer, binary.BigEndian, encoded)
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

				instr = Instruction{
					Op:   op,
					Dst:  RegEntry(dst),
					Ops1: RegEntry(ops1),
				}

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
				encoded, err := EncodeInstruction(instr)
				if err != nil {
					return nil, fmt.Errorf("%s: unexpected error: %v", op, err)
				}
				binary.Write(writer, binary.BigEndian, encoded)
			}
			/*
				p.parsedCode = append(p.parsedCode, byte(op))

				switch op {
				case Jmpt, Jmpf:
					p.parsedCode = append(p.parsedCode, p.parseLoc(p.pc, splitStr[1])...)
					p.parsedCode = append(p.parsedCode, []byte{Dec, 0}...)

					p.toFill[p.pc+3] = strings.TrimSpace(splitStr[2])

					p.pc += 4
				case Call:
					p.toFill[p.pc+1] = strings.TrimSpace(splitStr[1])

					p.parsedCode = append(p.parsedCode, []byte{Dec, 0}...)
					p.pc += 2
				default:
					for _, loc := range splitStr[1:] {
						p.parsedCode = append(p.parsedCode, p.parseLoc(p.pc, strings.TrimSpace(loc))...)
						p.pc += 2
					}
				}
			*/
		}
	}

	/*
		for pc, label := range p.toFill {
			p.parsedCode[pc] = byte(p.labels[label])
		}

		return p.parsedCode, nil
	*/
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
