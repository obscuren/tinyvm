// Copyright 2016 Jeffrey Wilcke
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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

	var code []byte
	if len(flag.Args()) > 0 {
		var err error
		code, err = ioutil.ReadFile(flag.Args()[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		code = asm.Parse(string(code))
	} else {
		code = asm.Parse(fibanocci)
	}
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
	addProgram = `
		jmp 	main
	add:    ; add taket two arguments
		add 	r0 r0 r1
		ret

	main:   ; main must be called with r0 and r1 set
		call 	add

		stop
	`
	stack = `
		push 	1
		pop
		push 	255
		mov 	r0 pop
		push 	1
		push 	2

		stop
	`

	call = `
	jmp main

	nop:
		ret
	main:
		call 	nop
	`

	example = `
		mov 	r0 0
		mov 	r10 0

	while_not_3:
		add 	r0 r0 1

		lt 	r10 r0 3
		jmpi 	r10 while_not_3

		mov 	r1 r0
		mov 	r10 0
	while_not_0:
		sub 	r1 r1 1

		gt 	r10 r1 0
		jmpi 	r10 while_not_0

	not_happening:
		eq 	1 0
		jmpi 	not_happening
	`

	mov = `
	mov 	r4 5
	mov 	r3 1
	`

	// r0 = c
	// r1 = next
	// r2 = first
	// r3 = second
	// r4 = n
	fibanocci = `
	mov	r4 5 	; find number 5
	mov	r3 1	; set r3 to 1

for_loop:
	lt 	r10 r0 r4
	jmpn 	r10 end
start_if:
	lteq 	r10 r0 1
	jmpn 	r10 else

	mov 	r1 r0
	jmp 	end_if
else:
	add 	r1 r2 r3
	mov 	r2 r3
	mov 	r3 r1
end_if:
	add 	r0 r0 1
	jmp 	for_loop
end:
	mov 	r0 r1
	
`
)
