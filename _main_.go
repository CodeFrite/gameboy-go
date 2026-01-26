package main_

import (
	"bufio"
	"fmt"
	"os"
)

/* create a new debugger, initialize it with a ROM and exit the program when done */
func main() {
	// Print the welcome message
	fmt.Println("Welcome to gameboy-go emulator")

	// Create a new reader to read the user input
	reader := bufio.NewReader(os.Stdin)

	// Select the mode: gameboy or debugger
	fmt.Println("Would you like to run the emulator in debug mode? (y/n)")

	// Read the user input
	var debugMode string
	var err error

	debugMode, err = reader.ReadString('\n')
	for debugMode != "y" && debugMode != "n" {
		if err != nil {
			fmt.Println("Error reading user input:", err)
		}
		fmt.Println("Would you like to run the emulator in debug mode? (y/n)")
		fmt.Scan(&debugMode)
	}
	/*
		// user does not want to run the emulator in debug mode
		if debugMode == "n" {
			// Create a new Gameboy
			gb := gameboy.NewGameboy(nil, nil, nil, nil)
			gb.LoadRom("Tetris.gb")
			gb.Run()

		} else if debugMode == "y" {
			// user wants to run the emulator in debug mode

			// Create a new Debugger
			db := gameboy.NewDebugger(nil, nil, nil, nil)

			// Initialize the Debugger with the ROM
			db.LoadRom("Tetris.gb")

			// main loop that asks the user to either if he wants to run, step, pause, add a breakpoint, delete a breakpoint, print the state or exit the debugger
			exit := false

			for !exit {
				// ask the user what he wants to do:
				var action string
				var err error

				fmt.Print("\nWhat do you want to do?? (run/R/r, step/S/s, pause/P/p, break/B/b, delete/D/d, list/L/l, state/X/x, exit/E/e): ")
				action, err = reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading user input:", err)
				}

				// remove the delimiter from the user input
				action = action[:len(action)-1]

				// - continue the execution (run/R/r)
				if action == "run" || action == "R" || action == "r" {
					fmt.Println("Running the debugger...")
					db.Run()
				}
				// - step into the next instruction (step/S/s)
				if action == "step" || action == "S" || action == "s" {
					fmt.Println("Stepping into the next instruction...")
					db.Step()
				}
				// - pause the execution (pause/P/p)
				if action == "pause" || action == "P" || action == "p" {
					fmt.Println("Pausing the execution...")
					continue // should be implemented
				}
				// - add a breakpoint (break/B/b)
				if action == "break" || action == "B" || action == "b" {
					fmt.Println("Adding a breakpoint...")
					// ask the user for the address of the breakpoint
					var addr string
					addr, err = reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error reading breakpoint address:", err)
					}
					addr16, err := strconv.ParseUint(addr[:len(addr)-1], 16, 16)
					if err != nil {
						fmt.Println("Error parsing breakpoint address to uint16:", err)
					}
					db.AddBreakPoint(uint16(addr16))
				}
				// - delete a breakpoint (delete/D/d)
				if action == "delete" || action == "D" || action == "d" {
					fmt.Println("Deleting a breakpoint...")
					// ask the user for the address of the breakpoint
					var addr uint16
					_, err = fmt.Scan("Enter the address of the breakpoint to delete: ", &addr)
					if err != nil {
						fmt.Println("Error reading breakpoint address:", err)
					}
					db.RemoveBreakPoint(addr)
				}
				// - list all breakpoints (list/L/l)
				if action == "list" || action == "L" || action == "l" {
					fmt.Println("Listing all breakpoints:", db.GetBreakPoints())
				}
				// - print the current state (state/X/x)
				if action == "state" || action == "X" || action == "x" {
					//db.PrintCPUState()
					//db.PrintInstruction()
					//db.PrintMemoryProperties()
				}
				// - exit the debugger (exit/E/e)
				if action == "exit" || action == "E" || action == "e" {
					fmt.Println("Exiting the debugger...")
					exit = true

				}
			}
		}
	*/
	// Exit program properly
	os.Exit(0)
}
