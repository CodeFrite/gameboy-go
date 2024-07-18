package gameboy

import "log"

type Gameboy struct {
	cpu          *CPU
	bootrom      *Memory
	cartridge    *Cartridge
	vram         *Memory
	wram         *Memory
	io_registers *Memory
	hram         *Memory
	bus          *Bus
	state        *GameboyState
}

func NewGameboy(romName string) *Gameboy {
	gb := &Gameboy{}
	gb.init(romName)
	return gb
}

// Initializers

// Memory Map:
// 0x0000-0x00FF (256 bytes) - Boot ROM
// 0x0100-0x7FFF (32KB switchable) - Cartridge ROM
// 0x8000-0x9FFF (8KB Video RAM) - VRAM
// 0xA000-0xBFFF (8KB switchable) - External RAM
// 0xC000-0xDFFF (8KB internal RAM) - WRAM
// 0xE000-0xFDFF (7.5KB Echo RAM) - Echo RAM
// 0xFE00-0xFE9F (160 bytes) - OAM
// (0xFEA0-0xFEFF (96 bytes) - Not used)
// 0xFF00-0xFF7F (128 bytes) - I/O Registers
// 0xFF80-0xFFFE (127 bytes) - High RAM
// 0xFFFF (1 byte) - Interrupt Enable Register

func (gb *Gameboy) init(romName string) *GameboyState {
	gb.initBus()
	gb.initBootRom()
	gb.initCPU()
	gb.initCartridge("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/", romName)
	gb.initMemory()
	gb.saveCurrentState()
	return gb.state
}

func (gb *Gameboy) initBus() {
	gb.bus = NewBus()
}

func (gb *Gameboy) initBootRom() {
	bootRomData := gb.getBootRomData()
	gb.bootrom = NewMemory(0x100)
	gb.bus.AttachMemory("Boot ROM", 0x0000, gb.bootrom)
	gb.bus.WriteBlob(0x0000, bootRomData)
}

func (gb *Gameboy) initCPU() {
	gb.cpu = NewCPU(gb.bus)
	gb.state = &GameboyState{
		PREV_CPU_STATE: nil,
		CURR_CPU_STATE: nil,
		INSTR:          nil,
		MEMORY_WRITES: []MemoryWrite{{
			Address: 0,
			Data:    []string{},
		}},
	}
	// IE register 1byte set by the CPU
	gb.bus.AttachMemory("IE", 0xFFFF, gb.cpu.IE)
}

func (gb *Gameboy) initCartridge(uri string, name string) {
	gb.cartridge = NewCartridge(uri, name)
	gb.bus.AttachMemory("Cartridge", 0x0100, gb.cartridge.rom)
}

func (gb *Gameboy) initMemory() {
	// initialize memories
	gb.hram = NewMemory(0x7F)                 // High RAM (127 bytes)
	gb.vram = NewMemoryWithRandomData(0x2000) // VRAM (8KB)
	gb.wram = NewMemory(0x2000)               // WRAM (8KB)
	gb.io_registers = NewMemory(0x0080)       // I/O Registers (128 bytes)

	// attach memories to the bus
	gb.bus.AttachMemory("High RAM (HRAM)", 0xFF80, gb.hram)
	gb.bus.AttachMemory("Video RAM (VRAM)", 0x8000, gb.vram)
	gb.bus.AttachMemory("Working RAM (WRAM)", 0xC000, gb.wram)
	gb.bus.AttachMemory("I/O Registers", 0xFF00, gb.io_registers)
}

//! Public interface

/*
 * Run the bootrom and then the game
 */
func (gb *Gameboy) Run() {
	for {
		gb.Step()
		if gb.cpu.halted {
			break
		}
	}
}

func (gb *Gameboy) Step() *GameboyState {
	// execute the instruction
	gb.cpu.step()
	// save the state
	gb.saveCurrentState()
	return gb.state
}

// Utility functions

// load data content as []uint8 from a rom file
func (gb *Gameboy) getBootRomData() []uint8 {
	romData, err := LoadRom("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return romData
}
