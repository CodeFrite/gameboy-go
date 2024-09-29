package gameboy

import (
	"fmt"
	"sync"
)

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
	ppuStateChannel chan<- *PpuState // v0.4.1
	//apuStateChannel chan<- *ApuState // v0.4.2
	//joypadStateChannel <-chan *JoypadState // v0.4.3
}

// creates a new gameboy struct
func NewGameboy(cpuStateChannel chan<- *CpuState, ppuStateChannel chan<- *PpuState) *Gameboy {
	gb := &Gameboy{
		cpuStateChannel: cpuStateChannel,
		ppuStateChannel: ppuStateChannel,
	}
	return gb
}

//! Public interface

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
	fmt.Println("Gameboy> tick")
	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	// tick the cpu
	wg.Add(1)
	go func() {
		gb.cpu.onTick()
		wg.Done()
	}()
	// tick the ppu
	wg.Add(1)
	go func() {
		gb.ppu.onTick()
		wg.Done()
	}()
	// tick the apu
	//gb.apu.onTick()

	// wait for all goroutines to finish
	fmt.Println("Gameboy.onTick> waiting for goroutines to finish")
	wg.Wait()

	// get the cpu, ppu and apu states and send them to the respective channels
	fmt.Println("Gameboy.onTick> sending cpu state back to the channel")
	gb.cpuStateChannel <- gb.cpu.getState()
	fmt.Println("Gameboy.onTick> terminating ...")
	gb.ppuStateChannel <- gb.ppu.getState()
	//gb.apuStateChannel <- gb.currApuState()

}
