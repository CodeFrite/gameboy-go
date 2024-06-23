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

	IR       uint8  // Instruction Register
	prefixed bool   // Is the current instruction prefixed with 0xCB
	Operand  uint16 // Current operand fetched from memory (this register doesn't physically exist in the CPU)

	HRAM [127]byte // 127bytes of High-Speed RAM
	bus  *Bus      // reference to the bus

	IME    bool // interrupt master enable
	halted bool
}

// Create a new CPU
func NewCPU(bus *Bus) *CPU {
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
	return uint16(c.B)<<8 | uint16(c.C)
}

func (c *CPU) setBC(value uint16) {
	c.B = byte(value >> 8)
	c.C = byte(value)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

func (c *CPU) setDE(value uint16) {
	c.D = byte(value >> 8)
	c.E = byte(value)
}

func (c *CPU) setHL(value uint16) {
	c.H = byte(value >> 8)
	c.L = byte(value)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

// Stack operations

// Push a value to the stack
func (c *CPU) push(value uint16) {
	// decrement the stack pointer
	c.SP -= 1
	// write the high byte to the stack
	c.bus.Write(c.SP, byte(value>>8))
	// decrement the stack pointer
	c.SP -= 1
	// write the low byte to the stack
	c.bus.Write(c.SP, byte(value))
}

// Pop a value to the stack
func (c *CPU) pop() uint16 {
	// pop the low byte to the stack
	low := c.bus.Read(c.SP)
	// increment the stack pointer
	c.SP += 1
	// write the high byte to the stack
	high := c.bus.Read(c.SP)
	// increment the stack pointer
	c.SP += 1
	return uint16(high)<<8 | uint16(low)
}

// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() {
	// Fetch the opcode from memory at the address in the program counter
	opcode := c.bus.Read(c.PC)

	// is it a prefixed instruction?
	if opcode == 0xCB {
		c.prefixed = true
		// fetch the next opcode
		opcode = c.bus.Read(c.PC + 1)
	}

	// Store the opcode in the instruction register
	c.IR = opcode
}

/*
 * Fetch the value of an operand
 * Save the result in cpu.Operand as an uint16 (must be casted to the correct type inside the different instruction handlers)
 */
func (c *CPU) fetchOperandValue(operand Operand) {
	var value, addr uint16
	switch operand.Name {
	case "n8": // always immediate
		value = uint16(c.bus.Read(c.PC + 1))
	case "n16": // always immediate
		// little-endian
		value = c.bus.Read16(c.PC + 1)
	case "a8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.PC + 1))
		} else {
			addr = 0xFF00 | uint16(c.bus.Read(c.PC+1))
			value = uint16(c.bus.Read(addr))
		}
	case "a16": // not always immediate
		if operand.Immediate {
			value = c.bus.Read16(c.PC + 1)
		} else {
			addr := c.bus.Read16(c.PC + 1)
			value = c.bus.Read16(addr)
		}
	case "e8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.PC + 1))
		} else {
			panic("e8 non immediate operand not implemented yet")
		}
	case "A":
		if operand.Immediate {
			value = uint16(c.A)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "B":
		if operand.Immediate {
			value = uint16(c.B)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "C":
		if operand.Immediate {
			value = uint16(c.C)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "D":
		if operand.Immediate {
			value = uint16(c.D)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "E":
		if operand.Immediate {
			value = uint16(c.E)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "H":
		if operand.Immediate {
			value = uint16(c.H)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "L":
		if operand.Immediate {
			value = uint16(c.L)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	default:
		err := fmt.Sprintf("Unknown operand name: %s (0x%02X)", operand.Name, c.IR)
		panic(err)
	}
	// save the current operand value into the cpu context
	c.Operand = value
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

/*
 * Print the current instruction being executed in the following format (examples):
 * PC: 0x00A0, Bytes: 00 			, ASM: NOP
 * PC: 0x00A1, Bytes: 40 			, ASM: LD B, $40
 * PC: 0x00A2, Bytes: 3E 01 	, ASM: LD A, $01
 * PC: 0x00A4, Bytes: F8 4E 	, ASM: LD HL, SP + $4E
 * PC: 0x00A6, Bytes: EA AB 01, ASM: LD [$01AB], HL
 */
func (c *CPU) printCurrentInstruction() {
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), c.prefixed)
	getBytes := func() []byte {
		var bytes []byte
		for i := 0; i < instruction.Bytes; i++ {
			bytes = append(bytes, c.bus.Read(c.PC+uint16(i)))
		}
		return bytes
	}

	getOperands := func() string {
		var operands []string
		for _, operand := range instruction.Operands {
			var value string
			if operand.Name == "n8" {
				value = fmt.Sprintf("$%02X", c.bus.Read(c.PC+1))
			} else if operand.Name == "n16" {
				value = fmt.Sprintf("$%04X", c.bus.Read16(c.PC+1))
			} else {
				value = operand.Name
			}

			if operand.Increment {
				value += "+"
			} else if operand.Decrement {
				value += "-"
			}

			if !operand.Immediate {
				value = "[" + value + "]"
			}

			operands = append(operands, value)
		}
		return strings.Join(operands, ", ")
	}

	fmt.Printf("PC: 0x%04X, SP: 0x%04X", c.PC, c.SP)
	fmt.Printf(", memory: %-6X", getBytes())
	fmt.Printf(", asm: %s %s\n", instruction.Mnemonic, getOperands())
}

func (c *CPU) printRegisters() {
	fmt.Printf("A: 0x%02X, B: 0x%02X, C: 0x%02X, D: 0x%02X, E: 0x%02X, H: 0x%02X, L: 0x%02X", c.A, c.B, c.C, c.D, c.E, c.H, c.L)
	fmt.Printf(", Z: %t, N: %t, H: %t, C: %t\n", c.getZFlag(), c.getNFlag(), c.getHFlag(), c.getCFlag())
}

// Execute one cycle of the CPU: fetch, decode and execute the next instruction
// TODO: i am supposed to return an error but i am always returning nil. Chose an error handling strategy and implement it
func (c *CPU) step() error {
	// 0. reset the prefixed flag
	c.prefixed = false

	// 1. Fetch the opcode from memory and save it to the instruction register IR
	c.fetchOpcode()

	// 2. Decode the instruction
	// get instruction from opcodes.json file with IR used as key
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), c.prefixed)
	// get the operands of the instruction
	operands := instruction.Operands
	// fetch the operand value
	if len(operands) == 1 {
		c.fetchOperandValue(operands[0])
	} else if len(operands) == 2 {
		// decode operand 2
		c.fetchOperandValue(operands[1])
	}

	// debug
	c.printCurrentInstruction()

	// 3. Execute the instruction
	if !c.prefixed {
		c.executeInstruction(instruction)
	} else {
		c.executeCBInstruction(instruction)
	}

	// debug
	c.printRegisters()
	fmt.Println()

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
