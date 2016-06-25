package main

import (
	"flag"
	"fmt"

	"github.com/obscuren/tinyvm/asm"
	"github.com/obscuren/tinyvm/vm"
)

var (
	statFlag  = flag.Bool("vmstats", false, "display virtual machine stats")
	printCode = flag.Bool("printcode", false, "prints executing code in hex")
)

func main() {
	flag.Parse()

	fmt.Println("TinyVM", vm.VersionString, "- (c) Jeffrey Wilcke")

	code := asm.Parse(stack)
	if *printCode {
		fmt.Printf("%x\n", code)
	}

	v := vm.New()
	if err := v.Exec(code); err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println("exit:", v.Get(asm.Reg, asm.R0))

	if *statFlag {
		v.Stats()
	}
}

const (
	stack = `
		push 1
		pop
		push 255
		mov r0 pop
		push 1
		push 2

		stop
	`
	call = `
	jmp main

	nop:
		ret
	main:
		call nop
	`
	minJmp = `
		jmp end
	end:
	`
	example = `
		mov r0 0

	while_not_3:
		add r0 r0 1

		lt r0 3
		jmpi while_not_3

		mov r1 r0
	while_not_0:
		sub r1 r1 1

		gt r1 0
		jmpi while_not_0

	not_happening:
		eq 1 0
		jmpi not_happening

		ret r0
	`

	// r0 = c
	// r1 = next
	// r2 = first
	// r3 = second
	// r4 = n
	fibanocci = `
	mov r4 5 	; find number 5
	mov r3 1	; set r3 to 1

for_loop:
	lt r0 r4
	jmpn end
start_if:
	lteq r0 1
	jmpn else

	mov r1 r0
	jmp end_if
else:
	add r1 r2 r3
	mov r2 r3
	mov r3 r1
end_if:
	add r0 r0 1
	jmp for_loop
end:
	ret r1
	
`
)
