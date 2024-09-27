package gameboy

import (
	"fmt"
	"reflect"
)

type CpuState struct {
	// Special registers
	CpuCycles uint64 `json:"cpuCycles"` // number of cycles the CPU has executed TODO: change to the correct type and implement the interrupt (overflow) handling
	PC        uint16 `json:"PC"`        // Program Counter
	SP        uint16 `json:"SP"`        // Stack Pointer
	A         uint8  `json:"A"`         // Accumulator
	F         uint8  `json:"F"`         // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	Z         bool   `json:"Z"`         // Zero flag
	N         bool   `json:"N"`         // Subtraction flag
	H         bool   `json:"H"`         // Half Carry flag
	C         bool   `json:"C"`         // Carry flag
	// 16-bits general purpose registers
	BC uint16 `json:"BC"`
	DE uint16 `json:"DE"`
	HL uint16 `json:"HL"`

	// Instruction
	PREFIXED      bool   `json:"prefixed"`     // Is the current instruction prefixed with 0xCB
	IR            uint8  `json:"IR"`           // Instruction Register
	OPERAND_VALUE uint16 `json:"operandValue"` // Current operand fetched from memory (this register doesn't physically exist in the CPU)

	// Interrupts
	IE  uint8 `json:"IE"`  // Interrupt Enable
	IME bool  `json:"IME"` // interrupt master enable

	// emulator state
	HALTED  bool `json:"HALTED"`  // is the CPU halted
	STOPPED bool `json:"STOPPED"` // is the CPU stopped
}

type GameboyState struct {
	PREV_CPU_STATE *CpuState     `json:"prevState"`
	CURR_CPU_STATE *CpuState     `json:"currState"`
	INSTR          *Instruction  `json:"instruction"`
	MEMORY_WRITES  []MemoryWrite `json:"memoryWrites"`
}

func (gbs *GameboyState) print() {
	gbs.printCPUState()
	gbs.printInstruction()
}

func (gbs *GameboyState) printCPUState() {
	fmt.Println("")
	fmt.Println("\n> CPU State:")
	fmt.Println("------------")
	// if previous and current states are nil, there is nothing to print
	if (gbs.PREV_CPU_STATE == nil) && (gbs.CURR_CPU_STATE == nil) {
		fmt.Println("> No CPU state to print")
		// if only the current state is available, print it
	} else if gbs.PREV_CPU_STATE == nil {
		curr := reflect.Indirect(reflect.ValueOf(gbs.CURR_CPU_STATE))
		typeOfCpu := curr.Type()
		for i := 0; i < curr.NumField(); i++ {
			if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
				fmt.Printf("- %s: %t\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
				fmt.Printf("- %s: 0x%02X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
				fmt.Printf("- %s: 0x%04X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
				fmt.Printf("- %s: %s\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			}
		}
	} else {
		prev := reflect.Indirect(reflect.ValueOf(gbs.PREV_CPU_STATE))
		curr := reflect.Indirect(reflect.ValueOf(gbs.CURR_CPU_STATE))
		typeOfCpu := prev.Type()

		for i := 0; i < prev.NumField(); i++ {
			if prev.Field(i).Interface() != curr.Field(i).Interface() {
				if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
					fmt.Printf("- %s: %t -> %t \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
					fmt.Printf("- %s: 0x%02X -> 0x%02X \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
					fmt.Printf("- %s: 0x%04X -> 0x%04X \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
					fmt.Printf("- %s: %s -> %s \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				}
			} else {
				if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
					fmt.Printf("- %s: %t \n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
					fmt.Printf("- %s: 0x%02X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
					fmt.Printf("- %s: 0x%04X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
					fmt.Printf("- %s: %s\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				}
			}
		}
	}
}

func (gbs *GameboyState) printInstruction() {
	fmt.Println("")
	fmt.Println("\n> Instruction:")
	fmt.Println("--------------")
	if gbs.INSTR == nil {
		fmt.Println("> No instruction to print")
	} else {
		fmt.Printf("- Opcode: 0x%02X\n", gbs.CURR_CPU_STATE.IR)
		fmt.Printf("- Mnemonic: %s\n", gbs.INSTR.Mnemonic)
		fmt.Printf("- Bytes: %d\n", gbs.INSTR.Bytes)
		fmt.Printf("- Cycles: %v\n", gbs.INSTR.Cycles)
		fmt.Printf("- Operands: %v\n", gbs.INSTR.Operands)
		fmt.Printf("- Immediate: %t\n", gbs.INSTR.Immediate)
		fmt.Printf("- Flags: %v\n", gbs.INSTR.Flags)
	}
}

// get the memories current content
func (d *Debugger) currCpuState() *CpuState {
	return &CpuState{
		CpuCycles:     d.gameboy.cpu.cpuCycles,
		PC:            d.gameboy.cpu.pc,
		SP:            d.gameboy.cpu.sp,
		A:             d.gameboy.cpu.a,
		F:             d.gameboy.cpu.f,
		Z:             d.gameboy.cpu.f&0x80 != 0,
		N:             d.gameboy.cpu.f&0x40 != 0,
		H:             d.gameboy.cpu.f&0x20 != 0,
		C:             d.gameboy.cpu.f&0x10 != 0,
		BC:            uint16(d.gameboy.cpu.b)<<8 | uint16(d.gameboy.cpu.c),
		DE:            uint16(d.gameboy.cpu.d)<<8 | uint16(d.gameboy.cpu.e),
		HL:            uint16(d.gameboy.cpu.h)<<8 | uint16(d.gameboy.cpu.l),
		PREFIXED:      d.gameboy.cpu.prefixed,
		IR:            d.gameboy.cpu.ir,
		OPERAND_VALUE: d.gameboy.cpu.operand,
		IE:            d.gameboy.cpu.bus.Read(0xFFFF),
		IME:           d.gameboy.cpu.ime,
		HALTED:        d.gameboy.cpu.halted,
		STOPPED:       d.gameboy.cpu.stopped,
	}
}

// currInstruction: returns the current instruction being processed based on cpu IR and prefix values
func (d *Debugger) currInstruction() *Instruction {
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", d.gameboy.cpu.ir)), d.gameboy.cpu.prefixed)
	return &instruction
}

// clear memory writes
func (d *Debugger) clearMemoryWrites() {
	d.gameboy.cpuBus.mmu.clearMemoryWrites()
}

// returns the current memory writes
func (d *Debugger) currMemoryWrites() []MemoryWrite {
	return d.gameboy.cpuBus.mmu.memoryWrites
}

/**
 * print the current state of the gameboy
 */
func (d *Debugger) PrintCPUState() {
	d.state.printCPUState()
}

/**
 * print the current instruction
 */
func (d *Debugger) PrintInstruction() {
	d.state.printInstruction()
}

/**
 * print the properties of the memories attached to the bus
 */
func (d *Debugger) PrintMemoryProperties() {
	memoryMaps := d.gameboy.cpuBus.mmu.GetMemoryMaps()
	fmt.Println("")
	fmt.Println("\n> Memory Mapping:")
	fmt.Println("-----------------")
	for _, memoryMap := range memoryMaps {
		fmt.Printf("> Memory %s: %d bytes @ 0x%04X->0x%04X\n", memoryMap.Name, len(memoryMap.Data), memoryMap.Address, memoryMap.Address+uint16(len(memoryMap.Data))-1)
	}
}
