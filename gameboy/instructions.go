package gameboy

import "fmt"

const debug = false

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
	0x11: {0x11, "LD DE, %X", 3, 12, LD_DE_n16},
	0x2C: {0x2C, "INC L", 1, 4, INC_L},
	0x4A: {0x4A, "LD C, D", 1, 4, LD_C_D},
	0x4B: {0x4B, "LD C, E", 1, 4, LD_C_E},
	0x53: {0x53, "LD D, E", 1, 4, LD_D_E},
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

// 0x00: No operation
func NOP(c *CPU, operand []byte) {
	// Do nothing
	c.incrementPC(1)
}

// 0x01: Load the next two bytes into the BC register
func LD_BC_n16(c *CPU, operand []byte) {
	// Load the next two bytes into the BC register
	if debug {
		fmt.Printf("LD BC, 0x%X\n (BC=0x%X => ", operand, c.BC)
	}
	c.BC = uint16(operand[1]) << 8 + uint16(operand[0])
	if debug {
		fmt.Printf("0x%X)\n", c.BC)
	}
	c.incrementPC(3)
}

// 0x11: Load the next two bytes into the DE register
func LD_DE_n16(c *CPU, operand []byte) {
	// Load the next two bytes into the DE register
	if debug {
		fmt.Printf("LD DE, 0x%X (DE=0x%X => ", operand, c.DE)
	}
	c.DE = bytesToUint16([2]byte{operand[0], operand[1]})
	if debug {
		fmt.Printf("0x%X)\n", c.DE)
	}
	c.incrementPC(3)
}

// 0x2C: Increment the value of register L
func INC_L(c *CPU, operand []byte) {
	if debug {
		fmt.Printf("INC L (HL=0x%X => ", c.HL)
	}
	c.HL = c.HL & 0xFF00 | (c.HL+1) & 0x00FF;
	if debug {
		fmt.Printf("0x%X)\n", c.HL)
	}
	c.incrementPC(1);
}

// 0x4A: Load the value of register E into register D
func LD_C_D(c *CPU, operand []byte) {
	if debug {
		fmt.Printf("LD C, D (DE=0x%X) (BC=0x%X => ", c.DE, c.BC)
	}
	c.BC = (c.BC & 0xFF00) |  (c.DE & 0xFF00) >> 8;
	if debug {
		fmt.Printf("0x%X)\n", c.BC)
	}
	c.incrementPC(1);
}

// 0x4B: Load the value of register E into register D
func LD_C_E(c *CPU, operand []byte) {
	if debug {
		fmt.Printf("LD C, E (DE=0x%X) (BC=0x%X => ", c.DE, c.BC)
	}
	c.BC = c.BC | (c.DE & 0x00FF);
	if debug {
		fmt.Printf("0x%X)\n", c.BC)
	}
	c.incrementPC(1);
}

// 0x53: Load the value of register E into register D
func LD_D_E(c *CPU, operand []byte) {
	if debug {
		fmt.Printf("LD D, E (DE=0x%X => ", c.DE)
	}
	c.DE = c.DE << 8 | (c.DE & 0x00FF);
	if debug {
		fmt.Printf("0x%X)\n", c.DE)
	}
	c.incrementPC(1);
}

// 0xC3: Load the next two bytes into the PC register
func JP_a16(c *CPU, operand []byte) {
	if debug {
		fmt.Printf("JP a16 (a16=0x%X)\n", operand)
	}
	c.PC = bytesToUint16([2]byte{operand[0], operand[1]})
}