package gameboy

import (
	"testing"

	"github.com/codefrite/gameboy-go/gameboy"
)

func TestNewMemory(t *testing.T) {
	memory := gameboy.NewMemory(0x2000)
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

func TestNewMemoryWithData(t *testing.T) {
	data := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		data[i] = uint8(i)
	}
	memory := gameboy.NewMemoryWithData(0x0100, data)
	for i := 0; i < 256; i++ {
		memoryCellData := memory.Read(uint16(i))
		if memoryCellData != uint8(i) {
			t.Errorf("Expected memory to be initialized with %02X at address %04X, got %02X", i, i, memoryCellData)
		}
	}
}

func TestNewMemoryWithDataUnmatchingSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected NewMemoryWithData to panic when data size does not match memory size")
		}
	}()
	data := make([]uint8, 256)
	gameboy.NewMemoryWithData(0xFFFF, data)
}

func TestNewMemoryWithRandomData(t *testing.T) {
	memorySize := uint16(0x2000)
	memory := gameboy.NewMemoryWithRandomData(memorySize)
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

func TestMemoryReadOutsideRange(t *testing.T) {
	// it shouldn't be possible to read outside the memory range
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Read to panic when reading outside the memory range")
		}
	}()
	memory := gameboy.NewMemory(10)
	memory.Read(11)
}

func TestMemoryWriteOutsideRange(t *testing.T) {
	// it shouldn't be possible to write outside the memory range
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected Write to panic when writing outside the memory range")
		}
	}()
	memory := gameboy.NewMemory(10)
	memory.Write(11, 0)
}

func TestMemoryDump(t *testing.T) {
	data := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		data[i] = uint8(i)
	}
	memory := gameboy.NewMemoryWithData(0x100, data)
	dump := memory.Dump(0, 256)
	for i := 0; i < 256; i++ {
		memoryCellData := dump[i]
		if memoryCellData != uint8(i) {
			t.Errorf("Expected memory to be initialized with %02X at address %04X, got %02X", i, i, memoryCellData)
		}
	}
}
