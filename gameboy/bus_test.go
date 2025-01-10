package gameboy

import (
	"testing"
)

// check if the bootrom is disabled after writing to 0xFF50 register
func TestDisableBootRom(t *testing.T) {
	bus := NewBus()
	// boot rom will contain only ones
	bootRom := NewMemory(0x0100)
	bootRom.ResetWithOnes()
	bus.AttachMemory("Boot ROM", 0x0000, bootRom)
	// cartridge rom will contain only zeros
	cartridgeRom := NewMemory(0x8000)
	for i := 0; i < len(cartridgeRom.data); i++ {
		cartridgeRom.data[i] = 0
	}
	bus.AttachMemory("Cartridge ROM", 0x0000, cartridgeRom)

	// before writing to 0xFF50, the boot rom should be read and return 0xFF for all addresses [0x0000, 0x00FF]
	for i := 0; i < 0x0100; i++ {
		if bus.Read(uint16(i)) != 0xFF {
			t.Errorf("Expected 0xFF, got %X", bus.Read(uint16(i)))
		}
	}
	// mmu memory maps should contain 2 memory maps
	if len(bus.memoryMaps) != 2 {
		t.Errorf("Expected 2 memory maps, got %d", len(bus.memoryMaps))
	}

	// write to 0xFF50
	bus.Write(0xFF50, 0x01)
	// after writing to 0xFF50, the cartridge rom should be read and return 0xFF for all addresses [0x0000, 0x7FFF]
	for i := 0; i < 0x0100; i++ {
		if bus.Read(uint16(i)) != 0x00 {
			t.Errorf("Expected 0xFF, got %X", bus.Read(uint16(i)))
		}
	}
	// mmu memory maps should contain 1 memory map
	if len(bus.memoryMaps) != 1 {
		t.Errorf("Expected 1 memory map, got %d", len(bus.memoryMaps))
	}
}
