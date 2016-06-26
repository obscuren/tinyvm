**NOTE: this project is much influx and should not be considered usabled**

TinyVM is a minimalistic 64bit Virtual Machine. The aim of TinyVM is to make it easy to embed
in other Go projects that require a Virtual Machine. TinyVM is **not** thread-safe, though
this is subject to change.

## Assembler

TinyVM comes with a small set of assembler instructions to make it easy to use. The `asm` package
contains an assembler language definition and a very simple compiler.

## VM

TinyVM comes with a small general purpose register (`r0..r15`, `pc`), unbounded memory (`[addr]`)
and a general purpose stack mechanism (`pop`, `push`). It supports arbitrary jumps `jmp(n|i)` and
a simply calling mechanism (`call`) and keeps an internel call stack to determine the positions for
returning (`ret`).

## Example

```asm
    jmp 	main
add:    ; add taket two arguments
	add 	r0 r0 r1
	ret

main:   ; main must be called with r0 and r1 set
	call 	add

	stop

```

```go
// parse the source code
code, err := asm.Parse(sourceCode)
if err != nil {
    panic(err)
}


v := vm.New()
// set the registers (as required for "main")
v.Set(asm.Reg, asm.R0, 3) // set r0 to 3
v.Set(asm.Reg, asm.R1, 2) // set r1 to 2

// execute the compiled code
if err := v.Exec(code); err != nil {
    fmt.Println("err", err)
        return
}

// exit: 5
fmt.Println("exit:", v.Get(asm.Reg, asm.R0))
```
