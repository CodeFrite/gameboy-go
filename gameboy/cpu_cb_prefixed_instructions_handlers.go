package gameboy

import "fmt"

// > instructions handlers (PREFIX CB)

// Route the execution to the corresponding instruction handler (PREFIX CB)
func (c *CPU) executeCBInstruction(instruction Instruction) {
	// Execute the corresponding instruction
	switch instruction.Mnemonic {
	case "RLC":
		c.RLC(&instruction)
	case "RRC":
		c.RRC(&instruction)
	case "RL":
		c.RL(&instruction)
	case "RR":
		c.RR(&instruction)
	case "SLA":
		c.SLA(&instruction)
	case "SRA":
		c.SRA(&instruction)
	case "SWAP":
		c.SWAP(&instruction)
	case "SRL":
		c.SRL(&instruction)
	case "BIT":
		c.BIT(&instruction)
	case "RES":
		c.RES(&instruction)
	case "SET":
		c.SET(&instruction)
	default:
		fmt.Println("Unknown instruction")
	}
}

func (c *CPU) RLC(instruction *Instruction) {
	panic("RLC not implemented")
}
func (c *CPU) RRC(instruction *Instruction) {
	panic("RRC not implemented")
}

/*
	 RL r8 / [HL]
	 Rotate r8 or [HL] left through carry: old bit 7 to Carry flag, new bit 0 to bit 7.
	 opcodes:
		- 0x10:	RL B
		- 0x11:	RL C
		- 0x12:	RL D
		- 0x13:	RL E
		- 0x14:	RL H
		- 0x15:	RL L
		- 0x16:	RL [HL]
		- 0x17:	RL A
		 flags: Z=Z N=0 H=0 C=C
*/
func (c *CPU) RL(instruction *Instruction) {
	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	carry := boolToInt(c.getCFlag())
	switch instruction.Operands[0].Name {
	case "A":
		if c.A&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.A = c.A<<1 | uint8(carry)
		if c.A == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "B":
		if c.B&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.B = c.B<<1 | uint8(carry)
		if c.B == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "C":
		if c.C&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.C = c.C<<1 | uint8(carry)
		if c.C == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "D":
		if c.D&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.D = c.D<<1 | uint8(carry)
		if c.D == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "E":
		if c.E&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.E = c.E<<1 | uint8(carry)
		if c.E == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "H":
		if c.H&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.H = c.H<<1 | uint8(carry)
		if c.H == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "L":
		if c.L&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.L = c.L<<1 | uint8(carry)
		if c.L == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "[HL]":
		val := c.bus.Read(c.getHL())
		if val&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		newVal := val<<1 | uint8(carry)
		c.bus.Write(c.getHL(), newVal)
		if newVal == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	}
	c.resetNFlag()
	c.resetHFlag()
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}
func (c *CPU) RR(instruction *Instruction) {
	panic("RR not implemented")
}
func (c *CPU) SLA(instruction *Instruction) {
	panic("SLA not implemented")
}
func (c *CPU) SRA(instruction *Instruction) {
	panic("SRA not implemented")
}
func (c *CPU) SWAP(instruction *Instruction) {
	panic("SWAP not implemented")
}
func (c *CPU) SRL(instruction *Instruction) {
	panic("SRL not implemented")
}

/*
BIT b, r8 / [HL]
Test bit b in register r8 or [HL]: If bit b is 0, Z is set.
opcodes: 0x40-0x7F
flags: Z=Z N=0 H=1 C=-
*/
func (c *CPU) BIT(instruction *Instruction) {
	// get the bit position to test
	opStr := instruction.Operands[0].Name // the bit position to test is given as a string
	var b uint16 = uint16(opStr[0] - '0')
	// check if bit b of operand is 0
	if c.Operand&uint16(1<<b) == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.setHFlag()
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}
func (c *CPU) RES(instruction *Instruction) {
	panic("RES not implemented")
}
func (c *CPU) SET(instruction *Instruction) {
	panic("SET not implemented")
}
