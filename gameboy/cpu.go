// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"fmt"
	"strings"
)

/*
 * CPU: executes instructions fetched from memory, reads and writes to memory (internal registers, flags & bus)
 */
type CPU struct {
	PC               uint16 // Program Counter
	SP               uint16 // Stack Pointer
	A                uint8  // Accumulator
	F                uint8  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	B, C, D, E, H, L uint8  // 16-bit general purpose registers
	IE               uint8  // Interrupt Enable

	IR      uint8  // Instruction Register
	Operand uint16 // Current operand fetched from memory (this register doesn't physically exist in the CPU)

	HRAM [127]byte // 127bytes of High-Speed RAM
	bus  *Bus      // reference to the bus

	IME    bool // interrupt master enable
	halted bool
}

func NewCPU(bus *Bus) *CPU {
	// initialize all registers to 0 except the program counter which starts at 0x100 (in cartridge ROM)
	return &CPU{
		bus: bus,
	}
}

/*
 * F flags register
 * 7 6 5 4 3 2 1 0 (position)
 * Z N H C 0 0 0 0 (flag)
 */

// Zero Flag operations
// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return c.F&0x80 == 0x80
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.F = c.F | 0x80
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.F = c.F & 0x7F
}

// Toggle the Z flag in the F register
func (c *CPU) toggleZFlag() {
	c.F = c.F ^ 0x80
}

// Carry Flag operations
// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return c.F&0x40 == 0x40
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.F = c.F | 0x40
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.F = c.F & 0xBF
}

// Toggle the N flag in the F register
func (c *CPU) toggleNFlag() {
	c.F = c.F ^ 0x40
}

// Half Carry Flag operations
// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return c.F&0x20 == 0x20
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.F = c.F | 0x20
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.F = c.F & 0xDF
}

// Toggle the H flag in the F register
func (c *CPU) toggleHFlag() {
	c.F = c.F ^ 0x20
}

// Carry Flag operations
// Get the C flag from the F register
func (c *CPU) getCFlag() bool {
	return c.F&0x10 == 0x10
}

// Set the C flag in the F register
func (c *CPU) setCFlag() {
	c.F = c.F | 0x10
}

// Reset the C flag in the F register
func (c *CPU) resetCFlag() {
	c.F = c.F & 0xEF
}

// Toggle the C flag in the F register
func (c *CPU) toggleCFlag() {
	c.F = c.F ^ 0x10
}

/*
 * 16-bit registers accessors
 */
func (c *CPU) getBC() uint16 {
	return uint16(c.B<<8 | c.C)
}

func (c *CPU) setBC(value uint16) {
	c.B = byte(value >> 8)
	c.C = byte(value)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.D<<8 | c.E)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.H<<8 | c.L)
}

// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() {
	// Fetch the opcode from memory at the address in the program counter
	opcode := c.bus.Read(c.PC)

	// Store the opcode in the instruction register
	c.IR = opcode
}

/*
 * Fetch the value of an operand
 * Returns an interface{} that can either be a uint8 or uint16
 */
func (c *CPU) fetchOperandValue(operand Operand) interface{} {
	var value interface{}
	switch operand.Name {
	case "n8": // always immediate
		value = c.bus.Read(c.PC + 1)
	case "n16": // always immediate
		// little-endian
		low := c.bus.Read(c.PC + 1)
		high := c.bus.Read(c.PC + 2)
		value = uint16(high)<<8 | uint16(low)
	case "a8": // not always immediate
		if operand.Immediate {
			value = c.bus.Read(c.PC + 1)
		}
		/* TODO: handle this case where i need to add 0xFF00 to the value (LDH instructions)
		else {
			addr := c.bus.Read(c.PC + 1)
			value = c.bus.Read(addr)
		}
		*/
	case "a16": // not always immediate
		if operand.Immediate {
			low := c.bus.Read(c.PC + 1)
			high := c.bus.Read(c.PC + 2)
			value = uint16(high)<<8 | uint16(low)
		} else {
			low := c.bus.Read(c.PC + 1)
			high := c.bus.Read(c.PC + 2)
			addr := uint16(high)<<8 | uint16(low)
			value = c.bus.Read(addr)
		}
	case "A":
		if operand.Immediate {
			value = c.A
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "B":
		if operand.Immediate {
			value = c.B
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "C":
		if operand.Immediate {
			value = c.C
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "D":
		if operand.Immediate {
			value = c.D
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "E":
		if operand.Immediate {
			value = c.E
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "H":
		if operand.Immediate {
			value = c.H
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "L":
		if operand.Immediate {
			value = c.L
		} else {
			panic("Non immediate operand not implemented yet")
		}
	default:
		panic("Unknown operand type")
	}
	return value
}

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
	case "PREFIX":
		c.PREFIX(&instruction)
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

func (c *CPU) incrementPC(offset uint16) {
	c.PC += uint16(offset)
}

// Execute one cycle of the CPU: fetch, decode and execute the next instruction
// TODO: i am supposed to return an error but i am always returning nil. Chose an error handling strategy and implement it
func (c *CPU) step() error {
	// 1. Fetch the opcode from memory and save it to the instruction register IR
	c.fetchOpcode()

	// 2. Decode the instruction

	// get instruction from opcodes.json file with IR used as key
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), false)

	// get the operands of the instruction
	operands := instruction.Operands

	// handle 0 operands instructions
	if len(operands) == 0 {
		// no operand to decode
		// execute the instruction
		c.executeInstruction(instruction)
	} else if len(operands) == 1 {
		fmt.Println(instruction)
		panic("CPU 1 operand instructions not implemented yet")
	} else if len(operands) == 2 {
		panic("CPU 2 operands instructions not implemented yet")
	}

	return nil
}

// Run the CPU
func (c *CPU) Run() {
	for {
		if c.halted {
			// if the CPU is halted, wiat for an interrupt to wake it up
			// TODO: implement the interrupt handling
			// ! for the moment, we will break the loop to avoid an infinite loop
			break
		}
		// Execute the next instruction
		if err := c.step(); err != nil {
			panic(err)
		}
	}
}
