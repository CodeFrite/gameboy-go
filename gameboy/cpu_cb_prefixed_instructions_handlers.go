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

// helper function that returns the left rotated value of a byte and the value of the carry flag
func rotateLeft(value uint8) (uint8, bool) {
	rotateValue := (value << 1) | (value >> 7)
	carry := (value & 0x80) == 0x80
	return rotateValue, carry
}

// Rotate r8 register left and save bit 7 to the Carry flag
// opcodes:
//   - 0x00 =	RLC B
//   - 0x01 =	RLC C
//   - 0x02 =	RLC D
//   - 0x03 =	RLC E
//   - 0x04 =	RLC H
//   - 0x05 =	RLC L
//   - 0x06 =	RLC [HL]
//   - 0x07 =	RLC A
//
// flags: Z=Z N=0 H=0 C=C
func (c *CPU) RLC(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "A":
		rotatedValue, carry := rotateLeft(c.a)
		c.a = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "B":
		rotatedValue, carry := rotateLeft(c.b)
		c.b = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "C":
		rotatedValue, carry := rotateLeft(c.c)
		c.c = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "D":
		rotatedValue, carry := rotateLeft(c.d)
		c.d = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "E":
		rotatedValue, carry := rotateLeft(c.e)
		c.e = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "H":
		rotatedValue, carry := rotateLeft(c.h)
		c.h = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "L":
		rotatedValue, carry := rotateLeft(c.l)
		c.l = rotatedValue
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "HL":
		valueAtHL := c.bus.Read(c.getHL())
		rotatedValue, carry := rotateLeft(valueAtHL)
		c.bus.Write(c.getHL(), rotatedValue)
		if carry {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		if rotatedValue == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	}

	c.resetNFlag()
	c.resetHFlag()
	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
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

	boolToUint8 := func(b bool) uint8 {
		if b {
			return 1
		}
		return 0
	}

	carry := boolToUint8(c.getCFlag())
	switch instruction.Operands[0].Name {
	case "A":
		if c.a&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.a = c.a<<1 | uint8(carry)
		if c.a == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "B":
		if c.b&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.b = c.b<<1 | uint8(carry)
		if c.b == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "C":
		if c.c&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.c = c.c<<1 | uint8(carry)
		if c.c == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "D":
		if c.d&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.d = c.d<<1 | uint8(carry)
		if c.d == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "E":
		if c.e&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.e = c.e<<1 | uint8(carry)
		if c.e == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "H":
		if c.h&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.h = c.h<<1 | uint8(carry)
		if c.h == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "L":
		if c.l&(1<<7) != 0 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.l = c.l<<1 | uint8(carry)
		if c.l == 0 {
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
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
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
	if c.operand&uint16(1<<b) == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.setHFlag()
	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
}
func (c *CPU) RES(instruction *Instruction) {
	panic("RES not implemented")
}
func (c *CPU) SET(instruction *Instruction) {
	panic("SET not implemented")
}
