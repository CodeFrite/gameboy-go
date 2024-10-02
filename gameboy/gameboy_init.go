package gameboy

const CRYSTAL_FREQUENCY = 1

/**
 * initializes the gameboy by creating the bus, bootrom, cpu, cartridge and the different memories
 */
func (gb *Gameboy) LoadRom(romName string) {
	// buses
	gb.cpuBus = NewBus()
	gb.ppuBus = NewBus()
	// cpu
	gb.cpu = NewCPU(gb.cpuBus)
	gb.cpu.Init()
	// ppu
	gb.ppu = NewPPU(gb.cpu, gb.ppuBus)
	// apu
	gb.apu = NewAPU()
	// cartridge
	gb.cartridge = NewCartridge("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms", romName)
	gb.cpuBus.AttachMemory("Cartridge ROM", 0x0000, gb.cartridge.rom)
	// vram & wram
	gb.initMemory()
	gb.initTimer()
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
	gb.vram = NewMemoryWithRandomData(0x2000) // VRAM (8KB)
	gb.wram = NewMemory(0x2000)               // WRAM (8KB)

	// attach memories to the bus
	gb.cpuBus.AttachMemory("Video RAM (VRAM)", 0x8000, gb.vram)
	gb.cpuBus.AttachMemory("Working RAM (WRAM)", 0xC000, gb.wram)
}

func (gb *Gameboy) initTimer() {
	gb.crystal = NewTimer(CRYSTAL_FREQUENCY)
	gb.crystal.Subscribe(gb)
}

// TODO: initialize the joypad channel to listen to incoming joypad events from user
