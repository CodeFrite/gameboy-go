// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"errors"
	"fmt"
)

type CPU struct {
	IR uint8 // Instruction Register
	IE uint8 // Interrupt Enable
	SP uint16 // Stack Pointer
	PC uint16 // Program Counter

	// 8-bit registers
	A uint8 // Accumulator
	F uint8 // Flags (Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4))

	// 16-bit general purpose registers
	BC, DE, HL uint16

	// 127bytes of High-Speed RAM
	HRAM [127]byte

	// reference to the bus
	bus *Bus
}

func NewCPU(bus *Bus) *CPU {
	// initialize all registers to 0 except the program counter which starts at 0x100 (in cartridge ROM)
	return &CPU{
		PC:0x0100, // Start at 0x100 and ignore the boot ROM for now
		bus: bus,
	}
}

func (c *CPU) PrintRegisters() {
	// Print the registers
	fmt.Println("> CPU Registers:")
	fmt.Printf("+ IR: 0x%X\n", c.IR)
	fmt.Printf("+ IE: 0x%X\n", c.IE)
	fmt.Printf("+ SP: 0x%X\n", c.SP)
	fmt.Printf("+ PC: 0x%X\n", c.PC)
	fmt.Printf("+ A: 0x%X\n", c.A)
	fmt.Printf("+ F: 0x%X\n", c.F)
	fmt.Printf("+ BC: 0x%X\n", c.BC)
	fmt.Printf("+ DE: 0x%X\n", c.DE)
	fmt.Printf("+ HL: 0x%X\n", c.HL)
	fmt.Printf("+ HRAM: 0x%X\n", c.HRAM)
}

func (c *CPU) PrintIR() {
	// Print the instruction register
	fmt.Printf("IR: 0x%X\n", c.IR)
}

func (c *CPU) PrintPC() {
	// Print the program counter
	fmt.Printf("PC: 0x%X\n", c.PC)
}

// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return c.F & 0x80 == 0x80
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.F = c.F | 0x80
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.F = c.F & 0x7F
}

// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return c.F & 0x40 == 0x40
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.F = c.F | 0x40
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.F = c.F & 0xBF
}

// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return c.F & 0x20 == 0x20
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.F = c.F | 0x20
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.F = c.F & 0xDF
}

// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() {
	// Fetch the opcode from memory at the address in the program counter
	opcode := c.bus.Read(uint16ToBytes(c.PC))
	// Store the opcode in the instruction register
	c.IR = opcode
}

func (c *CPU) getInstruction() (Instruction, error) {
	if instruction, ok := Instructions[c.IR]; ok {
		return instruction, nil
	}
	// Return a default instruction for unimplemented opcodes
	err := fmt.Sprintf("Unimplemented opcode: 0x%X", c.IR)
	return NotYetImplementedInstruction, errors.New(err)
}

func (c *CPU) executeInstruction(instruction Instruction) {
	// Fetch the operand from memory
	var operand []byte
	if instruction.Length == 2 {
		operand = make([]byte, 1)
		operand[0] = c.bus.Read(uint16ToBytes(c.PC+1))
	} else if instruction.Length == 3 {
		operand = make([]byte, 2)
		operand[0] = c.bus.Read(uint16ToBytes(c.PC+1))
		operand[1] = c.bus.Read(uint16ToBytes(c.PC+2))
	}
	// Execute the instruction
	instruction.Handler(c, operand)
}

func (c *CPU) incrementPC(offset uint16) {
	c.PC += uint16(offset)
}

func (c *CPU) Step() error {
	// Fetch the opcode from memory
	c.fetchOpcode()
	// Get the instruction corresponding to the opcode
	instruction, err := c.getInstruction()
	if err != nil {
		c.PrintPC()
		c.PrintIR()
		fmt.Println(err)
		return err
	}
	
	// Execute the instruction
	c.executeInstruction(instruction)
	
	return nil
}