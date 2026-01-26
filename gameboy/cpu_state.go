package gameboy

import (
	"fmt"
)

type CpuState struct {
	// Special registers
	CPU_CYCLES uint64 `json:"CPU_CYCLES"` // number of cycles the CPU has executed TODO: change to the correct type and implement the interrupt (overflow) handling
	PC         uint16 `json:"PC"`         // Program Counter
	SP         uint16 `json:"SP"`         // Stack Pointer
	A          uint8  `json:"A"`          // Accumulator
	F          uint8  `json:"F"`          // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	Z          bool   `json:"Z"`          // Zero flag
	N          bool   `json:"N"`          // Subtraction flag
	H          bool   `json:"H"`          // Half Carry flag
	C          bool   `json:"C"`          // Carry flag
	// 16-bits general purpose registers
	BC uint16 `json:"BC"`
	DE uint16 `json:"DE"`
	HL uint16 `json:"HL"`

	// Instruction
	INSTRUCTION   Instruction `json:"INSTRUCTION"`   // Current instruction
	PREFIXED      bool        `json:"PREFIXED"`      // Is the current instruction prefixed with 0xCB
	IR            uint8       `json:"IR"`            // Instruction Register
	OPERAND_VALUE uint16      `json:"OPERAND_VALUE"` // Current operand fetched from memory (this register doesn't physically exist in the CPU)

	// Interrupts
	IE  uint8 `json:"IE"`  // Interrupt Enable
	IME bool  `json:"IME"` // interrupt master enable

	// emulator state
	HALTED  bool `json:"HALTED"`  // is the CPU halted
	STOPPED bool `json:"STOPPED"` // is the CPU stopped
}

// get the memories current content
func (c *CPU) getState() CpuState {
	return CpuState{
		CPU_CYCLES:    c.cpuCycles,
		PC:            c.pc,
		SP:            c.sp,
		A:             c.a,
		F:             c.f,
		Z:             c.f&0x80 != 0,
		N:             c.f&0x40 != 0,
		H:             c.f&0x20 != 0,
		C:             c.f&0x10 != 0,
		BC:            uint16(c.b)<<8 | uint16(c.c),
		DE:            uint16(c.d)<<8 | uint16(c.e),
		HL:            uint16(c.h)<<8 | uint16(c.l),
		INSTRUCTION:   c.instruction,
		PREFIXED:      c.prefixed,
		IR:            c.ir,
		OPERAND_VALUE: c.operand,
		IE:            c.bus.Read(0xFFFF),
		IME:           c.ime,
		HALTED:        c.halted,
		STOPPED:       c.stopped,
	}
}

func (cs *CpuState) print() {
	fmt.Println("")
	fmt.Println("\n> CPU State:")
	fmt.Println("------------")
	fmt.Printf("Cycles: %d\n", cs.CPU_CYCLES)
	fmt.Printf("PC: 0x%04X\n", cs.PC)
	fmt.Printf("SP: 0x%04X\n", cs.SP)
	fmt.Printf("A: 0x%02X\n", cs.A)
	fmt.Printf("F: 0x%02X\n", cs.F)
	fmt.Printf("Z: %t\n", cs.F&0x80 != 0)
	fmt.Printf("N: %t\n", cs.F&0x40 != 0)
	fmt.Printf("H: %t\n", cs.F&0x20 != 0)
	fmt.Printf("C: %t\n", cs.F&0x10 != 0)
	fmt.Printf("BC: 0x%04X\n", cs.BC)
	fmt.Printf("DE: 0x%04X\n", cs.DE)
	fmt.Printf("HL: 0x%04X\n", cs.HL)
	fmt.Printf("PREFIXED: %t\n", cs.PREFIXED)
	fmt.Printf("IR: 0x%02X\n", cs.IR)
	fmt.Printf("INSTRUCTION: %s\n", cs.INSTRUCTION.Mnemonic)
	fmt.Printf("OPERAND_VALUE: 0x%02X\n", cs.OPERAND_VALUE)
	fmt.Printf("IE: 0x%02X\n", cs.IE)
	fmt.Printf("IME: %t\n", cs.IME)
	fmt.Printf("HALTED: %t\n", cs.HALTED)
	fmt.Printf("STOPPED: %t\n", cs.STOPPED)
	fmt.Println("")
}
