package gameboy

import (
	"errors"
	"fmt"
)

const (
	// Special Actions associated with certain memory addresses
	DISABLE_BOOT_ROM_REGISTER = 0xFF50 // On Write, disable the bootrom
)

type Accessible interface {
	Read(uint16) uint8
	Dump(uint16, uint16) []uint8
	Write(uint16, uint8)
	Size() uint16
}

// represents a memory mapped to a specific address
// Name: string name of the memory
// Address: uint16 address where the memory is mapped
// Memory: Accessible memory
type MemoryMap struct {
	Name    string
	Address uint16
	Memory  Accessible
}

// the memory management unit (MMU) is responsible for routing memory accesses to the correct memory
type MMU struct {
	memoryMaps   []MemoryMap
	memoryWrites []MemoryWrite
}

// constructor for the MMU struct
func NewMMU() *MMU {
	return &MMU{
		memoryMaps:   []MemoryMap{},
		memoryWrites: []MemoryWrite{},
	}
}

func (m *MMU) getMemoryWrites() *[]MemoryWrite {
	return &m.memoryWrites
}

// add a memory write to the memory writes
func (m *MMU) addMemoryWrite(memoryWrite MemoryWrite) {
	m.memoryWrites = append(m.memoryWrites, memoryWrite)
}

// clear the memory writes
func (m *MMU) clearMemoryWrites() {
	m.memoryWrites = []MemoryWrite{}
}

// Attach a memory to the MMU at the given address.
// At the moment, there is not check for overlapping memories: when reading from or writing to an address contained in
// multiple memories, the first one found will be used.
// name: string name of the memory
// address: uint16 address where the memory will be attached
// memory: Accessible memory to attach
// return void
func (m *MMU) AttachMemory(name string, address uint16, memory Accessible) {
	// TODO: check if the memory is already attached before attaching it (might need to pass a pointer to the memory instead of the memory itself or pool &memory in the AttachMemory method)
	for _, memoryMap := range m.memoryMaps {
		if &memoryMap.Memory == &memory {
			return
		}
	}
	m.memoryMaps = append(m.memoryMaps, MemoryMap{
		Name:    name,
		Address: address,
		Memory:  memory,
	})
}

// return the memory maps attached to the MMU as MemoryWrite[] (used by the debugger to display the memories)
// return []MemoryMap memory maps attached to the MMU
func (b *MMU) GetMemoryMaps() []MemoryWrite {
	memoryWrites := []MemoryWrite{}
	for _, memoryMap := range b.memoryMaps {
		memoryDump := memoryMap.Memory.Dump(0, memoryMap.Memory.Size())
		memoryWrite := MemoryWrite{
			Name:    memoryMap.Name,
			Address: memoryMap.Address,
			Data:    memoryDump,
		}
		memoryWrites = append(memoryWrites, memoryWrite)
	}
	return memoryWrites
}

// Return the memory map that contains the address or returns an error if the address is not found/mapped.
// address: uint16 address to look for
// return *MemoryMap memory map containing the address or an error if the address is not found
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

// Read the value at the given address.
// addr: uint16 address where the value will be read
// return uint8 value at the given address
// panic if the address is not found
func (b *MMU) Read(addr uint16) uint8 {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		return memoryMap.Memory.Read(addr - memoryMap.Address)
	} else {
		panic(err)
	}
}

// Dump memory from address 'from' to address 'to'
// from: uint16 start address
// to: uint16 end address
// return []uint8 (blob/memory dump)
// panic if the addresses are out of bounds.
func (b *MMU) Dump(from uint16, to uint16) []uint8 {
	memoryMap, err := b.findMemory(from)
	if err == nil {
		return memoryMap.Memory.Dump(from-memoryMap.Address, to-memoryMap.Address)
	} else {
		panic(err)
	}
}

// Reads the next 2 bytes from the given address and returns them as a uint16 little-endian value (HIGH LOW).
func (b *MMU) Read16(addr uint16) uint16 {
	return uint16(b.Read(addr+1))<<8 | uint16(b.Read(addr))
}

// Write the provided value at the given address.
// addr: uint16 address where the value will be written
// value: uint8 value to write
// return void
// panic if the address is not found
func (b *MMU) Write(addr uint16, value uint8) error {
	// on write to 0xFF50, disable the bootrom
	if addr == DISABLE_BOOT_ROM_REGISTER {
		b.DisableBootRom()
	} else if addr == REG_FF04_DIV {
		// if the address is the divider register, reset it to 0
		value = 0
	}

	// find the memory map containing the address
	memoryMap, err := b.findMemory(addr)

	// if the address is not found, return an error
	if err != nil {
		return err
	}

	// write the value to the memory

	memoryMap.Memory.Write(addr-memoryMap.Address, value)
	memoryWrite := MemoryWrite{
		Name:    memoryMap.Name,
		Address: addr - memoryMap.Address,
		Data:    []uint8{value},
	}
	b.addMemoryWrite(memoryWrite)
	return nil
}

// Write the provided blob at the given address.
// Please note that the entire blob should belong to the same memory otherwise,
// the write process will panic once it reaches the end of the first memory,
// resulting in the partial write of the blob.
// param addr: uint16 address where the blob will be written
// param blob: []uint8 blob to write
// return void
// panic if one of the addresses in the range (addr, addr+len(blob)) is not found
func (b *MMU) WriteBlob(addr uint16, blob []uint8) {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		// write the blob to the memory
		for i, value := range blob {
			memoryMap.Memory.Write(addr-memoryMap.Address+uint16(i), value)
		}
		// log the memory write
		memoryWrite := MemoryWrite{
			Name:    memoryMap.Name,
			Address: addr - memoryMap.Address,
			Data:    blob,
		}
		b.addMemoryWrite(memoryWrite)
	} else {
		panic(err)
	}
}

// Disable the bootrom by removing it from the memory maps if a write operation to 0xFF50 is detected.
func (mmu *MMU) DisableBootRom() {
	for idx, mem := range mmu.memoryMaps {
		if mem.Name == BOOT_ROM_MEMORY_NAME {
			// remove the bootrom from the memory maps
			mmu.memoryMaps = append(mmu.memoryMaps[0:idx], mmu.memoryMaps[idx+1:]...)
			return
		}
	}
}

// Special write operation for the timer registers when called by the Timer itself
// Allows the DIV register to be incremented, otherwise, any write for example from
// cartridge program to DIV register will reset it to 0.
func (mmu *MMU) timerWrite(addr uint16, value uint8) {
	// find the memory map containing the address
	memoryMap, err := mmu.findMemory(addr)
	// if the address is not found, return an error
	if err != nil {
		panic(err)
	}
	// write the value to the memory
	memoryMap.Memory.Write(addr-memoryMap.Address, value)
}
