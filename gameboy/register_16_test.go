package gameboy

import (
	"testing"
)

func TestInstantiation16(t *testing.T) {
	var r16 Register16 = Register16(0x0101)
	// set the register value to 0x0101
	if r16 != 0x0101 {
		t.Errorf("Expected r16 to be 0x00, got %04X", r16)
	}
	// set the register value to 0x2BE7
	r16 = Register16(0x2BE7)
	if r16 != 0x2BE7 {
		t.Errorf("Expected r16 to be 0x2BE7, got %04X", r16)
	}
	// set the register value to 0x1A77
	r16 = Register16(0x1A77)
	if r16.Get() != 0x1A77 {
		t.Errorf("Expected r16 to be 0x1A77, got %04X", r16)
	}
}

func TestGetterSetter16(t *testing.T) {
	var r16 Register16 = Register16(0)
	// set the register value to 0x2BE7
	r16.Set(0x2BE7)
	if r16.Get() != 0x2BE7 {
		t.Errorf("Expected r16 to be 0x2BE7, got %04X", r16.Get())
	}
	// set the register value to 0x1A77
	r16.Set(0x1A77)
	if r16.Get() != 0x1A77 {
		t.Errorf("Expected r16 to be 0x1A, got %04X", r16.Get())
	}
}

func TestGetBit16(t *testing.T) {
	// 0x0000 = b0000 0000 0000 0000
	var r16 Register16 = Register16(0x0000)
	for i := uint8(0); i < 16; i++ {
		if r16.GetBit(i) {
			t.Errorf("Expected r16 to have bit %d unset, got %v", i, r16.GetBit(i))
		}
	}

	// 0xFF = b1111 1111 1111 1111
	r16 = Register16(0xFFFF)
	for i := uint8(0); i < 16; i++ {
		if !r16.GetBit(i) {
			t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
		}
	}

	// 0xAA = b1010 1010 1010 1010
	r16 = Register16(0xAAAA)
	for i := uint8(0); i < 16; i++ {
		if i%2 == 0 {
			if r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
			}
		} else {
			if !r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d unset, got %v", i, r16.GetBit(i))
			}
		}
	}

	// 0x5555 = b0101 0101 0101 0101
	r16 = Register16(0x5555)
	for i := uint8(0); i < 16; i++ {
		if i%2 == 0 {
			if !r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d unset, got %v", i, r16.GetBit(i))
			}
		} else {
			if r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
			}
		}
	}
}

func TestSetBit16(t *testing.T) {
	// 0xFFFF = b1111 1111 1111 1111
	var r16 Register16 = Register16(0x00)
	for i := uint8(0); i < 16; i++ {
		r16.SetBit(i)
		if !r16.GetBit(i) {
			t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
		}
	}
	if r16 != 0xFFFF {
		t.Errorf("Expected r16 to be 0xFFFF, got %04X", r16)
	}

	// 0xAAAA = b1010 1010 1010 1010
	r16 = Register16(0x00)
	for i := uint8(0); i < 16; i++ {
		if i%2 != 0 {
			r16.SetBit(i)
			if !r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
			}
		}
	}
	if r16 != 0xAAAA {
		t.Errorf("Expected r16 to be 0xAAAA, got %04X", r16)
	}

	// 0x5555 = b0101 0101 0101 0101
	r16 = Register16(0x00)
	for i := uint8(0); i < 16; i++ {
		if i%2 == 0 {
			r16.SetBit(i)
			if !r16.GetBit(i) {
				t.Errorf("Expected r16 to have bit %d set, got %v", i, r16.GetBit(i))
			}
		}
	}
	if r16 != 0x5555 {
		t.Errorf("Expected r16 to be 0x5555, got %04X", r16)
	}
}

func TestResetBit16(t *testing.T) {
	// 0xFFFF = b1111 1111 1111 1111
	var r16 Register16 = Register16(0xFFFF)
	for i := uint8(0); i < 16; i++ {
		r16.ResetBit(i)
		if r16.GetBit(i) {
			t.Errorf("Expected r16 to have bit %d unset, got %v", i, r16.GetBit(i))
		}
	}
}
