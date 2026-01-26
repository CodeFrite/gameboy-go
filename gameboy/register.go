package gameboy

import "fmt"

// represents an 8-bit flag register
// it will allow the user to:
// - manipulate the register as a whole passing it a uint8 value
// - manipulate the register bit by bit using the GetBit, SetBit and ResetBit methods
// - have access to any bit using its shorthand notation (Z, N, H, C)
type Register8 struct {
	uint8
	address uint16
}

func NewRegister8(address uint16) *Register8 {
	return &Register8{address: address}
}

func (r8 *Register8) GetAddress() uint16 {
	return r8.address
}

func (r8 *Register8) Get() uint8 {
	return r8.uint8
}

func (r8 *Register8) Set(value uint8) {
	r8.uint8 = value
}

func (r8 *Register8) GetBit(bit uint8) bool {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> getBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	return (r8.uint8 & op) == op
}

func (r8 *Register8) SetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> setBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	r8.uint8 |= op
}

func (r8 *Register8) ResetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> resetBit: bit out of range: %v", bit))
	}
	r8.uint8 ^= 1 << bit
}

// 16-bit register

type Register16 uint16

func (r16 *Register16) Set(value uint16) {
	*r16 = Register16(value)
}

func (r16 *Register16) Get() uint16 {
	return uint16(*r16)
}

func (r16 *Register16) GetBit(bit uint8) bool {
	if bit > 15 {
		panic(fmt.Sprintf("Register16> getBit: bit out of range: %v", bit))
	}
	op := uint16(1 << bit)
	return (r16.Get() & op) == op
}

func (r16 *Register16) SetBit(bit uint8) {
	if bit > 15 {
		panic(fmt.Sprintf("Register16> setBit: bit out of range: %v", bit))
	}
	*r16 = Register16(r16.Get() | (1 << bit))
}

func (r16 *Register16) ResetBit(bit uint8) {
	if bit > 15 {
		panic(fmt.Sprintf("Register16> resetBit: bit out of range: %v", bit))
	}
	*r16 = Register16(r16.Get() & ^(1 << bit))
}
