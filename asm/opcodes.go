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
	And           // logical and operator (ops1 AND ops2)
	Xor           // exclusive or operator (ops1 XOR ops2)
	Orr           // or operator (ops1 OR ops2)
	Cmp

	Stop
	Nop

	// Data transfer op codes
	Ldr Op = iota // Load data
	Str           // Store data

	// Branching op codes
	Call Op = iota
	Ret
)

var OpString = map[string]Op{
	"mov": Mov,
	"add": Add,
	"sub": Sub,
	"rsb": Rsb,
	"and": And,
	"xor": Xor,
	"or":  Orr,
	"cmp": Cmp,

	"nop":  Nop,
	"stop": Stop,

	"ldr": Ldr,
	"str": Str,

	"call": Call,
	"ret":  Ret,
}

func (o Op) String() string {
	return OpToString[o]
}

var OpToString = map[Op]string{
	Mov: "mov",
	Add: "add",
	Sub: "sub",
	Rsb: "rsb",
	And: "and",
	Xor: "xor",
	Orr: "or",
	Cmp: "cmp",

	Nop:  "nop",
	Stop: "stop",

	Ldr: "ldr",
	Str: "str",

	Call: "call",
	Ret:  "ret",
}

type Cond byte

const (
	NoCond = iota
	Eq
	Ne
	Gt
	Gteq
	Lt
	Lteq
)

func (c Cond) String() string {
	return CondToString[c]
}

var StringToCond = map[string]Cond{
	"eq":   Eq,
	"ne":   Ne,
	"gt":   Gt,
	"lt":   Lt,
	"gteq": Gteq,
	"lteq": Lteq,
}

var CondToString = map[Cond]string{
	NoCond: "nocond",
	Eq:     "eq",
	Ne:     "ne",
	Gt:     "gt",
	Lt:     "lt",
	Gteq:   "gteq",
	Lteq:   "lteq",
}
