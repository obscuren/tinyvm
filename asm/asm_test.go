package asm

import "fmt"

func ExampleEncodeInstruction() {
	instr := Instruction{
		Instruction: Mov,
		Dst:         R1,
		Immediate:   true,
		Value:       260,
	}
	encoded, err := EncodeInstruction(instr)
	if err != err {
		fmt.Println(err)
	}
	fmt.Printf("%032b\n", encoded)

	instr = Instruction{
		Instruction: Mov,
		Dst:         R1,
		Ops1:        R2,
	}

	encoded, err = EncodeInstruction(instr)
	if err != err {
		fmt.Println(err)
	}
	fmt.Printf("%032b\n", encoded)
	// Output:
	// 00000001010100010000111101000001
	// 00000000010100010001000000000000
}
