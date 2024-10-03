package gameboy

import (
	"fmt"
	"testing"
	"time"
)

func TestGameboyStateChannels(t *testing.T) {
	cpuStateChannel := make(chan *CpuState)
	ppuStateChannel := make(chan *PpuState)
	apuStateChannel := make(chan *ApuState)
	memoryStateChannel := make(chan *[]MemoryWrite)
	joypadStateChannel := make(chan *JoypadState)

	gb := NewGameboy(cpuStateChannel, ppuStateChannel, apuStateChannel, memoryStateChannel, joypadStateChannel)

	if gb.cpuStateChannel != cpuStateChannel {
		t.Fatalf("Expected a gameboy cpu state channel to be initialized with address %p, got %p", cpuStateChannel, gb.cpuStateChannel)
	}

	if gb.ppuStateChannel != ppuStateChannel {
		t.Fatalf("Expected a gameboy ppu state channel to be initialized with address %p, got %p", ppuStateChannel, gb.ppuStateChannel)
	}

	if gb.apuStateChannel != apuStateChannel {
		t.Fatalf("Expected a gameboy apu state channel to be initialized with address %p, got %p", apuStateChannel, gb.apuStateChannel)
	}

	if gb.memoryStateChannel != memoryStateChannel {
		t.Fatalf("Expected a gameboy memory state channel to be initialized with address %p, got %p", memoryStateChannel, gb.memoryStateChannel)
	}
}

// When stepping the gameboy, the crystal oscillator should tick once, trigger the Gameboy.onTick() method,
// which in turn should tick the CPU, PPU and APU.
// In this test, we will tick the crystal 5 times and check that the CPU is executing once every tick
// We will then make sure that Gameboy.onTick() sends back the CPU state through the cpu state channel to the test case
func TestStepGameBoyDoTickCPU(t *testing.T) {
	// create a channel to listen to cpu state updates
	cpuStateChannel := make(chan *CpuState)
	// create a new gameboy
	gb := NewGameboy(cpuStateChannel, nil, nil, nil, nil)
	// initialize the gameboy
	gb.LoadRom("Tetris.gb")
	// start the gameboy
	fmt.Println("Gameboy> stepping ...")

	// step 5 times
	for i := 0; i < 5; i++ {
		// step in // routine
		go gb.Step()
		// listen to the cpu state channel waiting for one cpu state to arrive

		fmt.Println("Gameboy> listening to cpu state channel ...")
		select {
		case cpuState := <-cpuStateChannel:
			fmt.Println("Gameboy> cpu state received ...")
			cpuState.print()
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for cpu state")
		}
	}
}

// Scenario: the gameboy should run @1Hz and stop on a HALT instruction @PC=0x0007
// Given I run the gameboy
// And I replace instruction @PC=0x0007 by a HALT instruction (IR=0x76)
// Then the gameboy should start the crystal oscillator
// And the crystal oscillator should tick the CPU@v0.4.0, (PPU@v0.4.1 and APU@v0.4.2)
// And at each step i, the gameboy should wait for the CPU, PPU and APU to finish
// And at each step i, the gameboy should send the CPU, PPU and APU states to their respective channels
// And the test case should receive the CPU, PPU and APU states
func TestRunGameBoyDoTickCPUUntilHalt(t *testing.T) {
	// create a channel to listen to cpu state updates
	cpuStateChannel := make(chan *CpuState)
	// create a new gameboy
	gb := NewGameboy(cpuStateChannel, nil, nil, nil, nil)
	// initialize the gameboy
	gb.LoadRom("Tetris.gb")

	// replace memory location @0x0007 with a HALT instruction (IR=0x76)
	gb.cpuBus.Write(0x0007, 0x76)

	// run the gameboy
	fmt.Println("Gameboy> running ...")
	gb.Run()

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

// When stepping the gameboy, the crystal oscillator should tick once, trigger the Gameboy.onTick() method,
// which in turn should tick the CPU, PPU and APU.
// In this test, we will tick the crystal 5 times and check that the PPU is executing once every tick
// We will then make sure that Gameboy.onTick() sends back the PPU state through the ppu state channel to the test case
func TestStepGameBoyDoTickPPU(t *testing.T) {
	// create a channel to listen to ppu state updates
	ppuStateChannel := make(chan *PpuState)
	// create a new gameboy
	gb := NewGameboy(nil, ppuStateChannel, nil, nil, nil)
	// initialize the gameboy
	gb.LoadRom("Tetris.gb")
	// start the gameboy
	fmt.Println("Gameboy> stepping ...")

	// step 5 times
	for i := 0; i < 5; i++ {
		// step in // routine
		fmt.Println("Gameboy> stepping ...")
		go gb.Step()
		// listen to the ppu state channel waiting for one ppu state to arrive

		fmt.Println("Gameboy> listening to ppu state channel ...")
		select {
		case <-ppuStateChannel:
			fmt.Println("Gameboy> ppu state received ...")
			//ppuState.print()
		case <-time.After(5 * time.Second):
			t.Fatal("Test timed out waiting for ppu state")
		}
	}
}

// Scenario: the gameboy should run @1Hz and stop on a HALT instruction @PC=0x0007
// Given I run the gameboy
// And I replace instruction @PC=0x0007 by a HALT instruction (IR=0x76)
// Then the gameboy should start the crystal oscillator
// And the crystal oscillator should tick the CPU@v0.4.0, (PPU@v0.4.1 and APU@v0.4.2)
// And at each step i, the gameboy should wait for the CPU, PPU and APU to finish
// And at each step i, the gameboy should send the CPU, PPU and APU states to their respective channels
// And the test case should receive the CPU, PPU and APU states
func TestRunGameBoyDoTickPPUUntilHalt(t *testing.T) {
	// create a channel to listen to cpu state updates
	cpuStateChannel := make(chan *CpuState)
	ppuStateChannel := make(chan *PpuState)
	// create a new gameboy
	gb := NewGameboy(cpuStateChannel, ppuStateChannel, nil, nil, nil)
	// initialize the gameboy
	gb.LoadRom("Tetris.gb")

	// replace memory location @0x0007 with a HALT instruction (IR=0x76)
	gb.cpuBus.Write(0x0007, 0x76)

	// run the gameboy
	fmt.Println("Gameboy> running ...")
	go gb.Run()

	var cpuState *CpuState
loop:
	// step 5 times
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
