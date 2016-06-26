package asm

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

type parser struct {
	parsedCode []byte
	labels     map[string]int
	toFill     map[int]string
	pc         int
}

func Parse(code string) []byte {
	parser := &parser{
		labels: make(map[string]int),
		toFill: make(map[int]string),
	}
	return parser.parse(code)
}

func (p parser) parse(code string) []byte {
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
			op := OpString[strings.TrimSpace(splitStr[0])]

			p.parsedCode = append(p.parsedCode, byte(op))
			p.pc++

			switch op {
			case Jmpt, Jmpf:
				p.parsedCode = append(p.parsedCode, p.parseLoc(p.pc, splitStr[1])...)
				p.parsedCode = append(p.parsedCode, []byte{Dec, 0}...)

				p.toFill[p.pc+3] = strings.TrimSpace(splitStr[2])

				p.pc += 4
			case Jmp, Call:
				p.toFill[p.pc+1] = strings.TrimSpace(splitStr[1])

				p.parsedCode = append(p.parsedCode, []byte{Dec, 0}...)
				p.pc += 2
			default:
				for _, loc := range splitStr[1:] {
					p.parsedCode = append(p.parsedCode, p.parseLoc(p.pc, strings.TrimSpace(loc))...)
					p.pc += 2
				}
			}

		}
	}

	for pc, label := range p.toFill {
		p.parsedCode[pc] = byte(p.labels[label])
	}

	return p.parsedCode
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
