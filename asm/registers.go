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

type RegEntry byte

const (
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
