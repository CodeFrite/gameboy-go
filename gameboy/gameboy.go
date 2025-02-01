package gameboy

import (
	"log"
	"time"
)

// CONSTANTS
const (
	CRYSTAL_FREQUENCY    time.Duration = 4194304 // 4.194304MHz
	TICK_DURATION                      = time.Duration(1e9 / CRYSTAL_FREQUENCY)
	BOOT_ROM_MEMORY_NAME               = "Boot ROM"
	BOOT_ROM_START       uint16        = 0x0000
	BOOT_ROM_LEN         uint16        = 0x0100
	ROMS_URI                           = "/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms"

	// Gameboy states
	GB_STATE_NO_GAME_LOADED GameBoyState = "no game loaded" // no game loaded
	GB_STATE_PAUSED         GameBoyState = "paused"         // gameboy is paused
	GB_STATE_RUNNING        GameBoyState = "running"        // gameboy is running

	// Gameboy State's Actions
	GB_ACTION_LOAD_GAME GameBoyAction = "load"  // load a game
	GB_ACTION_RUN       GameBoyAction = "run"   // run the gameboy if it has a game loaded
	GB_ACTION_PAUSE     GameBoyAction = "pause" // pause the gameboy
	GB_ACTION_RESET     GameBoyAction = "reset" // reset the gameboy
)

type GameBoyState string
type GameBoyAction string
type GameboyActionMessage struct {
	Action  GameBoyAction
	payload interface{}
}

// the gameboy is composed out of a CPU, memories (ram & registers), a cartridge and a bus
type Gameboy struct {
	// state
	ticks uint64       // number of ticks since the gameboy started
	state GameBoyState // current state of the gameboy

	// components
	timer     *Timer // Gameboy Timer (DIV, TIMA, TMA, TAC)
	bus       *Bus
	cpu       *CPU
	ppu       *PPU
	apu       *APU
	bootrom   *Memory    // 0x0000-0x00FF: (256 bytes) - Boot ROM
	cartridge *Cartridge // Cartridge ROM (32KB) [0x0000-0x7FFF]
	vram      *Memory    // Video RAM (8KB) [0x8000-0x9FFF]
	wram      *Memory    // Working RAM (8KB) [0xC000-0xDFFF]
	joypad    *Joypad

	// state channels (sharing concrete types to avoid pointer values being changed before being sent to the frontend by the server)
	// TODO: now that i built my gameloop differently, i can pass pointers to the frontend instead of copying the state i guess
	gameboyActionChannel <-chan GameboyActionMessage // responsible to load a game, run, pause and stop the gameboy
	cpuStateChannel      chan<- CpuState
	ppuStateChannel      chan<- PpuState
	apuStateChannel      chan<- ApuState
	memoryStateChannel   chan<- []MemoryWrite
}

// create a new gameboy struct
func NewGameboy(
	gameboyActionChannel <-chan GameboyActionMessage,
	cpuStateChannel chan<- CpuState,
	ppuStateChannel chan<- PpuState,
	apuStateChannel chan<- ApuState,
	memoryStateChannel chan<- []MemoryWrite,
) *Gameboy {

	// components
	bus := NewBus()
	cpu := NewCPU(bus)
	ppu := NewPPU(bus)
	apu := NewAPU()

	// load the bootrom once for all
	bootrom := loadBootRom(ROMS_URI)
	bus.AttachMemory(BOOT_ROM_MEMORY_NAME, BOOT_ROM_START, bootrom)

	// create the gameboy struct
	gb := &Gameboy{
		bus:                  bus,
		cpu:                  cpu,
		ppu:                  ppu,
		apu:                  apu,
		bootrom:              bootrom,
		gameboyActionChannel: gameboyActionChannel,
		cpuStateChannel:      cpuStateChannel,
		ppuStateChannel:      ppuStateChannel,
		apuStateChannel:      apuStateChannel,
		memoryStateChannel:   memoryStateChannel,
		joypad:               NewJoypad(),
	}

	// initialize memories and timer
	gb.initMemory()
	gb.initTimer(bus)

	// start the gameboy state machine listener
	go gb.stateMachineListener()

	return gb
}

// reset the gameboy state:
// - reset the gameboy state to "GB_STATE_NO_GAME_LOADED"
// - reset the ticks count and the timer
// - reset the bus and initialize the memories
// - reset the cpu, ppu and apu states
func (gb *Gameboy) reset() {
	// reset the gameboy state
	gb.state = GB_STATE_NO_GAME_LOADED

	// reset the ticks count and the timer
	gb.ticks = 0
	gb.timer.reset()

	// reset the bus and initialize the memories
	gb.bus.reset()
	gb.initMemory()

	// reset the cpu, ppu and apu states
	gb.cpu.reset()
	gb.ppu.reset()
	gb.apu.reset()
}

// initialize the memories and attach them to the bus
//   - HRAM: 127 bytes @ 0xFF80
//   - VRAM: 8KB bytes @ 0x8000
//   - WRAM: 8KB @ 0xC000
//   - I/O Registers: 128 bytes @ 0xFF00
func (gb *Gameboy) initMemory() {
	// initialize memories
	gb.vram = NewMemoryWithRandomData(0x2000) // VRAM (8KB)
	gb.wram = NewMemoryWithRandomData(0x2000) // WRAM (8KB)

	// attach memories to the CPU bus
	gb.bus.AttachMemory("Video RAM (VRAM)", 0x8000, gb.vram)
	gb.bus.AttachMemory("Working RAM (WRAM)", 0xC000, gb.wram)
}

// instantiate the timer and subscribe the gameboy to it
func (gb *Gameboy) initTimer(bus *Bus) {
	gb.timer = NewTimer(bus)
}

// initialize the bootrom @ 0x0000 - 0x00FF
func loadBootRom(uri string) *Memory {
	bootromData, err := LoadRom(uri + "/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return NewMemoryWithData(BOOT_ROM_LEN, bootromData)
}

// initialize the gameboy by creating the bus, bootrom, cpu, cartridge and the different memories
func (gb *Gameboy) LoadRom(romName string) {
	// reset components cpu, ppu & apu
	gb.cpu.reset() // all registers are randomized apart from PC which is set to 0x100
	gb.ppu.reset()
	gb.apu.reset()

	// reset vram & wram
	gb.vram.ResetWithRandomData()
	gb.wram.ResetWithRandomData()

	// load the cartridge rom
	gb.cartridge = NewCartridge(ROMS_URI, romName)
	gb.bus.AttachMemory("Cartridge ROM", 0x0000, gb.cartridge.rom)

	// set the gameboy state to paused
	gb.state = GB_STATE_NO_GAME_LOADED
}

// send updated state on the respective channels
func (gb *Gameboy) sendState() {
	gb.cpuStateChannel <- gb.cpu.getState()
	gb.ppuStateChannel <- gb.ppu.getState()
	gb.apuStateChannel <- gb.apu.getState()
	gb.memoryStateChannel <- *gb.bus.getMemoryWrites()
}

// tick the gameboy once
func (gb *Gameboy) tick() {
	gb.timer.Tick()
	gb.cpu.Tick()
	gb.ppu.Tick()
	gb.apu.Tick()
	gb.ticks++
}

func (gb *Gameboy) Tick() {
	// clear memory writes
	gb.cpu.bus.clearMemoryWrites()
	// tick the gameboy
	gb.tick()
	// report state changes on their respective channels
	gb.cpuStateChannel <- gb.cpu.getState()
	gb.ppuStateChannel <- gb.ppu.getState()
	gb.apuStateChannel <- gb.apu.getState()
	gb.memoryStateChannel <- *gb.bus.getMemoryWrites()
}

// run the bootrom and then the game
// When the ppu finishes to draw a frame, it sends the whole state to the frontend (cpu, ppu, apu, memory)
func (gb *Gameboy) run() {
	// run the gameboy until it is paused or stopped
	for gb.state == GB_STATE_RUNNING {
		// timing the gameboy @4.194304MHz
		tickStartTime := time.Now()
		// tick the gameboy
		gb.tick()
		tickDuration := time.Since(tickStartTime)
		// send the state to the frontend when the ppu finishes to draw a frame or when it reaches pixel (0, 144)
		if gb.ppu.dotX == 0 && gb.ppu.dotY == LCD_Y_RESOLUTION {
			gb.ppuStateChannel <- gb.ppu.getState()
			gb.bus.clearMemoryWrites()
		}
		// wait for the next tick
		time.Sleep(TICK_DURATION - tickDuration)
	}
}

// GAMEBOY STATE MACHINE

// listen to the gameboy state actions channel
func (gb *Gameboy) stateMachineListener() {
	gb.state = GB_STATE_NO_GAME_LOADED
	for {
		select {
		case state := <-gb.gameboyActionChannel:
			switch state.Action {
			case GB_ACTION_LOAD_GAME:
				gb.LoadRom(state.payload.(string))
			case GB_ACTION_PAUSE:
				if gb.state == GB_STATE_RUNNING {
					gb.state = GB_STATE_PAUSED
				}
			case GB_ACTION_RUN:
				if gb.state == GB_STATE_PAUSED {
					gb.state = GB_STATE_RUNNING
					go gb.run()
				}
			case GB_ACTION_RESET:
				gb.reset()
				gb.state = GB_STATE_NO_GAME_LOADED
				return
			}
		}
	}
}
