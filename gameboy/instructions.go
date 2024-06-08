package gameboy

import "fmt"

type Instruction struct {
	Opcode byte;
	Mnemonic string;
	Length uint16;
	Cycles uint16;
	Handler func(*CPU, []byte);
}

var Instructions map[byte]Instruction = map[byte]Instruction{
	0x00: {0x00, "NOP", 1, 4, NOP},
	0x01: {0x01, "LD BC, %X", 3, 12, LD_BC_n16},
	0xC3:	{0xC3, "JP a16", 3, 16, JP_a16},
}

var NotYetImplementedInstruction = Instruction{
	Opcode:   0xFF,
	Mnemonic: "UNIMPLEMENTED",
	Length:   1,
	Cycles:   0,
	Handler: func(*CPU, []byte) {
		fmt.Println("Unimplemented instruction")
	},
}

func NOP(c *CPU, operand []byte) {
	// Do nothing
	c.incrementPC(1)
}

func LD_BC_n16(c *CPU, operand []byte) {
	// Load the next two bytes into the BC register
	c.BC = uint16(operand[1]) << 8 + uint16(operand[0])
}

func JP_a16(c *CPU, operand []byte) {
	// Load the next two bytes into the PC register
	c.PC = uint16(operand[1]) << 8 + uint16(operand[0])
}