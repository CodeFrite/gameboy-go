package gameboy

import "fmt"

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
	// return the current state
	return d.state
}

/**
 * run the gameboy until a breakpoint is reached or the gameboy is halted
 */
func (d *Debugger) Run() *GameboyState {
	// reset memory writes
	d.gameboy.clearMemoryWrites()
	// keep track of the number of steps executed during the run
	steps := 0
	// run the gameboy until a breakpoint is reached or the gameboy is halted
	for {
		d.state.PREV_CPU_STATE = d.gameboy.currCpuState() // since we do not know if this will be the last step before returning, we have to save the last state into the previous state at each iteration
		d.gameboy.Step()
		if contains(d.breakpoints, d.gameboy.cpu.PC) || d.gameboy.cpu.halted {
			break
		}
		steps++
	}
	// save the current state
	d.state.CURR_CPU_STATE = d.gameboy.currCpuState()
	d.state.INSTR = d.gameboy.currInstruction()
	d.state.MEMORY_WRITES = d.gameboy.currMemoryWrites()
	fmt.Printf("Executed %d steps\n", steps)
	return d.state
}

/**
 * print the current state of the gameboy
 */
func (d *Debugger) PrintCPUState() {
	d.state.printCPUState()
}

/**
 * print the properties of the memories attached to the bus
 */
func (d *Debugger) PrintMemoryProperties() {
	memoryMaps := d.gameboy.bus.mmu.GetMemoryMaps()
	fmt.Println("")
	fmt.Println("\n> Memory Mapping:")
	fmt.Println("-----------------")
	for _, memoryMap := range memoryMaps {
		fmt.Printf("> Memory %s: %d bytes @ 0x%04X->0x%04X\n", memoryMap.Name, len(memoryMap.Data), memoryMap.Address, memoryMap.Address+uint16(len(memoryMap.Data))-1)
	}
}

// return the list of memories attached to the mmu including their name, address and data
func (d *Debugger) GetAttachedMemories() []MemoryWrite {
	return d.gameboy.bus.mmu.GetMemoryMaps()
}
