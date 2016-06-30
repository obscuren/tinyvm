package asm

import "strings"

const (
	labelType      = ":" // label suffix
	comment        = ";" // comment prefix
	registerPrefix = "r" // register prefix
	numberPrefix   = "#" // number prefix
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
	return strings.HasPrefix(s, numberPrefix)
}
