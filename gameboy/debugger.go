package gameboy

const STATE_QUEUE_MAX_LENGTH = 100

/**
 * debugger struct: combination of a gameboy, its internal state and a list of breakpoints set by the user
 */
type Debugger struct {
	// state
	gameboy          *Gameboy
	cpuStateQueue    *fifo[CpuState]
	ppuStateQueue    *fifo[PpuState]
	apuStateQueue    *fifo[ApuState]
	memoryStateQueue fifo[[]MemoryWrite]
	joypadStateQueue *fifo[JoypadState]

	// break points
	breakpoints []uint16 // list of breakpoints addresses set by the user to pause the execution with a maximum of 100 breakpoints

	// state channels received from the client meant to listen to the gameboy state
	clientCpuStateChannel    chan<- *CpuState // v0.4.0
	clientPpuStateChannel    chan<- *PpuState // v0.4.1
	clientApuStateChannel    chan<- *ApuState // v0.4.2
	clientMemoryStateChannel chan<- []MemoryWrite
	clientJoypadStateChannel <-chan *JoypadState
	doneChannel              chan bool

	// internal channels corresponding to the channels received from the client and used to intercept, store in a queue, and then relay the state changes
	internalCpuStateChannel    chan *CpuState
	internalPpuStateChannel    chan *PpuState
	internalApuStateChannel    chan *ApuState
	internalMemoryStateChannel chan []MemoryWrite
	internalJoypadStateChannel chan *JoypadState
}

/**
 * creates a new debugger: instanciates a new gameboy and initializes the breakpoints list
 */
func NewDebugger(
	cpuStateChannel chan<- *CpuState,
	ppuStateChannel chan<- *PpuState,
	apuStateChannel chan<- *ApuState,
	memoryStateChannel chan<- []MemoryWrite,
	joypadStateChannel <-chan *JoypadState,
) *Debugger {

	internalCpuStateChannel := make(chan *CpuState)
	internalPpuStateChannel := make(chan *PpuState)
	internalApuStateChannel := make(chan *ApuState)
	internalMemoryStateChannel := make(chan []MemoryWrite)
	internalJoypadStateChannel := make(chan *JoypadState)
	doneChannel := make(chan bool)

	gb := NewGameboy(
		internalCpuStateChannel,
		internalPpuStateChannel,
		internalApuStateChannel,
		internalMemoryStateChannel,
		internalJoypadStateChannel,
	)

	return &Debugger{
		gameboy:                    gb,
		cpuStateQueue:              newFifo[CpuState](),
		ppuStateQueue:              newFifo[PpuState](),
		apuStateQueue:              newFifo[ApuState](),
		memoryStateQueue:           *newFifo[[]MemoryWrite](),
		joypadStateQueue:           newFifo[JoypadState](),
		clientCpuStateChannel:      cpuStateChannel,
		clientPpuStateChannel:      ppuStateChannel,
		clientApuStateChannel:      apuStateChannel,
		clientMemoryStateChannel:   memoryStateChannel,
		clientJoypadStateChannel:   joypadStateChannel,
		internalCpuStateChannel:    internalCpuStateChannel,
		internalPpuStateChannel:    internalPpuStateChannel,
		internalApuStateChannel:    internalApuStateChannel,
		internalMemoryStateChannel: internalMemoryStateChannel,
		internalJoypadStateChannel: internalJoypadStateChannel,
		doneChannel:                doneChannel, // used to notify client that crystal has stopped
		breakpoints:                make([]uint16, 0),
	}
}

/**
 * initializes the gameboy with the given ROM and returns a pointer to the gameboy state.
 */
func (d *Debugger) LoadRom(romName string) {
	d.gameboy.LoadRom(romName)
	initialCpuState := d.gameboy.cpu.getState()
	initialPpuState := d.gameboy.ppu.getState()
	initialApuState := d.gameboy.apu.getState()
	initialMemoryWrites := &d.gameboy.cpuBus.mmu.memoryWrites

	// save the initial state
	d.cpuStateQueue.push(initialCpuState)
	d.ppuStateQueue.push(initialPpuState)
	d.apuStateQueue.push(initialApuState)
	d.memoryStateQueue.push(initialMemoryWrites)
}

// run the next instruction and return the gameboy state
func (d *Debugger) Step() {
	// run the next instruction
	d.gameboy.Step()

	// add the new state to the queue and send it to the client if a channel is present
	cpuState := <-d.internalCpuStateChannel
	d.cpuStateQueue.push(cpuState)
	if d.clientCpuStateChannel != nil {
		d.clientCpuStateChannel <- cpuState
	}
	ppuState := <-d.internalPpuStateChannel
	d.ppuStateQueue.push(ppuState)
	if d.clientPpuStateChannel != nil {
		d.clientPpuStateChannel <- ppuState
	}
	apuState := <-d.internalApuStateChannel
	d.apuStateQueue.push(apuState)
	if d.clientApuStateChannel != nil {
		d.clientApuStateChannel <- apuState
	}
	memoryState := <-d.internalMemoryStateChannel
	d.memoryStateQueue.push(&memoryState)
	if d.clientMemoryStateChannel != nil {
		d.clientMemoryStateChannel <- memoryState
	}
}

// run the gameboy until a breakpoint is reached or the gameboy is halted
func (d *Debugger) Run() chan bool {
	d.listenToGameboyState()
	d.gameboy.Run()
	return d.doneChannel
}

func (d *Debugger) listenToGameboyState() {
	// launch the debugger internal channels in a go routine to listen to the gameboy state channels
	// TODO: add a done channel to stop the go routine
	go func() {
	loop:
		for {
			// listen to the gameboy state channels
			select {
			case cpuState := <-d.internalCpuStateChannel:

				d.cpuStateQueue.push(cpuState)
				if d.clientCpuStateChannel != nil {
					d.clientCpuStateChannel <- cpuState
				}
				// manage breakpoints
				if contains(d.breakpoints, cpuState.PC) || cpuState.HALTED || cpuState.STOPPED {
					// stop the gameboy

					d.gameboy.crystal.Stop()
					d.doneChannel <- true
					break loop
				}
			case ppuState := <-d.internalPpuStateChannel:

				d.ppuStateQueue.push(ppuState)
				if d.clientPpuStateChannel != nil {
					d.clientPpuStateChannel <- ppuState
				}
			case apuState := <-d.internalApuStateChannel:

				d.apuStateQueue.push(apuState)
				if d.clientApuStateChannel != nil {
					d.clientApuStateChannel <- apuState
				}
			case memoryState := <-d.internalMemoryStateChannel:

				d.memoryStateQueue.push(&memoryState)
				if d.clientMemoryStateChannel != nil {
					d.clientMemoryStateChannel <- memoryState
				}
				/*
					case joypadState := <-d.clientJoypadStateChannel:
						// TODO: use it for conditional breakpoints on joypad events
						d.internalJoypadStateChannel <- joypadState
				*/
			}
		}
	}()
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
 * saves the current state of the gameboy into the debugger state.
 */
/*
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
*/

/*
func (d *Debugger) saveState() {
	d.state.CURR_CPU_STATE = d.gameboy.cpu.getState()
	d.state.INSTR = d.gameboy.currInstruction()
	d.state.MEMORY_WRITES = d.gameboy.currMemoryWrites()
}*/

// return the list of memories attached to the mmu including their name, address and data
func (d *Debugger) GetAttachedMemories() []MemoryWrite {
	return d.gameboy.cpuBus.mmu.GetMemoryMaps()
}

// print the current state of the gameboy
func (d *Debugger) PrintCPUState() {
	d.cpuStateQueue.peek().print()
}

// print the current instruction
func (d *Debugger) PrintInstruction() {
}

func (d *Debugger) PrintMemoryProperties() {
	d.gameboy.PrintMemoryProperties()
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
