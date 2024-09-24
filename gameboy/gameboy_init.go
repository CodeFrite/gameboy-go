package gameboy

import "log"

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

/**
 * initializes the gameboy by creating the bus, bootrom, cpu, cartridge and the different memories
 */
func (gb *Gameboy) init(romName string) {
	gb.initCpuBus() // the bus is created first to be able to attach the memories to it
	gb.initPpuBus()
	gb.initBootRom()
	gb.initCPU()
	gb.initPPU()
	gb.initCartridge("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms", romName)
	gb.initMemory()
	gb.initTimer()
}

/**
 * initializes the bus
 */
func (gb *Gameboy) initCpuBus() {
	gb.cpuBus = NewBus()
}

func (gb *Gameboy) initPpuBus() {
	gb.ppuBus = NewBus()
}

/**
 * initializes the bootrom @ 0x0000
 */
func (gb *Gameboy) initBootRom() {
	bootRomData := gb.getBootRomData()
	gb.bootrom = NewMemory(0x100)
	gb.cpuBus.AttachMemory("Boot ROM", 0x0000, gb.bootrom)
	gb.cpuBus.WriteBlob(0x0000, bootRomData)
}

/**
 * initializes the CPU: creates a new CPU, passes it the bus and attaches the IE register to the bus @ 0xFFFF
 */
func (gb *Gameboy) initCPU() {
	gb.cpu = NewCPU(gb.cpuBus)
	// IE register 1byte set by the CPU
	gb.cpuBus.AttachMemory("IE", 0xFFFF, gb.cpu.IE)
}

func (gb *Gameboy) initPPU() {
	gb.ppu = NewPPU(gb.cpu, gb.ppuBus)
}

/**
 * initializes the cartridge: creates a new cartridge, attaches the ROM to the bus @ 0x0100
 * Please note that the first 0x100 correspond to the bootrom which i chosed to load separately from a different file (see cartridge.go)
 */
func (gb *Gameboy) initCartridge(uri string, name string) {
	gb.cartridge = NewCartridge(uri, name)
	gb.cpuBus.AttachMemory("Cartridge", 0x0100, gb.cartridge.rom)
}

/**
 * initializes the memories and attaches them to the bus
 * HRAM: 127 bytes @ 0xFF80
 * VRAM: 8KB bytes @ 0x8000
 * WRAM: 8KB @ 0xC000
 * I/O Registers: 128 bytes @ 0xFF00
 */
func (gb *Gameboy) initMemory() {
	// initialize memories
	gb.hram = NewMemory(0x7F)                 // High RAM (127 bytes)
	gb.vram = NewMemoryWithRandomData(0x2000) // VRAM (8KB)
	gb.wram = NewMemory(0x2000)               // WRAM (8KB)
	gb.io_registers = NewMemory(0x0080)       // I/O Registers (128 bytes)

	// attach memories to the bus
	gb.cpuBus.AttachMemory("High RAM (HRAM)", 0xFF80, gb.hram)
	gb.cpuBus.AttachMemory("Video RAM (VRAM)", 0x8000, gb.vram)
	gb.cpuBus.AttachMemory("Working RAM (WRAM)", 0xC000, gb.wram)
	gb.cpuBus.AttachMemory("I/O Registers", 0xFF00, gb.io_registers)
}

func (gb *Gameboy) initTimer() {
	gb.crystal = NewTimer(100)
	gb.crystal.Subscribe(gb.cpu)
	gb.crystal.Subscribe(gb.ppu)
}

/**
 * loads data content as []uint8 from a rom file
 * TODO: no need for this function: it already calls a utility function from utilies.go
 */
func (gb *Gameboy) getBootRomData() []uint8 {
	romData, err := LoadRom("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return romData
}
