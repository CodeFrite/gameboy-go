package gameboy

import "fmt"

func (instr *Instruction) print() {
	fmt.Println("")
	fmt.Println("\n> Instruction:")
	fmt.Println("--------------")
	if instr == nil {
		fmt.Println("> No instruction to print")
	} else {
		//fmt.Printf("- Opcode: 0x%02X\n", instr.IR) // find a way to bring this back
		fmt.Printf("- Mnemonic: %s\n", instr.Mnemonic)
		fmt.Printf("- Bytes: %d\n", instr.Bytes)
		fmt.Printf("- Cycles: %v\n", instr.Cycles)
		fmt.Printf("- Operands: %v\n", instr.Operands)
		fmt.Printf("- Immediate: %t\n", instr.Immediate)
		fmt.Printf("- Flags: %v\n", instr.Flags)
	}
}
