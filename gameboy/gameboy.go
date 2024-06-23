package gameboy

import (
	"log"
	"os"
)

type Gameboy struct {
	cpu          *CPU
	bootrom      *Memory
	cartridge    *Cartridge
	vram         *Memory
	wram         *Memory
	io_registers *Memory
	hram         *Memory
	bus          *Bus
}

func NewGameboy() *Gameboy {
	g := &Gameboy{}
	g.initMemory()
	g.initBus()
	g.initCPU()
	return g
}

// Initializers

func (g *Gameboy) initMemory() {
	// Bootrom
	g.bootrom = NewMemory(0x0100)
	// VRAM
	g.vram = NewMemory(0x2000)
	// WRAM
	g.wram = NewMemory(0x2000)
	// I/O Registers
	g.io_registers = NewMemory(0x007F)
	// high ram
	g.hram = NewMemory(0x007F)
}

func (g *Gameboy) initBus() {
	g.bus = NewBus()
	g.bus.AttachMemory(0x8000, g.vram)
	g.bus.AttachMemory(0xC000, g.wram)
	g.bus.AttachMemory(0xFF00, g.io_registers)
	g.bus.AttachMemory(0xFF80, g.hram)
}

func (g *Gameboy) initCPU() {
	g.cpu = NewCPU(g.bus)
}

// Utility functions
func getBootRomData() []uint8 {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rom, err := LoadRom(currentDir + "/roms/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return rom
}

func (g *Gameboy) loadCartridge(uri string, name string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	g.cartridge = NewCartridge(currentDir+uri, name)
	g.bus.AttachMemory(0x0100, g.cartridge)
}

/*
 * Run the bootrom and then the game
 */
func (g *Gameboy) Run(uri string, name string) {
	bootRomData := getBootRomData()
	g.bus.AttachMemory(0x0000, g.bootrom)
	g.bus.WriteBlob(0x0000, bootRomData)
	g.loadCartridge(uri, name)
	g.cpu.Boot()
	g.cpu.Run()
}
