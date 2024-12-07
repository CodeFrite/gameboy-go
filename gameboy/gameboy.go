package gameboy

import (
	"sync"
)

// the gameboy is composed out of a CPU, memories (ram & registers), a cartridge and a bus
type Gameboy struct {
	busyChannel chan bool
	crystal     *Timer // crystal oscillator running at 4.194304MHz
	ticks       uint64
	cpuBus      *Bus
	ppuBus      *Bus
	cpu         *CPU
	ppu         *PPU
	apu         *APU
	vram        *Memory
	wram        *Memory
	cartridge   *Cartridge // 0x0000-0x7FFF (32KB switchable) - Cartridge ROM
	joypad      *Joypad

	// state channels (sharing concrete types to avoid pointer values being changed before being sent to the frontend by the server)
	cpuStateChannel    chan<- CpuState
	ppuStateChannel    chan<- PpuState
	apuStateChannel    chan<- ApuState
	memoryStateChannel chan<- []MemoryWrite
	joypadStateChannel <-chan JoypadState
}

// creates a new gameboy struct
func NewGameboy(cpuStateChannel chan<- CpuState, ppuStateChannel chan<- PpuState, apuStateChannel chan<- ApuState, memoryStateChannel chan<- []MemoryWrite, joypadStateChannel <-chan JoypadState) *Gameboy {
	gb := &Gameboy{
		busyChannel:        make(chan bool, 1),
		cpuStateChannel:    cpuStateChannel,
		ppuStateChannel:    ppuStateChannel,
		apuStateChannel:    apuStateChannel,
		memoryStateChannel: memoryStateChannel,
		joypadStateChannel: joypadStateChannel,
		joypad:             NewJoypad(joypadStateChannel),
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
