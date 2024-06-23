package main

import (
	"os"

	"codefrite.dev/emulators/gameboy"
)

func main() {
	// Create a new Gameboy
	gb := gameboy.NewGameboy()
	gb.Run("/roms/", "tetris.gb")

	// Exit program properly
	os.Exit(0)
}
