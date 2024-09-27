package gameboy

import (
	"fmt"
	"reflect"
)

/**
 * debugger struct: combination of a gameboy, its internal state and a list of breakpoints set by the user
 */
type Debugger struct {
	gameboy     *Gameboy
	state       *GameboyState
	breakpoints []uint16 // list of breakpoints addresses set by the user to pause the execution with a maximum of 100 breakpoints
}

/**
 * creates a new debugger: instanciates a new gameboy and initializes the breakpoints list
 */
func NewDebugger() *Debugger {
	gb := NewGameboy()
	return &Debugger{
		gameboy:     gb,
		state:       &GameboyState{},
		breakpoints: make([]uint16, 100),
	}
}

/**
 * initializes the gameboy with the given ROM and returns a pointer to the gameboy state.
 */
func (d *Debugger) Init(romName string) *GameboyState {
	d.gameboy.init(romName)
	d.state = &GameboyState{
		PREV_CPU_STATE: nil,
		CURR_CPU_STATE: nil,
		INSTR:          nil,
		MEMORY_WRITES:  []MemoryWrite{},
	}
	return d.state
}

/**
 * helper function to check if a value is present in an uint16 array.
 */
func contains(arr []uint16, addr uint16) bool {
	for _, v := range arr {
		if v == addr {
			return true
		}
	}
	return false
}

/**
 * adds a breakpoint at the given address if not already present.
 */
func (d *Debugger) AddBreakPoint(addr uint16) {
	if contains(d.breakpoints, addr) {
		return
	} else {
		d.breakpoints = append(d.breakpoints, addr)
	}
}

/**
 * removes a breakpoint if present.
 */
func (d *Debugger) RemoveBreakPoint(addr uint16) {
	for i, v := range d.breakpoints {
		if v == addr {
			d.breakpoints = append(d.breakpoints[:i], d.breakpoints[i+1:]...)
			return
		}
	}
}

/**
 * retrieve the breakpoints list
 */
func (d *Debugger) GetBreakPoints() []uint16 {
	return d.breakpoints
}

/**
 * shifts the current state into the previous state
 */
func (d *Debugger) shiftState() {
	d.state.PREV_CPU_STATE = d.state.CURR_CPU_STATE
}

/**
 * saves the current state of the gameboy into the debugger state.
 */
func (d *Debugger) testInstructionExecution() {
	instr := d.gameboy.currInstruction()
	curr := d.state.CURR_CPU_STATE
	prev := d.state.PREV_CPU_STATE

	if (curr == nil) || (prev == nil) {
		return
	}

	// check the flags
	v := reflect.ValueOf(instr.Flags)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		key := typeOfS.Field(i).Name
		flagValue := field.String()

		// Use reflection to get the value of the field from curr
		currValue := reflect.ValueOf(curr).Elem().FieldByName(key)
		prevValue := reflect.ValueOf(prev).Elem().FieldByName(key)

		if flagValue == "0" {
			if currValue.Bool() {
				fmt.Printf("Debugger@0x%04X> 'testInstructionExecution' failed: flag should be 0 (%v)\n", d.gameboy.cpu.pc, instr.Mnemonic)
			}
		} else if flagValue == "1" {
			if !currValue.Bool() {
				fmt.Printf("Debugger@0x%04X> 'testInstructionExecution' failed: flag should be 1 (%v)\n", d.gameboy.cpu.pc, instr.Mnemonic)
			}
		} else if flagValue == "-" {
			// '-' means the flag is not relevant for this instruction, so we skip the check
			if currValue.Bool() != prevValue.Bool() {
				fmt.Printf("Debugger@0x%04X> 'testInstructionExecution' failed: flag should stay the same (%v) : %v(%v)=%v->%v\n", d.gameboy.cpu.pc, instr.Mnemonic, key, flagValue, currValue.Bool(), prevValue.Bool())

			}
		} else {
			//fmt.Printf("\nDebugger> 'testInstructionExecution' flag value %v unsupported", key)
		}
	}

	// check the registers
	// check the memory reads
	// check the memory writes
}

func (d *Debugger) saveState() {
	d.state.CURR_CPU_STATE = d.gameboy.currCpuState()
	d.state.INSTR = d.gameboy.currInstruction()
	d.state.MEMORY_WRITES = d.gameboy.currMemoryWrites()
}

/**
 * run the next instruction and return the gameboy state
 */
func (d *Debugger) Step() *GameboyState {
	// clear memory writes
	d.gameboy.clearMemoryWrites()
	// shift the current state into the previous state
	d.shiftState()
	// run the next instruction
	d.gameboy.Step()
	// save the current state
	d.saveState()

	// test the instruction execution
	// - read from memory locations named in the instruction operands
	// - write to memory locations named in the instruction operands
	// - changed the correct CPU registers, everytime or when the condition is met
	// - changed the correct flags, everytime or when the condition is met
	d.testInstructionExecution()

	// return the current state
	return d.state
}

/**
 * run the gameboy until a breakpoint is reached or the gameboy is halted
 */
func (d *Debugger) Run() *GameboyState {
	// TODO: I think that i can simply start the timer and the CPU will run until it is halted or stopped by a breakpoint but where ???
	// reset memory writes
	d.gameboy.clearMemoryWrites()

	// run the gameboy until a breakpoint is reached or the gameboy is halted
	for {
		d.state.PREV_CPU_STATE = d.gameboy.currCpuState() // since we do not know if this will be the last step before returning, we have to save the last state into the previous state at each iteration
		d.gameboy.Step()
		// save the current state
		d.shiftState()
		d.state.CURR_CPU_STATE = d.gameboy.currCpuState()
		d.state.INSTR = d.gameboy.currInstruction()
		d.state.MEMORY_WRITES = d.gameboy.currMemoryWrites()
		d.testInstructionExecution()
		// check if the current PC is a breakpoint or if the gameboy is halted
		if contains(d.breakpoints, d.gameboy.cpu.pc) || d.gameboy.cpu.halted || d.gameboy.cpu.stopped {
			// we do not halt the processor which is a feature used by the gameboy to save power, we just stop the execution from the debugger point of view by not executing any more instructions
			break
		}
	}
	return d.state
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

// return the list of memories attached to the mmu including their name, address and data
func (d *Debugger) GetAttachedMemories() []MemoryWrite {
	return d.gameboy.cpuBus.mmu.GetMemoryMaps()
}
