package asm

import (
	"errors"
	"fmt"
)

const (
	CondPos          = 28
	ImmediateFlagPos = 24
	SFlagPos         = 25
	InstrPos         = 20
	DstPos           = 16
	Ops1Pos          = 12
	Ops2Pos          = 8
	ImmediatePos     = 0
)

type Instruction struct {
	Raw uint32

	Cond Cond
	S    bool
	Op   Op
	Dst  RegEntry
	Ops1 RegEntry
	Ops2 RegEntry

	Immediate bool
	Value     uint32
}

func EncodeInstruction(instr Instruction) (uint32, error) {
	var encoded uint32
	encoded |= (uint32(instr.Cond) << CondPos)
	encoded |= (uint32(instr.Op) << InstrPos)
	encoded |= (uint32(instr.Dst) << DstPos)
	encoded |= (uint32(instr.Ops1) << Ops1Pos)
	if instr.S {
		encoded |= 1 << SFlagPos
	}
	if instr.Immediate {
		encoded |= 1 << ImmediateFlagPos
		encodedValue, err := encodeImmediate(instr.Value)
		if err != nil {
			return 0, fmt.Errorf("instruction encoder err: %v (value=%d)", err, instr.Value)
		}
		encoded |= encodedValue
	} else {
		encoded |= (uint32(instr.Ops2) << Ops2Pos)
	}
	return encoded, nil
}

func DecodeInstruction(instruction uint32) Instruction {
	var instr Instruction
	instr.Raw = instruction
	instr.Cond = Cond(getBits(instruction, CondPos, CondPos+3))
	instr.Op = Op(getBits(instruction, InstrPos, InstrPos+3))
	instr.Dst = RegEntry(getBits(instruction, DstPos, DstPos+3))
	instr.Ops1 = RegEntry(getBits(instruction, Ops1Pos, Ops1Pos+3))
	instr.S = isSet(instruction, SFlagPos)

	if isSet(instruction, ImmediateFlagPos) {
		instr.Immediate = true
		instr.Value = decodeImmediate(getBits(instruction, 0, 11))
	} else {
		instr.Ops2 = RegEntry(getBits(instruction, Ops2Pos, Ops2Pos+3))
	}

	return instr
}

func isSet(n uint32, bit uint32) bool {
	return (n >> bit & 1) == 1
}

func setBit(n *uint32, bit uint32) {
	*n |= 1 << bit
}

func getBits(n uint32, offset uint32, bits uint32) uint32 {
	var flag uint32
	for i := uint32(0); i <= bits-offset; i++ {
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
