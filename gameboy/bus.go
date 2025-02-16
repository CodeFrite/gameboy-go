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

// the BUS is responsible for managing memory accesses
type Bus struct {
	// state
	memoryMaps   []MemoryMap
	memoryWrites []MemoryWrite
	// memory access handlers
	writeHandlers map[uint16]func(uint8) uint8
}

// constructor for the MMU struct
func NewBus() *Bus {
	return &Bus{
		memoryMaps:   []MemoryMap{},
		memoryWrites: []MemoryWrite{},
	}
}

// reset the bus
func (bus *Bus) reset() {
	bus.memoryMaps = []MemoryMap{}
	bus.memoryWrites = []MemoryWrite{}
	// TODO: check if i need to reset the write handlers ???
}

// Memory Writes Operations

func (bus *Bus) getMemoryWrites() *[]MemoryWrite {
	return &bus.memoryWrites
}

// add a memory write to the memory writes
func (bus *Bus) addMemoryWrite(memoryWrite MemoryWrite) {
	bus.memoryWrites = append(bus.memoryWrites, memoryWrite)
}

// clear the memory writes
func (bus *Bus) clearMemoryWrites() {
	bus.memoryWrites = []MemoryWrite{}
}

// Memory Maps Operations

// Attach a memory to the BUS at the given address.
// At the moment, there is not check for overlapping memories: when reading from or writing to an address contained in
// multiple memories, the first one found will be used.
// name: string name of the memory
// address: uint16 address where the memory will be attached
// memory: Accessible memory to attach
// return void
func (bus *Bus) AttachMemory(name string, address uint16, memory Accessible) {
	// if memory is already attached, return
	for _, memoryMap := range bus.memoryMaps {
		if &memoryMap.Memory == &memory {
			return
		}
	}
	// otherwise, attach the memory
	bus.memoryMaps = append(bus.memoryMaps, MemoryMap{
		Name:    name,
		Address: address,
		Memory:  memory,
	})
}

// return the memory maps attached to the Bus as MemoryWrite[] (used by the debugger to display the memories)
// return []MemoryMap memory maps attached to the Bus
func (bus *Bus) GetMemoryMaps() []MemoryWrite {
	memoryWrites := []MemoryWrite{}
	for _, memoryMap := range bus.memoryMaps {
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
func (bus *Bus) findMemory(address uint16) (*MemoryMap, error) {
	for _, memoryMap := range bus.memoryMaps {
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

// Accessible Interface Implementation

// Read the value at the given address.
// addr: uint16 address where the value will be read
// return uint8 value at the given address
// panic if the address is not found
func (bus *Bus) Read(addr uint16) uint8 {

	// DEBUG: JOYPAD should return 0xFF until implemented
	if addr == REG_FF00_JOYP || (addr >= 0xFEA0 && addr <= 0xFEFF) {
		return 0xFF
	}

	memoryMap, err := bus.findMemory(addr)
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
func (bus *Bus) Dump(from uint16, to uint16) []uint8 {
	memoryMap, err := bus.findMemory(from)
	if err == nil {
		return memoryMap.Memory.Dump(from-memoryMap.Address, to-memoryMap.Address)
	} else {
		panic(err)
	}
}

// Reads the next 2 bytes from the given address and returns them as a uint16 little-endian value (HIGH LOW).
func (bus *Bus) Read16(addr uint16) uint16 {
	return uint16(bus.Read(addr+1))<<8 | uint16(bus.Read(addr))
}

func (bus *Bus) write(addr uint16, value uint8) error {

	// DEBUG: Trying to write to 0x2000 and 0xFEA0-0xFEFF should be ignored
	if addr == 0x2000 || (addr >= 0xFEA0 && addr <= 0xFEFF) {
		return nil
	}

	// find the memory map containing the address
	memoryMap, err := bus.findMemory(addr)

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
	bus.addMemoryWrite(memoryWrite)
	return nil
}

// Write the provided value at the given address.
// addr: uint16 address where the value will be written
// value: uint8 value to write
// return void
// panic if the address is not found
func (bus *Bus) Write(addr uint16, value uint8) error {
	// on write to 0xFF50, disable the bootrom
	if addr == DISABLE_BOOT_ROM_REGISTER {
		bus.DisableBootRom()
	} else if addr == REG_FF04_DIV {
		// if the address is the divider register, reset it to 0
		value = 0
	} else if addr == REG_FF00_JOYP {
		// if the address is the joypad register, leave the handling to the joypad

	}

	// check if there is a write handler for the address
	if writeHandler, ok := bus.writeHandlers[addr]; ok {
		writeHandler(value)
		return nil
	}

	return bus.write(addr, value)
}

// Write the provided blob at the given address.
// Please note that the entire blob should belong to the same memory otherwise,
// the write process will panic once it reaches the end of the first memory,
// resulting in the partial write of the blob.
// param addr: uint16 address where the blob will be written
// param blob: []uint8 blob to write
// return void
// panic if one of the addresses in the range (addr, addr+len(blob)) is not found
func (bus *Bus) WriteBlob(addr uint16, blob []uint8) {
	memoryMap, err := bus.findMemory(addr)
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
		bus.addMemoryWrite(memoryWrite)
	} else {
		panic(err)
	}
}

// Disable the bootrom by removing it from the memory maps if a write operation to 0xFF50 is detected.
func (bus *Bus) DisableBootRom() {
	for idx, mem := range bus.memoryMaps {
		if mem.Name == BOOT_ROM_MEMORY_NAME {
			// remove the bootrom from the memory maps
			bus.memoryMaps = append(bus.memoryMaps[0:idx], bus.memoryMaps[idx+1:]...)
			return
		}
	}
}

// Special write operation for the timer registers when called by the Timer itself
// Allows the DIV register to be incremented, otherwise, any write for example from
// cartridge program to DIV register will reset it to 0.
func (bus *Bus) timerWrite(addr uint16, value uint8) {
	// find the memory map containing the address
	memoryMap, err := bus.findMemory(addr)
	// if the address is not found, return an error
	if err != nil {
		panic(err)
	}
	// write the value to the memory
	memoryMap.Memory.Write(addr-memoryMap.Address, value)
}
