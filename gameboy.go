package main

import (
	"log"
	"os"

	"codefrite.dev/emulators/gameboy"
)

// TODO> Since the boot rom is not run on init, I should make sure to init the gameboy (CPU, BUS, WRAM, VRAM, etc) properly.
// TODO> There is more to it than setting the PC to 0x0100 and SP to 0xFFFE.
func initCPU(cpu *gameboy.CPU) {
	// set CPU PC to 0x100 to skip the boot rom and start executing the game
	cpu.PC = 0x0100

	// ! the SP is normally set to 0xFFFE by the boot rom. We will set it manually for the moment
	cpu.SP = 0xFFFE
}
func initGameboy() *gameboy.CPU {
	// 1.Init RAM
	// VRAM
	vram := gameboy.NewMemory(0x2000)

	// WRAM
	wram := gameboy.NewMemory(0x2000)

	// I/O Registers
	io_registers := gameboy.NewMemory(0x007F)

	// 2. Init Cartridge
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// load tetris for the moment
	c := gameboy.NewCartridge(currentDir+"/roms", "tetris.gb")

	// 3. init BUS
	bus := gameboy.NewBus()
	bus.AttachMemory(c, 0x0000)
	bus.AttachMemory(vram, 0x8000)
	bus.AttachMemory(wram, 0xC000)
	bus.AttachMemory(io_registers, 0xFF00)

	// 4. instantiate a new CPU
	cpu := gameboy.NewCPU(bus)

	return cpu
}

func main() {
	// init gameboy
	cpu := initGameboy()

	// init CPU since the boot rom is not run
	initCPU(cpu)

	// print cpu info
	//cpu.PrintRegisters()

	// set CPU PC to 0x100 to skip the boot rom and start executing the game
	cpu.PC = 0x0100

	// ! the SP is normally set to 0xFFFE by the boot rom. We will set it manually for the moment
	cpu.SP = 0xFFFE

	// main game loop
	cpu.Run()

	// Exit program properly
	os.Exit(0)
}
