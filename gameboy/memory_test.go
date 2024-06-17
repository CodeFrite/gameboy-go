package gameboy

import "testing"

func TestNewVRAM(t *testing.T) {
	v := NewMemory(0x2000)
	if v == nil {
		t.Error("Expected VRAM to be initialized")
	}
	if len(v.data) != 0x2000 {
		t.Errorf("Expected VRAM to have 0x2000 bytes, got %d", len(v.data))
	}
}

//func TestVRAMReadAndWrite(t *testing.T) {
