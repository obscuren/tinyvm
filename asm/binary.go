package asm

import (
	"errors"
	"fmt"
)

/*
   +--------------+---------+----------+----------+----------+----------+---------+---------+---------+
   | Bits         |31 .. 28 | 27 .. 24 | 23 .. 20 | 19 .. 16 | 15 .. 12 | 11 .. 8 | 7 ... 4 | 3 ... 0 |
   +--------------+---------+----------+----------+----------+----------+---------+---------+---------+
   | Description  |  COND   |    I     |    INS   |    Ds    |   Op1    |    Op2  |         |         |
   +--------------+---------+----------+----------+----------+----------+---------+---------+---------+
   | mov r1 #260  |  0000   |   0001   |   0101   |   0001   |   0000   |   1111  |  0100   |   0001  |
   | mov r1 r2    |  0000   |   0000   |   0101   |   0001   |   0002   |   0000  |  0000   |   0000  |
   +--------------+---------+----------+----------+----------+----------+---------+---------+---------+
*/

const (
	CondPos          = 27
	ImmediateFlagPos = 24
	InstrPos         = 19
	DestPos          = 16
	Ops1Pos          = 11
	Ops2Pos          = 7
	ImmediatePos     = 0
)

type Instruction struct {
	Ca, Cb, Cc, Cd bool
	Instruction    Op
	Dst            RegEntry
	Ops1           RegEntry
	Ops2           RegEntry

	Immediate bool
	Value     uint32
}

func EncodeInstruction(instr Instruction) (uint32, error) {
	var encoded uint32
	encoded |= (uint32(instr.Instruction) << InstrPos)
	encoded |= (uint32(instr.Dst) << DestPos)
	encoded |= (uint32(instr.Ops1) << Ops1Pos)
	if instr.Immediate {
		encoded |= 1 << ImmediateFlagPos
		encodedValue, err := encodeImmediate(instr.Value)
		if err != nil {
			return 0, fmt.Errorf("instruction encoder err: %v (value=%d)", err, instr.Value)
		}
		encoded |= encodedValue
	}
	return encoded, nil
}

func isSet(n uint32, bit uint32) bool {
	return (n >> bit & 1) == 1
}

func setBit(n *uint32, bit uint32) {
	*n |= 1 << bit
}

func getBits(n uint32, offset uint32, bits uint32) uint32 {
	var flag uint32
	for i := uint32(0); i < bits; i++ {
		setBit(&flag, i)
	}
	return n >> offset & flag
}

func rol(n, i uint32) uint32 {
	return (n << i) | (n >> (32 - i))
}

func ror(n, i uint32) uint32 {
	return (n >> i) | (n << (32 - i))
}

func encodeImmediate(n uint32) (uint32, error) {
	var m uint32
	for i := uint32(0); i < 16; i++ {
		m = rol(n, i*2)
		if m < 256 {
			return uint32(i<<8) | m, nil
		}
	}
	return 0, errors.New("unencodable constast")
}

func decodeImmediate(n uint32) uint32 {
	shift := (n >> 8) & 0xf
	immediate := n & 0xff
	return ror(immediate, shift*2)
}
