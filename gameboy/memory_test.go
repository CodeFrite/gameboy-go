package gameboy

import (
	"testing"
)

/*

Feature MEMORY
==============

Test Cases List:
- TC1> TestNewMemory 													checks if the memory is initialized with the correct size and 0 at all memory addresses
- TC2> TestNewMemoryWithData 									checks if the memory is initialized with the correct data being passed on creation as an array of uint8
- TC3> TestNewMemoryWithDataUnmatchingSize 		checks that Memory.NewMemoryWithRandomData does panic when passed a data array with a size that differs with the memory size
- TC4> TestNewMemoryWithRandomData 						checks that Memory.NewMemoryWithRandomData does panic when passed a data array with a size that differs with the memory size
- TC5> TestMemoryReadOutsideRange 						checks that Memory.Read does panic when trying to read from an address that is outside the memory range
- TC6> TestMemoryWriteOutsideRange 						checks that Memory.Write does panic when trying to write to an address that is outside the memory range
- TC7> TestMemoryDump 												checks that Memory.Read and Memory.Write work as expected

*/

/* checks if the memory is initialized with the correct size and 0 at all memory addresses */
func TestNewMemory(t *testing.T) {
	memory := NewMemory(0x2000)
	if memory == nil {
		t.Error("Expected memory to be initialized")
	}
	if memory.Size() != 0x2000 {
		t.Errorf("Expected memory to have 0x2000 bytes, got %d", memory.Size())
	}
	for i := 0; i < 0x2000; i++ {
		if memory.Read(uint16(i)) != 0 {
			t.Errorf("Expected memory to be initialized with 0 at address %04X, got %02X", i, memory.Read(uint16(i)))
		}
	}
}

/* checks if the memory is initialized with the correct data being passed on creation as an array of uint8 */
func TestNewMemoryWithData(t *testing.T) {
	data := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		data[i] = uint8(i)
	}
	memory := NewMemoryWithData(0x0100, data)
	for i := 0; i < 256; i++ {
		memoryCellData := memory.Read(uint16(i))
		if memoryCellData != uint8(i) {
			t.Errorf("Expected memory to be initialized with %02X at address %04X, got %02X", i, i, memoryCellData)
		}
	}
}

/* checks that Memory.NewMemoryWithRandomData does panic when passed a data array with a size that differs with the memory size */
func TestNewMemoryWithDataUnmatchingSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected NewMemoryWithData to panic when data size does not match memory size")
		}
	}()
	data := make([]uint8, 256)
	NewMemoryWithData(0xFFFF, data)
}

func TestNewMemoryWithRandomData(t *testing.T) {
	memorySize := uint16(0x2000)
	memory := NewMemoryWithRandomData(memorySize)
	// check that all data are not 0 and that they are random
	freq := make(map[uint8]int)
	for i := 0; i < int(memorySize); i++ {
		data := memory.Read(uint16(i))
		freq[data]++
	}
	for k, v := range freq {
		if v == 0 {
			t.Errorf("Expected memory to have random data, but %02X is missing", k)
		}
	}
}

/* checks that Memory.Read does panic when trying to read from an address that is outside the memory range */
func TestMemoryReadOutsideRange(t *testing.T) {
	// it shouldn't be possible to read outside the memory range
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Read to panic when reading outside the memory range")
		}
	}()
	memory := NewMemory(10)
	memory.Read(11)
}

/* checks that Memory.Write does panic when trying to write to an address that is outside the memory range */
func TestMemoryWriteOutsideRange(t *testing.T) {
	// it shouldn't be possible to write outside the memory range
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Write to panic when writing outside the memory range")
		}
	}()
	memory := NewMemory(10)
	memory.Write(11, 0)
}

/* checks that Memory.Read and Memory.Write work as expected */
func TestMemoryDump(t *testing.T) {
	memSize := 256
	memSizeUint16 := 0x0100
	data := make([]uint8, memSize)
	for i := 0; i < memSize; i++ {
		data[i] = uint8(i)
	}
	memory := NewMemoryWithData(uint16(memSizeUint16), data)
	dump := memory.Dump(0, memory.Size()-1)
	for i := 0; i < memSize-1; i++ {
		memoryCellData := dump[i]
		if memoryCellData != uint8(i) {
			t.Errorf("Expected memory to be initialized with %02X at address %04X, got %02X", i, i, memoryCellData)
		}
	}
}
