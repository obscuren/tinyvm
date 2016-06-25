package asm

type Op byte

const (
	Mov Op = iota
	Push
	Pop

	Jmp
	Jmpi
	Jmpn

	Gt
	Gteq
	Lt
	Lteq
	Eq
	Nq

	Call
	Ret

	Add
	Sub

	Dbg

	Stop
	Nop
)

var OpString = map[string]Op{
	"mov":  Mov,
	"push": Push,
	"pop":  Pop,
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
	"call": Call,
	"ret":  Ret,
	"nop":  Nop,
	"stop": Stop,
	"dbg":  Dbg,
}
