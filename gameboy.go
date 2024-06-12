package main

import (
	"log"
	"os"

	"codefrite.dev/emulators/gameboy"
)


func initGameboy() *gameboy.CPU {
	// 1. Init VRAM
	vram := gameboy.NewVRAM()

	// 2. Init WRAM
	wram := gameboy.NewWRAM()

	// 3. Init Cartridge
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// instantiate a new cartridge with tetris for the moment
	c := gameboy.NewCartridge(currentDir + "/roms", "tetris.gb")

	// 4. init BUS
	bus := gameboy.NewBus(c, vram, wram)

	// 4. instantiate a new CPU
	cpu := gameboy.NewCPU(bus)

	return cpu
}

func main() {
	// init gameboy
	cpu := initGameboy()

	// print cpu info
	//cpu.PrintRegisters()

	// set CPU PC to 0x100
	cpu.PC = 0x100
	
	// main game loop
	for {
		// cpu fetches the opcode located at the program counter
		err := cpu.Step()
		if err != nil {
			break
		}
	}
}