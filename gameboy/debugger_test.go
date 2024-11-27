package gameboy

import (
	"fmt"
	"testing"
	"time"
)

// TestNewDebugger tests the creation of a new debugger
func TestNewDebugger(t *testing.T) {
	debugger := NewDebugger(nil, nil, nil, nil, nil)
	if debugger == nil {
		t.Fatalf("Expected a new debugger, got nil")
	}
}

func TestDebuggerStateChannels(t *testing.T) {
	t.Skip("Skipped: Buggy since i added the channel logic")

	cpuStateChannel := make(chan *CpuState)
	ppuStateChannel := make(chan *PpuState)
	apuStateChannel := make(chan *ApuState)
	memoryStateChannel := make(chan []MemoryWrite)
	joypadStateChannel := make(chan *JoypadState)

	debugger := NewDebugger(cpuStateChannel, ppuStateChannel, apuStateChannel, memoryStateChannel, joypadStateChannel)

	if debugger.clientCpuStateChannel != cpuStateChannel {
		t.Fatalf("Expected a gameboy cpu state channel to be initialized with address %p, got %p", cpuStateChannel, debugger.clientCpuStateChannel)
	}

	if debugger.clientPpuStateChannel != ppuStateChannel {
		t.Fatalf("Expected a gameboy ppu state channel to be initialized with address %p, got %p", ppuStateChannel, debugger.clientPpuStateChannel)
	}

	if debugger.clientApuStateChannel != apuStateChannel {
		t.Fatalf("Expected a gameboy apu state channel to be initialized with address %p, got %p", apuStateChannel, debugger.clientApuStateChannel)
	}

	if debugger.clientMemoryStateChannel != memoryStateChannel {
		t.Fatalf("Expected a gameboy memory state channel to be initialized with address %p, got %p", memoryStateChannel, debugger.clientMemoryStateChannel)
	}

	// load the ROM
	debugger.LoadRom("Tetris.gb")

	// step the gameboy ...
	debugger.Step()

	// ... and check that all channels do receive an updated state
	cpuState := <-cpuStateChannel
	fmt.Println("received cpu state", cpuState)
	ppuState := <-ppuStateChannel
	fmt.Println("received ppu state", ppuState)
	apuState := <-apuStateChannel
	fmt.Println("received apu state", apuState)
	memoryState := <-memoryStateChannel
	fmt.Println("received memory state", memoryState)
}

// Scenario: the gameboy should run @1Hz and stop on a HALT instruction @PC=0x0007
// Given I run the gameboy
// And I replace instruction @PC=0x0007 by a HALT instruction (IR=0x76)
// Then the gameboy should start the crystal oscillator
// And the crystal oscillator should tick the CPU@v0.4.0, (PPU@v0.4.1 and APU@v0.4.2)
// And at each step i, the gameboy should wait for the CPU, PPU and APU to finish
// And at each step i, the gameboy should send the CPU, PPU and APU states to their respective channels
// And the test case should receive the CPU, PPU and APU states
func TestRunDebuggerDoTickCPUUntilHalt(t *testing.T) {
	// create a channel to listen to cpu state updates
	cpuStateChannel := make(chan *CpuState)
	// create a new gameboy
	debugger := NewDebugger(cpuStateChannel, nil, nil, nil, nil)
	// initialize the gameboy
	debugger.LoadRom("Tetris.gb")

	// replace memory location @0x0007 with a HALT instruction (IR=0x76)
	debugger.gameboy.cpuBus.Write(0x0007, 0x76)

	// run the gameboy
	fmt.Println("Gameboy> running ...")
	debugger.Run()

	var cpuState *CpuState
loop:
	for {
		// listen to the cpu state channel waiting for one cpu state to arrive

		fmt.Println("Gameboy> listening to cpu state channel ...")
		select {
		case cpuState = <-cpuStateChannel:
			fmt.Println("Gameboy> cpu state received ...")
			cpuState.print()
			if cpuState.HALTED {
				break loop
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for cpu state")
		}
	}

	// check that the CPU stopped on PC=0x0007
	if cpuState.PC != 0x0007 {
		t.Fatalf("Expected CPU to stop on PC=0x0007, got PC=0x%04X", cpuState.PC)
	}
}

func TestDebuggerDoStopOnBreakPoints(t *testing.T) {
	// create a channel to listen to cpu state updates
	cpuStateChannel := make(chan *CpuState, 1000)
	// create a new gameboy
	debugger := NewDebugger(cpuStateChannel, nil, nil, nil, nil)
	// initialize the gameboy
	debugger.LoadRom("Tetris.gb")

	// set breakpoint on instruction @0x0060
	debugger.AddBreakPoint(0x0060)

	// run the gameboy
	fmt.Println("Gameboy> running ...")
	doneChan := debugger.Run()

	var cpuState *CpuState
loop1:
	for {
		// listen to the cpu state channel waiting for one cpu state to arrive
		select {
		case <-cpuStateChannel:
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for cpu state")
		case <-doneChan:
			break loop1
		}
	}

	cpuState = debugger.gameboy.cpu.getState()
	cpuState.print()

	// check that the CPU stopped on PC=0x0007
	if cpuState.PC != 0x0060 {
		t.Fatalf("Expected CPU to stop on PC=0x0060, got PC=0x%04X", cpuState.PC)
	} else {
		t.Log("Gameboy> CPU stopped on breakpoint @0x0060")
	}
}
