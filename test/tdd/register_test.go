package test

import (
	"testing"

	"github.com/codefrite/gameboy-go/gameboy"
)

func TestInstantiation(t *testing.T) {
	var r8 gameboy.Register8 = gameboy.Register8(0x01)
	if r8 != 0x01 {
		t.Errorf("Expected r8 to be 0x00, got %02X", r8)
	}
	r8 = gameboy.Register8(0xE7)
	if r8 != 0xE7 {
		t.Errorf("Expected r8 to be 0xE7, got %02X", r8)
	}
	r8 = gameboy.Register8(0x1A)
	if r8.Get() != 0x1A {
		t.Errorf("Expected r8 to be 0xE7, got %02X", r8)
	}
}

func TestGetterSetter(t *testing.T) {
	var r8 gameboy.Register8 = gameboy.Register8(0)
	// set the value to 0xE7
	r8.Set(0xE7)
	if r8.Get() != 0xE7 {
		t.Errorf("Expected r8 to be 0xE7, got %02X", r8.Get())
	}
	r8.Set(0x1A)
	if r8.Get() != 0x1A {
		t.Errorf("Expected r8 to be 0x1A, got %02X", r8.Get())
	}
}

func TestGetBit(t *testing.T) {
	// 0x00 = b0000 0000
	var r8 gameboy.Register8 = gameboy.Register8(0x00)
	for i := uint8(0); i < 8; i++ {
		if r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
		}
	}

	// 0xFF = b1111 1111
	r8 = gameboy.Register8(0xFF)
	for i := uint8(0); i < 8; i++ {
		if !r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
		}
	}

	// 0xA0 = b1010 1010
	r8 = gameboy.Register8(0xAA)
	for i := uint8(0); i < 8; i++ {
		if i%2 == 0 {
			if r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		} else {
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
			}
		}
	}

	// 0x55 = b0101 0101
	r8 = gameboy.Register8(0x55)
	for i := uint8(0); i < 8; i++ {
		if i%2 == 0 {
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
			}
		} else {
			if r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		}
	}
}

func TestSetBit(t *testing.T) {
	// 0xFF = b1111 1111
	var r8 gameboy.Register8 = gameboy.Register8(0x00)
	for i := uint8(0); i < 8; i++ {
		r8.SetBit(i)
		if !r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
		}
	}

	// 0xAA = b1010 1010
	r8 = gameboy.Register8(0x00)
	for i := uint8(0); i < 8; i++ {
		if i%2 == 0 {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		} else {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		}
	}

	// 0x55 = b0101 0101
	r8 = gameboy.Register8(0x00)
	for i := uint8(0); i < 8; i++ {
		if i%2 == 0 {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		} else {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		}
	}
}

func TestResetBit(t *testing.T) {
	// 0xFF = b1111 1111
	var r8 gameboy.Register8 = gameboy.Register8(0xFF)
	for i := uint8(0); i < 8; i++ {
		r8.ResetBit(i)
		if r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
		}
	}
}
