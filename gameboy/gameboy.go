package gameboy

import (
	"fmt"
	"log"
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
	state        *GameboyState
}

func NewGameboy(romName string) *Gameboy {
	gb := &Gameboy{}
	gb.init(romName)
	return gb
}

// Initializers

func (gb *Gameboy) init(romName string) *GameboyState {
	gb.initBus()
	gb.initCPU()
	gb.initBootRom()
	gb.loadCartridge("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/", romName)
	gb.initMemory()
	gb.connectMemoryToBus()
	gb.cpu.prefetch()
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", gb.cpu.IR)), gb.cpu.Prefixed)
	gb.saveCurrentState(instruction)
	return gb.state
}

func (gb *Gameboy) initBus() {
	gb.bus = NewBus()
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
}

func (gb *Gameboy) initBootRom() {
	bootRomData := gb.getBootRomData()
	gb.bootrom = NewMemory(0x100)
	gb.bus.AttachMemory(0x0000, gb.bootrom)
	gb.bus.WriteBlob(0x0000, bootRomData)
}

func (gb *Gameboy) loadCartridge(uri string, name string) {
	gb.cartridge = NewCartridge(uri, name)
}

func (gb *Gameboy) initMemory() {
	// initialize memories
	gb.vram = NewMemory(0x2000)         // VRAM
	gb.wram = NewMemory(0x2000)         // WRAM
	gb.io_registers = NewMemory(0x007F) // I/O Registers
	gb.hram = NewMemory(0x007F)         // high ram
}

func (gb *Gameboy) connectMemoryToBus() {
	// attach memories to the bus
	gb.bus.AttachMemory(0x0100, gb.cartridge)
	gb.bus.AttachMemory(0x8000, gb.vram)
	gb.bus.AttachMemory(0xC000, gb.wram)
	gb.bus.AttachMemory(0xFF00, gb.io_registers)
	gb.bus.AttachMemory(0xFF80, gb.hram)
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
	// prefetch the instruction to get the correct output
	gb.cpu.prefetch()
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", gb.cpu.IR)), gb.cpu.Prefixed)
	// execute the instruction
	gb.cpu.step()
	// save the state
	gb.saveCurrentState(instruction)
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
