package gameboy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
 * JSONableSlice is a type alias for a slice of uint8 used to make it JSONable by defining a custom JSON marshalling method.
 */
type JSONableSlice []uint8

/**
 * represents a memory write operation used by the mmu to keep track of all memory changes between 2 states
 * memory writes can be either reset:
 * - before every step from the debugger
 * - before a new run from a breakpoint to another breakpoint
 */
type MemoryWrite struct {
	Name    string        `json:"name"`
	Address uint16        `json:"address"`
	Data    JSONableSlice `json:"data"`
}

// MarshalJSON is a custom JSON marshalling method for the JSONableSlice type
// It converts the slice of uint8 to a string of comma-separated values
func (j JSONableSlice) MarshalJSON() ([]byte, error) {
	strValues := make([]string, len(j))
	for i, v := range j {
		strValues[i] = fmt.Sprintf("%d", v)
	}
	result := strings.Join(strValues, ",")
	return json.Marshal(result)
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

/**
 * Attach a memory to the MMU at the given address.
 * At the moment, there is not check for overlapping memories: when reading from or writing to an address contained in
 * multiple memories, the first one found will be used. TODO: forbid overlapping memories to be attached.
 * @param name: string name of the memory
 * @param address: uint16 address where the memory will be attached
 * @param memory: Accessible memory to attach
 * @return void
 */
func (m *MMU) AttachMemory(name string, address uint16, memory Accessible) {
	m.memoryMaps = append(m.memoryMaps, MemoryMap{
		Name:    name,
		Address: address,
		Memory:  memory,
	})
}

/**
 * return the memory maps attached to the MMU as MemoryWrite[] (used by the debugger to display the memories)
 * @return []MemoryMap memory maps attached to the MMU
 */
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
func (b *MMU) Write(addr uint16, value uint8) error {
	memoryMap, err := b.findMemory(addr)
	if err == nil {
		memoryMap.Memory.Write(addr-memoryMap.Address, value)
		memoryWrite := MemoryWrite{
			Name:    memoryMap.Name,
			Address: addr - memoryMap.Address,
			Data:    []uint8{value},
		}
		b.addMemoryWrite(memoryWrite)
		return nil
	} else {
		return err
	}
}

/**
 * Write the provided blob at the given address.
 * Please note that the entire blob should belong to the same memory otherwise,
 * the write process will panic once it reaches the end of the first memory,
 * resulting in the partial write of the blob.
 * @param addr: uint16 address where the blob will be written
 * @param blob: []uint8 blob to write
 * @return void
 * @panic if one of the addresses in the range (addr, addr+len(blob)) is not found
 */
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
