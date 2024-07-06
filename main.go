package main

import (
	"os"

	"github.com/codefrite/gameboy-go/gameboy"
)

func main() {
	// Create a new Gameboy
	gb := gameboy.NewGameboy("tetris.gb")
	gb.Run()
	/*
		// Create a channel to communicate with the Gameboy
		ch := make(chan *gameboy.GameboyState)
		go gb.Run(ch)
		for range ch {
			fmt.Printf("+ state: %v\n", (<-ch).CURR_CPU_STATE)
		}
	*/
	// Exit program properly
	os.Exit(0)
}
