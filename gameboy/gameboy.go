package gameboy

import (
	"log"
	"sync"
)

// CONSTANTS
const (
	CRYSTAL_FREQUENCY           = 50000
	BOOT_ROM_MEMORY_NAME        = "Boot ROM"
	BOOT_ROM_START       uint16 = 0x0000
	BOOT_ROM_LEN         uint16 = 0x0100
	ROMS_URI                    = "/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms"
)

// the gameboy is composed out of a CPU, memories (ram & registers), a cartridge and a bus
type Gameboy struct {
	// state
	ticks       uint64    // number of ticks since the gameboy started
	busyChannel chan bool // used to prevent multiple clock ticks being processed by cpu/ppu/apu at the same time

	// components
	crystal   *Timer // crystal oscillator running at 4.194304MHz
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

// creates a new gameboy struct
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
		busyChannel:        make(chan bool, 1),
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
	gb.initTimer()

	return gb
}

// initializes the bootrom @ 0x0000 - 0x00FF
func loadBootRom(uri string) *Memory {
	bootromData, err := LoadRom(uri + "/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return NewMemoryWithData(BOOT_ROM_LEN, bootromData)
}

// initializes the memories and attaches them to the bus
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

// instantiates the timer and subscribes the gameboy to it
func (gb *Gameboy) initTimer() {
	gb.crystal = NewTimer(CRYSTAL_FREQUENCY)
	gb.crystal.Subscribe(gb)
}

// initializes the gameboy by creating the bus, bootrom, cpu, cartridge and the different memories
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

// runs the bootrom and then the game
func (gb *Gameboy) Run() {
	gb.crystal.Start()
}

// executes the next instruction
func (gb *Gameboy) Step() {
	// tick the crystal oscillator once
	gb.crystal.Tick()
}

// the gameboy ticks in parallel the cpu, ppu and apu and wait for these calls all to end using a wait group
func (gb *Gameboy) onTick() {
	// busy channel to prevent multiple ticks at the same time
	gb.busyChannel <- true

	// clear memory writes
	gb.cpu.bus.mmu.clearMemoryWrites()

	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	// tick the cpu 1 out of 3 ticks
	if gb.ticks%3 == 0 {
		wg.Add(1)
		go func() {
			gb.cpu.onTick()
			wg.Done()
		}()
	}

	// tick the ppu if FF40 bit 7 is set
	lcdc := gb.cpu.bus.Read(0xFF40)
	lcd_ppu_enabled := lcdc&0x80 == 0x80

	if lcd_ppu_enabled {
		wg.Add(1)
		go func() {
			gb.ppu.onTick()
			wg.Done()
		}()
	}

	// tick the apu
	wg.Add(1)
	go func() {
		gb.apu.onTick()
		wg.Done()
	}()

	// wait for all goroutines to finish
	wg.Wait()

	// now we can send the state to the respective channels

	// get the cpu, ppu and apu states and send them to the respective channels
	if gb.cpuStateChannel != nil {
		gb.cpuStateChannel <- gb.cpu.getState()
	}
	if gb.ppuStateChannel != nil {
		gb.ppuStateChannel <- gb.ppu.getState()
	}
	if gb.apuStateChannel != nil {
		gb.apuStateChannel <- gb.apu.getState()
	}
	if gb.memoryStateChannel != nil {
		gb.memoryStateChannel <- *gb.cpuBus.mmu.getMemoryWrites()
	}
	gb.ticks++
	<-gb.busyChannel
}
