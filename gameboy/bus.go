package gameboy

import "errors"

type Accessible interface {
	Read(uint16) uint8
	Write(uint16, uint8)
	Size() uint16
}

type MemoryMap struct {
	Address uint16
	Memory  Accessible
}

type Bus struct {
	router []MemoryMap
}

func NewBus() *Bus {
	return &Bus{
		router: []MemoryMap{},
	}
}

func (b *Bus) AttachMemory(memory Accessible, address uint16) {
	b.router = append(b.router, MemoryMap{
		Address: address,
		Memory:  memory,
	})
}

/*
 * Return the memory map that contains the address or return nil
 */
func (b *Bus) findMemory(address uint16) (*MemoryMap, error) {
	for _, memoryMap := range b.router {
		if address >= memoryMap.Address && address < memoryMap.Address+memoryMap.Memory.Size() {
			return &memoryMap, nil
		}
	}
	return nil, errors.New("Memory location not found")
}

func (b *Bus) Read(addr uint16) uint8 {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		return memoryMap.Memory.Read(addr - memoryMap.Address)
	} else {
		panic(err)
	}
}

func (b *Bus) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr+1))<<8 | uint16(b.Read(addr))
}

func (b *Bus) Write(addr uint16, value uint8) {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		memoryMap.Memory.Write(addr-memoryMap.Address, value)
	} else {
		panic(err)
	}
}
