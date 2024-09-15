package gameboy

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/codefrite/gameboy-go/gameboy"
)

/*

Feature non CB instructions
===========================

Test Cases List:

- TC1> NOP: should not change anything in the gameboy except the program counter and the clock
- TC2> STOP: should stop the gameboy (not implemented yet)
- TC3> HALT: should halt the gameboy by setting the HALT flag to true
- TC4> DI: should disable interrupts by setting the IME flag to false
- TC5> EI: should enable interrupts by setting the IME flag to true
- TC6> JP: should jump to the address specified in the instruction
- TC7> JR: should jump to the address specified in the instruction
- TC8> CALL: should call the address specified in the instruction
- TC9> RET: should return from a subroutine
- TC10> RETI: should return from a subroutine and enable interrupts
- TC11> RST: should call the address specified in the instruction
- TC12> LD: should load the value from the source into the destination
- TC13> LDH: should load the value from the source into the destination
- TC14> PUSH: should push the value from the source into the stack
- TC15> POP: should pop the value from the stack into the destination
- TC16> ADD: should add the value from the source to the destination
- TC17> ADC: should add the value from the source to the destination with the carry
- TC18> AND: should perform a bitwise AND between the source and the destination
- TC19> INC: should increment the value of the destination
- TC20> CCF: should flip the carry flag
- TC21> CP: should compare the value from the source to the destination
- TC22> CPL: should flip all bits of the destination
- TC23> DAA: should adjust the destination to be a valid BCD number
- TC24> DEC: should decrement the value of the destination
- TC25> SUB: should subtract the value from the source to the destination
- TC26> SBC: should subtract the value from the source to the destination with the carry
- TC27> SCF: should set the carry flag
- TC28> OR: should perform a bitwise OR between the source and the destination
- TC29> XOR: should perform a bitwise XOR between the source and the destination
- TC30> RLA: should rotate the destination left through the carry
- TC31> RLCA: should rotate the destination left
- TC32> RRA: should rotate the destination right through the carry
- TC33> RRCA: should rotate the destination right

*/

// global variables
var (
	bus     *gameboy.Bus
	memory1 *gameboy.Memory
	memory2 *gameboy.Memory
	cpu     *gameboy.CPU

	cpuState1 *gameboy.CpuState
	cpuState2 *gameboy.CpuState
)

/* SUPPORT FUNCs */

// allow func to panic without stopping the test
func mayPanic(f func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	f()
	return
}

// initialize the test environment with the following preconditions:
// create a bus /
// create two 8KB memory and attach them to the bus /
// create a cpu
// initialize the cpu states
func preconditions() {
	// create a bus
	bus = gameboy.NewBus()
	// create a first memory and attach it to the bus
	memory1 = gameboy.NewMemory(0x2000)
	bus.AttachMemory("RAM 1", 0x0000, memory1)
	// create a second memory and attach it to the bus
	memory2 = gameboy.NewMemory(0x2000)
	bus.AttachMemory("RAM 2", 0x2000, memory2)
	// create a cpu
	cpu = gameboy.NewCPU(bus)
	cpu.PC = 0x0000
	cpu.SP = 0xFFFE
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
func getCpuState() *gameboy.CpuState {
	return &gameboy.CpuState{
		PC:            cpu.PC,
		SP:            cpu.SP,
		A:             cpu.A,
		F:             cpu.F,
		Z:             cpu.F&0x80 != 0,
		N:             cpu.F&0x40 != 0,
		H:             cpu.F&0x20 != 0,
		C:             cpu.F&0x10 != 0,
		BC:            uint16(cpu.B)<<8 | uint16(cpu.C),
		DE:            uint16(cpu.D)<<8 | uint16(cpu.E),
		HL:            uint16(cpu.H)<<8 | uint16(cpu.L),
		PREFIXED:      cpu.Prefixed,
		IR:            cpu.IR,
		OPERAND_VALUE: cpu.Operand,
		IE:            cpu.IE.Read(0),
		IME:           cpu.IME,
		HALTED:        cpu.Halted,
	}
}

func printCpuState(cpuState *gameboy.CpuState) {
	fmt.Println(" ***   *** *** ***   *** ***   ***   *** *** ***   ***   *** ***   *** *** ***   ***")
	fmt.Println("CPU STATE:")
	fmt.Printf("PC: 0x%4X\n", cpuState.PC)
	fmt.Printf("SP: 0x%4X\n", cpuState.SP)
	fmt.Printf("A: 0x%2X\n", cpuState.A)
	fmt.Printf("F: 0x%2X\n", cpuState.F)
	fmt.Printf("Z: %t\n", cpuState.Z)
	fmt.Printf("N: %t\n", cpuState.N)
	fmt.Printf("H: %t\n", cpuState.H)
	fmt.Printf("C: %t\n", cpuState.C)
	fmt.Printf("BC: 0x%4X\n", cpuState.BC)
	fmt.Printf("DE: 0x%4X\n", cpuState.DE)
	fmt.Printf("HL: 0x%4X\n", cpuState.HL)
	fmt.Printf("PREFIXED: %t\n", cpuState.PREFIXED)
	fmt.Printf("IR: 0x%2X\n", cpuState.IR)
	fmt.Printf("OPERAND_VALUE: 0x%2X\n", cpuState.OPERAND_VALUE)
	fmt.Println("IE:", cpuState.IE)
	fmt.Println("IME:", cpuState.IME)
	fmt.Println("HALTED:", cpuState.HALTED)
}

// shift the state of the cpu from mem1 to mem2
func shiftCpuState(mem1 *gameboy.CpuState, mem2 *gameboy.CpuState) {
	*mem2 = gameboy.CpuState{
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
func loadProgramIntoMemory(memory *gameboy.Memory, program []uint8) {
	for idx, val := range program {
		memory.Write(uint16(idx), val)
	}
}

func compareCpuState(mem1 *gameboy.CpuState, mem2 *gameboy.CpuState) []string {
	result := make([]string, 0)
	// Loop over the fields of the CpuState struct
	v1 := reflect.ValueOf(*mem1)
	v2 := reflect.ValueOf(*mem2)
	typeOfS := v1.Type()

	for i := 0; i < v1.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		val1 := v1.Field(i).Interface()
		val2 := v2.Field(i).Interface()
		if val1 != val2 {
			result = append(result, fieldName)
		}
	}
	return result
}

/* TEST CASES */

/* TC1: should not change anything in the gameboy except the program counter and the clock */

func TestNOP(t *testing.T) {
	preconditions()

	// test program : 0xFF NOP instructions
	var testData []uint8 = make([]uint8, 0xFF)
	for i := 0; i < 0x0F; i++ {
		testData[i] = 0x00
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// shift mem1 to mem2
	cpu.Step()
	shiftCpuState(cpuState1, cpuState2)

	for i := 0; i < len(testData); i++ {
		// execute the next instruction and shift the saved state mem1 to mem2 and save the new state in mem1
		shiftCpuState(cpuState1, cpuState2)
		cpu.Step()
		cpuState1 = getCpuState()
		cmp := compareCpuState(cpuState1, cpuState2)
		// check if there is only one difference between the two states: PC incremented to i
		if len(cmp) != 1 {
			t.Errorf("Error> NOP instruction should change exactly one field, the PC. Here got %v", cmp)
		} else {

			// the key should be PC
			if cmp[0] != "PC" {
				t.Errorf("Error> NOP instruction should change the PC field, here got %v\n", cmp[0])
			}
		}
	}

	postconditions()
}

/* TC2: should stop the gameboy */
func TestSTOP(t *testing.T) {
	preconditions()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check if the gameboy is halted on the STOP instruction
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}
	finalState := getCpuState()
	// check if the last PC = 0x0005, position of the STOP instruction in the test data program
	if finalState.PC != 0x0005 {
		t.Errorf("Error> STOP instruction: the program counter should have stopped at the STOP instruction @0x0005, got @0x%X4 \n", finalState.PC)
	}

	/*
		// debugging output
		printCpuState(finalState)
	*/

	postconditions()
}

/* TC3: should halt the gameboy by setting the HALT flag to true */
func TestHALT(t *testing.T) {
	preconditions()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x76, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check if the gameboy is halted on the HALT instruction
	if !cpu.Halted {
		t.Errorf("Error> HALT instruction should halt the gameboy\n")
	}
	finalState := getCpuState()
	// check if the last PC = 0x0001, position of the HALT instruction in the test data program
	if finalState.PC != 0x0005 {
		t.Errorf("Error> HALT instruction: the program counter should have stopped at the HALT instruction\n")
	}

	postconditions()
}

/* TC4: should disable interrupts by setting the IME flag to false */
func TestDI(t *testing.T) {
	preconditions()
	cpu.IME = true

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xF3, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program and control step by step the IME flag
	for i := 0; i < len(testData); i++ {
		cpu.Step()
		/*
			// debugging output
			fmt.Println()
			fmt.Println(" ***   *** *** ***   *** ***   ***   *** *** ***   ***   *** ***   *** *** ***   ***")
			fmt.Println(i, ">")
			printCpuState(getCpuState())
		*/

		// the IME flag should stay up until the end of the execution after the DI instruction
		if i >= 0 && i <= 5 {
			if !cpu.IME {
				t.Errorf("Error> DI instruction should disable the IME flag after the execution of the next instruction\n")
			}
		} else if i >= 6 {
			if cpu.IME {
				t.Errorf("Error> DI instruction should disable the IME flag\n")
			}
		}
	}

	postconditions()
}

/* TC5: should enable interrupts by setting the IME flag to true */
func TestEI(t *testing.T) {
	preconditions()
	cpu.IME = false

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0xFB, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program and control step by step the IME flag
	for i := 0; i < len(testData); i++ {
		cpu.Step()
		/*
			// debugging output
			fmt.Println()
			fmt.Println(" ***   *** *** ***   *** ***   ***   *** *** ***   ***   *** ***   *** *** ***   ***")
			fmt.Println(i, ">")
			printCpuState(getCpuState())
		*/
		// the IME flag should stay down until the end of the execution after the EI instruction
		if i >= 0 && i <= 5 {
			if cpu.IME {
				t.Errorf("Error> EI instruction should enable the IME flag after the execution of the next instruction\n")
			}
		} else if i >= 6 {
			if !cpu.IME {
				t.Errorf("Error> EI instruction should enable the IME flag\n")
			}
		}
	}

	postconditions()
}

// TC6: should jump to the address specified in the instruction
// opcodes:
//   - 0xC3 = JP     a16
//   - 0xE9 = JP HL
//   - 0xCA = JP  Z, a16
//   - 0xC2 = JP NZ, a16
//   - 0xDA = JP  C, a16
//   - 0xD2 = JP NC, a16
//   - flags: -
func TestJP(t *testing.T) {

	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JP a16 / JP HL / JP Z, a16 / JP C, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.F

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC3, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JP a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xE9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JP HL ; precondition: HL = 0x00D0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xDA, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 4	;		JP C, a16(0x00B0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	; step 5	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0xCA, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	; step 3	;		JP Z, a16(0x0040)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xC3, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 1	;		JP a16(0x0020)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00B0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00B0 {
		t.Errorf("Error> JP instruction: the program counter should have stopped at the STOP instruction @0x00B0, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> JP instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JP a16 / JP HL / JP NZ, a16 / JP NC, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.F

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC3, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JP a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xE9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JP HL ; precondition: HL = 0x00D0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xD2, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 4	;		JP NC, a16(0x00B0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	; step 5	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0xC2, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	; step 3	;		JP NZ, a16(0x0040)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xC3, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 1	;		JP a16(0x0020)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00B0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00B0 {
		t.Errorf("Error> JP instruction: the program counter should have stopped at the STOP instruction @0x00B0, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> JP instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

}

// TC7: should jump to the address specified in the instruction
// Jumps to a relative address from the next instruction
//
//	opcodes:
//	- 0x18 = JR r8
//	- 0x20 = JR NZ, r8
//	- 0x28 = JR Z, r8
//	- 0x30 = JR NC, r8
//	- 0x38 = JR C, r8
//	- flags: -
func TestJR(t *testing.T) {
	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JR a16 / JR HL / JR Z, a16 / JR C, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.F

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x18, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x28, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR Z, r8(0x40) ; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x18, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8(255 - 32 = 223 = DF) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x28, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR Z, r8(0x40)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x38, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR C, r8(0x40) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x38, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR C, r8(0x60) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x38, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR C, r8(255-80 = 175 = AF) ; precondition: C = 0
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JR a16 / JR HL / JR NZ, a16 / JR NC, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.F

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x18, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x20, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR NZ, r8(0x40) ; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x18, 0xE0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8(255 - 32 = 223 = DF) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x20, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR NZ, r8(0x40)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x30, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR NC, r8(0x40) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x30, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR NC, r8(0x60) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x30, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR NC, r8(255-80 = 175 = AF) ; precondition: C = 0
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}

// TC8: should call the address specified in the instruction
// opcodes:
//   - 0xCD = CALL a16
//   - 0xCC = CALL Z, a16
//   - 0xC4 = CALL NZ, a16
//   - 0xDC = CALL C, a16
//   - 0xD4 = CALL NC, a16
//   - flags: -
func TestCALL(t *testing.T) {
	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions CALL a16 / CALL HL / CALL Z, a16 / CALL C, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.F

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCD, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xCC, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL Z, a16(0x00D0) ; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xDC, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 4	;		CALL C, a16(0x00B0) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 6	;		STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0xDC, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	; step 5	;   CALL C, a16(0x0060) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0xCC, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	; step 3	;		CALL Z, a16(0x0040)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 1	;		CALL a16(0x0020)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00B0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x0060 {
		t.Errorf("Error> CALL instruction: the program counter should have stopped at the STOP instruction @0x0060, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> CALL instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions CALL a16 / CALL HL / CALL NZ, a16 / CALL NC, a16
	 */

	// preconditions
	preconditions()
	cpu.H = 0x00
	cpu.L = 0xD0
	cpu.F = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.F

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCD, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xC4, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL NZ, a16(0x00D0) ; precondition: Z = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xD4, 0xB0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 4	;		CALL NC, a16(0x00B0) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 6	;		STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0xD4, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	; step 5	;   CALL NC, a16(0x0060) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0xC4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	; step 3	;		CALL NZ, a16(0x0040) ; precondition: Z = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 1	;		CALL a16(0x0020)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00B0
	if !cpu.Stopped {
		t.Errorf("Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x0060 {
		t.Errorf("Error> CALL instruction: the program counter should have stopped at the STOP instruction @0x0060, got @0x%X4 \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("Error> CALL instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}

/* TC9: should return from a subroutine */
func TestRET(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC10: should return from a subroutine and enable interrupts */
func TestRETI(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC11: should call the address specified in the instruction */
func TestRST(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC12: should load the value from the source into the destination */
func TestLD(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC13: should load the value from the source into the destination */
func TestLDH(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC14: should push the value from the source into the stack */
func TestPUSH(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC15: should pop the value from the stack into the destination */
func TestPOP(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC16: should add the value from the source to the destination */
func TestADD(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC17: should add the value from the source to the destination with the carry */
func TestADC(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC18: should perform a bitwise AND between the source and the destination */
func TestAND(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC19: should increment the value of the destination */
func TestINC(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC20: should flip the carry flag */
func TestCCF(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC21: should compare the value from the source to the destination */
func TestCP(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC22: should flip all bits of the destination */
func TestCPL(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC23: should adjust the destination to be a valid BCD number */
func TestDAA(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC24: should decrement the value of the destination */
func TestDEC(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC25: should subtract the value from the source to the destination */
func TestSUB(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC26: should subtract the value from the source to the destination with the carry */
func TestSBC(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC27: should set the carry flag */
func TestSCF(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC28: should perform a bitwise OR between the source and the destination */
func TestOR(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC29: should perform a bitwise XOR between the source and the destination */
func TestXOR(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC30: should rotate the destination left through the carry */
func TestRLA(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC31: should rotate the destination left */
func TestRLCA(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC32: should rotate the destination right through the carry */
func TestRRA(t *testing.T) {
	t.Error("not implemented yet")
}

/* TC33: should rotate the destination right */
func TestRRCA(t *testing.T) {
	t.Error("not implemented yet")
}