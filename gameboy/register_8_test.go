package gameboy

import (
	"testing"
)

func TestInstantiation8(t *testing.T) {
	var r8 Register8
	r8.Set(0x01)
	if r8.Get() != 0x01 {
		t.Errorf("Expected r8 to be 0x00, got %02X", r8)
	}
	r8.Set(0xE7)
	if r8.Get() != 0xE7 {
		t.Errorf("Expected r8 to be 0xE7, got %02X", r8)
	}
	r8.Set(0x1A)
	if r8.Get() != 0x1A {
		t.Errorf("Expected r8 to be 0x1A, got %02X", r8)
	}
}

func TestGetterSetter8(t *testing.T) {
	var r8 Register8
	r8.Set(0x00)
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

func TestGetBit8(t *testing.T) {
	// 0x00 = b0000 0000
	var r8 Register8
	r8.Set(0x00)
	for i := uint8(0); i < 8; i++ {
		if r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
		}
	}

	// 0xFF = b1111 1111
	r8.Set(0xFF)
	for i := uint8(0); i < 8; i++ {
		if !r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
		}
	}

	// 0xA0 = b1010 1010
	r8.Set(0xAA)
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
	r8.Set(0x55)
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

func TestSetBit8(t *testing.T) {
	// 0xFF = b1111 1111
	var r8 Register8
	r8.Set(0x00)
	for i := uint8(0); i < 8; i++ {
		r8.SetBit(i)
		if !r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
		}
	}

	// 0xAA = b1010 1010
	r8.Set(0x00)
	for i := uint8(0); i < 8; i++ {
		if i%2 != 0 {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		}
	}
	if r8.Get() != 0xAA {
		t.Errorf("Expected r8 to be 0xAA, got %02X", r8)
	}

	// 0x55 = b0101 0101
	r8.Set(0x00)
	for i := uint8(0); i < 8; i++ {
		if i%2 == 0 {
			r8.SetBit(i)
			if !r8.GetBit(i) {
				t.Errorf("Expected r8 to have bit %d set, got %v", i, r8.GetBit(i))
			}
		}
	}
	if r8.Get() != 0x55 {
		t.Errorf("Expected r8 to be 0x55, got %02X", r8)
	}
}

func TestResetBit8(t *testing.T) {
	// 0xFF = b1111 1111
	var r8 Register8
	r8.Set(0xFF)
	for i := uint8(0); i < 8; i++ {
		r8.ResetBit(i)
		if r8.GetBit(i) {
			t.Errorf("Expected r8 to have bit %d unset, got %v", i, r8.GetBit(i))
		}
	}
}
