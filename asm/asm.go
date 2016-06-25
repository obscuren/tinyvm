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
