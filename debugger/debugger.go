package debugger

import (
	"time"

	ds "github.com/codefrite/gameboy-go/datastructure"
	"github.com/codefrite/gameboy-go/gameboy"
)

const STATE_QUEUE_MAX_LENGTH = 100

// debugger struct: combination of a gameboy, its internal state and a list of breakpoints set by the user
type Debugger struct {
	// state
	gameboy     *gameboy.Gameboy
	programFlow *ds.Fifo[uint16] // queue of program counter positions to render a diagram of the program flow
	breakpoints []uint16         // list of breakpoints addresses set by the user to pause the execution with a maximum of 100 breakpoints

	cpuStateQueue    *ds.Fifo[gameboy.CpuState]
	memoryStateQueue *ds.Fifo[[]gameboy.MemoryWrite]

	// state channels received from the client meant to listen to the gameboy state
	clientCpuStateChannel    chan<- gameboy.CpuState
	clientMemoryStateChannel chan<- []gameboy.MemoryWrite
	doneChannel              chan bool

	// internal channels corresponding to the channels received from the client and used to intercept, store in a queue, and then relay the state changes
	internalCpuStateChannel    chan gameboy.CpuState
	internalMemoryStateChannel chan []gameboy.MemoryWrite
}

// instantiate a new debugger:
// - instanciates a new gameboy
// - initializes the internal channels to listen to the gameboy state
// - initializes the breakpoints list
// - initializes the program flow queue
// - initializes the state queues (cpu, ppu, apu, memory, joypad)
func NewDebugger(
	cpuStateChannel chan<- gameboy.CpuState,
	memoryStateChannel chan<- []gameboy.MemoryWrite,
) *Debugger {

	// instantiate an empty debugger
	debugger := &Debugger{
		internalCpuStateChannel:  make(chan gameboy.CpuState), // we always need to listen to the cpu state to handle breakpoints
		clientCpuStateChannel:    cpuStateChannel,
		clientMemoryStateChannel: memoryStateChannel,
		doneChannel:              make(chan bool), // used to notify client that crystal has stopped
		breakpoints:              make([]uint16, 0),
	}

	// create the internal channels to listen to the gameboy state if they are used by the client
	if memoryStateChannel != nil {
		debugger.internalMemoryStateChannel = make(chan []gameboy.MemoryWrite)
	}

	// instantiate a new gameboy with the debugger internal channels
	gb := gameboy.NewGameboy(
		nil, // gameboyActionChannel is not used as the debugger ticks the gameboy manually
		debugger.internalCpuStateChannel,
		nil,
		nil,
		debugger.internalMemoryStateChannel,
	)

	// attach the gameboy to the debugger
	debugger.gameboy = gb

	// initializes the debugger state with empty state queues and breakpoints list
	debugger.reset()

	// launch the debugger state channels listener
	go debugger.stateChannelsListener()

	return debugger
}

// resets the debugger state (queues, breakpoints, etc)
func (d *Debugger) reset() {
	d.programFlow = ds.NewFifo[uint16](STATE_QUEUE_MAX_LENGTH)
	d.cpuStateQueue = ds.NewFifo[gameboy.CpuState](STATE_QUEUE_MAX_LENGTH)
	d.memoryStateQueue = ds.NewFifo[[]gameboy.MemoryWrite](STATE_QUEUE_MAX_LENGTH)
	d.breakpoints = make([]uint16, 0)
}

// initializes the gameboy with the given ROM and returns a pointer to the gameboy state
func (d *Debugger) LoadRom(romName string) {
	d.reset()
	go d.gameboy.LoadRom(romName)
}

func (d *Debugger) GetMemoryMaps() []gameboy.MemoryWrite {
	return d.gameboy.GetMemoryMaps()
}

// returns the gameboy state
func (d *Debugger) GetCpuState() gameboy.CpuState {
	return d.gameboy.GetCpuState()
}

func (d *Debugger) GetMemoryWrites() []gameboy.MemoryWrite {
	return d.gameboy.GetMemoryWrites()
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

// Execution Control

// Tick the gameboy once and send the state to the client
func (d *Debugger) Tick() {
	d.gameboy.Tick() // must run in a goroutine to avoid deadlock (gameboy.sendState() is blocking until we read from the internal channels)
}

// Run the gameboy until we reach a breakpoint
func (d *Debugger) Run() {
	for {
		// tick the gameboy @ 1 Hz
		startTime := time.Now()
		d.gameboy.Tick()
		tickDuration := time.Since(startTime)
		time.Sleep(time.Second - tickDuration)

		/*
			// check if we reached a breakpoint
			if contains(d.breakpoints, cpuState.PC) {
				break
			}
		*/
	}
}

// listen to the internal channels and relay the state to the client
func (d *Debugger) stateChannelsListener() {
	for {
		select {
		case cpuState := <-d.internalCpuStateChannel:
			d.clientCpuStateChannel <- cpuState
		case memoryState := <-d.internalMemoryStateChannel:
			d.clientMemoryStateChannel <- memoryState
		case <-d.doneChannel:
			return
		}
	}
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
