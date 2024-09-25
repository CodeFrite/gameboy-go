package gameboy

import "fmt"

// 8-bit register

type Register8 uint8

func (r8 *Register8) Get() uint8 {
	return uint8(*r8)
}

func (r8 *Register8) Set(value uint8) {
	*r8 = Register8(value)
}

func (r8 *Register8) GetBit(bit uint8) bool {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> getBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	return (r8.Get() & op) == op
}

func (r8 *Register8) SetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> setBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	*r8 = Register8(r8.Get() | op)
}

func (r8 *Register8) ResetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("Register8> resetBit: bit out of range: %v", bit))
	}
	*r8 = Register8(r8.Get() & ^(1 << bit))
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
