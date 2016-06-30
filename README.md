**NOTE: this project is much influx and should not be considered usabled**

TinyVM is a minimalistic 64bit Virtual Machine. The aim of TinyVM is to make it easy to embed
in other Go projects that require a Virtual Machine. TinyVM is **not** thread-safe, though
this is subject to change.

## Installation

To install the TinyVM binary please make sure you've got Go properly set up, then run
`go get github.com/obscuren/tinyvm`

### Usage

Basic: `tinyvm <flags> file`. TinyVM allows you to set the registers from the command line using the
`-r#`. Where `#` is the register number (0 to 15). Please take extra care when setting register 15.
This register is used for the program counter and allows you to control the flow of execution. Please
refer to the `-help` option for more information.

## Assembler

TinyVM comes with a small set of assembler instructions to make it easy to use. The `asm` package
contains an assembler language definition and a very simple compiler.

## VM

TinyVM comes with a small general purpose register (`r0..r15`), unbounded memory (`[addr]`)
and a general purpose stack mechanism (`pop`, `push`). `r15` is a special register for the
program counter and can be set to jump to arbitrary position in code. TinyVM also has a very
simple calling mechanism (`call`) and keeps an internal call stack to determine the positions
for returning (`ret`).

Setting register `r15` to anything other than the default (`0`) means execution will start from
that position and onward. In the future we'll allow labels to be specified in the form of
`v.Set(asm.Reg, asm.R15, "my_label")`, but this has to be implemented in both the vm as well as
the compiler who does not yet emit label information during the assembly stage of the "compiler".

 See Appendix I for a list op assembly operations.

## Conditional execution

TinyVM supports (like ARM) conditional execution e.g. `moveq` would only be executed if the
conditional value were to be set to zero. The conditional value can be set by appending `s`
to the mnemonic (e.g. `movs`, which sets the 25th bit) or by using the comparison instructions
`cmp` and `tst`. By default data processing instructions do not set the condition code flag.

An conditional execution must always be preceeded by a comparison instruction or an instruction
with the conditional code bit set.

## Instruction encoding

The instruction scheme used by TinyVM is based on the ARM instruction scheme though with
some arrangement in the ops order. TinyVM uses 4-bit rotation value (9th to 12th bit) to
shift the 8-bit immediate value. This clever trick allows us to encode a lot of values in
the range of `2^0..32` in only 12-bits. When an immediate value is encoded in `Ops2` the
25th bit is set to 1, indicating an immediate value is encoded in the lower 12 bits of the
instruction.

The last 4 bits will be used for conditional execution (like ARM). Any instruction can be
form of `operation[condition]` e.g. `addeq` for *add if equal* or `movgt` for *mov if
greater than*.

```
+--------------+---------+----------+----------+----------+----------+---------+---------+---------+
| Bits         |31 .. 28 | 27 .. 24 | 23 .. 20 | 19 .. 16 | 15 .. 12 | 11 .. 8 | 7 ... 4 | 3 ... 0 |
+--------------+---------+----------+----------+----------+----------+---------+---------+---------+
| Description  |  COND   |     SI   |    INS   |    Ds    |   Ops1   |   Ops2  |         |         |
+--------------+---------+----------+----------+----------+----------+---------+---------+---------+
| mov r1 #260  |  0000   |   0001   |   1010   |   0001   |   0000   |   1111  |  0100   |   0001  |
| mov r1 r2    |  0000   |   0000   |   1010   |   0001   |   0002   |   0000  |  0000   |   0000  |
+--------------+---------+----------+----------+----------+----------+---------+---------+---------+
```

## Example

### Integraton

The following example is an integration example on how you could embed TinyVM in to your
own project.

```asm
    mov     r15 main
add:    ; add taket two arguments
	add 	r0 r0 r1
	ret

main:   ; main must be called with r0 and r1 set
	call 	add

	stop

```

```go
// parse the source code
code, err := asm.Assemble(sourceCode)
if err != nil {
    panic(err)
}


v := vm.New(false) // pass "true" for debug info
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

### ASM samples

#### Counter

The following code fragment is a loop which runs until the counter in `r0` hits zero
When it hits zero the condition code `ne` (not equal to zero) controling branch becomes
false and exits the loop.

```asm
	mov     r0   #10
loop:
	subs	r0   r0 #1
	movne	r15  loop
```


## Appendix I

All operations take at least 2 argument. The first argument (dst=destination) must be a register (`r#`).

| Opcode | Argument count | Description |
|:------:|:--------------:|:-----------:|
| `mov`  | 2              | Moves `ops1` in to register `dst`
| `add`  | 3              | `ops1 + ops2` and sets the result to register `dst`
| `sub`  | 3              | `ops1 - ops2` and sets the result to register `dst`
| `rsb`  | 3              | `ops2 - ops1` and sets the result to register `dst`
| `and`  | 3              | `ops1 & ops2` and sets the result to register `dst`
| `xor`  | 3              | `ops1 ^ ops2` and sets the result to register `dst`
| `orr`  | 3              | `ops1 | ops2` and sets the result to register `dst`
| `cmp`  | 2              | `ops1 - ops2` and sets the result to the condition value
| `call` | 1              | sets `r15` to `dst` and pushes pc to the pc stack
| `ret`  | 0              | pops the pc of the pc stack and sets `r15`. `len(stack)==0` halt execution
| `stop` | 0              | halts execution

