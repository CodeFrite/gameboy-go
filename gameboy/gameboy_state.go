package gameboy

import "fmt"

type CpuState struct {
	// Special registers
	PC uint16 `json:"PC"` // Program Counter
	SP uint16 `json:"SP"` // Stack Pointer
	A  uint8  `json:"A"`  // Accumulator
	F  uint8  `json:"F"`  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	Z  bool   `json:"Z"`  // Zero flag
	N  bool   `json:"N"`  // Subtraction flag
	H  bool   `json:"H"`  // Half Carry flag
	C  bool   `json:"C"`  // Carry flag
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
	HALTED bool `json:"HALTED"` // is the CPU halted
}

type MemoryWrite struct {
	Address uint16  `json:"address"`
	Data    []uint8 `json:"data"`
}

type GameboyState struct {
	PREV_CPU_STATE *CpuState     `json:"prevState"`
	CURR_CPU_STATE *CpuState     `json:"currState"`
	INSTR          *Instruction  `json:"instruction"`
	MEMORY_WRITES  []MemoryWrite `json:"memoryWrites"`
}

// shift the current state to the previous state and reset the current state
func (gb *Gameboy) shiftState() {
	gb.state.PREV_CPU_STATE = gb.state.CURR_CPU_STATE
}

func (gb *Gameboy) getCurrentState() *GameboyState {
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", gb.cpu.IR)), gb.cpu.Prefixed)
	data := gb.bus.Dump(0, gb.bootrom.Size())
	memoryWrites := []MemoryWrite{}
	memoryWrites = append(memoryWrites, MemoryWrite{
		Address: 0x0000,
		Data:    data,
	})
	return &GameboyState{
		PREV_CPU_STATE: gb.state.CURR_CPU_STATE,
		CURR_CPU_STATE: &CpuState{
			PC:            gb.cpu.PC,
			SP:            gb.cpu.SP,
			A:             gb.cpu.A,
			F:             gb.cpu.F,
			Z:             gb.cpu.F&0x80 != 0,
			N:             gb.cpu.F&0x40 != 0,
			H:             gb.cpu.F&0x20 != 0,
			C:             gb.cpu.F&0x10 != 0,
			BC:            uint16(gb.cpu.B)<<8 | uint16(gb.cpu.C),
			DE:            uint16(gb.cpu.D)<<8 | uint16(gb.cpu.E),
			HL:            uint16(gb.cpu.H)<<8 | uint16(gb.cpu.L),
			PREFIXED:      gb.cpu.Prefixed,
			IR:            gb.cpu.IR,
			OPERAND_VALUE: gb.cpu.Operand,
			IE:            gb.cpu.IE,
			IME:           gb.cpu.IME,
			HALTED:        gb.cpu.halted,
		},
		INSTR:         &instruction,
		MEMORY_WRITES: memoryWrites,
	}
}

func (gb *Gameboy) saveCurrentState() {
	gb.state = gb.getCurrentState()
}

func (gb *Gameboy) State() *GameboyState {
	return gb.state
}
