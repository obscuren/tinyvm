package asm

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
