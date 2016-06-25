package vm

import (
	"strconv"
	"strings"
)

const (
	labelType = ":"
	comment   = ";"
)

func isLabel(s string) bool {
	return strings.HasSuffix(s, labelType)
}

func Parse(code string) []byte {
	var (
		parsedCode []byte
		labels     = make(map[string]int)

		toFill = make(map[int]string)
	)

	pc := 0
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

		switch {
		case isLabel(line):
			line = strings.TrimSuffix(line, labelType)
			labels[line] = pc
			parsedCode = append(parsedCode, byte(Nop))
			pc++
		default:
			var (
				splitStr = strings.Split(line, " ")
				op       = OpString[splitStr[0]]
			)
			parsedCode = append(parsedCode, byte(op))
			pc++

			switch op {
			case Jmpi, Jmpn, Jmp:
				toFill[pc+1] = splitStr[1]

				parsedCode = append(parsedCode, []byte{10, 0}...)
				pc += 2
			default:
				for _, loc := range splitStr[1:] {
					parsedCode = append(parsedCode, parseLoc(loc)...)
					pc += 2
				}
			}

		}
	}

	for pc, label := range toFill {
		parsedCode[pc] = byte(labels[label])
	}

	return parsedCode
}

func parseLoc(s string) []byte {
	switch {
	case s == "pc":
		return []byte{0, byte(Pc)}
	case strings.HasPrefix(s, "r"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{0, byte(R0 + n)}

	case strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{1, byte(n)}
	default:
		n, _ := strconv.Atoi(s)
		return []byte{10, byte(n)}
	}
}

type Op byte

const (
	Mov Op = iota

	Jmp
	Jmpi
	Jmpn

	Gt
	Gteq
	Lt
	Lteq
	Eq
	Nq

	Ret

	Add
	Sub

	Nop

	Dbg
)

var OpString = map[string]Op{
	"mov":  Mov,
	"jmp":  Jmp,
	"jmpi": Jmpi,
	"jmpn": Jmpn,
	"gt":   Gt,
	"gteq": Gteq,
	"lt":   Lt,
	"lteq": Lteq,
	"eq":   Eq,
	"nq":   Nq,
	"add":  Add,
	"sub":  Sub,
	"ret":  Ret,
	"nop":  Nop,
	"dbg":  Dbg,
}

type Reg byte

const (
	Pc Reg = 0 // program counter register

	// general purpose registers
	R0 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7
	R8
	R9
	R10
	R11
	R12
	R13
	R14
	R15
	MaxRegister
)

var RegToString = map[Reg]string{
	Pc:  "pc",
	R0:  "r0",
	R1:  "r1",
	R2:  "r2",
	R3:  "r3",
	R4:  "r4",
	R5:  "r5",
	R6:  "r6",
	R7:  "r7",
	R8:  "r8",
	R9:  "r9",
	R10: "r10",
	R11: "r11",
	R12: "r12",
	R13: "r13",
	R14: "r14",
	R15: "r15",
}
