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

// helper function that returns the right rotated value of a byte and the value of the carry flag
func rotateRight(value uint8) (uint8, bool) {
	rotateValue := (value >> 1) | (value << 7)
	carry := (value & 0x01) == 0x01
	return rotateValue, carry
}

// Rotate r8 register right and save bit 0 to the Carry flag
// opcodes:
//   - 0x08 =	RRC B
//   - 0x09 =	RRC C
//   - 0x0A =	RRC D
//   - 0x0B =	RRC E
//   - 0x0C =	RRC H
//   - 0x0D =	RRC L
//   - 0x0E =	RRC [HL]
//   - 0x0F =	RRC A
//
// flags: Z=Z N=0 H=0 C=C
func (c *CPU) RRC(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "A":
		rotatedValue, carry := rotateRight(c.a)
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
		rotatedValue, carry := rotateRight(c.b)
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
		rotatedValue, carry := rotateRight(c.c)
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
		rotatedValue, carry := rotateRight(c.d)
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
		rotatedValue, carry := rotateRight(c.e)
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
		rotatedValue, carry := rotateRight(c.h)
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
		rotatedValue, carry := rotateRight(c.l)
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
		rotatedValue, carry := rotateRight(valueAtHL)
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

// RL r8 / [HL]
// Rotate r8 or [HL] left through carry: old bit 7 to Carry flag, new bit 0 from carry flag
// opcodes:
//   - 0x10:	RL B
//   - 0x11:	RL C
//   - 0x12:	RL D
//   - 0x13:	RL E
//   - 0x14:	RL H
//   - 0x15:	RL L
//   - 0x16:	RL [HL]
//   - 0x17:	RL A
//
// flags: Z=Z N=0 H=0 C=C
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
	case "HL":
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

// RR r8 / [HL]
// Rotate r8 or [HL] right through carry: old bit 0 to Carry flag, new bit 7 from carry flag
// opcodes:
//   - 0x18:	RR B
//   - 0x19:	RR C
//   - 0x1A:	RR D
//   - 0x1B:	RR E
//   - 0x1C:	RR H
//   - 0x1D:	RR L
//   - 0x1E:	RR [HL]
//   - 0x1F:	RR A
//
// flags: Z=Z N=0 H=0 C=C
func (c *CPU) RR(instruction *Instruction) {
	boolToUint8 := func(b bool) uint8 {
		if b {
			return 1
		}
		return 0
	}

	carry := boolToUint8(c.getCFlag())
	switch instruction.Operands[0].Name {
	case "A":
		if c.a&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.a = c.a>>1 | uint8(carry)<<7
		if c.a == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "B":
		if c.b&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.b = c.b>>1 | uint8(carry)<<7
		if c.b == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "C":
		if c.c&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.c = c.c>>1 | uint8(carry)<<7
		if c.c == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "D":
		if c.d&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.d = c.d>>1 | uint8(carry)<<7
		if c.d == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "E":
		if c.e&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.e = c.e>>1 | uint8(carry)<<7
		if c.e == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "H":
		if c.h&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.h = c.h>>1 | uint8(carry)<<7
		if c.h == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "L":
		if c.l&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		c.l = c.l>>1 | uint8(carry)<<7
		if c.l == 0 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
	case "HL":
		val := c.bus.Read(c.getHL())
		if val&0x01 == 0x01 {
			c.setCFlag()
		} else {
			c.resetCFlag()
		}
		newVal := val>>1 | uint8(carry)<<7
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

// Shift Left Arithmetically register r8.
// Carry ←╂─ b7 ← ... ← b0 ←╂─ 0
// opcodes:
//   - 0x20:	SLA B
//   - 0x21:	SLA C
//   - 0x22:	SLA D
//   - 0x23:	SLA E
//   - 0x24:	SLA H
//   - 0x25:	SLA L
//   - 0x26:	SLA [HL]
//   - 0x27:	SLA A
//
// flags: Z=Z N=0 H=0 C=b7
func (c *CPU) SLA(instruction *Instruction) {
	if c.operand&0x80 == 0x80 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	if c.operand == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.resetHFlag()
	switch instruction.Operands[0].Name {
	case "A":
		c.a = c.a << 1
	case "B":
		c.b = c.b << 1
	case "C":
		c.c = c.c << 1
	case "D":
		c.d = c.d << 1
	case "E":
		c.e = c.e << 1
	case "H":
		c.h = c.h << 1
	case "L":
		c.l = c.l << 1
	case "HL":
		valueAtHL := c.bus.Read(c.getHL())
		c.bus.Write(c.getHL(), valueAtHL<<1)
	}

	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
}

// Shift Right Arithmetically register r8 (bit 7 of r8 is unchanged).
// ┃ b7 → ... → b0 ─╂→ Carry
// opcodes:
//   - 0x28:	SRA B
//   - 0x29:	SRA C
//   - 0x2A:	SRA D
//   - 0x2B:	SRA E
//   - 0x2C:	SRA H
//   - 0x2D:	SRA L
//   - 0x2E:	SRA [HL]
//   - 0x2F:	SRA A
//
// flags: Z=Z N=0 H=0 C=b0
func (c *CPU) SRA(instruction *Instruction) {
	msb := uint8(c.operand) & 0x80
	if c.operand&0x01 == 0x01 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	if c.operand == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.resetHFlag()
	switch instruction.Operands[0].Name {
	case "A":
		c.a = c.a>>1 | msb
	case "B":
		c.b = c.b>>1 | msb
	case "C":
		c.c = c.c>>1 | msb
	case "D":
		c.d = c.d>>1 | msb
	case "E":
		c.e = c.e>>1 | msb
	case "H":
		c.h = c.h>>1 | msb
	case "L":
		c.l = c.l>>1 | msb
	case "HL":
		valueAtHL := c.bus.Read(c.getHL())
		c.bus.Write(c.getHL(), valueAtHL>>1|msb)
	}
	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
}

// Swap upper & lower nibles of r8 register
// opcodes:
//   - 0x30:	SWAP B
//   - 0x31:	SWAP C
//   - 0x32:	SWAP D
//   - 0x33:	SWAP E
//   - 0x34:	SWAP H
//   - 0x35:	SWAP L
//   - 0x36:	SWAP [HL]
//   - 0x37:	SWAP A
//
// flags: Z=Z N=0 H=0 C=0
func (c *CPU) SWAP(instruction *Instruction) {
	if c.operand == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.resetHFlag()
	c.resetCFlag()
	switch instruction.Operands[0].Name {
	case "A":
		c.a = c.a<<4 | c.a>>4
	case "B":
		c.b = c.b<<4 | c.b>>4
	case "C":
		c.c = c.c<<4 | c.c>>4
	case "D":
		c.d = c.d<<4 | c.d>>4
	case "E":
		c.e = c.e<<4 | c.e>>4
	case "H":
		c.h = c.h<<4 | c.h>>4
	case "L":
		c.l = c.l<<4 | c.l>>4
	case "HL":
		valueAtHL := c.bus.Read(c.getHL())
		c.bus.Write(c.getHL(), valueAtHL<<4|valueAtHL>>4)
	}
	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
}

// Shift Right Logically register r8.
// 0 ─╂→ b7 → ... → b0 ─╂→ Carry
// opcodes:
//   - 0x38:	SRL B
//   - 0x39:	SRL C
//   - 0x3A:	SRL D
//   - 0x3B:	SRL E
//   - 0x3C:	SRL H
//   - 0x3D:	SRL L
//   - 0x3E:	SRL [HL]
//   - 0x3F:	SRL A
//
// flags: Z=Z N=0 H=0 C=b0
func (c *CPU) SRL(instruction *Instruction) {
	shiftedValue := uint8(c.operand) >> 1
	if shiftedValue == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	if c.operand&0x01 == 0x01 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	c.resetNFlag()
	c.resetHFlag()
	switch instruction.Operands[0].Name {
	case "A":
		c.a = shiftedValue
	case "B":
		c.b = shiftedValue
	case "C":
		c.c = shiftedValue
	case "D":
		c.d = shiftedValue
	case "E":
		c.e = shiftedValue
	case "H":
		c.h = shiftedValue
	case "L":
		c.l = shiftedValue
	case "HL":
		c.bus.Write(c.getHL(), shiftedValue)
	}
	// update the program counter offset
	c.offset = c.pc + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.cpuCycles += uint64(instruction.Cycles[0])
}

// Test bit b in register r8 or [HL]: If bit b is 0, Z is set.
// opcodes:
//   - 0x40:	BIT 0, B
//   - 0x41:	BIT 0, C
//   - 0x42:	BIT 0, D
//   - 0x43:	BIT 0, E
//   - 0x44:	BIT 0, H
//   - 0x45:	BIT 0, L
//   - 0x46:	BIT 0, [HL]
//   - 0x47:	BIT 0, A
//   - 0x48:	BIT 1, B
//   - 0x49:	BIT 1, C
//   - 0x4A:	BIT 1, D
//   - 0x4B:	BIT 1, E
//   - 0x4C:	BIT 1, H
//   - 0x4D:	BIT 1, L
//   - 0x4E:	BIT 1, [HL]
//   - 0x4F:	BIT 1, A
//   - 0x50:	BIT 2, B
//   - 0x51:	BIT 2, C
//   - 0x52:	BIT 2, D
//   - 0x53:	BIT 2, E
//   - 0x54:	BIT 2, H
//   - 0x55:	BIT 2, L
//   - 0x56:	BIT 2, [HL]
//   - 0x57:	BIT 2, A
//   - 0x58:	BIT 3, B
//   - 0x59:	BIT 3, C
//   - 0x5A:	BIT 3, D
//   - 0x5B:	BIT 3, E
//   - 0x5C:	BIT 3, H
//   - 0x5D:	BIT 3, L
//   - 0x5E:	BIT 3, [HL]
//   - 0x5F:	BIT 3, A
//   - 0x60:	BIT 4, B
//   - 0x61:	BIT 4, C
//   - 0x62:	BIT 4, D
//   - 0x63:	BIT 4, E
//   - 0x64:	BIT 4, H
//   - 0x65:	BIT 4, L
//   - 0x66:	BIT 4, [HL]
//   - 0x67:	BIT 4, A
//   - 0x68:	BIT 5, B
//   - 0x69:	BIT 5, C
//   - 0x6A:	BIT 5, D
//   - 0x6B:	BIT 5, E
//   - 0x6C:	BIT 5, H
//   - 0x6D:	BIT 5, L
//   - 0x6E:	BIT 5, [HL]
//   - 0x6F:	BIT 5, A
//   - 0x70:	BIT 6, B
//   - 0x71:	BIT 6, C
//   - 0x72:	BIT 6, D
//   - 0x73:	BIT 6, E
//   - 0x74:	BIT 6, H
//   - 0x75:	BIT 6, L
//   - 0x76:	BIT 6, [HL]
//   - 0x77:	BIT 6, A
//   - 0x78:	BIT 7, B
//   - 0x79:	BIT 7, C
//   - 0x7A:	BIT 7, D
//   - 0x7B:	BIT 7, E
//   - 0x7C:	BIT 7, H
//   - 0x7D:	BIT 7, L
//   - 0x7E:	BIT 7, [HL]
//   - 0x7F:	BIT 7, A
//
// flags: Z=Z N=0 H=1 C=-
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
