package gameboy

import (
	"fmt"
	"math/rand"
	"reflect"
)

// global variables
var (
	bus     *Bus
	memory1 *Memory
	memory2 *Memory
	cpu     *CPU

	cpuState1 *CpuState
	cpuState2 *CpuState
)

// initialize the test environment with the following preconditions:
// create a bus /
// create two 8KB memory and attach them to the bus /
// create a cpu
// initialize the cpu states
func preconditions() {
	// create a bus
	bus = NewBus()
	// create a first memory and attach it to the bus
	memory1 = NewMemory(0x2000)
	bus.AttachMemory("RAM 1", 0x0000, memory1)
	// create a second memory and attach it to the bus
	memory2 = NewMemory(0xDFFF)
	bus.AttachMemory("RAM 2", 0x2000, memory2)
	// create IE memory and attach it to the bus
	ie := NewMemory(0x0001)
	bus.AttachMemory("IE", 0xFFFF, ie)

	// create a cpu
	cpu = NewCPU(bus)
	cpu.pc = 0x0000
	cpu.sp = 0xFFFE
	cpu.halted = false
	cpu.stopped = false

	// initialize the cpu states
	cpuState := getCpuState()
	cpuState1 = cpuState
	cpuState2 = cpuState
}

// clean up the test environment by setting all the variables to nil
func postconditions() {
	// clean up
	bus = nil
	memory1 = nil
	memory2 = nil
	cpu = nil
	cpuState1 = nil
	cpuState2 = nil
}

// save the state of the cpu
func getCpuState() *CpuState {
	return &CpuState{
		PC:            cpu.pc,
		SP:            cpu.sp,
		A:             cpu.a,
		F:             cpu.f,
		Z:             cpu.f&0x80 != 0,
		N:             cpu.f&0x40 != 0,
		H:             cpu.f&0x20 != 0,
		C:             cpu.f&0x10 != 0,
		BC:            uint16(cpu.b)<<8 | uint16(cpu.c),
		DE:            uint16(cpu.d)<<8 | uint16(cpu.e),
		HL:            uint16(cpu.h)<<8 | uint16(cpu.l),
		PREFIXED:      cpu.prefixed,
		IR:            cpu.ir,
		OPERAND_VALUE: cpu.operand,
		IE:            cpu.GetIEFlag(),
		IME:           cpu.ime,
		HALTED:        cpu.halted,
	}
}

func printCpuState(cpuState *CpuState) {
	fmt.Println(" ***   *** *** ***   *** ***   ***   *** *** ***   ***   *** ***   *** *** ***   ***")
	fmt.Println("CPU STATE:")
	fmt.Printf("PC: 0x%04X\n", cpuState.PC)
	fmt.Printf("SP: 0x%04X\n", cpuState.SP)
	fmt.Printf("A: 0x%02X\n", cpuState.A)
	fmt.Printf("F: 0x%02X\n", cpuState.F)
	fmt.Printf("Z: %t\n", cpuState.Z)
	fmt.Printf("N: %t\n", cpuState.N)
	fmt.Printf("H: %t\n", cpuState.H)
	fmt.Printf("C: %t\n", cpuState.C)
	fmt.Printf("BC: 0x%04X\n", cpuState.BC)
	fmt.Printf("DE: 0x%04X\n", cpuState.DE)
	fmt.Printf("HL: 0x%04X\n", cpuState.HL)
	fmt.Printf("PREFIXED: %t\n", cpuState.PREFIXED)
	fmt.Printf("IR: 0x%02X\n", cpuState.IR)
	fmt.Printf("OPERAND_VALUE: 0x%02X\n", cpuState.OPERAND_VALUE)
	fmt.Println("IE:", cpuState.IE)
	fmt.Println("IME:", cpuState.IME)
	fmt.Println("HALTED:", cpuState.HALTED)
}

// shift the state of the cpu from mem1 to mem2
func shiftCpuState(mem1 *CpuState, mem2 *CpuState) {
	*mem2 = CpuState{
		PC:            mem1.PC,
		SP:            mem1.SP,
		A:             mem1.A,
		F:             mem1.F,
		Z:             mem1.Z,
		N:             mem1.N,
		H:             mem1.H,
		C:             mem1.C,
		BC:            mem1.BC,
		DE:            mem1.DE,
		HL:            mem1.HL,
		PREFIXED:      mem1.PREFIXED,
		IR:            mem1.IR,
		OPERAND_VALUE: mem1.OPERAND_VALUE,
		IE:            mem1.IE,
		IME:           mem1.IME,
		HALTED:        mem1.HALTED,
	}
}

// load program into the memory starting from the address 0x0000
func loadProgramIntoMemory(memory *Memory, program []uint8) {
	for idx, val := range program {
		memory.Write(uint16(idx), val)
	}
}

func compareCpuState(mem1 *CpuState, mem2 *CpuState) []string {
	result := make([]string, 0)
	// Loop over the fields of the CpuState struct
	v1 := reflect.ValueOf(*mem1)
	v2 := reflect.ValueOf(*mem2)
	typeOfS := v1.Type()

	for i := 0; i < v1.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		val1 := v1.Field(i).Interface()
		val2 := v2.Field(i).Interface()
		if fieldName != "INSTRUCTION" && val1 != val2 {
			result = append(result, fieldName)
		}
	}
	return result
}

func printMemoryProperties() {
	memoryMaps := bus.mmu.GetMemoryMaps()
	fmt.Println("\n> Memory Mapping:")
	fmt.Println("-----------------")
	for _, memoryMap := range memoryMaps {
		fmt.Printf("> Memory %s: %d bytes @ 0x%04X->0x%04X\n", memoryMap.Name, len(memoryMap.Data), memoryMap.Address, memoryMap.Address+uint16(len(memoryMap.Data))-1)
	}
}

func randomizeFlags() {
	randomBit := rand.Intn(2)
	if randomBit == 0 {
		cpu.resetZFlag()
	} else {
		cpu.setZFlag()
	}
	randomBit = rand.Intn(2)
	if randomBit == 0 {
		cpu.resetNFlag()
	} else {
		cpu.setNFlag()
	}
	randomBit = rand.Intn(2)
	if randomBit == 0 {
		cpu.resetHFlag()
	} else {
		cpu.setHFlag()
	}
	randomBit = rand.Intn(2)
	if randomBit == 0 {
		cpu.resetCFlag()
	} else {
		cpu.setCFlag()
	}
}
