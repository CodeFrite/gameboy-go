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

/**
 * represents a memory mapped to a specific address
 * @param Name: string name of the memory
 * @param Address: uint16 address where the memory is mapped
 * @param Memory: Accessible memory
 */
type MemoryMap struct {
	Name    string
	Address uint16
	Memory  Accessible
}

/**
 * represents a memory write operation used by the mmu to keep track of all memory changes between 2 states
 * memory writes can be either reset:
 * - before every step from the debugger
 * - before a new run from a breakpoint to another breakpoint
 */
type MemoryWrite struct {
	Name    string   `json:"name"`
	Address uint16   `json:"address"`
	Data    []string `json:"data"`
}

/**
 * the memory management unit (MMU) is responsible for routing memory accesses to the correct memory.
 */
type MMU struct {
	memoryMaps   []MemoryMap
	memoryWrites []MemoryWrite
}

/**
 * constructor for the MMU struct
 */
func NewMMU() *MMU {
	return &MMU{
		memoryMaps:   []MemoryMap{},
		memoryWrites: []MemoryWrite{},
	}
}

/**
 * clear memory writes
 */
func (m *MMU) clearMemoryWrites() {
	m.memoryWrites = []MemoryWrite{}
}

/**
 * Attach a memory to the MMU at the given address.
 * At the moment, there is not check for overlapping memories: when reading from or writing to an address contained in
 * multiple memories, the first one found will be used. TODO: forbid overlapping memories to be attached.
 * @param name: string name of the memory
 * @param address: uint16 address where the memory will be attached
 * @param memory: Accessible memory to attach
 * @return void
 */
func (b *MMU) AttachMemory(name string, address uint16, memory Accessible) {
	b.memoryMaps = append(b.memoryMaps, MemoryMap{
		Name:    name,
		Address: address,
		Memory:  memory,
	})
}

/**
 * return the memory maps attached to the MMU (used by the debugger to display the memories)
 * @return []MemoryMap memory maps attached to the MMU
 */
func (b *MMU) GetMemoryMaps() []MemoryMap {
	return b.memoryMaps
}

// TODO: i could implement a detachMemory method to remove a memory from the MMU like the cartrige from the gameboy

/**
 * Return the memory map that contains the address or returns an error if the address is not found/mapped.
 * @param address: uint16 address to look for
 * @return *MemoryMap memory map containing the address or an error if the address is not found
 */
func (b *MMU) findMemory(address uint16) (*MemoryMap, error) {
	for _, memoryMap := range b.memoryMaps {
		memoryMapSize := memoryMap.Memory.Size() - 1
		condition1 := address >= memoryMap.Address
		condition2 := address <= memoryMap.Address+memoryMapSize
		if condition1 && condition2 {
			return &memoryMap, nil
		}
	}
	errMessage := fmt.Sprintf("Memory location 0x%04X not found", address)
	return nil, errors.New(errMessage)
}

/**
 * Read the value at the given address.
 * @param addr: uint16 address where the value will be read
 * @return uint8 value at the given address
 * @panic if the address is not found
 */
func (b *MMU) Read(addr uint16) uint8 {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		return memoryMap.Memory.Read(addr - memoryMap.Address)
	} else {
		panic(err)
	}
}

/**
 * Dump memory from address 'from' to address 'to'
 * @param from: uint16 start address
 * @param to: uint16 end address
 * @return []uint8 (blob/memory dump)
 * @panic if the addresses are out of bounds.
 */
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

/**
 * Write the provided value at the given address.
 * @param addr: uint16 address where the value will be written
 * @param value: uint8 value to write
 * @return void
 * @panic if the address is not found
 */
func (b *MMU) Write(addr uint16, value uint8) {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		memoryMap.Memory.Write(addr-memoryMap.Address, value)
		memoryWrite := MemoryWrite{
			Name:    memoryMap.Name,
			Address: addr,
			Data:    []string{fmt.Sprintf("0x%02X", value)},
		}
		b.memoryWrites = append(b.memoryWrites, memoryWrite)
	} else {
		panic(err)
	}
}

/**
 * Write the provided blob at the given address.
 * @param addr: uint16 address where the blob will be written
 * @param blob: []uint8 blob to write
 * @return void
 * @panic if one of the addresses in the range (addr, addr+len(blob)) is not found
 */
func (b *MMU) WriteBlob(addr uint16, blob []uint8) {
	for i, value := range blob {
		b.Write(addr+uint16(i), uint8(value))
		// memory writes are already logged in the Write method
	}
}
