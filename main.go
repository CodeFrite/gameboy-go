package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/codefrite/gameboy-go/gameboy"
	"github.com/codefrite/gameboy-go/gui"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// Start CPU profiling
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer cpuProfile.Close()

	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start CPU profile: %v\n", err)
		os.Exit(1)
	}
	defer pprof.StopCPUProfile()

	// Start memory profiling at the end
	defer func() {
		memProfile, err := os.Create("mem.prof")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create memory profile: %v\n", err)
			return
		}
		defer memProfile.Close()

		if err := pprof.WriteHeapProfile(memProfile); err != nil {
			fmt.Fprintf(os.Stderr, "Could not write memory profile: %v\n", err)
		}
	}()

	// Instantiate a GUI
	gui, err := gui.NewGUI()

	// On error simply exit
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Instantiate a new gameboy
	gbActionMessageChannel := make(chan gameboy.GameboyActionMessage, 1)
	gbCpuStateChannel := make(chan gameboy.CpuState, 1)
	gbPpuStateChannel := make(chan gameboy.PpuState, 1)
	//gbApuStateChannel := make(chan gameboy.ApuState, 1)
	//gbMemoryStateChannel := make(chan []gameboy.MemoryWrite, 1)

	gb := gameboy.NewGameboy(gbActionMessageChannel, gbCpuStateChannel, gbPpuStateChannel, nil, nil)

	// loading tetris rom
	gb.LoadRom("tetris.gb")

	// running the gameboy
	gbActionMessageChannel <- gameboy.GameboyActionMessage{
		Action:  gameboy.GB_ACTION_RUN,
		Payload: nil,
	}

	running := true

	loopCount := 0
	seconds := 0

	for running {
		loopCount++
		if loopCount%(1398101) == 0 {
			seconds++
		}
		// Handle events (keyboard, quit, etc.)
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					fmt.Printf("Key pressed: %s\n", sdl.GetKeyName(e.Keysym.Sym))
					if e.Keysym.Sym == sdl.K_ESCAPE {
						running = false
					}
				}
			}
		}

		// Wait for a signal from the gameboy (non-blocking)
		var cycle int64 = 0
		select {
		case cpuState := <-gbCpuStateChannel:
			if cycle++; cycle%100000 == 0 {
				fmt.Println("CPU TIC (", seconds, "):", cpuState.PC)
			}
		case ppuState := <-gbPpuStateChannel:
			// Clear screen and draw every frame
			gui.LCDClear()

			//fmt.Println("Received PPU state:", ppuState.IMAGE)
			// Drawing the current PPU image to the screen
			gui.LCDDrawImage(ppuState.IMAGE)
			// nothing new to draw
			gui.LCDPresent()
		default:
			time.Sleep(time.Microsecond) // Prevent busy-wait
		}
	}
}
