package gameboy

import "testing"

func TestBytesToUint16(t *testing.T) {
	// Test the bytesToUint16 function
	var b [2]byte
	b[0] = 0x12
	b[1] = 0x34
	if bytesToUint16(b) != 0x1234 {
		t.Errorf("bytesToUint16 failed. Expected 0x1234, got 0x%X", bytesToUint16(b))
	}
}

func TestUint16ToBytes(t *testing.T) {
	// Test the uint16ToBytes function
	var u uint16 = 0x1234
	b := uint16ToBytes(u)
	if b[0] != 0x12 || b[1] != 0x34 {
		t.Errorf("uint16ToBytes failed. Expected [0x12, 0x34], got [0x%X, 0x%X]", b[1], b[0])
	}
}

func TestBytesToBytes(t *testing.T) {
	// Test the bytesToBytes function
	var b [2]byte
	b[0] = 0x12
	b[1] = 0x34

	res := uint16ToBytes(bytesToUint16(b))
	if res[0] != 0x12 || res[1] != 0x34 {
		t.Errorf("uint16ToBytes(bytesToUint16(b)) failed. Expected 0x1234, got [0x%X,0x%X]", res[0], res[1])
	}
}

func TestUint16ToUint16(t *testing.T) {
	// Test the uint16ToUint16 function
	var u uint16 = 0x1234

	res := bytesToUint16(uint16ToBytes(u))
	if res != 0x1234 {
		t.Errorf("bytesToUint16(uint16ToBytes(u)) failed. Expected 0x1234, got 0x%X", res)
	}
}