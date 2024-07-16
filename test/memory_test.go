package gameboy

import (
	"testing"

	"github.com/codefrite/gameboy-go/gameboy"
)

func TestNewMemory(t *testing.T) {
	memory := gameboy.NewMemory(0x2000)
	if memory == nil {
		t.Error("Expected VRAM to be initialized")
	}
	if memory.Size() != 0x2000 {
		t.Errorf("Expected VRAM to have 0x2000 bytes, got %d", memory.Size())
	}
}

// func TestVRAMReadAndWrite(t *testing.T) {
func TestReadAndWrite(t *testing.T) {
	memory := gameboy.NewMemory(0x2000)
	for i := 0; i < 0x2000; i++ {
		memory.Write(uint16(i), uint8(i))
	}
	for i := 0; i < 0x2000; i++ {
		data := memory.Read(uint16(i))
		if data != uint8(i) {
			t.Errorf("Expected VRAM to have value %02X at address %04X, got %02X", i, i, data)
		}
	}
}
