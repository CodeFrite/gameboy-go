package gameboy

import "math/rand"

/**
 * Memory struct represents any memory type RAM or ROM in the GameBoy.
 * It is used to hold uint8 values and is indexed by a uint16 address.
 * It is used here to represent the BOOTROM, cartridge ROM and RAM, Video RAM, Working RAM, High RAM, OAM, ...
 * While CPU registers are defined inside the CPU struct as []uint8/[]uint16, the IE register is defined as a memory cell.
 */
type Memory struct {
	data []uint8
}

/**
 * NewMemory creates a new Memory object with the given size.
 * All memory locations are initialized to 0.
 */
func NewMemory(size uint16) *Memory {
	return &Memory{data: make([]uint8, size)}
}

/**
 * NewMemoryWithData creates a new Memory object with the given size and data.
 * The data is copied into the memory. The data size must match the memory size.
 */
func NewMemoryWithData(size uint16, data []uint8) *Memory {
	if len(data) != int(size) {
		panic("Data size does not match memory size")
	}
	return &Memory{data: data}
}

/**
 * NewMemoryWithRandomData creates a new Memory object with the given size and random data.
 * It is used to simulate the initial state of the memories in the gameboy which state cannot be predicted.
 * It allows us to see the bootrom clearing the VRAM on startup.
 */
func NewMemoryWithRandomData(size uint16) *Memory {
	data := make([]uint8, size)
	for i := 0; i < len(data); i++ {
		data[i] = uint8(rand.Intn(256))
	}
	return &Memory{data: data}
}

/**
 * returns the size of the memory in bytes.
 */
func (m *Memory) Size() uint16 {
	return uint16(len(m.data))
}

/**
 * Read returns the uint8 value at the given address.
 * panics if the address is out of bounds.
 */
func (m *Memory) Read(addr uint16) uint8 {
	if addr >= uint16(len(m.data)) {
		panic("Memory address out of bounds while reading")
	}
	return m.data[addr]
}

/**
 * Dump memory from address 'from' to address 'to'
 * panics if the addresses are out of bounds.
 */
func (m *Memory) Dump(from uint16, to uint16) []uint8 {
	if from >= m.Size() || to >= m.Size() {
		panic("Memory address out of bounds while dumping")
	}
	return m.data[from:to]
}

/**
 * Write the provided uint8 value at the given address.
 * panics if the address is out of bounds.
 */
func (m *Memory) Write(addr uint16, value uint8) {
	if addr >= uint16(len(m.data)) {
		panic("Memory address out of bounds while writing")
	}
	m.data[addr] = value
}
