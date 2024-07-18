// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"fmt"
)

/*
 * CPU: executes instructions fetched from memory, reads and writes to memory (internal registers, flags & bus)
 */
type CPU struct {
	// Registers
	PC               uint16 // Program Counter
	SP               uint16 // Stack Pointer
	A                uint8  // Accumulator
	F                uint8  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	B, C, D, E, H, L uint8  // 16-bit general purpose registers

	// Instruction
	IR       uint8  // Instruction Register
	Prefixed bool   // Is the current instruction prefixed with 0xCB
	Operand  uint16 // Current operand fetched from memory (this register doesn't physically exist in the CPU)
	offset   uint16 // offset used in some instructions

	// Memory
	bus *Bus // reference to the bus

	// Interrupts
	IE     *Memory // Interrupt Enable
	IME    bool    // interrupt master enable
	halted bool
}

// Create a new CPU
func NewCPU(bus *Bus) *CPU {
	return &CPU{
		bus: bus,
		IE:  NewMemory(1),
	}
}

// Increment the Program Counter by the given offset
func (c *CPU) incrementPC() {
	c.PC = c.offset
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
func (c *CPU) fetchOpcode() (opcode uint8, prefixed bool) {
	// Fetch the opcode from memory at the address in the program counter
	opcode = c.bus.Read(c.PC)

	// is it a prefixed instruction?
	if opcode == 0xCB {
		prefixed = true
		// fetch the next opcode
		opcode = c.bus.Read(c.PC + 1)
	}
	return opcode, prefixed
}

/*
 * Fetch the value of an operand
 * Save the result in cpu.Operand as an uint16 (must be casted to the correct type inside the different instruction handlers)
 */
func (c *CPU) fetchOperandValue(operand Operand) uint16 {
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
			addr = 0xFF00 + uint16(c.bus.Read(c.PC+1))
			value = uint16(c.bus.Read(addr))
			fmt.Println("addr", addr)
			fmt.Println("a8", value)
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

	case "BC":
		if operand.Immediate {
			value = c.getBC()
		} else {
			value = c.bus.Read16(c.getBC())
		}
	case "DE":
		if operand.Immediate {
			value = c.getDE()
		} else {
			value = c.bus.Read16(c.getDE())
		}
	case "HL":
		if operand.Immediate {
			value = c.getHL()
		} else {
			value = c.bus.Read16(c.getHL())
		}
	case "SP": // always immediate
		value = c.SP
	case "$00": // RST $00
		value = 0x00
	case "$08": // RST $08
		value = 0x08
	case "$10": // RST $10
		value = 0x10
	case "$18": // RST $18
		value = 0x18
	case "$20": // RST $20
		value = 0x20
	case "$28": // RST $28
		value = 0x28
	case "$30": // RST $30
		value = 0x30
	case "$38": // RST $38
		value = 0x38

	default:
		err := fmt.Sprintf("Unknown operand name: %s (0x%02X)", operand.Name, c.IR)
		panic(err)
	}
	return value
}

// Execute one cycle of the CPU: fetch, decode and execute the next instruction
// TODO: i am supposed to return an error but i am always returning nil. Chose an error handling strategy and implement it
func (c *CPU) step() error {
	// update the pc
	c.incrementPC()
	// reset the offset
	c.offset = 0

	// 0. reset the prefixed flag
	c.Prefixed = false

	// 1. Store the opcode in the instruction register & prefix flag
	opCode, prefixed := c.fetchOpcode()
	c.Prefixed = prefixed
	c.IR = opCode

	// 2. Decode the instruction
	// get instruction from opcodes.json file with IR used as key
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), c.Prefixed)
	// get the operands of the instruction
	operands := instruction.Operands
	// fetch the operand value
	if len(operands) == 1 {
		c.Operand = c.fetchOperandValue(operands[0])
	} else if len(operands) == 2 {
		// decode operand 2
		c.Operand = c.fetchOperandValue(operands[1])
	}

	// 3. Execute the instruction
	if !c.Prefixed {
		c.executeInstruction(instruction)
	} else {
		c.executeCBInstruction(instruction)
	}

	if c.offset == 0 {
		c.offset = c.PC + uint16(instruction.Bytes)
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

// Boot the CPU and returns when the boot process is done (PC=0x0100)
func (c *CPU) Boot() {
	c.PC = 0x0000
	for c.PC != 0x0100 {
		c.step()
	}
}
