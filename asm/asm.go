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
		line = strings.Replace(line, "\t", " ", -1)

		switch {
		case isLabel(line):
			line = strings.TrimSuffix(line, labelType)
			labels[line] = pc
			parsedCode = append(parsedCode, byte(Nop))
			pc++
		default:
			var splitStr []string
			for _, str := range strings.Split(line, " ") {
				if len(str) > 0 {
					splitStr = append(splitStr, str)
				}
			}
			op := OpString[strings.TrimSpace(splitStr[0])]

			parsedCode = append(parsedCode, byte(op))
			pc++

			switch op {
			case Jmpi, Jmpn:
				parsedCode = append(parsedCode, parseLoc(splitStr[1])...)
				parsedCode = append(parsedCode, []byte{Dec, 0}...)

				toFill[pc+3] = strings.TrimSpace(splitStr[2])

				pc += 4
			case Jmp, Call:
				toFill[pc+1] = strings.TrimSpace(splitStr[1])

				parsedCode = append(parsedCode, []byte{Dec, 0}...)
				pc += 2
			default:
				for _, loc := range splitStr[1:] {
					parsedCode = append(parsedCode, parseLoc(strings.TrimSpace(loc))...)
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

const (
	Reg byte = iota
	Mem
	Dec
	Stack
)

func parseLoc(s string) []byte {
	switch {
	case s == "pc":
		return []byte{Reg, byte(Pc)}
	case s == "pop":
		return []byte{Stack, 0}
	case strings.HasPrefix(s, "r"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{Reg, byte(R0 + n)}

	case strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"):
		n, _ := strconv.Atoi(s[1:])
		return []byte{Mem, byte(n)}
	default:
		n, _ := strconv.Atoi(s)
		return []byte{Dec, byte(n)}
	}
}
