package asm

type RegEntry byte

const (
	Pc RegEntry = 0 // program counter register

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

var RegToString = map[RegEntry]string{
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
