package gameboy

const STATE_QUEUE_MAX_LENGTH = 100

// debugger struct: combination of a gameboy, its internal state and a list of breakpoints set by the user
type Debugger struct {
	// state
	gameboy     *Gameboy
	programFlow *fifo[uint16] // queue of program counter positions to render a diagram of the program flow
	breakpoints []uint16      // list of breakpoints addresses set by the user to pause the execution with a maximum of 100 breakpoints

	cpuStateQueue    *fifo[CpuState]
	ppuStateQueue    *fifo[PpuState]
	apuStateQueue    *fifo[ApuState]
	memoryStateQueue *fifo[[]MemoryWrite]
	joypadStateQueue *fifo[JoypadState]

	// state channels received from the client meant to listen to the gameboy state
	clientCpuStateChannel    chan<- CpuState
	clientPpuStateChannel    chan<- PpuState
	clientApuStateChannel    chan<- ApuState
	clientMemoryStateChannel chan<- []MemoryWrite
	clientJoypadStateChannel <-chan JoypadState
	doneChannel              chan bool

	// internal channels corresponding to the channels received from the client and used to intercept, store in a queue, and then relay the state changes
	internalCpuStateChannel    chan CpuState
	internalPpuStateChannel    chan PpuState
	internalApuStateChannel    chan ApuState
	internalMemoryStateChannel chan []MemoryWrite
	internalJoypadStateChannel chan JoypadState
}

// instantiate a new debugger:
// - instanciates a new gameboy
// - initializes the internal channels to listen to the gameboy state
// - initializes the breakpoints list
// - initializes the program flow queue
// - initializes the state queues (cpu, ppu, apu, memory, joypad)
func NewDebugger(
	cpuStateChannel chan<- CpuState,
	ppuStateChannel chan<- PpuState,
	apuStateChannel chan<- ApuState,
	memoryStateChannel chan<- []MemoryWrite,
	joypadStateChannel <-chan JoypadState,
) *Debugger {
	// create the internal channels to listen to the gameboy state
	internalCpuStateChannel := make(chan CpuState)
	internalPpuStateChannel := make(chan PpuState)
	internalApuStateChannel := make(chan ApuState)
	internalMemoryStateChannel := make(chan []MemoryWrite)
	internalJoypadStateChannel := make(chan JoypadState)
	doneChannel := make(chan bool)

	gb := NewGameboy(
		internalCpuStateChannel,
		internalPpuStateChannel,
		internalApuStateChannel,
		internalMemoryStateChannel,
		internalJoypadStateChannel,
	)

	debugger := Debugger{
		gameboy:                    gb,
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

	// initializes the debugger state with empty state queues and breakpoints list
	debugger.reset()

	return &debugger
}

// initializes the gameboy with the given ROM and returns a pointer to the gameboy state
func (d *Debugger) LoadRom(romName string) {
	d.gameboy.LoadRom(romName)

	d.reset()

	initialCpuState := d.gameboy.cpu.getState()
	initialPpuState := d.gameboy.ppu.getState()
	initialApuState := d.gameboy.apu.getState()
	initialMemoryWrites := d.gameboy.cpuBus.mmu.memoryWrites

	// save the initial state
	d.programFlow.push(initialCpuState.PC)
	d.cpuStateQueue.push(initialCpuState)
	d.ppuStateQueue.push(initialPpuState)
	d.apuStateQueue.push(initialApuState)
	d.memoryStateQueue.push(initialMemoryWrites)
}

// tick the gameboy and return its state
func (d *Debugger) Tick() {
	// run the next instruction
	d.gameboy.Tick()

	// add the new state to the queue and send it to the client if a channel is present
	cpuState := <-d.internalCpuStateChannel
	d.cpuStateQueue.push(cpuState)
	d.programFlow.push(cpuState.PC)
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
	d.memoryStateQueue.push(memoryState)
	if d.clientMemoryStateChannel != nil {
		d.clientMemoryStateChannel <- memoryState
	}
}

// run the gameboy until a breakpoint is reached or the gameboy is halted
func (d *Debugger) Run() chan bool {
	go d.listenToGameboyState()
	d.gameboy.Run()
	return d.doneChannel
}

// pause the gameboy
func (d *Debugger) Pause() {
	d.gameboy.Pause()
}

// resume the gameboy
func (d *Debugger) Resume() {
	d.gameboy.Resume()
}

// stop the gameboy
func (d *Debugger) Stop() {
	d.gameboy.Stop()
}

// adds a breakpoint at the given address if not already present
func (d *Debugger) AddBreakPoint(addr uint16) {
	if contains(d.breakpoints, addr) {
		return
	} else {
		d.breakpoints = append(d.breakpoints, addr)
	}
}

// removes a breakpoint if present
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

// return the list of memories attached to the mmu including their name, address and data
func (d *Debugger) GetAttachedMemories() []MemoryWrite {
	return d.gameboy.cpuBus.mmu.GetMemoryMaps()
}

// HELPER FUNCS

// helper function to check if a value is present in an uint16 array
func contains(arr []uint16, addr uint16) bool {
	for _, v := range arr {
		if v == addr {
			return true
		}
	}
	return false
}

// resets the debugger state (queues, breakpoints, etc)
func (d *Debugger) reset() {
	d.programFlow = newFifo[uint16](INFINITE_MAX_NODE_COUNT)
	d.cpuStateQueue = newFifo[CpuState](STATE_QUEUE_MAX_LENGTH)
	d.ppuStateQueue = newFifo[PpuState](STATE_QUEUE_MAX_LENGTH)
	d.apuStateQueue = newFifo[ApuState](STATE_QUEUE_MAX_LENGTH)
	d.memoryStateQueue = newFifo[[]MemoryWrite](STATE_QUEUE_MAX_LENGTH)
	d.joypadStateQueue = newFifo[JoypadState](STATE_QUEUE_MAX_LENGTH)
	d.breakpoints = make([]uint16, 0)
}

// when the gameboy is run, launches a go routine that:
// - listens to the gameboy state channels and relay the state changes to the client
// - pushes the state changes to the queues (cpu, ppu, apu, memory, joypad)
// - listens to the done channel to stop the go routine
func (d *Debugger) listenToGameboyState() {
	// launch the debugger internal channels in a go routine to listen to the gameboy state channels
	// TODO: add a done channel to stop the go routine
	go func() {
	loop:
		for {
			// listen to the gameboy state channels
			select {
			case cpuState := <-d.internalCpuStateChannel:

				// manage breakpoints
				if contains(d.breakpoints, cpuState.PC) || cpuState.HALTED || cpuState.STOPPED {
					// stop the gameboy crystal from ticking
					//d.gameboy.Stop()

					// we must send the ppu, apu and memory states to the client
					if d.clientPpuStateChannel != nil {
						ppuState := <-d.internalPpuStateChannel
						d.ppuStateQueue.push(ppuState)
						d.clientPpuStateChannel <- ppuState
					}

					if d.clientApuStateChannel != nil {
						apuState := <-d.internalApuStateChannel
						d.apuStateQueue.push(apuState)
						d.clientApuStateChannel <- apuState
					}

					if d.clientMemoryStateChannel != nil {
						memoryState := <-d.internalMemoryStateChannel
						d.memoryStateQueue.push(memoryState)
						d.clientMemoryStateChannel <- memoryState
					}

					d.cpuStateQueue.push(cpuState)
					d.programFlow.push(cpuState.PC)
					if d.clientCpuStateChannel != nil {
						d.clientCpuStateChannel <- cpuState
					}

					break loop
				} else {
					d.cpuStateQueue.push(cpuState)
					d.programFlow.push(cpuState.PC)
					if d.clientCpuStateChannel != nil {
						d.clientCpuStateChannel <- cpuState
					}
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

				d.memoryStateQueue.push(memoryState)
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
		d.doneChannel <- true
	}()
}
