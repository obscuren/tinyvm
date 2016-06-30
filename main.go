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
	versionFlag = flag.Bool("version", false, "outputs version string")
	statFlag    = flag.Bool("vmstats", false, "display virtual machine stats")
	printCode   = flag.Bool("printcode", false, "prints executing code in hex")
	debug       = flag.Bool("debug", false, "prints debug information during execution")
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Println("TinyVM", vm.VersionString, "- (c) Jeffrey Wilcke")
		os.Exit(0)
	}

	var (
		code []byte
		err  error
	)
	if len(flag.Args()) > 0 {
		var err error
		code, err = ioutil.ReadFile(flag.Args()[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		code, err = asm.Assemble(string(code))
	} else {
		err = fmt.Errorf("Usage: tinyvm <flags> filename")
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if *printCode {
		fmt.Printf("(len=%d) %x\n", len(code), code)
		for i := 0; i < len(code); i += 4 {
			for _, b := range code[i : i+4] {
				fmt.Printf("%08b", b)
			}
			fmt.Printf(" ")
		}
		fmt.Println()
	}

	v := vm.New(*debug)
	for i, registerFlag := range registerFlags {
		v.Set(asm.Reg, uint32(i), uint32(*registerFlag))
	}

	if err := v.Exec(code); err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	if *statFlag {
		v.Stats()
	}

	fmt.Println(v.Get(asm.Reg, asm.R0))
}

var registerFlags [asm.MaxRegister]*int

func init() {
	for i := 0; i < asm.MaxRegister; i++ {
		registerFlags[i] = flag.Int(fmt.Sprintf("r%d", i), 0, fmt.Sprintf("sets the r%d register", i))
	}
}
