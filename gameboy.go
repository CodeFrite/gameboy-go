package main

import (
	"log"
	"os"

	"codefrite.dev/emulators/gameboy"
)

func initGameboy() *gameboy.CPU {
	// 1. Init VRAM
	vram := gameboy.NewMemory(0x2000)

	// 2. Init WRAM
	wram := gameboy.NewMemory(0x2000)

	// 3. Init Cartridge
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// load tetris for the moment
	c := gameboy.NewCartridge(currentDir+"/roms", "tetris.gb")

	// 4. init BUS
	bus := gameboy.NewBus()
	bus.AttachMemory(c, 0x0000)
	bus.AttachMemory(vram, 0x8000)
	bus.AttachMemory(wram, 0xC000)

	// 5. instantiate a new CPU
	cpu := gameboy.NewCPU(bus)

	return cpu
}

func main() {
	// init gameboy
	cpu := initGameboy()

	// print cpu info
	//cpu.PrintRegisters()

	// set CPU PC to 0x100 to skip the boot rom and start executing the game
	cpu.PC = 0x0100

	// main game loop
	cpu.Run()

	return
}
