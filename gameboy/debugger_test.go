package gameboy

import (
	"fmt"
	"testing"
)

// TestNewDebugger tests the creation of a new debugger
func TestNewDebugger(t *testing.T) {
	debugger := NewDebugger(nil, nil, nil, nil, nil)
	if debugger == nil {
		t.Fatalf("Expected a new debugger, got nil")
	}
}

func TestDebuggerStateChannels(t *testing.T) {
	cpuStateChannel := make(chan *CpuState)
	ppuStateChannel := make(chan *PpuState)
	apuStateChannel := make(chan *ApuState)
	memoryStateChannel := make(chan *[]MemoryWrite)
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
