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

type Op byte

const (
	Mov Op = iota
	Add
	Sub

	Cmp

	Call
	Ret

	Push
	Pop

	Jmpt

	Stop
	Nop
)

var OpString = map[string]Op{
	"mov":  Mov,
	"add":  Add,
	"sub":  Sub,
	"cmp":  Cmp,
	"push": Push,
	"pop":  Pop,
	"call": Call,
	"ret":  Ret,
	"nop":  Nop,
	"stop": Stop,
}

func (o Op) String() string {
	return OpToString[o]
}

var OpToString = map[Op]string{
	Mov:  "mov",
	Add:  "add",
	Sub:  "sub",
	Cmp:  "cmp",
	Push: "push",
	Pop:  "pop",
	Call: "call",
	Ret:  "ret",
	Nop:  "nop",
	Stop: "stop",
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
