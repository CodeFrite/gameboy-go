package gameboy

import (
	"errors"
	"fmt"
)

type Accessible interface {
	Read(uint16) uint8
	Dump(uint16, uint16) []uint8
	Write(uint16, uint8)
	Size() uint16
}

type MemoryMap struct {
	Name    string
	Address uint16
	Memory  Accessible
}

type MMU struct {
	router []MemoryMap
}

func NewMMU() *MMU {
	return &MMU{
		router: []MemoryMap{},
	}
}

func (b *MMU) AttachMemory(name string, address uint16, memory Accessible) {
	b.router = append(b.router, MemoryMap{
		Name:    name,
		Address: address,
		Memory:  memory,
	})
}

/*
 * Return the memory map that contains the address or return nil
 */
func (b *MMU) findMemory(address uint16) (*MemoryMap, error) {
	for _, memoryMap := range b.router {
		if address >= memoryMap.Address && address < memoryMap.Address+memoryMap.Memory.Size() {
			return &memoryMap, nil
		}
	}
	errMessage := fmt.Sprintf("Memory location 0x%04X not found", address)
	return nil, errors.New(errMessage)
}

func (b *MMU) Read(addr uint16) uint8 {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		return memoryMap.Memory.Read(addr - memoryMap.Address)
	} else {
		panic(err)
	}
}

func (b *MMU) Dump(from uint16, to uint16) []uint8 {
	memoryMap, err := b.findMemory(from)
	if err == nil {
		return memoryMap.Memory.Dump(from-memoryMap.Address, to-memoryMap.Address)
	} else {
		panic(err)
	}
}
func (b *MMU) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr+1))<<8 | uint16(b.Read(addr))
}

func (b *MMU) Write(addr uint16, value uint8) {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		memoryMap.Memory.Write(addr-memoryMap.Address, value)
	} else {
		panic(err)
	}
}

func (b *MMU) WriteBlob(addr uint16, blob []uint8) {
	for i, value := range blob {
		b.Write(addr+uint16(i), uint8(value))
	}
}
