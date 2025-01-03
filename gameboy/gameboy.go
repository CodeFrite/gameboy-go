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
	STATE_RUNNING GameBoyState = "running"
	STATE_PAUSED  GameBoyState = "paused"
	STATE_STOPPED GameBoyState = "stopped"
)

type GameBoyState string

// the gameboy is composed out of a CPU, memories (ram & registers), a cartridge and a bus
type Gameboy struct {
	// state
	ticks uint64       // number of ticks since the gameboy started
	state GameBoyState // current state of the gameboy

	// components
	timer     *Timer // Gameboy Timer (DIV, TIMA, TMA, TAC)
	cpuBus    *Bus
	ppuBus    *Bus
	cpu       *CPU
	ppu       *PPU
	apu       *APU
	bootrom   *Memory    // 0x0000-0x00FF: (256 bytes) - Boot ROM
	cartridge *Cartridge // Cartridge ROM (32KB) [0x0000-0x7FFF]
	vram      *Memory    // Video RAM (8KB) [0x8000-0x9FFF]
	wram      *Memory    // Working RAM (8KB) [0xC000-0xDFFF]
	joypad    *Joypad

	// state channels (sharing concrete types to avoid pointer values being changed before being sent to the frontend by the server)
	cpuStateChannel    chan<- CpuState
	ppuStateChannel    chan<- PpuState
	apuStateChannel    chan<- ApuState
	memoryStateChannel chan<- []MemoryWrite
	joypadStateChannel <-chan JoypadState
}

// create a new gameboy struct
func NewGameboy(cpuStateChannel chan<- CpuState, ppuStateChannel chan<- PpuState, apuStateChannel chan<- ApuState, memoryStateChannel chan<- []MemoryWrite, joypadStateChannel <-chan JoypadState) *Gameboy {
	// components
	cpuBus := NewBus()
	ppuBus := NewBus()
	cpu := NewCPU(cpuBus)
	ppu := NewPPU(cpu, ppuBus)
	apu := NewAPU()

	// load the bootrom once for all
	bootrom := loadBootRom(ROMS_URI)
	cpuBus.AttachMemory(BOOT_ROM_MEMORY_NAME, BOOT_ROM_START, bootrom)

	// create the gameboy struct
	gb := &Gameboy{
		cpuBus:             cpuBus,
		ppuBus:             ppuBus,
		cpu:                cpu,
		ppu:                ppu,
		apu:                apu,
		bootrom:            bootrom,
		cpuStateChannel:    cpuStateChannel,
		ppuStateChannel:    ppuStateChannel,
		apuStateChannel:    apuStateChannel,
		memoryStateChannel: memoryStateChannel,
		joypadStateChannel: joypadStateChannel,
		joypad:             NewJoypad(joypadStateChannel),
	}

	// initialize memories and timer
	gb.initMemory()
	gb.initTimer(cpuBus)

	return gb
}

// initialize the bootrom @ 0x0000 - 0x00FF
func loadBootRom(uri string) *Memory {
	bootromData, err := LoadRom(uri + "/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return NewMemoryWithData(BOOT_ROM_LEN, bootromData)
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

	// attach memories to the bus
	gb.cpuBus.AttachMemory("Video RAM (VRAM)", 0x8000, gb.vram)
	gb.cpuBus.AttachMemory("Working RAM (WRAM)", 0xC000, gb.wram)
}

// instantiate the timer and subscribe the gameboy to it
func (gb *Gameboy) initTimer(bus *Bus) {
	gb.timer = NewTimer(bus)
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
	gb.cpuBus.AttachMemory("Cartridge ROM", 0x0000, gb.cartridge.rom)
}

// send state
func (gb *Gameboy) sendState() {
	gb.cpuStateChannel <- gb.cpu.getState()
	gb.ppuStateChannel <- gb.ppu.getState()
	gb.apuStateChannel <- gb.apu.getState()
	gb.memoryStateChannel <- *gb.cpuBus.mmu.getMemoryWrites()
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
	gb.cpu.bus.mmu.clearMemoryWrites()
	// tick the gameboy
	gb.tick()
	// report state changes on their respective channels
	gb.cpuStateChannel <- gb.cpu.getState()
	gb.ppuStateChannel <- gb.ppu.getState()
	gb.apuStateChannel <- gb.apu.getState()
	gb.memoryStateChannel <- *gb.cpuBus.mmu.getMemoryWrites()
}

// run the bootrom and then the game
// When the ppu finishes to draw a frame, it sends the whole state to the frontend (cpu, ppu, apu, memory)
func (gb *Gameboy) Run() {
	// clear memory writes
	gb.cpu.bus.mmu.clearMemoryWrites()
	// run the gameboy until it is paused or stopped
	for gb.state == STATE_RUNNING {
		// timing the gameboy @4.194304MHz
		tickStartTime := time.Now()
		// tick the gameboy
		gb.tick()
		tickDuration := time.Since(tickStartTime)
		// send the state to the frontend when the ppu finishes to draw a frame or when it reaches pixel (0, 144)
		if gb.ppu.dotX == 0 && gb.ppu.dotY == LCD_Y_RESOLUTION {
			gb.sendState()
			gb.cpu.bus.mmu.clearMemoryWrites()
		}
		// wait for the next tick
		time.Sleep(TICK_DURATION - tickDuration)
	}
	// report state changes on their respective channels
	gb.sendState()
}

// pauses the gameboy
func (gb *Gameboy) Pause() {
	if gb.state == STATE_RUNNING {
		gb.state = STATE_PAUSED
	}
}

// resumes the gameboy
func (gb *Gameboy) Resume() {
	if gb.state == STATE_PAUSED {
		gb.state = STATE_RUNNING
	}
}

// stop the gameboy
func (gb *Gameboy) Stop() {
	if gb.state == STATE_RUNNING || gb.state == STATE_PAUSED {
		gb.state = STATE_STOPPED
	}
}
