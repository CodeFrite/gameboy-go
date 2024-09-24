package gameboy

import (
	"fmt"
	"strings"
)

// > instructions handlers (NO PREFIX)

// Route the execution to the corresponding instruction handler
func (c *CPU) executeInstruction(instruction Instruction) {
	// Execute the corresponding instruction
	switch instruction.Mnemonic {
	case "NOP":
		c.NOP(&instruction)
	case "STOP":
		c.STOP(&instruction)
	case "HALT":
		c.HALT(&instruction)
	case "DI":
		c.DI(&instruction)
	case "EI":
		c.EI(&instruction)
	case "JP":
		c.JP(&instruction)
	case "JR":
		c.JR(&instruction)
	case "CALL":
		c.CALL(&instruction)
	case "RET":
		c.RET(&instruction)
	case "RETI":
		c.RETI(&instruction)
	case "RST":
		c.RST(&instruction)
	case "LD":
		c.LD(&instruction)
	case "LDH":
		c.LDH(&instruction)
	case "PUSH":
		c.PUSH(&instruction)
	case "POP":
		c.POP(&instruction)
	case "ADD":
		c.ADD(&instruction)
	case "ADC":
		c.ADC(&instruction)
	case "AND":
		c.AND(&instruction)
	case "INC":
		c.INC(&instruction)
	case "CCF":
		c.CCF(&instruction)
	case "CP":
		c.CP(&instruction)
	case "CPL":
		c.CPL(&instruction)
	case "DAA":
		c.DAA(&instruction)
	case "DEC":
		c.DEC(&instruction)
	case "SUB":
		c.SUB(&instruction)
	case "SBC":
		c.SBC(&instruction)
	case "SCF":
		c.SCF(&instruction)
	case "OR":
		c.OR(&instruction)
	case "XOR":
		c.XOR(&instruction)
	case "RLA":
		c.RLA(&instruction)
	case "RLCA":
		c.RLCA(&instruction)
	case "RRA":
		c.RRA(&instruction)
	case "RRCA":
		c.RRCA(&instruction)
	default:
		// Handle illegal instructions first
		if strings.HasPrefix(instruction.Mnemonic, "ILLEGAL_") {
			c.ILLEGAL(&instruction)
		} else {
			err := fmt.Sprintf("Unknown instruction: 0x%02X= %s", c.IR, instruction.Mnemonic)
			panic(err)
		}
	}
}

// Misc / Control instructions

/*
Disable Interrupts (DI)
Disables the IME flag to prevent the CPU from responding to interrupts
opcodes: 0xF3
flags: -
*/
func (c *CPU) DI(instruction *Instruction) {
	// ask the CPU to disable interrupts after the next instruction
	c.IME_DISABLE_NEXT_CYCLE = true
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
}

/*
Enable Interrupts (EI)
Enables the IME flag to allow the CPU to respond to interrupts
Does not enable interrupts immediately, the next instruction will be executed before the interrupts are enabled
opcodes: 0xFB
flags: -
*/
func (c *CPU) EI(instruction *Instruction) {
	// ask the CPU to enable interrupts after the next instruction
	c.IME_ENABLE_NEXT_CYCLE = true
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
}

/*
 * HALT: Halt the CPU until an interrupt occurs
 * opcodes: 0x76
 * flags: -
 */
func (c *CPU) HALT(instruction *Instruction) {
	c.Halted = true
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
}

/*
NOP: No operation, does nothing apart from incrementing the program counter
opcodes: 0x00=NOP
flags impacted: -
*/
func (c *CPU) NOP(instruction *Instruction) {
	// do nothing (no business logic)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
}
func (c *CPU) STOP(instruction *Instruction) {
	// stop the CPU
	c.Stopped = true
	// TODO: update the 0xFF04 register (DIV) to 0

	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
}

// Jump / Call instructions
// only the conditional instructions increment the PC if the condition is not met since they are meant to position the PC at the operand
/*
	CALL: Call a subroutine = if condition is met, push the address of the next instruction to the stack and jump to the address
	Otherwise, continue with the next instruction
	opcodes:
		- 0xC4 = CALL NZ, a16
		- 0xCC = CALL Z, a16
		- 0xCD = CALL a16
		- 0xD4 = CALL NC, a16
		- 0xDC = CALL C, a16
	flags: -
*/
func (c *CPU) CALL(instruction *Instruction) {
	offset := c.PC + uint16(instruction.Bytes)
	switch instruction.Operands[0].Name {
	case "Z":
		if c.getZFlag() {
			c.push(offset) // push the address of the next instruction onto the stack
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "NZ":
		if !c.getZFlag() {
			c.push(offset)
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "C":
		if c.getCFlag() {
			c.push(offset)
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "NC":
		if !c.getCFlag() {
			c.push(offset)
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "a16":
		c.push(offset)
		// update the number of cycles executed by the CPU
		c.CpuCycles += instruction.Cycles[0]
		c.Offset = c.Operand
	default:
		panic("CALL: unknown operand")
	}
}

/*
JP: Jumps to an address
opcodes:
  - 0xC3 = JP     a16
  - 0xE9 = JP HL
  - 0xCA = JP  Z, a16
  - 0xC2 = JP NZ, a16
  - 0xDA = JP  C, a16
  - 0xD2 = JP NC, a16

flags: -
*/
func (c *CPU) JP(instruction *Instruction) {
	offset := c.PC + uint16(instruction.Bytes)
	switch instruction.Operands[0].Name {
	case "Z":
		if c.getZFlag() {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "NZ":
		if !c.getZFlag() {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "C":
		if c.getCFlag() {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "NC":
		if !c.getCFlag() {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
			c.Offset = c.Operand
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = offset
		}
	case "a16":
		// update the number of cycles executed by the CPU
		c.CpuCycles += instruction.Cycles[0]
		c.Offset = c.Operand
	case "HL":
		// update the number of cycles executed by the CPU
		c.CpuCycles += instruction.Cycles[0]
		c.Offset = c.Operand
	default:
		panic("JP: unknown operand")
	}
}

/*
JR: Jump relative
Jumps to a relative address from the next instruction
opcodes:
  - 0x18 = JR r8
  - 0x20 = JR NZ, r8
  - 0x28 = JR Z, r8
  - 0x30 = JR NC, r8
  - 0x38 = JR C, r8

flags: -
! the int8 r8 operand has already been casted to uint16, safely because uint16 > int8 retains the sign value
! by converting for -1 to 0xFF, -2 to 0xFE, etc which is the expected behavior
*/
func (c *CPU) JR(instruction *Instruction) {
	offset := uint16(int(c.PC) + int(int8(c.Operand)) + int(instruction.Bytes))
	switch instruction.Operands[0].Name {
	case "Z":
		if c.getZFlag() {
			// update the program counter offset
			c.Offset = offset
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = c.PC + uint16(instruction.Bytes)
		}
	case "NZ":
		if !c.getZFlag() {
			// update the program counter offset
			c.Offset = offset
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = c.PC + uint16(instruction.Bytes)
		}
	case "C":
		if c.getCFlag() {
			// update the program counter offset
			c.Offset = offset
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = c.PC + uint16(instruction.Bytes)
		}
	case "NC":
		if !c.getCFlag() {
			// update the program counter offset
			c.Offset = offset
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[0]
		} else {
			// update the number of cycles executed by the CPU
			c.CpuCycles += instruction.Cycles[1]
			c.Offset = c.PC + uint16(instruction.Bytes)
		}
	case "e8":
		c.CpuCycles += instruction.Cycles[0]
		// update the program counter offset
		c.Offset = offset
	default:
		errMessage := fmt.Sprint("JR: unknown operand, got ", instruction.Operands[0].Name)
		panic(errMessage)
	}
}

/*
RET: Return from a subroutine
This intruction pops the address from the stack and jumps to it
opcodes:
  - 0xC9 = RET
  - 0xC8 = RET Z
  - 0xC0 = RET NZ
  - 0xD8 = RET C
  - 0xD0 = RET NC

flags: -
*/
func (c *CPU) RET(instruction *Instruction) {
	if len(instruction.Operands) == 0 {
		c.Offset = c.pop()
	} else {
		switch instruction.Operands[0].Name {
		case "Z":
			if c.getZFlag() {
				c.Offset = c.pop()
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[0]
			} else {
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[1]
			}
		case "NZ":
			if !c.getZFlag() {
				c.Offset = c.pop()
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[0]
			} else {
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[1]
			}
		case "C":
			if c.getCFlag() {
				c.Offset = c.pop()
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[0]
			} else {
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[1]
			}
		case "NC":
			if !c.getCFlag() {
				c.Offset = c.pop()
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[0]
			} else {
				// update the number of cycles executed by the CPU
				c.CpuCycles += instruction.Cycles[1]
			}
		default:
			panic("RET: unknown operand")
		}
	}
}

/*
RETI: Return from interrupt
Return from subroutine and enable interrupts.
This is basically equivalent to executing EI then RET, meaning that IME is set right after this instruction.
opcodes: 0xD9
flags: -
*/
func (c *CPU) RETI(instruction *Instruction) {
	c.Offset = c.pop()
	c.IME = true
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
RST: Restart
Restart the CPU at a specific address by pushing the current address to the stack and jumping to the specified address
opcodes:
  - 0xC7 = RST $00
  - 0xCF = RST $08
  - 0xD7 = RST $10
  - 0xDF = RST $18
  - 0xE7 = RST $20
  - 0xEF = RST $28
  - 0xF7 = RST $30
  - 0xFF = RST $38

flags: -
*/
func (c *CPU) RST(instruction *Instruction) {
	c.push(c.PC + uint16(instruction.Bytes))
	c.Offset = c.Operand

	// Update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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
flags: - (except for 0xF8 where Z->0 N->0 H->H C->C)

NOTE: all LD instructions have 2 operands, the first one is always the destination and the second one is always the source (except for LD HL, SP+r8)
=> we will 'automate' the process of fetching the operands expect for LD HL, SP+r8 that will be handled manually
*/
func (c *CPU) LD(instruction *Instruction) {
	var address uint16
	var err error

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
			err = c.bus.Write(address, uint8(c.Operand))
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
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
			err := c.bus.Write(c.getBC(), uint8(c.Operand))
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
		}
	case "DE":
		if instruction.Operands[0].Immediate {
			c.setDE(c.Operand)
		} else {
			err := c.bus.Write(c.getDE(), uint8(c.Operand))
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
		}
	case "HL":
		if instruction.Operands[0].Immediate {
			// LD HL, SP+e8 (0xF8)
			if len(instruction.Operands) == 3 {
				newValue := c.SP + c.Operand
				// set or reset the H flag
				if newValue > 0x0F {
					c.setHFlag()
				} else {
					c.resetHFlag()
				}
				// set or reset the C flag
				if newValue > 0xFF {
					c.setCFlag()
				} else {
					c.resetCFlag()
				}
				// load the result into HL
				c.setHL(newValue)
				// update flags
				c.resetZFlag()
				c.resetNFlag()
			} else {
				c.setHL(c.Operand)
				// no flags are impacted
			}
		} else {
			err = c.bus.Write(c.getHL(), uint8(c.Operand))
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
			// no flags are impacted
		}
		if instruction.Operands[0].Increment {
			c.setHL(c.getHL() + 1)
		} else if instruction.Operands[0].Decrement {
			c.setHL(c.getHL() - 1)
		}
	case "SP":
		c.SP = c.Operand
	case "a16":
		low := c.bus.Read(c.PC + 1)
		high := c.bus.Read(c.PC + 2)
		addr := uint16(high)<<8 | uint16(low)
		err = c.bus.Write(addr, uint8(c.SP))
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
		err = c.bus.Write(addr+1, uint8(c.SP>>8))
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	default:
		panic("LD: unknown operand")
	}

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)

	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]

}

/*
LDH: Load data from memory address 0xFF00+a8 to A or from A to memory address 0xFF00+a8
opcodes:
  - 0xE0 = LDH [a8], A
  - 0xF0 = LDH A, [a8]

flags: -
*/
func (c *CPU) LDH(instruction *Instruction) {
	var err error

	switch instruction.Operands[0].Name {
	case "A":
		c.A = uint8(c.Operand)
	case "a8":
		a8 := 0xFF00 + c.Operand
		err = c.bus.Write(a8, c.A)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	default:
		panic("LDH: unknown operand")
	}
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
PUSH: Push a 16-bit register pair onto the stack
opcodes:
  - 0xC5 = PUSH BC
  - 0xD5 = PUSH DE
  - 0xE5 = PUSH HL
  - 0xF5 = PUSH AF

flags: -
*/
func (c *CPU) PUSH(instruction *Instruction) {
	c.SP--
	switch instruction.Operands[0].Name {
	case "AF":
		err := c.bus.Write(c.SP, c.A)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
		c.SP--
		err = c.bus.Write(c.SP, c.F)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	case "BC":
		err := c.bus.Write(c.SP, c.B)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
		c.SP--
		err = c.bus.Write(c.SP, c.C)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	case "DE":
		err := c.bus.Write(c.SP, c.D)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
		c.SP--
		err = c.bus.Write(c.SP, c.E)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	case "HL":
		err := c.bus.Write(c.SP, c.H)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
		c.SP--
		err = c.bus.Write(c.SP, c.L)
		if err != nil {
			fmt.Printf("\n> Panic @0x%04X\n", c.PC)
			panic(err)
		}
	default:
		panic("PUSH: unknown operand")
	}
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
POP: Pop a 16-bit register pair from the stack
opcodes:
  - 0xC1 = POP BC
  - 0xD1 = POP DE
  - 0xE1 = POP HL
  - 0xF1 = POP AF (flags are restored from the stack)

flags: - except for 0xF1 where Z->Z N->N H->H C->C
*/
func (c *CPU) POP(instruction *Instruction) {
	low := c.bus.Read(c.SP)
	c.SP++
	high := c.bus.Read(c.SP)
	c.SP++
	switch instruction.Operands[0].Name {
	case "AF":
		c.A = high
		c.F = low
	case "BC":
		c.B = high
		c.C = low
	case "DE":
		c.D = high
		c.E = low
	case "HL":
		c.H = high
		c.L = low
	default:
		panic("POP: unknown operand")
	}
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
	 CP: compare 2 memory locations and/or registers by subtracting them without storing the result
	 opcodes:
		- B8 = CP A, B
		- B9 = CP A, C
		- BA = CP A, D
		- BB = CP A, E
		- BC = CP A, H
		- BD = CP A, L
		- BE = CP A, [HL]
		- BF = CP A, A
		- FE = CP A, n8
	 flags: Z:Z N:1 H:H C:C
*/
func (c *CPU) CP(instruction *Instruction) {
	val := c.A - uint8(c.Operand)
	// update flags
	if val == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.setNFlag()
	if c.A&0x0F < uint8(c.Operand)&0x0F {
		c.setHFlag()
	} else {
		c.resetHFlag()
	}
	if uint8(c.Operand) > c.A {
		c.setCFlag()
	} else {
		c.resetCFlag()
	}
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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
	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
	 DEC: Decrement register or memory
	 opcodes:
		- 0x05=DEC B
		- 0x0B=DEC BC
		- 0x0D=DEC C
		- 0x15=DEC D
		- 0x1B=DEC DE
		- 0x1D=DEC E
		- 0x25=DEC H
		- 0x2B=DEC HL
		- 0x2D=DEC L
		- 0x35=DEC [HL]
		- 0x3B=DEC SP
		- 0x3D=DEC A

flags: Z:Z N:1 H:H C:- for all but the 16-bits registers

When to set H ? There will be a borrow from bit 4 if the lower nibble is 0
*/
func (c *CPU) DEC(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "A":
		// check H before DEC
		if c.A&0x0F == 0x00 {
			c.setHFlag()
		}
		c.A--
		if c.A == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "B":
		if c.B&0x0F == 0x00 {
			c.setHFlag()
		}
		c.B--
		if c.B == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "C":
		if c.C&0x0F == 0x00 {
			c.setHFlag()
		}
		c.C--
		if c.C == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "D":
		if c.D&0x0F == 0x00 {
			c.setHFlag()
		}
		c.D--
		if c.D == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "E":
		if c.E&0x0F == 0x00 {
			c.setHFlag()
		}
		c.E--
		if c.E == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "H":
		if c.H&0x0F == 0x00 {
			c.setHFlag()
		}
		c.H--
		if c.H == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "L":
		if c.L&0x0F == 0x00 {
			c.setHFlag()
		}
		c.L--
		if c.L == 0x00 {
			c.setZFlag()
		}
		c.setNFlag()
	case "BC":
		c.setBC(c.getBC() - 1)
	case "DE":
		c.setDE(c.getDE() - 1)
	case "HL":
		if instruction.Operands[0].Immediate {
			c.setHL(c.getHL() - 1)
		} else {
			addr := c.getHL()
			val := c.bus.Read(addr)
			err := c.bus.Write(addr, val-1)
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
		}
	case "SP":
		c.SP--
	default:
		panic("DEC: unknown operand")
	}

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
}

/*
	 INC: Increment register or memory
	 opcodes:
	 	- 0x04=INC B
		- 0x0C=INC C
		- 0x14=INC D
		- 0x1C=INC E
		- 0x24=INC H
		- 0x2C=INC L
		- 0x34=INC [HL]
		- 0x3C=INC A

	 flags: Z:Z N:0 H:H C:-
*/
func (c *CPU) INC(instruction *Instruction) {
	switch instruction.Operands[0].Name {
	case "A":
		c.A++
		if c.A == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.A & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "B":
		c.B++
		if c.B == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.B & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "C":
		c.C++
		if c.C == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.C & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "D":
		c.D++
		if c.D == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.D & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "E":
		c.E++
		if c.E == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.E & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "H":
		c.H++
		if c.H == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.H & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "L":
		c.L++
		if c.L == 0x00 {
			c.setZFlag()
		} else {
			c.resetZFlag()
		}
		if (c.L & 0x0F) == 0x00 {
			c.setHFlag()
		} else {
			c.resetHFlag()
		}
		// reset the N flag
		c.resetNFlag()
	case "HL":
		if instruction.Operands[0].Immediate {
			c.setHL(c.getHL() + 1)
		} else if instruction.Operands[0].Increment {
			addr := c.getHL()
			val := c.bus.Read(addr) + 1
			err := c.bus.Write(addr, val)
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}

			if val == 0x00 {
				c.setZFlag()
			} else {
				c.resetZFlag()
			}
			if (val & 0x0F) == 0x00 {
				c.setHFlag()
			} else {
				c.resetHFlag()
			}
			// reset the N flag
			c.resetNFlag()
		}
	case "DE":
		if instruction.Operands[0].Immediate {
			c.setDE(c.getDE() + 1)
		} else if instruction.Operands[0].Increment {
			addr := c.getDE()
			val := c.bus.Read(addr) + 1
			err := c.bus.Write(addr, val)
			if err != nil {
				fmt.Printf("\n> Panic @0x%04X\n", c.PC)
				panic(err)
			}
		}
	default:
		panic(fmt.Sprintf(">> PANIC >> INC instruction: unknown operand %s", instruction.Operands[0].Name))
	}

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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
	// update flags
	if c.A == 0x00 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	// reset N, H and C flags
	c.resetNFlag()
	c.resetHFlag()
	c.resetCFlag()

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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

	// update the program counter offset
	c.Offset = c.PC + uint16(instruction.Bytes)
	// update the number of cycles executed by the CPU
	c.CpuCycles += instruction.Cycles[0]
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
