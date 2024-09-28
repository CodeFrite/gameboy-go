package gameboy

import "log"

// the gameboy is composed out of a CPU, memories (ram & registers), a cartridge and a bus
type Gameboy struct {
	cpu       *CPU
	ppu       *PPU
	cartridge *Cartridge // 0x0000-0x7FFF (32KB switchable) - Cartridge ROM
	vram      *Memory
	wram      *Memory
	cpuBus    *Bus
	ppuBus    *Bus
	crystal   *Timer // crystal oscillator running at 4.194304MHz

	// state channels
	cpuStateChannel chan<- *CpuState // v0.4.0
	//ppuStateChannel chan<- *PpuState // v0.4.1
	//apuStateChannel chan<- *ApuState // v0.4.2
	//joypadStateChannel <-chan *JoypadState // v0.4.3
}

// creates a new gameboy struct
func NewGameboy(cpuStateChannel chan<- *CpuState) *Gameboy {
	gb := &Gameboy{
		cpuStateChannel: cpuStateChannel,
	}
	return gb
}

//! Public interface

// runs the bootrom and then the game
func (gb *Gameboy) Run() {
	for {
		if gb.cpu.halted || gb.cpu.stopped {
			break
		}
		gb.Step()
	}
}

// executes the next instruction
func (gb *Gameboy) Step() {
	// check if the CPU is halted or stopped
	if gb.cpu.halted || gb.cpu.stopped {
		return
	}
	// execute the instruction
	gb.cpu.Step()
}

func (gb *Gameboy) Crash(err error) {
	log.Fatal("Gameboy crashed")

}
