// Copyright 2016 Jeffrey Wilcke
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package asm

type Mode byte

const (
	DataProcessing Mode = iota
	DataTransfer
	Branching
)

type Op byte

const (
	// Data processing op codes
	Mov Op = iota // move data to register
	Add           // addition operator (ops1 + ops2)
	Sub           // subtraction operator (ops1 - ops2)
	Rsb           // reverse subtraction operator (ops2 - ops1)
	Mul           // multiplicaton operator (ops1 x ops2)
	Div           // division operator (ops1 / ops2)
	And           // logical and operator (ops1 AND ops2)
	Xor           // exclusive or operator (ops1 XOR ops2)
	Orr           // or operator (ops1 OR ops2)
	Lsl           // logical shift left (ops1 << ops2)
	Lsr           // logical shift right (ops1 >> ops2)
	Cmp

	// Data transfer op codes
	Ldm Op = iota // Load memory
	Stm           // Store memory

	// Branching op codes
	Call Op = iota
	Ret

	// Pseudo opcodes
	Push Op = 64 + iota
	Pop
)

var OpString = map[string]Op{
	"mov": Mov,
	"add": Add,
	"sub": Sub,
	"rsb": Rsb,
	"mul": Mul,
	"div": Div,
	"and": And,
	"xor": Xor,
	"orr": Orr,
	"lsl": Lsl,
	"lsr": Lsr,
	"cmp": Cmp,

	"ldm": Ldm,
	"stm": Stm,

	"call": Call,
	"ret":  Ret,

	// pseudo codes
	"push": Push,
	"pop":  Pop,
}

var PseudoOpcodes = map[Op]bool{Push: true, Pop: true}

func (o Op) String() string {
	return OpToString[o]
}

var OpToString = map[Op]string{
	Mov: "mov",
	Add: "add",
	Sub: "sub",
	Rsb: "rsb",
	Mul: "mul",
	Div: "div",
	And: "and",
	Xor: "xor",
	Orr: "orr",
	Lsl: "lsl",
	Lsr: "lsr",
	Cmp: "cmp",

	Ldm: "ldm",
	Stm: "stm",

	Call: "call",
	Ret:  "ret",

	Push: "push",
	Pop:  "pop",
}

type Cond byte

const (
	NoCond = iota
	Eq     // equal
	Ne     // not equal
	Gt     // greater than
	Lt     // less than
	Gte    // greater than or equal
	Lte    // less than or equal
)

func (c Cond) String() string {
	return CondToString[c]
}

var StringToCond = map[string]Cond{
	"eq":  Eq,
	"ne":  Ne,
	"gt":  Gt,
	"lt":  Lt,
	"gte": Gte,
	"lte": Lte,
}

var CondToString = map[Cond]string{
	NoCond: "nocond",
	Eq:     "eq",
	Ne:     "ne",
	Gt:     "gt",
	Lt:     "lt",
	Gte:    "gte",
	Lte:    "lte",
}
