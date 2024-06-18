package gameboy

import (
	"fmt"
)

// > instructions handlers (NO PREFIX)

// Misc / Control instructions

/*
 Disable Interrupts (DI)
 Disables the IME flag to prevent the CPU from responding to interrupts
 opcodes: 0xF3
 flags: -
*/
func (c *CPU) DI(instruction *Instruction) {
	c.IME = false
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 Enable Interrupts (EI)
 Enables the IME flag to allow the CPU to respond to interrupts
 Does not enable interrupts immediately, the next instruction will be executed before the interrupts are enabled
 opcodes: 0xFB
 flags: -
*/
func (c *CPU) EI(instruction *Instruction) {
	// execute the next instruction before enabling interrupts
	c.incrementPC(uint16(instruction.Bytes))
	c.step()
	c.IME = true
}

/*
 * HALT: Halt the CPU until an interrupt occurs
 * opcodes: 0x76
 * flags: -
 */
func (c *CPU) HALT(instruction *Instruction) {
	c.halted = true
}

/*
	NOP: No operation, does nothing apart from incrementing the program counter
	opcodes: 0x00=NOP
	flags impacted: -
*/
func (c *CPU) NOP(instruction *Instruction) {
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 0xCB = PREFIX CB
 Indicates that the next instruction is a CB instruction
 flags: -
*/
func (c *CPU) PREFIX(instruction *Instruction) {
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
	// fetch the next opcode
	c.fetchOpcode()
	// get the instruction from the opcodes.json file
	cbInstruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), true)
	// execute the instruction
	c.executeCBInstruction(cbInstruction)
}
func (c *CPU) STOP(instruction *Instruction) {
	panic("STOP not implemented")
}

// Jump / Call instructions
// only the conditional instructions increment the PC if the condition is not met since they are meant to position the PC at the operand
/*
	CALL: Call a subroutine
	opcodes:
		- 0xC4 = CALL NZ, a16
		- 0xCC = CALL Z, a16
		- 0xCD = CALL a16
		- 0xD4 = CALL NC, a16
		- 0xDC = CALL C, a16
*/
func (c *CPU) CALL(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "Z":
		if c.getZFlag() {
			c.push(c.PC)
			c.PC = uint16(c.Operand)
		}
	case "NZ":
		if !c.getZFlag() {
			c.push(c.PC)
			c.PC = uint16(c.Operand)
		}
	case "C":
		if c.getCFlag() {
			c.push(c.PC)
			c.PC = uint16(c.Operand)
		}
	case "NC":
		if !c.getCFlag() {
			c.push(c.PC)
			c.PC = uint16(c.Operand)
		}
	case "a16":
		c.push(c.PC)
		c.PC = uint16(c.Operand)
	default:
		panic("CALL: unknown operand")
	}
}

/*
	JP: Jumps to an address
	opcodes:
		- 0xC2 = JP NZ, a16
		- 0xC3 = JP a16
		- 0xCA = JP Z, a16
		- 0xD2 = JP NC, a16
		- 0xDA = JP C, a16
		- 0xE9 = JP HL
	flags: -
*/
func (c *CPU) JP(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "Z":
		if c.getZFlag() {
			c.PC = uint16(c.Operand)
		} else {
			c.incrementPC(uint16(instruction.Bytes))
		}
	case "NZ":
		if !c.getZFlag() {
			c.PC = uint16(c.Operand)
		}
	case "C":
		if c.getCFlag() {
			c.PC = uint16(c.Operand)
		}
	case "NC":
		if !c.getCFlag() {
			c.PC = uint16(c.Operand)
		}
	case "a16":
		c.PC = uint16(c.Operand)
	case "HL":
		c.PC = uint16(c.Operand)
	default:
		panic("JP: unknown operand")
	}
}
func (c *CPU) JR(instruction *Instruction) {
	panic("JR not implemented")
}

/*
	RET: Return from a subroutine
	This intruction pops the address from the stack and jumps to it
	opcodes: 0xC9
	flags: -
*/
func (c *CPU) RET(instruction *Instruction) {
	c.PC = c.bus.Read16(c.SP)
}

/*
	RETI: Return from interrupt
	Return from subroutine and enable interrupts.
	This is basically equivalent to executing EI then RET, meaning that IME is set right after this instruction.
	opcodes: 0xD9
	flags: -
*/
func (c *CPU) RETI(instruction *Instruction) {
	panic("RETI not implemented")
}
func (c *CPU) RST(instruction *Instruction) {
	panic("RST not implemented")
}

// Load / Store instructions

/*
	LD: Load data from one location to another
	opcodes:

	LD r8, n8 = LD A/B/C/D/E/H/L, n8 								[7]
	LD r8, r8 = LD A/B/C/D/E/H/L, A/B/C/D/E/H/L 		[49]
	LD r8, [r8] + 0xFF00 = LD A, [C]								[1]
	LD r8, [a16] = LD A, [a16]											[1]
	LD r8, [r16] = LD A, [BC]/[DE]/[HL]/[HL+]/[HL-]	[5]
							 = LD B/C/D/E/L/H, [HL]							[6]
	LD [r8], r8 = LD [C], A													[1]

	LD r16, n16 = LD BC/DE/HL/SP, n16								[4]
	LD r16, r16 + e8 = LD HL, SP+e8									[1]
	LD r16, r16 = LD SP, HL													[1]
	LD [r16], r8 = LD [BC]/[DE]/[HL+]/[HL-],  A			[4]
							 = LD [HL],  A/B/C/D/E/H/L					[7]
	LD [r16], n8 = LD [HL], n8											[1]

	LD [a16], r8 = LD [a16], A											[1]
	LD [a16], r16 = LD [a16], SP										[1]


	flags: - except for 0xF8 where Z->0 N->0 H->C C->C

	NOTE: all LD instructions have 2 operands, the first one is always the destination and the second one is always the source (except for LD HL, SP+r8)
	=> we will 'automate' the process of fetching the operands expect for LD HL, SP+r8 that will be handled manually
*/
func (c *CPU) LD(instruction *Instruction) {
	var address uint16
	switch instruction.Operands[0].Name {
	case "A":
		c.A = uint8(c.Operand)
	case "B":
		c.B = uint8(c.Operand)
	case "C":
		if instruction.Operands[0].Immediate {
			c.C = uint8(c.Operand)
		} else {
			address = 0xFF00 | uint16(c.C)
			c.bus.Write(address, uint8(c.Operand))
		}
	case "D":
		c.D = uint8(c.Operand)
	case "E":
		c.E = uint8(c.Operand)
	case "H":
		c.H = uint8(c.Operand)
	case "L":
		c.L = uint8(c.Operand)
	case "BC":
		if instruction.Operands[0].Immediate {
			c.setBC(c.Operand)
		} else {
			c.bus.Write(c.getBC(), uint8(c.Operand))
		}
	case "DE":
		if instruction.Operands[0].Immediate {
			c.setDE(c.Operand)
		} else {
			c.bus.Write(c.getDE(), uint8(c.Operand))
		}
	case "HL":
		if instruction.Operands[0].Immediate {
			c.setHL(c.Operand)
		} else {
			c.bus.Write(c.getHL(), uint8(c.Operand))
		}
		if instruction.Operands[0].Increment {
			c.setHL(c.getHL() + 1)
		} else if instruction.Operands[0].Decrement {
			c.setHL(c.getHL() - 1)
		}
	case "SP":
		c.SP = uint16(c.Operand)
	case "a16":
		panic("LD [a16], r8/r16 not implemented")
	}
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}
func (c *CPU) LDH(instruction *Instruction) {
	panic("LDH not implemented")
}
func (c *CPU) PUSH(instruction *Instruction) {
	panic("PUSH not implemented")
}
func (c *CPU) POP(instruction *Instruction) {
	panic("POP not implemented")
}

// Arithmetic / Logical instructions
func (c *CPU) ADC(instruction *Instruction) {
	panic("ADC not implemented")
}
func (c *CPU) ADD(instruction *Instruction) {
	panic("ADD not implemented")
}
func (c *CPU) AND(instruction *Instruction) {
	panic("AND not implemented")
}

/*
 * CCF: Complement Carry Flag
 * opcodes: 0x3F
 * flags: Z:- N:0 H:0 C:~C
 */
func (c *CPU) CCF(instruction *Instruction) {
	c.toggleCFlag()
	// reset N and H flags
	c.resetNFlag()
	c.resetHFlag()
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}
func (c *CPU) CP(instruction *Instruction) {
	panic("CP not implemented")
}

/*
 * CPL: Complement A (flip all bits)
 * opcodes: 0x2F=CPL
 * flags: Z:- N:1 H:1 C:-
 */
func (c *CPU) CPL(instruction *Instruction) {
	// flip all bits of the accumulator
	c.A = ^c.A
	// update flags
	c.setNFlag()
	c.setHFlag()
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 * DAA: Decimal Adjust Accumulator
 * This instruction adjusts the contents of the accumulator to form a binary-coded decimal (BCD) representation.
 * The DAA instruction adjusts the result of an addition or subtraction operation so that the correct representation of the result is obtained.
 * It only relies on the content of the accumulator and the flags to correct the result.
 * opcodes: 0x27=DAA
 * flags: Z:Z N:- H:0 C:C
 */
func (c *CPU) DAA(instruction *Instruction) {
	offset := uint8(0)
	// if the last operation was an addition
	if !c.getNFlag() {
		// lower nibble correction
		if c.A&0x0F > 0x09 || c.getHFlag() {
			offset = 0x06
		}
		// upper nibble correction
		if c.A > 0x99 || c.getCFlag() {
			offset |= 0x60
			c.setCFlag() // set the carry flag
		}
		// apply the correction
		c.A += offset

	} else {
		// if the last operation was subtraction
		// lower nibble correction
		if c.A&0x0F > 0x09 || c.getHFlag() {
			offset = 0xFA // adjust for subtraction in BCD
		}
		// upper nibble correction
		if c.A&0xF0 > 0x90 || c.getCFlag() {
			offset |= 0xA0 // adjust for subtraction in BCD
			c.setCFlag()   // set the carry flag for subtraction
		}
		// apply the correction
		c.A -= offset
	}
	// update Z flag
	if c.A == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	// reset the H flag
	c.resetHFlag()
	// N flag is not modified
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

func (c *CPU) DEC(instruction *Instruction) {
	panic("DEC not implemented")
}
func (c *CPU) INC(instruction *Instruction) {
	panic("INC not implemented")
}
func (c *CPU) SUB(instruction *Instruction) {
	panic("SUB not implemented")
}
func (c *CPU) SBC(instruction *Instruction) {
	panic("SBC not implemented")
}

/*
 * SCF: Set Carry Flag
 * opcodes: 0x37=SCF
 * flags: Z:- N:0 H:0 C:1
 */
func (c *CPU) SCF(instruction *Instruction) {
	c.F |= 0b00010000

	// reset the N and H flags and leave the Z flag unchanged
	c.resetNFlag()
	c.resetHFlag()
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}
func (c *CPU) OR(instruction *Instruction) {
	panic("OR not implemented")
}

/*
	XOR: Bitwise XOR
	opcodes:
		- 0xA8 = XOR A, B
		- 0xA9 = XOR A, C
		- 0xAA = XOR A, D
		- 0xAB = XOR A, E
		- 0xAC = XOR A, H
		- 0xAD = XOR A, L
		- 0xAE = XOR A, [HL]
		- 0xAF = XOR A, A
		- 0xEE = XOR A, n8
	flags: Z->Z N->0 H->0 C->0
	note: 0xAF XOR 0xAF = 0x00 (Z flag is always set)
*/
func (c *CPU) XOR(instruction *Instruction) {
	c.A = c.A ^ uint8(c.Operand)
	c.incrementPC(uint16(instruction.Bytes))
}

// Shift / Rotate and Bit instructions

/*
 * RLA: Rotate A left through carry
 * opcodes: 0x17=RLA
 * flags: Z:0 N:0 H:0 C:C = bit 7 of A before rotation
 */
func (c *CPU) RLA(instruction *Instruction) {
	// save the carry flag value
	carry := c.getCFlag()

	// update the carry flag with accumulator MSB
	if c.A&0x80 == 0x80 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}

	// rotate the accumulator left and replace LSB with old carry value
	if carry {
		c.A = c.A<<1 | 0x01
	} else {
		c.A = c.A << 1
	}

	// update flags
	c.resetZFlag()
	c.resetNFlag()
	c.resetHFlag()

	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 * RLCA: Rotate A left
 * opcodes: 0x07=RLCA
 * flags: Z:0 N:0 H:0 C:C (bit 7 of A)
 */
func (c *CPU) RLCA(instruction *Instruction) {
	// update the carry flag with accumulator MSB
	if c.A&0x80 == 0x80 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	// rotate the accumulator left
	c.A = (c.A << 1) | (c.A >> 7)

	// update flags
	c.resetZFlag()
	c.resetNFlag()
	c.resetHFlag()

	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 * RRA: Rotate A right through carry
 * opcodes: 0x17=RRA
 * flags: Z:0 N:0 H:0 C:C = bit 0 of A before rotation
 */
func (c *CPU) RRA(instruction *Instruction) {
	// save the carry flag value
	carry := c.getCFlag()

	// update the carry flag with accumulator MSB
	if c.A&0x01 == 0x01 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}

	// rotate the accumulator left and replace LSB with old carry value
	if carry {
		c.A = c.A>>1 | 0x80
	} else {
		c.A = c.A >> 1
	}

	// update flags
	c.resetZFlag()
	c.resetNFlag()
	c.resetHFlag()

	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

/*
 * RRCA: Rotate A right
 * opcodes: 0x0F=RRCA
 * flags: Z:0 N:0 H:0 C:C (bit 0 of A)
 */
func (c *CPU) RRCA(instruction *Instruction) {
	// update the carry flag with Accumulator LSB
	if c.A&0x01 == 0x01 {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	// rotate the accumulator right
	c.A = (c.A >> 1) | (c.A << 7)
	// update flags
	c.resetZFlag()
	c.resetNFlag()
	c.resetHFlag()
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}

// Illegal instructions
/*
 panic when an illegal instruction is encountered
 opcodes: 0xD3, 0xDB, 0xDD, 0xE3, 0xE4, 0xEB, 0xEC, 0xED, 0xF4, 0xFC, 0xFD
*/
func (c *CPU) ILLEGAL(instruction *Instruction) {
	err := fmt.Sprintf("Illegal instruction encountered: 0x%02X=%v", c.IR, instruction.Mnemonic)
	panic(err)
}
