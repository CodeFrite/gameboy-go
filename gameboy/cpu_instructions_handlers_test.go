package gameboy

import (
	"fmt"
	"testing"
)

/*

Testing non CB instructions
===========================

Test Cases List:

- NOP: should not change anything in the gameboy except the program counter and the clock
- STOP: should stop the gameboy (not implemented yet)
- HALT: should halt the gameboy by setting the HALT flag to true
- DI: should disable interrupts by setting the IME flag to false
- EI: should enable interrupts by setting the IME flag to true
- JP: should jump to the address specified in the instruction
- JR: should jump to the address specified in the instruction
- CALL: should call the address specified in the instruction
- RET: should return from a subroutine
- RETI: should return from a subroutine and enable interrupts
- RST: should call the address specified in the instruction
- LD: should load the value from the source into the destination
- LDH: should load the value from the source into the destination
- PUSH: should push the value from the source into the stack
- POP: should pop the value from the stack into the destination
- ADD: should add the value from the source to the destination
- ADC: should add the value from the source to the destination with the carry
- AND: should perform a bitwise AND between the source and the destination
- INC: should increment the value of the destination
- CCF: should flip the carry flag
- CP: should compare the value from the source to the destination
- CPL: should flip all bits of the destination
- DAA: should adjust the destination to be a valid BCD number
- DEC: should decrement the value of the destination
- SUB: should subtract the value from the source to the destination
- SBC: should subtract the value from the source to the destination with the carry
- SCF: should set the carry flag
- OR: should perform a bitwise OR between the source and the destination
- XOR: should perform a bitwise XOR between the source and the destination
- RLA: should rotate the destination left through the carry
- RLCA: should rotate the destination left
- RRA: should rotate the destination right through the carry
- RRCA: should rotate the destination right

*/

// TEST CASES

// NOP: should not change anything in the gameboy except the program counter and the clock
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
			t.Errorf("[TestNOP_CHK_1] Error> NOP instruction should change exactly one field, the PC. Here got %v", cmp)
		} else {

			// the key should be PC
			if cmp[0] != "PC" {
				t.Errorf("[TestNOP_CHK_2] Error> NOP instruction should change the PC field, here got %v\n", cmp[0])
			}

			// PC should be equal to 0x0E
			if cpuState1.PC != uint16(i+1) {
				t.Errorf("[TestNOP_CHK_3] Error> NOP instruction should increment the PC by 1, here got %v\n", cpuState1.PC)
			}
		}
		// NOP instruction shouldn't change the flags
		if cpuState1.F != cpuState2.F {
			t.Errorf("[TestNOP_CHK_4] Error> NOP instruction shouldn't change the flags\n")
		}
	}

	postconditions()
}

// STOP: should stop the gameboy
func TestSTOP(t *testing.T) {
	preconditions()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check if the gameboy is halted on the STOP instruction
	if !cpu.stopped {
		t.Errorf("[TestSTOP_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}
	finalState := getCpuState()
	// check if the last PC = 0x0005, position of the STOP instruction in the test data program
	if finalState.PC != 0x0005 {
		t.Errorf("[TestSTOP_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x0005, got @0x%04X \n", finalState.PC)
	}

	/*
		// debugging output
		printCpuState(finalState)
	*/

	postconditions()
}

// HALT: should halt the gameboy by setting the HALT flag to true
func TestHALT(t *testing.T) {
	preconditions()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x76, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check if the gameboy is halted on the HALT instruction
	if !cpu.halted {
		t.Errorf("[TestHALT_CHK_1] Error> HALT instruction should halt the gameboy\n")
	}
	finalState := getCpuState()
	// check if the last PC = 0x0001, position of the HALT instruction in the test data program
	if finalState.PC != 0x0005 {
		t.Errorf("[TestHALT_CHK_2] Error> HALT instruction: the program counter should have stopped at the HALT instruction\n")
	}

	postconditions()
}

// DI: should disable interrupts by setting the IME flag to false
func TestDI(t *testing.T) {
	preconditions()
	cpu.ime = true

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
			if !cpu.ime {
				t.Errorf("[TestDI_CHK_1] Error> DI instruction should disable the IME flag after the execution of the next instruction\n")
			}
		} else if i >= 6 {
			if cpu.ime {
				t.Errorf("[TestDI_CHK_2] Error> DI instruction should disable the IME flag\n")
			}
		}
	}

	postconditions()
}

// EI: should enable interrupts by setting the IME flag to true
func TestEI(t *testing.T) {
	preconditions()
	cpu.ime = false

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
			if cpu.ime {
				t.Errorf("[TestEI_CHK_1] Error> EI instruction should enable the IME flag after the execution of the next instruction\n")
			}
		} else if i >= 6 {
			if !cpu.ime {
				t.Errorf("[TestEI_CHK_2] Error> EI instruction should enable the IME flag\n")
			}
		}
	}

	postconditions()
}

// JP: conditional jump to the address specified in the instruction
// opcodes:
//   - 0xC2 = JP NZ, a16
//   - 0xC3 = JP     a16
//   - 0xCA = JP  Z, a16
//   - 0xD2 = JP NC, a16
//   - 0xDA = JP  C, a16
//   - 0xE9 = JP HL
//   - flags: -
func TestJP(t *testing.T) {
	t.Run("0xC2: JP NZ, a16", test_0xC2_JP_NZ_a16)
	t.Run("0xC3: JP a16", test_0xC3_JP_a16)
	t.Run("0xCA: JP Z, a16", test_0xCA_JP_Z_a16)
	t.Run("0xD2: JP NC, a16", test_0xD2_JP_NC_a16)
	t.Run("0xDA: JP C, a16", test_0xDA_JP_C_a16)
	t.Run("0xE9: JP HL", test_0xE9_JP_HL)
	t.Run("JP integration", test_JP_integration)
}
func test_0xC3_JP_a16(t *testing.T) {
	// TC1: positive case: should always jump to the a16 address specified in the 3bytes instruction
	preconditions()
	cpu.resetZFlag()
	saveFlags := cpu.f

	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC3, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0xC3, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xC3_JP_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState.PC != 0x001F {
		t.Errorf("[test_0xC3_JP_a16_CHK_2] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x001F, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xC3_JP_a16_CHK_3] Error> JP instruction: no flags should have changed\n")
	}
}
func test_0xE9_JP_HL(t *testing.T) {
	// TC1: positive case: should always jump to the address stored in HL
	preconditions()
	cpu.setHL(0x0020)
	saveFlags := cpu.f

	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xE9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xE9_JP_HL_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState.PC != 0x002F {
		t.Errorf("[test_0xE9_JP_HL_CHK_2] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x002F, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xE9_JP_HL_CHK_3] Error> JP instruction: no flags should have changed\n")
	}
}
func test_0xCA_JP_Z_a16(t *testing.T) {
	// TC1: positive case: should jump to a16 operand address if Z=1
	preconditions()
	cpu.setZFlag()
	saveFlags1 := cpu.f

	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCA, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0xCA, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState1 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState1.PC != 0x001F {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x001F, got @0x%04X \n", finalState1.PC)
	}

	// no flags should have changed
	if finalState1.F != saveFlags1 {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_3] Error> no flags should have changed\n")
	}

	// TC2: positive case: should not jump if Z=0
	preconditions()
	cpu.resetZFlag()
	saveFlags2 := cpu.f

	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCA, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // no jump // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, //
		0x00, 0x00, 0x00, 0x00, 0xCA, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //
	}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState2 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState2.PC != 0x000F {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_5] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x000F, got @0x%04X \n", finalState2.PC)
	}

	// no flags should have changed
	if finalState2.F != saveFlags2 {
		t.Errorf("[test_0xCA_JP_Z_a16_CHK_6] Error> JP instruction: no flags should have changed\n")
	}
}
func test_0xC2_JP_NZ_a16(t *testing.T) {
	// TC1: positive case: should jump if Z=0
	preconditions()
	cpu.resetZFlag()
	saveFlags1 := cpu.f

	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC2, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0xC2, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState1 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState1.PC != 0x001F {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_2] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x001F, got @0x%04X \n", finalState1.PC)
	}

	// no flags should have changed
	if finalState1.F != saveFlags1 {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_3] Error> JP instruction: no flags should have changed\n")
	}

	// TC2: positive case: should not jump if Z=1
	preconditions()
	cpu.setZFlag()
	saveFlags2 := cpu.f

	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC2, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // no jump // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, //
		0x00, 0x00, 0x00, 0x00, 0xC2, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //
	}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState2 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState2.PC != 0x000F {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_5] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x000F, got @0x%04X \n", finalState2.PC)
	}

	// no flags should have changed
	if finalState2.F != saveFlags2 {
		t.Errorf("[test_0xC2_JP_NZ_a16_CHK_6] Error> JP instruction: no flags should have changed\n")
	}
}
func test_0xDA_JP_C_a16(t *testing.T) {
	// TC1: positive case: should jump to a16 operand address if C flag is set
	preconditions()
	cpu.setCFlag()
	saveFlags1 := cpu.f

	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xDA, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0xDA, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState1 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState1.PC != 0x001F {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x001F, got @0x%04X \n", finalState1.PC)
	}

	// no flags should have changed
	if finalState1.F != saveFlags1 {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_3] Error> no flags should have changed\n")
	}

	// TC2: positive case: should not jump if C flag is reset
	preconditions()
	cpu.resetCFlag()
	saveFlags2 := cpu.f

	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xDA, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // no jump // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, //
		0x00, 0x00, 0x00, 0x00, 0xDA, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //
	}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState2 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState2.PC != 0x000F {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_5] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x000F, got @0x%04X \n", finalState2.PC)
	}

	// no flags should have changed
	if finalState2.F != saveFlags2 {
		t.Errorf("[test_0xDA_JP_C_a16_CHK_6] Error> JP instruction: no flags should have changed\n")
	}
}
func test_0xD2_JP_NC_a16(t *testing.T) {
	// TC1: positive case: should jump to a16 operand address if C flag is reset
	preconditions()
	cpu.resetCFlag()
	saveFlags1 := cpu.f

	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xD2, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // jump to 0x0020
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // STOP @0x1F
		0x00, 0x00, 0x00, 0x00, 0xD2, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // jump to 0x1A
	}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState1 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState1.PC != 0x001F {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_2] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x001F, got @0x%04X \n", finalState1.PC)
	}

	// no flags should have changed
	if finalState1.F != saveFlags1 {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_3] Error> JP instruction: no flags should have changed\n")
	}

	// TC2: positive case: should not jump to a16 operand address if C flag is set
	preconditions()
	cpu.setCFlag()
	saveFlags2 := cpu.f

	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xD2, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // no jump // STOP @0x0F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, //
		0x00, 0x00, 0x00, 0x00, 0xD2, 0x1A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //
	}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x001F
	if !cpu.stopped {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState2 := getCpuState()

	// check if the last PC = 0x001F, position of the STOP instruction in the test data program
	if finalState2.PC != 0x000F {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_5] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x000F, got @0x%04X \n", finalState2.PC)
	}

	// no flags should have changed
	if finalState2.F != saveFlags2 {
		t.Errorf("[test_0xD2_JP_NC_a16_CHK_6] Error> JP instruction: no flags should have changed\n")
	}
}
func test_JP_integration(t *testing.T) {

	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JP a16 / JP HL / JP Z, a16 / JP C, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.f

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
	if !cpu.stopped {
		t.Errorf("[test_JP_integration_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00B0 {
		t.Errorf("[test_JP_integration_CHK_2] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x00B0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_JP_integration_CHK_3] Error> JP instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JP a16 / JP HL / JP NZ, a16 / JP NC, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.f

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
	if !cpu.stopped {
		t.Errorf("[test_JP_integration_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00B0 {
		t.Errorf("[test_JP_integration_CHK_5] Error> JP instruction: the program counter should have stopped at the STOP instruction @0x00B0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_JP_integration_CHK_6] Error> JP instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}

// JR: should jump to the address specified in the instruction
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
	t.Run("0x18: JR r8", test_0x18_JR_r8)
	t.Run("0x20: JR NZ, r8", test_0x20_JR_NZ_r8)
	t.Run("0x28: JR Z, r8", test_0x28_JR_Z_r8)
	t.Run("0x30: JR NC, r8", test_0x30_JR_NC_r8)
	t.Run("0x38: JR C, r8", test_0x38_JR_C_r8)
	t.Run("JR integration", test_JR_integration)
}
func test_0x18_JR_r8(t *testing.T) {
	// TC: jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0
	// preconditions
	preconditions()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x18, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x18, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x18, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x18, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x18, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x18, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x18, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x18_JR_r8_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x18_JR_r8_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x18_JR_r8_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0x20_JR_NZ_r8(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0
	// preconditions
	preconditions()
	cpu.resetZFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x20, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x20, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x20, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x20, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x20, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x20, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x20, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// should not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.resetZFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x20, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x20, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x20, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x20, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x20, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x20, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x20, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_5] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x20_JR_NZ_r8_CHK_6] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0x28_JR_Z_r8(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0
	// preconditions
	preconditions()
	cpu.setZFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x28, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x28, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x28, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x28, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x28, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x28, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x28, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// should not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.setZFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x28, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x28, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x28, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x28, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x28, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x28, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x28, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_5] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x28_JR_Z_r8_CHK_6] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0x30_JR_NC_r8(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0
	// preconditions
	preconditions()
	cpu.resetCFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x30, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x30, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x30, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x30, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x30, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// should not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.resetCFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x30, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x30, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x30, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x30, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_5] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x30_JR_NC_r8_CHK_6] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0x38_JR_C_r8(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0
	// preconditions
	preconditions()
	cpu.setCFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x38, 0x3E, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x38, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x38, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x38, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x38, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x38_JR_C_r8_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x38_JR_C_r8_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x38_JR_C_r8_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// should not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.setCFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40 - 2 bytes for the instruction length)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x38, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8 (255 - 32 = 223 = DF)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x38, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR r8(0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x38, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR r8(0x60)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x38, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR r8(255-80 = 175 = AF)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0x38_JR_C_r8_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0x38_JR_C_r8_CHK_5] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0x38_JR_C_r8_CHK_6] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_JR_integration(t *testing.T) {
	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JR a16 / JR HL / JR Z, a16 / JR C, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x18, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x28, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR Z, r8(0x40) ; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x18, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  	; step 1	;		JR r8(255 - 32 = 223 = DF) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x28, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  	; step 3	;		JR Z, r8(0x40)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x38, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  	; step 4	;		JR C, r8(0x40) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x38, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  	; step 5	;   JR C, r8(0x60) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x38, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  	; step 6	;   JR C, r8(255-80 = 175 = AF) ; precondition: C = 0
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[JP_TC1_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[JP_TC1_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[JP_TC1_CHK_3] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions JR a16 / JR HL / JR NZ, a16 / JR NC, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x18, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		JR r8(+64 => 0x40)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0x20, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	JR NZ, r8(0x40) ; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0x18, 0xDE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		JR r8(255 - 32 = 223 = DF) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0x20, 0x1E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		JR NZ, r8(0x40)	; precondition: Z = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0x30, 0x3E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		JR NC, r8(0x40) ; precondition: C = 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0x30, 0x2E, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   JR NC, r8(0x60) ; precondition: C = 0
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0x30, 0xAE, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   JR NC, r8(255-80 = 175 = AF) ; precondition: C = 0
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[JP_TC2_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[JP_TC2_CHK_2] Error> JR instruction: the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[JP_TC2_CHK_2] Error> JR instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}

// CALL: should call the address specified in the instruction
// opcodes:
//   - 0xCD = CALL a16
//   - 0xCC = CALL Z, a16
//   - 0xC4 = CALL NZ, a16
//   - 0xDC = CALL C, a16
//   - 0xD4 = CALL NC, a16
//   - flags: -
func TestCALL(t *testing.T) {
	test_0xCD_CALL_a16(t)
	test_0xCC_CALL_Z_a16(t)
	test_0xC4_CALL_NZ_a16(t)
	test_0xDC_CALL_C_a16(t)
	test_0xD4_CALL_NC_a16(t)
	test_CALL_integration(t)
}
func test_0xCD_CALL_a16(t *testing.T) {
	// TC: jumps from 0x0000 => 0x0040 => 0x002F => 0x0060 => 0x0080 => 0x00C0 => 0x00F0 => 0x00A0 using the CALL a16 instruction
	// preconditions
	preconditions()
	saveFlags := cpu.f

	// test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCD, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL a16(0x0040)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xCD, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL a16(0x0060)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		CALL a16(0x002F)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xCD, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		CALL a16(0x0080)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xCD, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		CALL a16(0x00C0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xCD, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   CALL a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xCD, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   CALL a16(0x00A0)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xCD_CALL_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0xCD_CALL_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xCD_CALL_a16_CHK_3] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0xCC_CALL_Z_a16(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x0040 => 0x002F => 0x0060 => 0x0080 => 0x00C0 => 0x00F0 => 0x00A0 using the CALL a16 instruction
	// preconditions
	preconditions()
	cpu.setZFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCC, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL Z, a16(0x0040)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xCC, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL Z, a16(0x0060)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xCC, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		CALL Z, a16(0x002F)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xCC, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		CALL Z, a16(0x0080)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xCC, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		CALL Z, a16(0x00C0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xCC, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   CALL Z, a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xCC, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   CALL Z, a16(0x00A0)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_3] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// does not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.resetZFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xCC, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // @000X	; step 0	;		do not CALL Z, a16(0x0040) and stops @0x000F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xCC, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xCC, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xCC, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xCC, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xCC, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xCC, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x000F {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_5] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xCC_CALL_Z_a16_CHK_6] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0xC4_CALL_NZ_a16(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x0040 => 0x002F => 0x0060 => 0x0080 => 0x00C0 => 0x00F0 => 0x00A0 using the CALL a16 instruction
	// preconditions
	preconditions()
	cpu.resetZFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL Z, a16(0x0040)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xC4, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL Z, a16(0x0060)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xC4, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		CALL Z, a16(0x002F)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xC4, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		CALL Z, a16(0x0080)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xC4, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		CALL Z, a16(0x00C0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xC4, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   CALL Z, a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xC4, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   CALL Z, a16(0x00A0)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_3] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// does not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.setZFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xC4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // @000X	; step 0	;		do not CALL Z, a16(0x0040) and stops @0x000F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xC4, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xC4, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xC4, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xC4, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xC4, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xC4, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x000F {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_5] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xC4_CALL_NZ_a16_CHK_6] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0xDC_CALL_C_a16(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x0040 => 0x002F => 0x0060 => 0x0080 => 0x00C0 => 0x00F0 => 0x00A0 using the CALL a16 instruction
	// preconditions
	preconditions()
	cpu.setCFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xDC, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL Z, a16(0x0040)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xDC, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL Z, a16(0x0060)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xDC, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		CALL Z, a16(0x002F)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xDC, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		CALL Z, a16(0x0080)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xDC, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		CALL Z, a16(0x00C0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xDC, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   CALL Z, a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xDC, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   CALL Z, a16(0x00A0)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_3] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// does not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.resetCFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xDC, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // @000X	; step 0	;		do not CALL Z, a16(0x0040) and stops @0x000F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xDC, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xDC, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xDC, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xDC, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xDC, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xDC, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x000F {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_5] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xDC_CALL_C_a16_CHK_6] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_0xD4_CALL_NC_a16(t *testing.T) {
	// TC1: positive case
	// jumps from 0x0000 => 0x0040 => 0x002F => 0x0060 => 0x0080 => 0x00C0 => 0x00F0 => 0x00A0 using the CALL a16 instruction
	// preconditions
	preconditions()
	cpu.resetCFlag()
	saveFlags := cpu.f

	// test data
	testData1 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xD4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 0	;		CALL Z, a16(0x0040)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xD4, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 2	; 	CALL Z, a16(0x0060)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xD4, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X  ; step 1	;		CALL Z, a16(0x002F)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xD4, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X  ; step 3	;		CALL Z, a16(0x0080)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xD4, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X  ; step 4	;		CALL Z, a16(0x00C0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	; step 7	;   STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xD4, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX  ; step 5	;   CALL Z, a16(0x00F0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xD4, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX  ; step 6	;   CALL Z, a16(0x00A0)
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData1)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x00A0 {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_2] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_3] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()

	// TC2: negative case
	// does not jump and stop @0x000F
	// preconditions
	preconditions()
	cpu.setCFlag()
	saveFlags = cpu.f

	// test data
	testData2 := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0xD4, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, // @000X	; step 0	;		do not CALL Z, a16(0x0040) and stops @0x000F
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X
		0xD4, 0x60, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X
		0xD4, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X
		0xD4, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X
		0xD4, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X
		0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX
		0xD4, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX
		0xD4, 0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData2)

	// run the program
	cpu.Run()

	// check if the gameboy is stopped on the STOP instruction @0x00A0
	if !cpu.stopped {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_4] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x000F {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_5] Error> the program counter should have stopped at the STOP instruction @0x00A0, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[test_0xD4_CALL_NC_a16_CHK_6] Error> no flags should have changed\n")
	}

	// postconditions
	postconditions()
}
func test_CALL_integration(t *testing.T) {
	/*
	 * TC1: positive cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions CALL a16 / CALL HL / CALL Z, a16 / CALL C, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1
	saveFlags := cpu.f

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
	if !cpu.stopped {
		t.Errorf("[CALL_TC1_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x0060 {
		t.Errorf("[CALL_TC1_CHK_2] Error> CALL instruction: the program counter should have stopped at the STOP instruction @0x0060, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[CALL_TC1_CHK_3] Error> CALL instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()

	/*
	 * TC2: negative cases
	 * jumps from 0x0000 => 0x00F0 => 0x0020 => 0x00D0 => 0x0040 => 0x00B0 using the different positive conditions CALL a16 / CALL HL / CALL NZ, a16 / CALL NC, a16
	 */

	// preconditions
	preconditions()
	cpu.h = 0x00
	cpu.l = 0xD0
	cpu.f = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0
	saveFlags = cpu.f

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
	if !cpu.stopped {
		t.Errorf("[CALL_TC2_CHK_1] Error> STOP instruction should stop the gameboy\n")
	}

	// check the final state of the cpu
	finalState = getCpuState()

	// check if the last PC = 0x00B0, position of the STOP instruction in the test data program
	if finalState.PC != 0x0060 {
		t.Errorf("[CALL_TC2_CHK_2] Error> CALL instruction: the program counter should have stopped at the STOP instruction @0x0060, got @0x%04X \n", finalState.PC)
	}

	// no flags should have changed
	if finalState.F != saveFlags {
		t.Errorf("[CALL_TC2_CHK_3] Error> CALL instruction: no flags should have changed\n")
	}

	// postconditions
	postconditions()
}

// RET: should return from a subroutine : Will call 3 func and return from them in reverse order
// opcodes:
//   - 0xC9 = RET
//   - 0xC8 = RET Z
//   - 0xC0 = RET NZ
//   - 0xD8 = RET C
//   - 0xD0 = RET NC
//
// flags: - - - -
func TestRET(t *testing.T) {
	t.Run("0xC9: RET", test_0xC9_RET)
	t.Run("0xC8: RET Z", test_0xC8_RET_Z)
	t.Run("0xC0: RET NZ", test_0xC0_RET_NZ)
	t.Run("0xD8: RET C", test_0xD8_RET_C)
	t.Run("0xD0: RET NC", test_0xD0_RET_NC)
}
func test_0xC9_RET(t *testing.T) {

	// preconditions
	preconditions()

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP 				; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP 				; CALL 0x00C9 : RET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP					; RET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[test_0xC9_RET_CHK_1] Error> RET instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[test_0xC9_RET_CHK_2] Error> RET instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}
func test_0xC8_RET_Z(t *testing.T) {

	// preconditions
	preconditions()
	cpu.f = 0x80 // Z = 1 / C = 0 / H = 0 / N = 0

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP 				; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xC8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RET Z
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xC8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP 				; CALL 0x00C9 : RET Z
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xC8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RET Z
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xC8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP					; RET Z
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[test_0xC8_RET_Z_CHK_1] Error> RET instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[test_0xC8_RET_Z_CHK_2] Error> RET instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}
func test_0xC0_RET_NZ(t *testing.T) {

	/* TC3: RET NZ */

	// preconditions
	preconditions()
	cpu.f = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP 				; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RET NZ
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP 				; CALL 0x00C9 : RET NZ
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RET NZ
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP					; RET NZ
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[test_0xC0_RET_NZ_CHK_1] Error> RET instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[test_0xC0_RET_NZ_CHK_2] Error> RET instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}
func test_0xD8_RET_C(t *testing.T) {

	/* TC4: RET C */

	// preconditions
	preconditions()
	cpu.f = 0xFF // Z = 1 / C = 1 / H = 1 / N = 1

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP 				; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RET C
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP 				; CALL 0x00C9 : RET C
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RET C
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xD8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP					; RET C
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[test_0xD8_RET_C_CHK_1] Error> RET instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[test_0xD8_RET_C_CHK_2] Error> RET instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}
func test_0xD0_RET_NC(t *testing.T) {

	/* TC5: RET NC */

	// preconditions
	preconditions()
	cpu.f = 0x00 // Z = 0 / C = 0 / H = 0 / N = 0

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP 				; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RET NC
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP 				; CALL 0x00C9 : RET NC
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RET NC
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xD0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP					; RET NC
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[test_0xD0_RET_NC_CHK_1] Error> RET instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[test_0xD0_RET_NC_CHK_2] Error> RET instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}

// RETI: should return from a subroutine and enable interrupts
// 0xD9 = RETI
func TestRETI(t *testing.T) {
	// preconditions
	preconditions()
	bus.Write(0xFFFF, 0x00) // disable interrupts IME=0

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7	0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0xCD, 0x20, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @000X	; step 4	; NOP ; NOP ; CALL 0x0020 ; NOP ; NOP ; NOP ; STOP
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @001X	;
		0x00, 0xCD, 0x40, 0x00, 0xD9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @002X	; step 3	; NOP ; CALL 0x0040 ; RETI
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @003X	;
		0x00, 0x00, 0xCD, 0x60, 0x00, 0xD9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	; step 2	; NOP ; NOP ; CALL 0x00C9 : RETI
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0xCD, 0x80, 0x00, 0xD9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	; step 1	; NOP ; CALL 0x0080 ; RETI
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0xD9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	; step 0	; NOP : NOP; start program here with RETI ; precondition: SP = 0x0080 [0x60, 0x00]
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00FX	;
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x0008 {
		t.Errorf("[TestRETI_CHK_1] Error> RETI instruction: the program counter should have stopped at the STOP instruction @0x0008, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[TestRETI_CHK_2] Error> RETI instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// check if the IME = 1
	if !finalState.IME {
		t.Errorf("[TestRETI_CHK_3] Error> RETI instruction: the interrupt master enable flag should have been set\n")
	}

	// postconditions
	postconditions()
}

// RST: should call the address specified in the instruction
// opcodes:
//   - 0xC7 = RST $00
//   - 0xCF = RST $08
//   - 0xD7 = RST $10
//   - 0xDF = RST $18
//   - 0xE7 = RST $20
//   - 0xEF = RST $28
//   - 0xF7 = RST $30
//   - 0xFF = RST $38
//   - flags: - - - -
func TestRST(t *testing.T) {

	/* TC1: RST - we start @0x00F0 and execute:
	 * - @0xF0 RST $38
	 */

	// preconditions
	preconditions()
	cpu.offset = 0x00F0

	// prepare the test data
	testData := []uint8{
		//X0	0xX1	0xX2	0xX3	0xX4	0xX5	0xX6	0xX7 *0xX8	0xX9	0xXA	0xXB	0xXC	0xXD	0xXE	0xXF
		0x00, 0x00, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0xC7, 0x00, 0x00, 0x00, 0x00, 0xC9, 0x00, // @000X	; 				; RET ; RST $00 ; RET
		0x00, 0xCF, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0xD7, 0x00, 0x00, 0x00, 0x00, 0xC9, 0x00, // @001X	; RST $00 ; RET ; RST $10	; RET
		0x00, 0xDF, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0xE7, 0x00, 0x00, 0x00, 0x00, 0xC9, 0x00, // @002X	; RST $18 ; RET ; RST $20	; RET
		0x00, 0xEF, 0x00, 0xC9, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF7, 0x00, 0x00, 0x00, 0x00, 0xC9, 0x00, // @003X	; RST $28 ; RET ; RST $30	; RET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @004X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @005X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @006X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @007X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @008X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @009X	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00AX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00BX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00CX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00DX	;
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // @00EX	;
		0x00, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, // @00FX	; RST $38 ; STOP
	}

	// load the program into the memory
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the last PC = 0x0008, position of the RET instruction in the test data program
	if finalState.PC != 0x00FE {
		t.Errorf("[TestRST_CHK_1] Error> RST instruction: the program counter should have stopped at the STOP instruction @0x00FE, got @0x%04X \n", finalState.PC)
	}

	// check if the SP = 0xFFFE
	if finalState.SP != 0xFFFE {
		t.Errorf("[TestRST_CHK_2] Error> RST instruction: the stack pointer should have stopped at 0xFFFE, got @0x%04X \n", finalState.SP)
	}

	// postconditions
	postconditions()
}

// LD: should load the value from the source into the destination
// opcodes:
// > LD r8, n8
//   - 0x3E = LD A, n8
//   - 0x06 = LD B, n8
//   - 0x0E = LD C, n8
//   - 0x16 = LD D, n8
//   - 0x1E = LD E, n8
//   - 0x26 = LD H, n8
//   - 0x2E = LD L, n8
//
// > LD r16, n16
//   - 0x01 = LD BC, n16
//   - 0x11 = LD DE, n16
//   - 0x21 = LD HL, n16
//   - 0x31 = LD SP, n16
//
// > LD r8 (A/B/C/D/E/H/L), r8
//   - 0x7F = LD A, A
//   - 0x78 = LD A, B
//   - 0x79 = LD A, C
//   - 0x7A = LD A, D
//   - 0x7B = LD A, E
//   - 0x7C = LD A, H
//   - 0x7D = LD A, L
//   - 0x47 = LD B, A
//   - 0x40 = LD B, B
//   - 0x41 = LD B, C
//   - 0x42 = LD B, D
//   - 0x43 = LD B, E
//   - 0x44 = LD B, H
//   - 0x45 = LD B, L
//   - 0x4F = LD C, A
//   - 0x48 = LD C, B
//   - 0x49 = LD C, C
//   - 0x4A = LD C, D
//   - 0x4B = LD C, E
//   - 0x4C = LD C, H
//   - 0x4D = LD C, L
//   - 0x57 = LD D, A
//   - 0x50 = LD D, B
//   - 0x51 = LD D, C
//   - 0x52 = LD D, D
//   - 0x53 = LD D, E
//   - 0x54 = LD D, H
//   - 0x55 = LD D, L
//   - 0x5F = LD E, A
//   - 0x58 = LD E, B
//   - 0x59 = LD E, C
//   - 0x5A = LD E, D
//   - 0x5B = LD E, E
//   - 0x5C = LD E, H
//   - 0x5D = LD E, L
//   - 0x67 = LD H, A
//   - 0x60 = LD H, B
//   - 0x61 = LD H, C
//   - 0x62 = LD H, D
//   - 0x63 = LD H, E
//   - 0x64 = LD H, H
//   - 0x65 = LD H, L
//   - 0x6F = LD L, A
//   - 0x68 = LD L, B
//   - 0x69 = LD L, C
//   - 0x6A = LD L, D
//   - 0x6B = LD L, E
//   - 0x6C = LD L, H
//   - 0x6D = LD L, L
//
// > LD r8, [HL]
//   - 0x7E = LD A, [HL]
//   - 0x46 = LD B, [HL]
//   - 0x4E = LD C, [HL]
//   - 0x56 = LD D, [HL]
//   - 0x5E = LD E, [HL]
//   - 0x66 = LD H, [HL]
//   - 0x6E = LD L, [HL]
//
// > LD [HL], n8/r8
//   - 0x36 = LD [HL], n8
//   - 0x77 = LD [HL], A
//   - 0x70 = LD [HL], B
//   - 0x71 = LD [HL], C
//   - 0x72 = LD [HL], D
//   - 0x73 = LD [HL], E
//   - 0x74 = LD [HL], H
//   - 0x75 = LD [HL], L
//
// > LD A, from address
//   - 0xFA = LD A, [n16]
//   - 0xF2 = LD A, [C]
//   - 0x0A = LD A, [BC]
//   - 0x1A = LD A, [DE]
//   - 0x2A = LD A, [HL+]
//   - 0x3A = LD A, [HL-]
//
// > LD to address, A
//   - 0xEA = LD [n16], A
//   - 0xE2 = LD [C], A
//   - 0x02 = LD [BC], A
//   - 0x12 = LD [DE], A
//   - 0x22 = LD [HL+], A
//   - 0x32 = LD [HL-], A
//
// > LD Stack Pointer
//   - 0xF9 = LD SP, HL
//   - 0x08 = LD [n16], SP
//   - 0xF8 = LD HL, SP+r8
//
// flags: - - - - except for 0xF8 where Z:0 N:0 H:H C:C
func TestLD(t *testing.T) {
	// > LD r8, n8
	t.Run("0x3E_LD_A_n8", test_0x3E_LD_A_n8)
	t.Run("0x06_LD_B_n8", test_0x06_LD_B_n8)
	t.Run("0x0E_LD_C_n8", test_0x0E_LD_C_n8)
	t.Run("0x16_LD_D_n8", test_0x16_LD_D_n8)
	t.Run("0x1E_LD_E_n8", test_0x1E_LD_E_n8)
	t.Run("0x26_LD_H_n8", test_0x26_LD_H_n8)
	t.Run("0x2E_LD_L_n8", test_0x2E_LD_L_n8)

	// > LD r16, n16
	t.Run("0x01_LD_BC_n16", test_0x01_LD_BC_n16)
	t.Run("0x11_LD_DE_n16", test_0x11_LD_DE_n16)
	t.Run("0x21_LD_HL_n16", test_0x21_LD_HL_n16)
	t.Run("0x31_LD_SP_n16", test_0x31_LD_SP_n16)

	// > LD r8 (A/B/C/D/E/H/L), r8
	t.Run("0x7F_LD_A_A", test_0x7F_LD_A_A)
	t.Run("0x78_LD_A_B", test_0x78_LD_A_B)
	t.Run("0x79_LD_A_C", test_0x79_LD_A_C)
	t.Run("0x7A_LD_A_D", test_0x7A_LD_A_D)
	t.Run("0x7B_LD_A_E", test_0x7B_LD_A_E)
	t.Run("0x7C_LD_A_H", test_0x7C_LD_A_H)
	t.Run("0x7D_LD_A_L", test_0x7D_LD_A_L)

	t.Run("0x47_LD_B_A", test_0x47_LD_B_A)
	t.Run("0x40_LD_B_B", test_0x40_LD_B_B)
	t.Run("0x41_LD_B_C", test_0x41_LD_B_C)
	t.Run("0x42_LD_B_D", test_0x42_LD_B_D)
	t.Run("0x43_LD_B_E", test_0x43_LD_B_E)
	t.Run("0x44_LD_B_H", test_0x44_LD_B_H)
	t.Run("0x45_LD_B_L", test_0x45_LD_B_L)

	t.Run("0x4F_LD_C_A", test_0x4F_LD_C_A)
	t.Run("0x48_LD_C_B", test_0x48_LD_C_B)
	t.Run("0x49_LD_C_C", test_0x49_LD_C_C)
	t.Run("0x4A_LD_C_D", test_0x4A_LD_C_D)
	t.Run("0x4B_LD_C_E", test_0x4B_LD_C_E)
	t.Run("0x4C_LD_C_H", test_0x4C_LD_C_H)
	t.Run("0x4D_LD_C_L", test_0x4D_LD_C_L)

	t.Run("0x57_LD_D_A", test_0x57_LD_D_A)
	t.Run("0x50_LD_D_B", test_0x50_LD_D_B)
	t.Run("0x51_LD_D_C", test_0x51_LD_D_C)
	t.Run("0x52_LD_D_D", test_0x52_LD_D_D)
	t.Run("0x53_LD_D_E", test_0x53_LD_D_E)
	t.Run("0x54_LD_D_H", test_0x54_LD_D_H)
	t.Run("0x55_LD_D_L", test_0x55_LD_D_L)

	t.Run("0x5F_LD_E_A", test_0x5F_LD_E_A)
	t.Run("0x58_LD_E_B", test_0x58_LD_E_B)
	t.Run("0x59_LD_E_C", test_0x59_LD_E_C)
	t.Run("0x5A_LD_E_D", test_0x5A_LD_E_D)
	t.Run("0x5B_LD_E_E", test_0x5B_LD_E_E)
	t.Run("0x5C_LD_E_H", test_0x5C_LD_E_H)
	t.Run("0x5D_LD_E_L", test_0x5D_LD_E_L)

	t.Run("0x67_LD_H_A", test_0x67_LD_H_A)
	t.Run("0x60_LD_H_B", test_0x60_LD_H_B)
	t.Run("0x61_LD_H_C", test_0x61_LD_H_C)
	t.Run("0x62_LD_H_D", test_0x62_LD_H_D)
	t.Run("0x63_LD_H_E", test_0x63_LD_H_E)
	t.Run("0x64_LD_H_H", test_0x64_LD_H_H)
	t.Run("0x65_LD_H_L", test_0x65_LD_H_L)

	t.Run("0x6F_LD_L_A", test_0x6F_LD_L_A)
	t.Run("0x68_LD_L_B", test_0x68_LD_L_B)
	t.Run("0x69_LD_L_C", test_0x69_LD_L_C)
	t.Run("0x6A_LD_L_D", test_0x6A_LD_L_D)
	t.Run("0x6B_LD_L_E", test_0x6B_LD_L_E)
	t.Run("0x6C_LD_L_H", test_0x6C_LD_L_H)
	t.Run("0x6D_LD_L_L", test_0x6D_LD_L_L)

	// > LD r8, [HL]
	t.Run("0x7E_LD_A__HL", test_0x7E_LD_A__HL)
	t.Run("0x46_LD_B__HL", test_0x46_LD_B__HL)
	t.Run("0x4E_LD_C__HL", test_0x4E_LD_C__HL)
	t.Run("0x56_LD_D__HL", test_0x56_LD_D__HL)
	t.Run("0x5E_LD_E__HL", test_0x5E_LD_E__HL)
	t.Run("0x66_LD_H__HL", test_0x66_LD_H__HL)
	t.Run("0x6E_LD_L__HL", test_0x6E_LD_L__HL)

	// > LD [HL], n8/r8
	t.Run("0x36_LD__HL_n8", test_0x36_LD__HL_n8)
	t.Run("0x77_LD__HL_A", test_0x77_LD__HL_A)
	t.Run("0x70_LD__HL_B", test_0x70_LD__HL_B)
	t.Run("0x71_LD__HL_C", test_0x71_LD__HL_C)
	t.Run("0x72_LD__HL_D", test_0x72_LD__HL_D)
	t.Run("0x73_LD__HL_E", test_0x73_LD__HL_E)
	t.Run("0x74_LD__HL_H", test_0x74_LD__HL_H)
	t.Run("0x75_LD__HL_L", test_0x75_LD__HL_L)

	// > LD A, from address
	t.Run("0xFA_LD_A__a16", test_0xFA_LD_A__a16)
	t.Run("0xF2_LD_A__C", test_0xF2_LD_A__C)
	t.Run("0x0A_LD_A__BC", test_0x0A_LD_A__BC)
	t.Run("0x1A_LD_A__DE", test_0x1A_LD_A__DE)
	t.Run("0x2A_LD_A__HL_", test_0x2A_LD_A__HLp)
	t.Run("0x3A_LD_A__HL_", test_0x3A_LD_A__HLm)

	// > LD to address, A
	t.Run("0xEA_LD__a16_A", test_0xEA_LD__a16_A)
	t.Run("0xE2_LD__C_A", test_0xE2_LD__C_A)
	t.Run("0x02_LD__BC_A", test_0x02_LD__BC_A)
	t.Run("0x12_LD__DE_A", test_0x12_LD__DE_A)
	t.Run("0x22_LD__HLp_A", test_0x22_LD__HLp_A)
	t.Run("0x32_LD__HLm_A", test_0x32_LD__HLm_A)

	// > LD Stack Pointer
	t.Run("0xF9_LD_SP_HL", test_0xF9_LD_SP_HL)
	t.Run("0x08_LD__a16_SP", test_0x08_LD__a16_SP)
	t.Run("0xF8_LD_HL_SP_e8", test_0xF8_LD_HL_SP_e8)
}

// > LD r8, n8
var testData_LD_r8_n8 = []uint8{0x00, 0xFF, 0x0F, 0xF0, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}

func test_0x3E_LD_A_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = 0x77
		testProgram := []uint8{0x3E, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x3E_LD_A_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x3E_LD_A_n8] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x3E_LD_A_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x06_LD_B_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = 0x77
		testProgram := []uint8{0x06, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x06_LD_B_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x06_LD_B_n8] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x06_LD_B_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x0E_LD_C_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = 0x77
		testProgram := []uint8{0x0E, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x0E_LD_C_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x0E_LD_C_n8] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x0E_LD_C_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x16_LD_D_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = 0x77
		testProgram := []uint8{0x16, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x16_LD_D_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x16_LD_D_n8] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x16_LD_D_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x1E_LD_E_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = 0x77
		testProgram := []uint8{0x1E, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x1E_LD_E_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x1E_LD_E_n8] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x1E_LD_E_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x26_LD_H_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = 0x77
		testProgram := []uint8{0x26, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x26_LD_H_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x26_LD_H_n8] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x26_LD_H_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x2E_LD_L_n8(t *testing.T) {
	for idx, data := range testData_LD_r8_n8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = 0x77
		testProgram := []uint8{0x2E, data, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x2E_LD_L_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x2E_LD_L_n8] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x2E_LD_L_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

var testData_LD_r16_n16 = []uint16{0x0000, 0xFFFF, 0x00FF, 0xFF00, 0x00AA, 0x00BB, 0x00CC, 0x00DD, 0x00EE, 0x00FF, 0x1234, 0x3456, 0x5678, 0x789A, 0x9ABC, 0xBCDE, 0xDEFF}

// > LD r16, n16
func test_0x01_LD_BC_n16(t *testing.T) {
	for idx, data := range testData_LD_r16_n16 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setBC(0x7777)
		testProgram := []uint8{0x01, uint8(data & 0x00FF), uint8((data & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		bc := cpu.getBC()
		if bc != data {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected register BC to be 0x%04X, got 0x%04X\n", idx, data, bc)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x11_LD_DE_n16(t *testing.T) {
	for idx, data := range testData_LD_r16_n16 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setDE(0x7777)
		testProgram := []uint8{0x11, uint8(data & 0x00FF), uint8((data & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0x11_LD_DE_n16] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		de := cpu.getDE()
		if de != data {
			t.Errorf("[test_0x11_LD_DE_n16] %v> expected register DE to be 0x%04X, got 0x%04X\n", idx, data, de)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x11_LD_DE_n16] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x21_LD_HL_n16(t *testing.T) {
	for idx, data := range testData_LD_r16_n16 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x7777)
		testProgram := []uint8{0x21, uint8(data & 0x00FF), uint8((data & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0x21_LD_HL_n16] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		hl := cpu.getHL()
		if hl != data {
			t.Errorf("[test_0x21_LD_HL_n16] %v> expected register HL to be 0x%04X, got 0x%04X\n", idx, data, hl)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x21_LD_HL_n16] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x31_LD_SP_n16(t *testing.T) {
	for idx, data := range testData_LD_r16_n16 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.sp = 0x7777
		testProgram := []uint8{0x31, uint8(data & 0x00FF), uint8((data & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		sp := cpu.sp
		if sp != data {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected register SP to be 0x%04X, got 0x%04X\n", idx, data, sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x01_LD_BC_n16] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

var testData_LD_r8_r8 = []uint8{0x00, 0xFF, 0x0F, 0xF0, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}

// > LD r8 (A/B/C/D/E/H/L), r8
func test_0x7F_LD_A_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x7F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7F_LD_A_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7F_LD_A_A] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x7F_LD_A_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7F_LD_A_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x78_LD_A_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x78, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x78_LD_A_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x78_LD_A_B] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x78_LD_A_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x78_LD_A_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x79_LD_A_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x79, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x79_LD_A_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x79_LD_A_C] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x79_LD_A_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x79_LD_A_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x7A_LD_A_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x7A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7A_LD_A_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7A_LD_A_D] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x7A_LD_A_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7A_LD_A_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x7B_LD_A_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x7B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7B_LD_A_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7B_LD_A_E] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x7B_LD_A_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7B_LD_A_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x7C_LD_A_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x7C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7C_LD_A_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7C_LD_A_H] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x7C_LD_A_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7C_LD_A_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x7D_LD_A_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x7D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7D_LD_A_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7D_LD_A_L] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x7D_LD_A_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7D_LD_A_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x47_LD_B_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x47, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x47_LD_B_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x47_LD_B_A] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x47_LD_B_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x47_LD_B_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x40_LD_B_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x47, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x40_LD_B_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x40_LD_B_B] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x40_LD_B_B] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x40_LD_B_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x41_LD_B_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x41, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x41_LD_B_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x41_LD_B_C] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x41_LD_B_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x41_LD_B_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x42_LD_B_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x42, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x42_LD_B_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x42_LD_B_D] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x42_LD_B_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x42_LD_B_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x43_LD_B_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x43, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x43_LD_B_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x43_LD_B_E] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x43_LD_B_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x43_LD_B_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x44_LD_B_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x44, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x44_LD_B_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x44_LD_B_H] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x44_LD_B_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x44_LD_B_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x45_LD_B_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x45, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x45_LD_B_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x45_LD_B_L] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x45_LD_B_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x45_LD_B_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x4F_LD_C_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x4F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4F_LD_C_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4F_LD_C_A] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x4F_LD_C_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4F_LD_C_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x48_LD_C_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x48, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x48_LD_C_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x48_LD_C_B] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x48_LD_C_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x48_LD_C_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x49_LD_C_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x49, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x49_LD_C_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x49_LD_C_C] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x49_LD_C_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x49_LD_C_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x4A_LD_C_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x4A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4A_LD_C_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4A_LD_C_D] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x4A_LD_C_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4A_LD_C_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x4B_LD_C_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x4B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4B_LD_C_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4B_LD_C_E] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x4B_LD_C_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4B_LD_C_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x4C_LD_C_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x4C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4C_LD_C_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4C_LD_C_H] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x4C_LD_C_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4C_LD_C_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x4D_LD_C_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x4D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4D_LD_C_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4D_LD_C_L] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x4D_LD_C_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4D_LD_C_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x57_LD_D_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x57, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x57_LD_D_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x57_LD_D_A] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x57_LD_D_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x57_LD_D_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x50_LD_D_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x50, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x50_LD_D_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x50_LD_D_B] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x50_LD_D_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x50_LD_D_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x51_LD_D_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x51, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x51_LD_D_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x51_LD_D_C] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x51_LD_D_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x51_LD_D_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x52_LD_D_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x52, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x52_LD_D_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x52_LD_D_D] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x52_LD_D_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x52_LD_D_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x53_LD_D_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x53, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x53_LD_D_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x53_LD_D_E] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x53_LD_D_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x53_LD_D_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x54_LD_D_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x54, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x54_LD_D_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x54_LD_D_H] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x54_LD_D_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x54_LD_D_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x55_LD_D_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x55, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x55_LD_D_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x55_LD_D_L] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x55_LD_D_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x55_LD_D_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x5F_LD_E_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x5F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5F_LD_E_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5F_LD_E_A] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x5F_LD_E_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5F_LD_E_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x58_LD_E_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x58, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x58_LD_E_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x58_LD_E_B] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x58_LD_E_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x58_LD_E_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x59_LD_E_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x59, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x59_LD_E_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x59_LD_E_C] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x59_LD_E_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x59_LD_E_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x5A_LD_E_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x5A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5A_LD_E_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5A_LD_E_D] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x5A_LD_E_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5A_LD_E_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x5B_LD_E_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x5B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5B_LD_E_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5B_LD_E_E] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x5B_LD_E_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5B_LD_E_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x5C_LD_E_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x5C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5C_LD_E_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5C_LD_E_H] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x5C_LD_E_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5C_LD_E_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x5D_LD_E_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x5D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5D_LD_E_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5D_LD_E_L] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x5D_LD_E_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5D_LD_E_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x67_LD_H_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x67, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x67_LD_H_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x67_LD_H_A] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x67_LD_H_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x67_LD_H_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x60_LD_H_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x60, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x60_LD_H_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x60_LD_H_B] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x60_LD_H_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x60_LD_H_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x61_LD_H_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x61, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x61_LD_H_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x61_LD_H_C] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x61_LD_H_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x61_LD_H_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x62_LD_H_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x62, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x62_LD_H_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x62_LD_H_D] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x62_LD_H_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x62_LD_H_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x63_LD_H_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x63, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x63_LD_H_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x63_LD_H_E] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x63_LD_H_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x63_LD_H_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x64_LD_H_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x64, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x64_LD_H_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x64_LD_H_H] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x64_LD_H_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x64_LD_H_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x65_LD_H_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x65, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x65_LD_H_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x65_LD_H_L] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x65_LD_H_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x65_LD_H_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

func test_0x6F_LD_L_A(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = data
		testProgram := []uint8{0x6F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6F_LD_L_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6F_LD_L_A] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.a != data {
			t.Errorf("[test_0x6F_LD_L_A] %v> expected source register A to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6F_LD_L_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x68_LD_L_B(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = data
		testProgram := []uint8{0x68, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x68_LD_L_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x68_LD_L_B] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.b != data {
			t.Errorf("[test_0x68_LD_L_B] %v> expected source register B to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x68_LD_L_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x69_LD_L_C(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = data
		testProgram := []uint8{0x69, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x69_LD_L_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x69_LD_L_C] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.c != data {
			t.Errorf("[test_0x69_LD_L_C] %v> expected source register C to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x69_LD_L_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x6A_LD_L_D(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = data
		testProgram := []uint8{0x6A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6A_LD_L_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6A_LD_L_D] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.d != data {
			t.Errorf("[test_0x6A_LD_L_D] %v> expected source register D to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6A_LD_L_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x6B_LD_L_E(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = data
		testProgram := []uint8{0x6B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6B_LD_L_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6B_LD_L_E] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.e != data {
			t.Errorf("[test_0x6B_LD_L_E] %v> expected source register E to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6B_LD_L_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x6C_LD_L_H(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = data
		testProgram := []uint8{0x6C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6C_LD_L_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6C_LD_L_H] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.h != data {
			t.Errorf("[test_0x6C_LD_L_H] %v> expected source register H to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6C_LD_L_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x6D_LD_L_L(t *testing.T) {
	for idx, data := range testData_LD_r8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = data
		testProgram := []uint8{0x6D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6D_LD_L_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6D_LD_L_L] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check that source register is unaffected
		if cpu.l != data {
			t.Errorf("[test_0x6D_LD_L_L] %v> expected source register L to be unaffected 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6D_LD_L_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

// > LD r8, [HL]
var testData_LD_r8__HL = []uint8{0x00, 0xFF, 0x0F, 0xF0, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}

func test_0x7E_LD_A__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x7E, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x7E_LD_A__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.a != data {
			t.Errorf("[test_0x7E_LD_A__HL] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, data, cpu.a)
		}
		// check that source register is unaffected
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x7E_LD_A__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x7E_LD_A__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x46_LD_B__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x46, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x46_LD_B__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.b != data {
			t.Errorf("[test_0x46_LD_B__HL] %v> expected register B to be 0x%02X, got 0x%02X\n", idx, data, cpu.b)
		}
		// check that source register is unaffected
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x46_LD_B__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x46_LD_B__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x4E_LD_C__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x4E, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x4E_LD_C__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.c != data {
			t.Errorf("[test_0x4E_LD_C__HL] %v> expected register C to be 0x%02X, got 0x%02X\n", idx, data, cpu.c)
		}
		// check that source register is unaffected
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x4E_LD_C__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x4E_LD_C__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x56_LD_D__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x56, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x56_LD_D__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.d != data {
			t.Errorf("[test_0x56_LD_D__HL] %v> expected register D to be 0x%02X, got 0x%02X\n", idx, data, cpu.d)
		}
		// check that source register is unaffected
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x56_LD_D__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x56_LD_D__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x5E_LD_E__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x5E, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x5E_LD_E__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.e != data {
			t.Errorf("[test_0x5E_LD_E__HL] %v> expected register E to be 0x%02X, got 0x%02X\n", idx, data, cpu.e)
		}
		// check that source register is unaffected
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x5E_LD_E__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x5E_LD_E__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x66_LD_H__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x66, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x66_LD_H__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.h != data {
			t.Errorf("[test_0x66_LD_H__HL] %v> expected register H to be 0x%02X, got 0x%02X\n", idx, data, cpu.h)
		}
		/* This test is not relevant here since we are loading data into the H register which was used as HL to locate the source data
		 		// check that source register is unaffected
				valueAtHL := cpu.bus.Read(cpu.getHL())
				if valueAtHL != data {
					t.Errorf("[test_0x66_LD_H__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
				}
		*/
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x66_LD_H__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x6E_LD_L__HL(t *testing.T) {
	for idx, data := range testData_LD_r8__HL {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		testProgram := []uint8{0x6E, 0x10, data}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x6E_LD_L__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register
		if cpu.l != data {
			t.Errorf("[test_0x6E_LD_L__HL] %v> expected register L to be 0x%02X, got 0x%02X\n", idx, data, cpu.l)
		}
		/* This test is not relevant here since we are loading data into the L register which was used as HL to locate the source data
		 		// check that source register is unaffected
				valueAtHL := cpu.bus.Read(cpu.getHL())
				if valueAtHL != data {
					t.Errorf("[test_0x66_LD_H__HL] %v> expected source value at memory location HL to be unaffected 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
				}
		*/
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x6E_LD_L__HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

// > LD [HL], n8/r8
var testData_LD__HL_n8_r8 = []uint8{0x00, 0xFF, 0x0F, 0xF0, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE}

func test_0x36_LD__HL_n8(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0x36, data, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x36_LD__HL_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x36_LD__HL_n8] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0003 {
			t.Errorf("[test_0x36_LD__HL_n8] %v> expected address in register HL to be unaffected 0x%02X, got 0x%02X\n", idx, 0x0003, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x36_LD__HL_n8] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x77_LD__HL_A(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.a = data
		testProgram := []uint8{0x77, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x77_LD__HL_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x77_LD__HL_A] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x77_LD__HL_A] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x77_LD__HL_A] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x70_LD__HL_B(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.b = data
		testProgram := []uint8{0x70, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x70_LD__HL_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x70_LD__HL_B] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x70_LD__HL_B] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x70_LD__HL_B] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x71_LD__HL_C(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.c = data
		testProgram := []uint8{0x71, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x71_LD__HL_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x71_LD__HL_C] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x71_LD__HL_C] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x71_LD__HL_C] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x72_LD__HL_D(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.d = data
		testProgram := []uint8{0x72, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x72_LD__HL_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x72_LD__HL_D] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x72_LD__HL_D] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x72_LD__HL_D] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x73_LD__HL_E(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.e = data
		testProgram := []uint8{0x73, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x73_LD__HL_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x73_LD__HL_E] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x73_LD__HL_E] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x73_LD__HL_E] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x74_LD__HL_H(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.h = data
		testProgram := []uint8{0x74, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x74_LD__HL_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x74_LD__HL_H] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		/* This test is not relevant here since we are loading data into the H register which was used as HL to locate the source data
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x74_LD__HL_H] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		*/
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x74_LD__HL_H] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x75_LD__HL_L(t *testing.T) {
	for idx, data := range testData_LD__HL_n8_r8 {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0002)
		cpu.l = data
		testProgram := []uint8{0x75, 0x10, 0x77}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x75_LD__HL_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check that data has been correctly loaded into memory location pointed by HL
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data {
			t.Errorf("[test_0x75_LD__HL_L] %v> expected memory location pointed by HL to be 0x%02X, got 0x%02X\n", idx, data, valueAtHL)
		}
		/* This test is not relevant here since we are loading data into the L register which was used as HL to locate the source data
		// check that register HL is unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0x75_LD__HL_L] %v> expected address in register HL to be unaffected 0x0002, got 0x%02X\n", idx, cpu.getHL())
		}
		*/
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x75_LD__HL_L] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

// > LD A, from address
var fromToAddress = []uint16{0x0010, 0x002F, 0x0035, 0x004E, 0x1F5F, 0x3F6A, 0x0273, 0xFFFF}
var value8bitAtAddress = []uint8{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xFF}

func test_0xFA_LD_A__a16(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.bus.Write(addr, value)
		testProgram := []uint8{0xFA, uint8(addr & 0x00FF), uint8((addr & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected PC to be 0x0003, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
	}
}
func test_0xF2_LD_A__C(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(addr & 0x00FF)
		value := value8bitAtAddress[idx]
		cpu.bus.Write(0xFF00+uint16(cpu.c), value)
		testProgram := []uint8{0xF2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xF2_LD_A__C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0xF2_LD_A__C] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xF2_LD_A__C] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x0A_LD_A__BC(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setBC(addr)
		value := value8bitAtAddress[idx]
		cpu.bus.Write(cpu.getBC(), value)
		testProgram := []uint8{0x0A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x0A_LD_A__BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0x0A_LD_A__BC] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x0A_LD_A__BC] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x1A_LD_A__DE(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setDE(addr)
		value := value8bitAtAddress[idx]
		cpu.bus.Write(cpu.getDE(), value)
		testProgram := []uint8{0x1A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x1A_LD_A__DE] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0x1A_LD_A__DE] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x1A_LD_A__DE] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x2A_LD_A__HLp(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(addr)
		value := value8bitAtAddress[idx]
		cpu.bus.Write(cpu.getHL(), value)
		testProgram := []uint8{0x2A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x2A_LD_A__HLp] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0x2A_LD_A__HLp] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check that register HL has been incremented
		if cpu.getHL() != addr+1 {
			t.Errorf("[test_0x2A_LD_A__HLp] %v> expected register HL to be incremented to 0x%04X, got 0x%04X\n", idx, addr+1, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x2A_LD_A__HLp] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x3A_LD_A__HLm(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(addr)
		value := value8bitAtAddress[idx]
		cpu.bus.Write(cpu.getHL(), value)
		testProgram := []uint8{0x3A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x3A_LD_A__HLm] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into register A
		if cpu.a != value {
			t.Errorf("[test_0x3A_LD_A__HLm] %v> expected register A to be 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check that register HL has been decremented
		if cpu.getHL() != addr-1 {
			t.Errorf("[test_0x3A_LD_A__HLm] %v> expected register HL to be decremented to 0x%04X, got 0x%04X\n", idx, addr-1, cpu.getHL())
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x3A_LD_A__HLm] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}

// > LD to address, A
func test_0xEA_LD__a16_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		testProgram := []uint8{0xEA, uint8(addr & 0x00FF), uint8((addr & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected PC to be 0x0003, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected memory location pointed by n16 operand to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xFA_LD_A__a16] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0xE2_LD__C_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		cpu.c = uint8(addr & 0x00FF)
		testProgram := []uint8{0xE2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xE2_LD__C_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(0xFF00 + (addr & 0x00FF))
		if valueAtAddr != value {
			t.Errorf("[test_0xE2_LD__C_A] %v> expected memory location pointed by 0xFF00 + C register to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0xE2_LD__C_A] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xE2_LD__C_A] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x02_LD__BC_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		cpu.setBC(addr)
		testProgram := []uint8{0x02, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x02_LD__BC_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0x02_LD__BC_A] %v> expected memory location pointed by LD register to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0x02_LD__BC_A] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x02_LD__BC_A] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x12_LD__DE_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		cpu.setDE(addr)
		testProgram := []uint8{0x12, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x12_LD__DE_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0x12_LD__DE_A] %v> expected memory location pointed by DE register to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0x12_LD__DE_A] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x12_LD__DE_A] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x22_LD__HLp_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		cpu.setHL(addr)
		testProgram := []uint8{0x22, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x22_LD__HLp_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0x22_LD__HLp_A] %v> expected memory location pointed by HL register to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check that register HL has been incremented
		if cpu.getHL() != addr+1 {
			t.Errorf("[test_0x22_LD__HLp_A] %v> expected register HL to be incremented to 0x%04X, got 0x%04X\n", idx, addr+1, cpu.getHL())
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0x22_LD__HLp_A] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x22_LD__HLp_A] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x32_LD__HLm_A(t *testing.T) {
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value8bitAtAddress[idx]
		cpu.a = value
		cpu.setHL(addr)
		testProgram := []uint8{0x32, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x32_LD__HLm_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0x32_LD__HLm_A] %v> expected memory location pointed by HL register to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check that register HL has been decremented
		if cpu.getHL() != addr-1 {
			t.Errorf("[test_0x32_LD__HLm_A] %v> expected register HL to be decremented to 0x%04X, got 0x%04X\n", idx, addr-1, cpu.getHL())
		}
		// check if A register is unaffected
		if cpu.a != value {
			t.Errorf("[test_0x32_LD__HLm_A] %v> expected register A to be unaffected 0x%02X, got 0x%02X\n", idx, value, cpu.a)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x32_LD__HLm_A] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}

// > LD Stack Pointer
var value16bitAtAddress = []uint16{0x0123, 0x4567, 0x89AB, 0xCDEF, 0xFFAA, 0x19A7, 0xD65C, 0x71B9}

func test_0xF9_LD_SP_HL(t *testing.T) {
	for idx, value := range value16bitAtAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(value)
		testProgram := []uint8{0xF9, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack pointer
		if cpu.sp != value {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected stack pointer to be 0x%02X, got 0x%02X\n", idx, value, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0x08_LD__a16_SP(t *testing.T) {
	// Stores SP & $FF at address a16 and SP >> 8 at address a16 + 1
	for idx, addr := range fromToAddress {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		value := value16bitAtAddress[idx]
		cpu.sp = value
		testProgram := []uint8{0x08, uint8(addr & 0x00FF), uint8((addr & 0xFF00) >> 8), 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0003 {
			t.Errorf("[test_0x08_LD__a16_SP] %v> expected PC to be 0x0003, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into memory location
		valueAtAddr := cpu.bus.Read16(addr)
		if valueAtAddr != value {
			t.Errorf("[test_0x08_LD__a16_SP] %v> expected memory location pointed by a16 operand to be 0x%02X, got 0x%02X\n", idx, value, valueAtAddr)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0x08_LD__a16_SP] %v> expected flags to be unaffected 0x00, got 0x%02X\n", idx, cpu.f)
		}
		postconditions()
	}
}
func test_0xF8_LD_HL_SP_e8(t *testing.T) {
	// this is the only LD instruction that affects the flags
	// this is also the only LD instruction that uses an 8-bit signed operand
	type TestDataEntry struct {
		sp            uint16
		e8            uint8
		expectedHL    uint16
		expectedHFlag bool
		expectedCFlag bool
	}
	var testData = []TestDataEntry{
		{0xFFF8, 0x10, 0x0008, false, true},  // 0> 0xFFF8 + 0x10 (+16) = 0x0008 - H = 0, C = 1
		{0xFFF8, 0xF0, 0xFFE8, false, true},  // 1> 0xFFF8 + 0xF0 (-16) = 0xFFE8 - H = 1, C = 0
		{0x0001, 0xFF, 0x0000, true, true},   // 2> 0x0001 + 0xFF (-1)  = 0x0000 - H = 1, C = 1
		{0x00FF, 0x01, 0x0100, true, true},   // 3> 0x00FF + 0x01 (+1)  = 0x0100 - H = 0, C = 1
		{0x0100, 0xFE, 0x00FE, false, false}, // 4> 0x0100 + 0xFE (-2)  = 0x00FE - H = 1, C = 1
		{0x7FFF, 0x01, 0x8000, true, true},   // 5> 0x7FFF + 0x01 (+1)  = 0x8000 - H = 0, C = 0
		{0x8000, 0xFF, 0x7FFF, false, false}, // 6> 0x8000 + 0xFF (-1)  = 0x7FFF - H = 1, C = 0
		{0x1234, 0x20, 0x1254, false, false}, // 7> 0x1234 + 0x20 (+32) = 0x1254 - H = 0, C = 0
		{0x00F0, 0x10, 0x0100, false, true},  // 8> 0x00F0 + 0x10 (+16) = 0x0100 - H = 1, C = 1
		{0xABCD, 0x30, 0xABFD, false, false}, // 9> 0xABCD + 0x30 (+48) = 0xABFD - H = 1, C = 0
	}

	for idx, entry := range testData {
		preconditions()
		randomizeFlags()
		cpu.sp = entry.sp
		e8 := entry.e8
		testProgram := []uint8{0xF8, e8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack pointer
		if cpu.getHL() != entry.expectedHL {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected HL register to be 0x%02X, got 0x%02X\n", idx, entry.expectedHL, cpu.getHL())
		}
		// check flags
		if cpu.getZFlag() {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected Z flag to be 0, got 1\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected N flag to be 0, got 1\n", idx)
		}
		if cpu.getHFlag() != entry.expectedHFlag {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected H flag to be %v, got %v\n", idx, entry.expectedHFlag, cpu.getHFlag())
		}
		if cpu.getCFlag() != entry.expectedCFlag {
			t.Errorf("[test_0xF9_LD_SP_HL] %v> expected C flag to be %v, got %v\n", idx, entry.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}

// LDH: should load the value from the source into the destination and not impact the flags
func TestLDH(t *testing.T) {
	t.Run("0xE0: LDH__a8_A", test_0xE0_LDH__a8_A)
	t.Run("0xF0: LDH_A__a8", test_0xF0_LDH_A__a8)
}
func test_0xF0_LDH_A__a8(t *testing.T) {
	preconditions()

	// print A initial value
	//getCpuState().print()

	// set flags to some arbitrary values to check if they are not impacted by the instruction
	cpu.f = 0xE5

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0xF0, 0x77, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	// the instruction LDH will look @0xFF77 for the value to load into A. Let's set this value to 0xB5
	bus.Write(0xFF77, 0xB5)

	//printMemoryProperties()

	// run the program and control step by step the IME flag
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	//finalState.print()

	// program should have stopped at 0x0005
	if finalState.PC != 0x0005 {
		t.Errorf("[0xF0_LDH_A__a8_TC13_CHK_0] Error> LDH A, (a8) instruction: the program counter should have stopped at 0x0005, got 0x%04X \n", finalState.PC)
	}

	// A should be 0xB5
	if finalState.A != 0xB5 {
		t.Errorf("[0xF0_LDH_A__a8_TC13_CHK_1] Error> LDH A, (a8) instruction: the A register should have been set to 0xB5, got 0x%02X \n", finalState.A)
	}

	// check if the flags are not impacted
	if finalState.F != 0xE5 {
		t.Errorf("[0xF0_LDH_A__a8_TC13_CHK_2] Error> LDH A, (a8) instruction: the flags should not have been impacted, got 0x%02X \n", finalState.F)
	}

	postconditions()
}
func test_0xE0_LDH__a8_A(t *testing.T) {
	preconditions()

	// set flags to some arbitrary values to check if they are not impacted by the instruction
	cpu.f = 0xE5

	// set the value of A to 0xB5
	cpu.a = 0xB5

	// print A initial value
	//getCpuState().print()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0xE0, 0x77, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData)

	//printMemoryProperties()

	// run the program and control step by step the IME flag
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	//finalState.print()

	// program should have stopped at 0x0005
	if finalState.PC != 0x0005 {
		t.Errorf("[0xF0_LDH_A__a8_TC13_CHK_0] Error> LDH (a8), A instruction: the program counter should have stopped at 0x0005, got 0x%04X \n", finalState.PC)
	}

	// [FF77] should be 0xB5
	inMemoryValue := bus.Read(0xFF77)
	if inMemoryValue != 0xB5 {
		t.Errorf("[0xF0_LDH__a8_A_TC13_CHK_1] Error> LDH (a8), A instruction: the memory location @0x77 should have been set to 0xB5, got 0x%02X \n", inMemoryValue)
	}

	// check if the flags are not impacted
	if finalState.F != 0xE5 {
		t.Errorf("[0xF0_LDH_A__a8_TC13_CHK_2] Error> LDH (a8), A instruction: the flags should not have been impacted, got 0x%02X \n", finalState.F)
	}

	postconditions()
}

// PUSH: Push a 16-bit register pair onto the stack
// opcodes:
//   - 0xC5 = PUSH BC
//   - 0xD5 = PUSH DE
//   - 0xE5 = PUSH HL
//   - 0xF5 = PUSH AF
//
// flags: -
func TestPUSH(t *testing.T) {
	t.Run("0xC5_PUSH_BC", test_0xC5_PUSH_BC)
	t.Run("0xD5_PUSH_DE", test_0xD5_PUSH_DE)
	t.Run("0xE5_PUSH_HL", test_0xE5_PUSH_HL)
	t.Run("0xF5_PUSH_AF", test_0xF5_PUSH_AF)
}

var pushTestData = []uint16{0x1234, 0x5678, 0x9ABC, 0xDEF0, 0xFFAA, 0x19A7, 0xD65C, 0x71B9, 0x0000, 0xFFFF}

// pushing all the test data values one after the other into the stack without resetting the gameboy
// and performing checks after each push
func test_0xC5_PUSH_BC(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	for idx, value := range pushTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.setBC(value)
		testProgram := []uint8{0xC5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xC5_PUSH_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack
		valueAtSP := cpu.bus.Read16(cpu.sp)
		if valueAtSP != value {
			t.Errorf("[test_0xC5_PUSH_BC] %v> expected value saved @SP to be 0x%04X, got 0x%04X\n", idx, value, valueAtSP)
		}
		// check if stack pointer has been decremented
		expectedSP := 0xFFFE - uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xC5_PUSH_BC] %v> expected stack pointer to be decremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xC5_PUSH_BC] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xD5_PUSH_DE(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	for idx, value := range pushTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.setDE(value)
		testProgram := []uint8{0xD5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xD5_PUSH_DE] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack
		valueAtSP := cpu.bus.Read16(cpu.sp)
		if valueAtSP != value {
			t.Errorf("[test_0xD5_PUSH_DE] %v> expected value saved @SP to be 0x%04X, got 0x%04X\n", idx, value, valueAtSP)
		}
		// check if stack pointer has been decremented
		expectedSP := 0xFFFE - uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xD5_PUSH_DE] %v> expected stack pointer to be decremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xD5_PUSH_DE] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xE5_PUSH_HL(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	for idx, value := range pushTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.setHL(value)
		testProgram := []uint8{0xE5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xE5_PUSH_HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack
		valueAtSP := cpu.bus.Read16(cpu.sp)
		if valueAtSP != value {
			t.Errorf("[test_0xE5_PUSH_HL] %v> expected value saved @SP to be 0x%04X, got 0x%04X\n", idx, value, valueAtSP)
		}
		// check if stack pointer has been decremented
		expectedSP := 0xFFFE - uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xE5_PUSH_HL] %v> expected stack pointer to be decremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xE5_PUSH_HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xF5_PUSH_AF(t *testing.T) {
	preconditions()
	randomizeFlags()

	for idx, value := range pushTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.a = uint8(value >> 8)
		cpu.f = uint8(value)
		testProgram := []uint8{0xF5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xF5_PUSH_AF] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into stack
		valueAtSP := cpu.bus.Read16(cpu.sp)
		if valueAtSP != value {
			t.Errorf("[test_0xF5_PUSH_AF] %v> expected value saved @SP to be 0x%04X, got 0x%04X\n", idx, value, valueAtSP)
		}
		// check if stack pointer has been decremented
		expectedSP := 0xFFFE - uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xF5_PUSH_AF] %v> expected stack pointer to be decremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
	}
	postconditions()
}

// POP: Pop a 16-bit register pair from the stack
// opcodes:
//   - 0xC1 = POP BC
//   - 0xD1 = POP DE
//   - 0xE1 = POP HL
//   - 0xF1 = POP AF (flags are restored from the stack)
//
// flags: - except for 0xF1 where Z->Z N->N H->H C->C
func TestPOP(t *testing.T) {
	t.Run("0xC1_POP_BC", test_0xC1_POP_BC)
	t.Run("0xD1_POP_DE", test_0xD1_POP_DE)
	t.Run("0xE1_POP_HL", test_0xE1_POP_HL)
	t.Run("0xF1_POP_AF", test_0xF1_POP_AF)
}

var popTestData = []uint16{0x1234, 0x5678, 0x9ABC, 0xDEF0, 0xFFAA, 0x19A7, 0xD65C, 0x71B9, 0x0000, 0xFFFF}

func test_0xC1_POP_BC(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	// to isolate this test from the PUSH BC test, we will push the values into the HRAM and update the stack pointer manually
	// we will push them in reverse order to POP them in the correct order
	for idx, _ := range popTestData {
		value := popTestData[len(popTestData)-1-idx]
		// extract the high and low bytes
		high := uint8(value >> 8)
		low := uint8(value)
		// push the value into the HRAM
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, high)
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, low)
	}

	// execute the test
	testProgram := []uint8{0xC1, 0x10}
	loadProgramIntoMemory(memory1, testProgram)

	for idx, value := range popTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xC1_POP_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into BC register
		if cpu.getBC() != value {
			t.Errorf("[test_0xC1_POP_BC] %v> expected BC register to be 0x%04X, got 0x%04X\n", idx, value, cpu.getBC())
		}
		// check if stack pointer has been incremented
		expectedSP := 0xFFFE - uint16(len(popTestData)*2) + uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xC1_POP_BC] %v> expected stack pointer to be incremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xC1_POP_BC] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xD1_POP_DE(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	// to isolate this test from the PUSH BC test, we will push the values into the HRAM and update the stack pointer manually
	// we will push them in reverse order to POP them in the correct order
	for idx, _ := range popTestData {
		value := popTestData[len(popTestData)-1-idx]
		// extract the high and low bytes
		high := uint8(value >> 8)
		low := uint8(value)
		// push the value into the HRAM
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, high)
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, low)
	}

	// execute the test
	testProgram := []uint8{0xD1, 0x10}
	loadProgramIntoMemory(memory1, testProgram)

	for idx, value := range popTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xD1_POP_DE] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into DE register
		if cpu.getDE() != value {
			t.Errorf("[test_0xD1_POP_DE] %v> expected DE register to be 0x%04X, got 0x%04X\n", idx, value, cpu.getDE())
		}
		// check if stack pointer has been incremented
		expectedSP := 0xFFFE - uint16(len(popTestData)*2) + uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xD1_POP_DE] %v> expected stack pointer to be incremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xD1_POP_DE] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xE1_POP_HL(t *testing.T) {
	preconditions()
	randomizeFlags()
	saveFlags := cpu.f

	// to isolate this test from the PUSH BC test, we will push the values into the HRAM and update the stack pointer manually
	// we will push them in reverse order to POP them in the correct order
	for idx, _ := range popTestData {
		value := popTestData[len(popTestData)-1-idx]
		// extract the high and low bytes
		high := uint8(value >> 8)
		low := uint8(value)
		// push the value into the HRAM
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, high)
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, low)
	}

	// execute the test
	testProgram := []uint8{0xE1, 0x10}
	loadProgramIntoMemory(memory1, testProgram)

	for idx, value := range popTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xE1_POP_HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into HL register
		if cpu.getHL() != value {
			t.Errorf("[test_0xE1_POP_HL] %v> expected HL register to be 0x%04X, got 0x%04X\n", idx, value, cpu.getHL())
		}
		// check if stack pointer has been incremented
		expectedSP := 0xFFFE - uint16(len(popTestData)*2) + uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xE1_POP_HL] %v> expected stack pointer to be incremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
		// check if flags are unaffected
		if cpu.f != saveFlags {
			t.Errorf("[test_0xE1_POP_HL] %v> expected flags to be unaffected 0x%02X, got 0x%02X\n", idx, saveFlags, cpu.f)
		}
	}
	postconditions()
}
func test_0xF1_POP_AF(t *testing.T) {
	preconditions()
	randomizeFlags()

	// to isolate this test from the PUSH BC test, we will push the values into the HRAM and update the stack pointer manually
	// we will push them in reverse order to POP them in the correct order
	for idx, _ := range popTestData {
		value := popTestData[len(popTestData)-1-idx]
		// extract the high and low bytes
		high := uint8(value >> 8)
		low := uint8(value)
		// push the value into the HRAM
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, high)
		cpu.sp -= 1
		cpu.bus.Write(cpu.sp, low)
	}

	// execute the test
	testProgram := []uint8{0xF1, 0x10}
	loadProgramIntoMemory(memory1, testProgram)

	for idx, value := range popTestData {
		// soft reset of the cpu after STOP instruction which blocks the cpu execution and set the offset and pc
		cpu.stopped = false
		cpu.offset = 0x0000
		cpu.pc = 0x0000

		cpu.Run()

		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xF1_POP_AF] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != uint8(value>>8) {
			t.Errorf("[test_0xF1_POP_AF] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, uint8(value>>8), cpu.a)
		}
		// check data correctly loaded into F register
		if cpu.f != uint8(value) {
			t.Errorf("[test_0xF1_POP_AF] %v> expected F register to be 0x%02X, got 0x%02X\n", idx, uint8(value), cpu.f)
		}
		// check if stack pointer has been incremented
		expectedSP := 0xFFFE - uint16(len(popTestData)*2) + uint16(2*(idx+1))
		if cpu.sp != expectedSP {
			t.Errorf("[test_0xF1_POP_AF] %v> expected stack pointer to be incremented to 0x%04X, got 0x%04X\n", idx, expectedSP, cpu.sp)
		}
	}
	postconditions()
}

// Add both operands together (8/16 bits, direct/indirect) and store back to operand 1 location (direct/indirect)
// opcodes:
//   - 0x09 = ADD HL, BC
//   - 0x19 = ADD HL, DE
//   - 0x29 = ADD HL, HL
//   - 0x39 = ADD HL, SP
//
// flags: Z:- N:0 H:H C:C
//
//   - 0x80 = ADD A, B
//   - 0x81 = ADD A, C
//   - 0x82 = ADD A, D
//   - 0x83 = ADD A, E
//   - 0x84 = ADD A, H
//   - 0x85 = ADD A, L
//   - 0x86 = ADD A, [HL]
//   - 0x87 = ADD A, A
//   - 0xC6 = ADD A, n8
//
// flags: Z:Z N:0 H:H C:C
//
//   - 0xE8 = ADD SP, e8
//
// flags: Z:0 N:0 H:H C:C
func TestADD(t *testing.T) {
	t.Run("0x09_ADD_HL_BC", test_0x09_ADD_HL_BC)
	t.Run("0x19_ADD_HL_DE", test_0x19_ADD_HL_DE)
	t.Run("0x29_ADD_HL_HL", test_0x29_ADD_HL_HL)
	t.Run("0x39_ADD_HL_SP", test_0x39_ADD_HL_SP)
	t.Run("0xC6_ADD_A_n8", test_0xC6_ADD_A_n8)
	t.Run("0x87_ADD_A_A", test_0x87_ADD_A_A)
	t.Run("0x80_ADD_A_B", test_0x80_ADD_A_B)
	t.Run("0x81_ADD_A_C", test_0x81_ADD_A_C)
	t.Run("0x82_ADD_A_D", test_0x82_ADD_A_D)
	t.Run("0x83_ADD_A_E", test_0x83_ADD_A_E)
	t.Run("0x84_ADD_A_H", test_0x84_ADD_A_H)
	t.Run("0x85_ADD_A_L", test_0x85_ADD_A_L)
	t.Run("0x86_ADD_A__HL", test_0x86_ADD_A__HL)
	t.Run("0xE8_ADD_SP_e8", test_0xE8_ADD_SP_e8)
}

var testData_ADD_operand1_16bit = []uint16{0x0000, 0x00F0, 0xABCD, 0x1234, 0x5678, 0x9ABC, 0xDEF0, 0xFFAA, 0x19A7, 0xD65C, 0x71B9}
var testData_ADD_operand2_16bit = []uint16{0xFFFF, 0x00F0, 0x1549, 0xF987, 0xB500, 0x0044, 0x5426, 0x0F0F, 0xFF0F, 0x0F0F, 0x0F0F}
var expected_ADD_sum_16bit = []uint16{0xFFFF, 0x01E0, 0xC116, 0x0BBB, 0x0B78, 0x9B00, 0x3316, 0x0EB9, 0x18B6, 0xE56B, 0x80C8}
var expected_ADD_HFlag = []bool{false, false, true, false, false, false, true, true, true, true, true}
var expected_ADD_CFlag = []bool{false, false, false, true, true, false, true, true, true, false, false}

func test_0x09_ADD_HL_BC(t *testing.T) {
	for idx, data := range testData_ADD_operand1_16bit {
		preconditions()
		randomizeFlags()
		saveZFlag := cpu.getZFlag()
		cpu.setHL(data)
		cpu.setBC(testData_ADD_operand2_16bit[idx])
		testProgram := []uint8{0x09, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into HL register
		if cpu.getHL() != expected_ADD_sum_16bit[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected HL register to be 0x%04X, got 0x%04X\n", idx, expected_ADD_sum_16bit[idx], cpu.getHL())
		}
		// check flags
		if cpu.getZFlag() != saveZFlag {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected Z flag to be unaffected %t, got %t\n", idx, saveZFlag, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag[idx], cpu.getCFlag())
		}
	}
}
func test_0x19_ADD_HL_DE(t *testing.T) {
	for idx, data := range testData_ADD_operand1_16bit {
		preconditions()
		randomizeFlags()
		saveZFlag := cpu.getZFlag()
		cpu.setHL(data)
		cpu.setDE(testData_ADD_operand2_16bit[idx])
		testProgram := []uint8{0x19, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into HL register
		if cpu.getHL() != expected_ADD_sum_16bit[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected HL register to be 0x%04X, got 0x%04X\n", idx, expected_ADD_sum_16bit[idx], cpu.getHL())
		}
		// check flags
		if cpu.getZFlag() != saveZFlag {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected Z flag to be unaffected %t, got %t\n", idx, saveZFlag, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag[idx], cpu.getCFlag())
		}
	}
}
func test_0x29_ADD_HL_HL(t *testing.T) {
	var expectedHL = []uint16{0x0000, 0x01E0, 0x579A, 0x2468, 0xACF0, 0x3578, 0xBDE0, 0xFF54, 0x334E, 0xACB8, 0xE372}
	var expectedHFlag = []bool{false, false, true, false, false, true, true, true, true, false, false}
	var expectedCFlag = []bool{false, false, true, false, false, true, true, true, false, true, false}

	for idx, data := range testData_ADD_operand1_16bit {
		preconditions()
		randomizeFlags()
		saveZFlag := cpu.getZFlag()
		cpu.setHL(data)
		testProgram := []uint8{0x29, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into HL register
		if cpu.getHL() != expectedHL[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected HL register to be 0x%04X, got 0x%04X\n", idx, expectedHL[idx], cpu.getHL())
		}
		// check flags
		if cpu.getZFlag() != saveZFlag {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected Z flag to be unaffected %t, got %t\n", idx, saveZFlag, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expectedHFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected H flag to be %t, got %t\n", idx, expectedHFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expectedCFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected C flag to be %t, got %t\n", idx, expectedCFlag[idx], cpu.getCFlag())
		}
	}
}
func test_0x39_ADD_HL_SP(t *testing.T) {
	for idx, data := range testData_ADD_operand1_16bit {
		preconditions()
		randomizeFlags()
		saveZFlag := cpu.getZFlag()
		cpu.setHL(data)
		cpu.sp = testData_ADD_operand2_16bit[idx]
		testProgram := []uint8{0x39, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into HL register
		if cpu.getHL() != expected_ADD_sum_16bit[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected HL register to be 0x%04X, got 0x%04X\n", idx, expected_ADD_sum_16bit[idx], cpu.getHL())
		}
		// check flags
		if cpu.getZFlag() != saveZFlag {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected Z flag to be unaffected %t, got %t\n", idx, saveZFlag, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag[idx] {
			t.Errorf("[test_0x09_ADD_HL_BC] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag[idx], cpu.getCFlag())
		}
	}
}

var testData_ADD_operand1_8bit = []uint8{0x00, 0xF0, 0xAB, 0x12, 0x56, 0xBC, 0xDE, 0xFF, 0xA7, 0xD6, 0x71}
var testData_ADD_operand2_8bit = []uint8{0xFF, 0xF0, 0x15, 0xF9, 0xB5, 0x44, 0x54, 0x01, 0x0F, 0x0F, 0x0F}
var expected_ADD_sum_8bit = []uint8{0xFF, 0xE0, 0xC0, 0x0B, 0x0B, 0x00, 0x32, 0x00, 0xB6, 0xE5, 0x80}
var expected_ADD_HFlag_8bit = []bool{false, false, true, false, false, true, true, true, true, true, true}
var expected_ADD_CFlag_8bit = []bool{false, true, false, true, true, true, true, true, false, false, false}

func test_0xC6_ADD_A_n8(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0xC6, testData_ADD_operand2_8bit[idx], 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0xC6_ADD_A_n8] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x87_ADD_A_A(t *testing.T) {
	var expected_A = []uint8{0x00, 0xE0, 0x56, 0x24, 0xAC, 0x78, 0xBC, 0xFE, 0x4E, 0xAC, 0xE2}
	var expected_HFlag = []bool{false, false, true, false, false, true, true, true, false, false, false}
	var expected_CFlag = []bool{false, true, true, false, false, true, true, true, true, true, false}
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0x87, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_A[idx] {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_A[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_A[idx] != 0 {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_HFlag[idx] {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected H flag to be %t, got %t\n", idx, expected_HFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_CFlag[idx] {
			t.Errorf("[test_0x87_ADD_A_A] %v> expected C flag to be %t, got %t\n", idx, expected_CFlag[idx], cpu.getCFlag())
		}
	}
}
func test_0x80_ADD_A_B(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.b = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x80, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x80_ADD_A_B] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x81_ADD_A_C(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.c = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x81, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x81_ADD_A_C] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x82_ADD_A_D(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.d = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x82, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x82_ADD_A_D] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x83_ADD_A_E(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.e = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x83, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x83_ADD_A_E] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x84_ADD_A_H(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.h = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x84, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x84_ADD_A_H] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x85_ADD_A_L(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.l = testData_ADD_operand2_8bit[idx]
		testProgram := []uint8{0x85, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x85_ADD_A_L] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x86_ADD_A__HL(t *testing.T) {
	for idx, data := range testData_ADD_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.setHL(0x0002)
		testProgram := []uint8{0x86, 0x10, testData_ADD_operand2_8bit[idx]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADD_sum_8bit[idx] {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADD_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADD_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADD_HFlag_8bit[idx] {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected H flag to be %t, got %t\n", idx, expected_ADD_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADD_CFlag_8bit[idx] {
			t.Errorf("[test_0x86_ADD_A__HL] %v> expected C flag to be %t, got %t\n", idx, expected_ADD_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0xE8_ADD_SP_e8(t *testing.T) {
	type TestDataEntry struct {
		sp            uint16
		e8            uint8
		expectedSP    uint16
		expectedHFlag bool
		expectedCFlag bool
	}
	var testData = []TestDataEntry{
		{0xFFF8, 0x10, 0x0008, true, true},   // 0> 0xFFF8 + 0x10 (+16) = 0x0008 - H = 1, C = 1
		{0xFFF8, 0xF0, 0xFFE8, true, true},   // 1> 0xFFF8 + 0xF0 (-16) = 0xFFE8 - H = 1, C = 1
		{0x0001, 0xFF, 0x0000, false, false}, // 2> 0x0001 + 0xFF (-1)  = 0x0000 - H = 0, C = 0
		{0x00FF, 0x01, 0x0100, false, false}, // 3> 0x00FF + 0x01 (+1)  = 0x0100 - H = 0, C = 0
		{0x0100, 0xFE, 0x00FE, false, false}, // 4> 0x0100 + 0xFE (-2)  = 0x00FE - H = 1, C = 1
		{0x7FFF, 0x01, 0x8000, true, false},  // 5> 0x7FFF + 0x01 (+1)  = 0x8000 - H = 1, C = 0
		{0x8000, 0xFF, 0x7FFF, false, false}, // 6> 0x8000 + 0xFF (-1)  = 0x7FFF - H = 0, C = 0
		{0x1234, 0x20, 0x1254, false, false}, // 7> 0x1234 + 0x20 (+32) = 0x1254 - H = 0, C = 0
		{0x00F0, 0x10, 0x0100, false, false}, // 8> 0x00F0 + 0x10 (+16) = 0x0100 - H = 0, C = 0
		{0xABCD, 0x30, 0xABFD, false, false}, // 9> 0xABCD + 0x30 (+48) = 0xABFD - H = 0, C = 0
	}

	for idx, data := range testData {
		preconditions()
		randomizeFlags()
		cpu.sp = data.sp
		testProgram := []uint8{0xE8, data.e8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.sp != data.expectedSP {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected SP register to be 0x%02X, got 0x%02X\n", idx, data.expectedSP, cpu.a)
		}
		// check flags
		if cpu.getZFlag() {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected Z flag to be reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected N flag to be reset\n", idx)
		}
		if cpu.getHFlag() != data.expectedHFlag {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected H flag to be %t, got %t\n", idx, data.expectedHFlag, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.expectedCFlag {
			t.Errorf("[test_0xE8_ADD_SP_e8] %v> expected C flag to be %t, got %t\n", idx, data.expectedCFlag, cpu.getCFlag())
		}
	}
}

var testData_ADC_operand1_8bit = []uint8{0x00, 0xF0, 0xAB, 0x12, 0x56, 0xBC, 0xDE, 0xFF, 0xA7, 0xD6, 0x71}
var testData_ADC_operand2_8bit = []uint8{0xFE, 0xF0, 0x14, 0xF9, 0xB4, 0x44, 0x53, 0x01, 0x0F, 0x0E, 0x0F}
var testData_ADC_carry = []bool{true, false, true, false, true, false, true, false, true, true, false}
var expected_ADC_sum_8bit = []uint8{0xFF, 0xE0, 0xC0, 0x0B, 0x0B, 0x00, 0x32, 0x00, 0xB7, 0xE5, 0x80}
var expected_ADC_HFlag_8bit = []bool{false, false, true, false, false, true, true, true, true, true, true}
var expected_ADC_CFlag_8bit = []bool{false, true, false, true, true, true, true, true, false, false, false}

// Add both operand and carry flag to A register (8 bits, direct/indirect) and store back to register A
// opcodes:
//   - 0x88 = ADC A, B
//   - 0x89 = ADC A, C
//   - 0x8A = ADC A, D
//   - 0x8B = ADC A, E
//   - 0x8C = ADC A, H
//   - 0x8D = ADC A, L
//   - 0x8E = ADC A, [HL]
//   - 0x8F = ADC A, A
//   - 0xCE = ADC A, n8
//
// flags: Z:Z N:0 H:H C:C
func TestADC(t *testing.T) {
	t.Run("0x88: ADC A, B", test_0x88_ADC_A_B)
	t.Run("0x89: ADC A, C", test_0x89_ADC_A_C)
	t.Run("0x8A: ADC A, D", test_0x8A_ADC_A_D)
	t.Run("0x8B: ADC A, E", test_0x8B_ADC_A_E)
	t.Run("0x8C: ADC A, H", test_0x8C_ADC_A_H)
	t.Run("0x8D: ADC A, L", test_0x8D_ADC_A_L)
	t.Run("0x8E: ADC A, [HL]", test_0x8E_ADC_A__HL)
	t.Run("0x8F: ADC A, A", test_0x8F_ADC_A_A)
	t.Run("0xCE: ADC A, n8", test_0xCE_ADC_A_n8)
}
func test_0x88_ADC_A_B(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.b = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x88, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x88_ADC_A_B] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x89_ADC_A_C(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.c = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x89, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8A_ADC_A_D(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.d = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x8A_ADC_A_D] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8B_ADC_A_E(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.e = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x8B_ADC_A_E] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8C_ADC_A_H(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.h = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x8C_ADC_A_H] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8D_ADC_A_L(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.l = testData_ADC_operand2_8bit[idx]
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x8D_ADC_A_L] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8E_ADC_A__HL(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.setHL(0x0002)
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8E, 0x10, testData_ADC_operand2_8bit[idx]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x8E_ADC_A__HL] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}
func test_0x8F_ADC_A_A(t *testing.T) {
	expectedA := []uint8{0x01, 0xE0, 0x57, 0x24, 0xAD, 0x78, 0xBD, 0xFE, 0x4F, 0xAD, 0xE2}
	expectedHFlag := []bool{false, false, true, false, false, true, true, true, false, false, false}
	expectedCFlag := []bool{false, true, true, false, false, true, true, true, true, true, false}
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0x8F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expectedA[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expectedA[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expectedA[idx] != 0 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expectedHFlag[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected H flag to be %t, got %t\n", idx, expectedHFlag[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expectedCFlag[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected C flag to be %t, got %t\n", idx, expectedCFlag[idx], cpu.getCFlag())
		}
	}
}
func test_0xCE_ADC_A_n8(t *testing.T) {
	for idx, data := range testData_ADC_operand1_8bit {
		preconditions()
		randomizeFlags()
		cpu.a = data
		if testData_ADC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCE, testData_ADC_operand2_8bit[idx], 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_ADC_sum_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_ADC_sum_8bit[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() && expected_ADC_sum_8bit[idx] != 0 {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected Z flag to be false, got true\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected N flag to be reset, got set\n", idx)
		}
		if cpu.getHFlag() != expected_ADC_HFlag_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected H flag to be %t, got %t\n", idx, expected_ADC_HFlag_8bit[idx], cpu.getHFlag())
		}
		if cpu.getCFlag() != expected_ADC_CFlag_8bit[idx] {
			t.Errorf("[test_0x89_ADC_A_C] %v> expected C flag to be %t, got %t\n", idx, expected_ADC_CFlag_8bit[idx], cpu.getCFlag())
		}
	}
}

// Bitwise AND operation between A register and operand (8 bits, direct/indirect) and store back to register A
// opcodes:
//   - 0xA0 = AND A, B
//   - 0xA1 = AND A, C
//   - 0xA2 = AND A, D
//   - 0xA3 = AND A, E
//   - 0xA4 = AND A, H
//   - 0xA5 = AND A, L
//   - 0xA6 = AND A, [HL]
//   - 0xA7 = AND A, A
//   - 0xE6 = AND A, n8
//
// flags: Z:Z N:0 H:1 C:0
func TestAND(t *testing.T) {
	t.Run("0xA0: AND A, B", test_0xA0_AND_A_B)
	t.Run("0xA1: AND A, C", test_0xA1_AND_A_C)
	t.Run("0xA2: AND A, D", test_0xA2_AND_A_D)
	t.Run("0xA3: AND A, E", test_0xA3_AND_A_E)
	t.Run("0xA4: AND A, H", test_0xA4_AND_A_H)
	t.Run("0xA5: AND A, L", test_0xA5_AND_A_L)
	t.Run("0xA6: AND A, [HL]", test_0xA6_AND_A__HL)
	t.Run("0xA7: AND A, A", test_0xA7_AND_A_A)
	t.Run("0xE6: AND A, n8", test_0xE6_AND_A_n8)
}

var testData_AND_operand1 = []uint8{0b00000000, 0b11111111, 0b00001111, 0b11110000, 0b11001100, 0b00110011, 0b10101010, 0b01010101}
var testData_AND_operand2 = []uint8{0b11111111, 0b00001111, 0b11001100, 0b00110011, 0b10101010, 0b01010101, 0b10101010, 0b00110011}
var expected_AND_result = []uint8{0b00000000, 0b00001111, 0b00001100, 0b00110000, 0b10001000, 0b00010001, 0b10101010, 0b00010001}

func test_0xA0_AND_A_B(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.b = testData_AND_operand2[idx]
		testProgram := []uint8{0xA0, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA0_AND_A_B] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA1_AND_A_C(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.c = testData_AND_operand2[idx]
		testProgram := []uint8{0xA1, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA1_AND_A_C] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA2_AND_A_D(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.d = testData_AND_operand2[idx]
		testProgram := []uint8{0xA2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA2_AND_A_D] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA3_AND_A_E(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.e = testData_AND_operand2[idx]
		testProgram := []uint8{0xA3, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA3_AND_A_E] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA4_AND_A_H(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.h = testData_AND_operand2[idx]
		testProgram := []uint8{0xA4, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA4_AND_A_H] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA5_AND_A_L(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.l = testData_AND_operand2[idx]
		testProgram := []uint8{0xA5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA5_AND_A_L] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA6_AND_A__HL(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		cpu.setHL(0x0002)
		testProgram := []uint8{0xA6, 0x10, testData_AND_operand2[idx]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA6_AND_A__HL] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xA7_AND_A_A(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0xA7, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected PC to be 0x0001, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != testData_AND_operand1[idx] {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, testData_AND_operand1[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xA7_AND_A_A] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}
func test_0xE6_AND_A_n8(t *testing.T) {
	for idx, data := range testData_AND_operand1 {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0xE6, testData_AND_operand2[idx], 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected PC to be 0x0002, got 0x%04X\n", idx, cpu.pc)
		}
		// check data correctly loaded into A register
		if cpu.a != expected_AND_result[idx] {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected A register to be 0x%02X, got 0x%02X\n", idx, expected_AND_result[idx], cpu.a)
		}
		// check flags
		if expected_AND_result[idx] == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected Z flag to be set, got reset\n", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected N flag to be reset, got set\n", idx)
		}
		if !cpu.getHFlag() {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected H flag to be set, got reset\n", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xE6_AND_A_n8] %v> expected C flag to be reset, got set\n", idx)
		}
	}
}

// INC: Increment register or memory @[HL]
// opcodes:
//   - 0x3C=INC A
//   - 0x04=INC B
//   - 0x0C=INC C
//   - 0x14=INC D
//   - 0x1C=INC E
//   - 0x24=INC H
//   - 0x2C=INC L
//   - 0x34=INC [HL]
//   - 0x03=INC BC
//   - 0x13=INC DE
//   - 0x23=INC HL
//   - 0x33=INC SP
//
// flags: Z:Z N:0 H:H C:-
func TestINC(t *testing.T) {
	t.Run("0x3C: INC A", test_0x3C_INC_A)
	t.Run("0x04: INC B", test_0x04_INC_B)
	t.Run("0x0C: INC C", test_0x0C_INC_C)
	t.Run("0x14: INC D", test_0x14_INC_D)
	t.Run("0x1C: INC E", test_0x1C_INC_E)
	t.Run("0x24: INC H", test_0x24_INC_H)
	t.Run("0x2C: INC L", test_0x2C_INC_L)
	t.Run("0x34: INC [HL]", test_0x34_INC_HL)
	t.Run("0x03: INC BC", test_0x03_INC_BC)
	t.Run("0x13: INC DE", test_0x13_INC_DE)
	t.Run("0x23: INC HL", test_0x23_INC_HL)
	t.Run("0x33: INC SP", test_0x33_INC_SP)
}
func test_0x3C_INC_A(t *testing.T) {
	// TC 1: increment register A 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.a = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x3C, 0x3C, 0x3C, 0x3C, 0x3C, 0x3C, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x3C_INC_A] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register A = 0x77
	if cpu.a != 0x77 {
		t.Errorf("[test_0x3C_INC_A] Error> the value of register A should have been set to 0x77, got 0x%02X \n", cpu.a)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register A from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.a = 0x000F

	testData2 := []uint8{0x3C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x3C_INC_A] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value at [HL] = 0x10
	if cpu.a != 0x10 {
		t.Errorf("[test_0x3C_INC_A] Error> the value of register A should have been set to 0x10, got 0x%02X \n", cpu.a)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing [HL] from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.a = 0x00FF

	testData3 := []uint8{0x3C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x3C_INC_A] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register A = 0x00
	if cpu.a != 0x00 {
		t.Errorf("[test_0x3C_INC_A] Error> the valueof register A should have been set to 0x00, got 0x%02X \n", cpu.a)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x3C_INC_A] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x04_INC_B(t *testing.T) {
	// TC 1: increment register B 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.b = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x04_INC_B] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x77
	if cpu.b != 0x77 {
		t.Errorf("[test_0x04_INC_B] Error> the value of register B should have been set to 0x77, got 0x%02X \n", cpu.b)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register B from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.b = 0x000F

	testData2 := []uint8{0x04, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x04_INC_B] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x10
	if cpu.b != 0x10 {
		t.Errorf("[test_0x04_INC_B] Error> the value of register B should have been set to 0x10, got 0x%02X \n", cpu.b)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register B from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.b = 0x00FF

	testData3 := []uint8{0x04, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x04_INC_B] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x00
	if cpu.b != 0x00 {
		t.Errorf("[test_0x04_INC_B] Error> the valueof register B should have been set to 0x00, got 0x%02X \n", cpu.b)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x04_INC_B] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x0C_INC_C(t *testing.T) {
	// TC 1: increment register C 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.c = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x0C, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x0C_INC_C] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register C = 0x77
	if cpu.c != 0x77 {
		t.Errorf("[test_0x0C_INC_C] Error> the value of register C should have been set to 0x77, got 0x%02X \n", cpu.c)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register C from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.c = 0x000F

	testData2 := []uint8{0x0C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x0C_INC_C] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register C = 0x10
	if cpu.c != 0x10 {
		t.Errorf("[test_0x0C_INC_C] Error> the value of register C should have been set to 0x10, got 0x%02X \n", cpu.c)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register C from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.c = 0x00FF

	testData3 := []uint8{0x0C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x0C_INC_C] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register C = 0x00
	if cpu.c != 0x00 {
		t.Errorf("[test_0x0C_INC_C] Error> the valueof register C should have been set to 0x00, got 0x%02X \n", cpu.c)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x0C_INC_C] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x14_INC_D(t *testing.T) {
	// TC 1: increment register D 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.d = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x14_INC_D] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register D = 0x77
	if cpu.d != 0x77 {
		t.Errorf("[test_0x14_INC_D] Error> the value of register D should have been set to 0x77, got 0x%02X \n", cpu.d)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register D from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.d = 0x000F

	testData2 := []uint8{0x14, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x14_INC_D] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register D = 0x10
	if cpu.d != 0x10 {
		t.Errorf("[test_0x14_INC_D] Error> the value of register D should have been set to 0x10, got 0x%02X \n", cpu.d)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register D from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.d = 0x00FF

	testData3 := []uint8{0x14, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x14_INC_D] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register D = 0x00
	if cpu.d != 0x00 {
		t.Errorf("[test_0x14_INC_D] Error> the valueof register D should have been set to 0x00, got 0x%02X \n", cpu.d)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x14_INC_D] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x1C_INC_E(t *testing.T) {
	// TC 1: increment register E 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.e = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x1C, 0x1C, 0x1C, 0x1C, 0x1C, 0x1C, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x1C_INC_E] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register E = 0x77
	if cpu.e != 0x77 {
		t.Errorf("[test_0x1C_INC_E] Error> the value of register E should have been set to 0x77, got 0x%02X \n", cpu.e)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register E from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.e = 0x000F

	testData2 := []uint8{0x1C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x1C_INC_E] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register E = 0x10
	if cpu.e != 0x10 {
		t.Errorf("[test_0x1C_INC_E] Error> the value of register E should have been set to 0x10, got 0x%02X \n", cpu.e)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register E from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.e = 0x00FF

	testData3 := []uint8{0x1C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x1C_INC_E] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register E = 0x00
	if cpu.e != 0x00 {
		t.Errorf("[test_0x1C_INC_E] Error> the valueof register E should have been set to 0x00, got 0x%02X \n", cpu.e)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x1C_INC_E] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x24_INC_H(t *testing.T) {
	// TC 1: increment register H 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.h = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x24, 0x24, 0x24, 0x24, 0x24, 0x24, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x24_INC_H] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x77
	if cpu.h != 0x77 {
		t.Errorf("[test_0x24_INC_H] Error> the value of register H should have been set to 0x77, got 0x%02X \n", cpu.h)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register H from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.h = 0x000F

	testData2 := []uint8{0x24, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x24_INC_H] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x10
	if cpu.h != 0x10 {
		t.Errorf("[test_0x24_INC_H] Error> the value of register H should have been set to 0x10, got 0x%02X \n", cpu.h)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register H from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.h = 0x00FF

	testData3 := []uint8{0x24, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x24_INC_H] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x00
	if cpu.h != 0x00 {
		t.Errorf("[test_0x24_INC_H] Error> the valueof register H should have been set to 0x00, got 0x%02X \n", cpu.h)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x2C_INC_L(t *testing.T) {
	// TC 1: increment register L 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.l = 0x0071
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x2C, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x2C_INC_L] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x77
	if cpu.l != 0x77 {
		t.Errorf("[test_0x2C_INC_L] Error> the value of register L should have been set to 0x77, got 0x%02X \n", cpu.l)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register L from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.l = 0x000F

	testData2 := []uint8{0x2C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x2C_INC_L] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x10
	if cpu.l != 0x10 {
		t.Errorf("[test_0x2C_INC_L] Error> the value of register L should have been set to 0x10, got 0x%02X \n", cpu.l)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register L from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.l = 0x00FF

	testData3 := []uint8{0x2C, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x2C_INC_L] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register H = 0x00
	if cpu.l != 0x00 {
		t.Errorf("[test_0x2C_INC_L] Error> the valueof register L should have been set to 0x00, got 0x%02X \n", cpu.l)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x2C_INC_L] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x24_INC_H] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x34_INC_HL(t *testing.T) {
	// TC 1: increment [HL] 6 times from 0x71 to 0x77 and stop @0x0006
	preconditions()

	cpu.setHL(0x0007)
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x34, 0x34, 0x34, 0x34, 0x34, 0x34, 0x10, 0x71}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x34_INC_HL] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value at [HL] = 0x77
	if bus.Read(0x0007) != 0x77 {
		t.Errorf("[test_0x34_INC_HL] Error> the value at [HL] should have been set to 0x77, got 0x%02X \n", bus.Read(0x0007))
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing [HL] from 0x0F to 0x10 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setHL(0x0007)

	testData2 := []uint8{0x34, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x34_INC_HL] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value at [HL] = 0x10
	if bus.Read(0x0007) != 0x10 {
		t.Errorf("[test_0x34_INC_HL] Error> the value at [HL] should have been set to 0x10, got 0x%02X \n", bus.Read(0x0007))
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing [HL] from 0xFF to 0x00 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setHL(0x0007)

	testData3 := []uint8{0x34, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x34_INC_HL] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value at [HL] = 0x00
	if bus.Read(0x0007) != 0x00 {
		t.Errorf("[test_0x34_INC_HL] Error> the value at [HL] should have been set to 0x00, got 0x%02X \n", bus.Read(0x0007))
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x34_INC_HL] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x03_INC_BC(t *testing.T) {
	// TC 1: increment register BC 6 times from 0xFF71 to 0xFF77 and stop @0x0006
	preconditions()

	cpu.setBC(0xFF71)
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x03_INC_BC] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0xFF77
	if cpu.getBC() != 0xFF77 {
		t.Errorf("[test_0x03_INC_BC] Error> the value of register BC should have been set to 0xFF77, got 0x%04X \n", cpu.getBC())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register BC from 0x00FF to 0x0100 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setBC(0x00FF)

	testData2 := []uint8{0x03, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x03_INC_BC] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x0100
	if cpu.getBC() != 0x0100 {
		t.Errorf("[test_0x03_INC_BC] Error> the value of register BC should have been set to 0x10, got 0x%04X \n", cpu.getBC())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register BC from 0xFFFF to 0x0000 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setBC(0xFFFF)

	testData3 := []uint8{0x03, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x03_INC_BC] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0x0000
	if cpu.getBC() != 0x0000 {
		t.Errorf("[test_0x03_INC_BC] Error> the valueof register BC should have been set to 0x00, got 0x%04X \n", cpu.getBC())
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x03_INC_BC] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x13_INC_DE(t *testing.T) {
	// TC 1: increment register DE 6 times from 0xFF71 to 0xFF77 and stop @0x0006
	preconditions()

	cpu.setDE(0xFF71)
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x13, 0x13, 0x13, 0x13, 0x13, 0x13, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x13_INC_DE] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0xFF77
	if cpu.getDE() != 0xFF77 {
		t.Errorf("[test_0x13_INC_DE] Error> the value of register DE should have been set to 0xFF77, got 0x%04X \n", cpu.getDE())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register DE from 0x00FF to 0x0100 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setDE(0x00FF)

	testData2 := []uint8{0x13, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x13_INC_DE] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x0100
	if cpu.getDE() != 0x0100 {
		t.Errorf("[test_0x13_INC_DE] Error> the value of register DE should have been set to 0x10, got 0x%04X \n", cpu.getDE())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register DE from 0xFFFF to 0x0000 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setDE(0xFFFF)

	testData3 := []uint8{0x13, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x13_INC_DE] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0x0000
	if cpu.getDE() != 0x0000 {
		t.Errorf("[test_0x13_INC_DE] Error> the valueof register DE should have been set to 0x00, got 0x%04X \n", cpu.getDE())
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x13_INC_DE] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x23_INC_HL(t *testing.T) {
	// TC 1: increment register HL 6 times from 0xFF71 to 0xFF77 and stop @0x0006
	preconditions()

	cpu.setHL(0xFF71)
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x23, 0x23, 0x23, 0x23, 0x23, 0x23, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x23_INC_HL] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0xFF77
	if cpu.getHL() != 0xFF77 {
		t.Errorf("[test_0x23_INC_HL] Error> the value of register DE should have been set to 0xFF77, got 0x%04X \n", cpu.getHL())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register HL from 0x00FF to 0x0100 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setHL(0x00FF)

	testData2 := []uint8{0x23, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x23_INC_HL] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x0100
	if cpu.getHL() != 0x0100 {
		t.Errorf("[test_0x23_INC_HL] Error> the value of register HL should have been set to 0x10, got 0x%04X \n", cpu.getHL())
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register HL from 0xFFFF to 0x0000 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.setHL(0xFFFF)

	testData3 := []uint8{0x23, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x23_INC_HL] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register BC = 0x0000
	if cpu.getHL() != 0x0000 {
		t.Errorf("[test_0x23_INC_HL] Error> the valueof register HL should have been set to 0x00, got 0x%04X \n", cpu.getHL())
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x23_INC_HL] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}
func test_0x33_INC_SP(t *testing.T) {
	// TC 1: increment register SP 6 times from 0xFF71 to 0xFF77 and stop @0x0006
	preconditions()

	cpu.sp = 0xFF71
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()

	testData1 := []uint8{0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x10, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData1)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0006 {
		t.Errorf("[test_0x33_INC_SP] Error> the program counter should have stopped at 0x0006, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register SP = 0xFF77
	if cpu.sp != 0xFF77 {
		t.Errorf("[test_0x33_INC_SP] Error> the value of register SP should have been set to 0xFF77, got 0x%04X \n", cpu.sp)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is reset
	if cpu.getHFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the H flag should have been reset, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC2: check H flag (Z:Z N:0 H:H C:-) by incrementing register SP from 0x00FF to 0x0100 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.sp = 0x00FF

	testData2 := []uint8{0x33, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData2)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x33_INC_SP] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register B = 0x0100
	if cpu.sp != 0x0100 {
		t.Errorf("[test_0x33_INC_SP] Error> the value of register SP should have been set to 0x10, got 0x%04X \n", cpu.sp)
	}

	// check that Z flag is not set
	if cpu.getZFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the Z flag should have been reset, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the N flag should have been reset, got %t \n", cpu.getNFlag())
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the C flag should have been left untouched, got %t \n", cpu.getCFlag())
	}

	postconditions()

	// TC3: check Z & H flags (Z:Z N:0 H:H C:-) by incrementing register SP from 0xFFFF to 0x0000 and stop @0x0001
	preconditions()

	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()
	cpu.setCFlag()
	cpu.sp = 0xFFFF

	testData3 := []uint8{0x33, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	loadProgramIntoMemory(memory1, testData3)
	cpu.Run()

	// check the final state of the cpu
	if cpu.pc != 0x0001 {
		t.Errorf("[test_0x33_INC_SP] Error> the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
	}

	// check if the value of register SP = 0x0000
	if cpu.sp != 0x0000 {
		t.Errorf("[test_0x33_INC_SP] Error> the valueof register SP should have been set to 0x00, got 0x%04X \n", cpu.sp)
	}

	// check that Z flag is set
	if !cpu.getZFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the Z flag should have been set, got %t \n", cpu.getZFlag())
	}

	// check that N flag is reset
	if cpu.getNFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the N flag should have been reset \n")
	}

	// check if the H flag is set
	if !cpu.getHFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the H flag should have been set, got %t \n", cpu.getHFlag())
	}

	// check if C flag is left untouched
	if !cpu.getCFlag() {
		t.Errorf("[test_0x33_INC_SP] Error> the C flag should have been left untouched \n")
	}
	postconditions()
}

// CCF: should flip the carry flag
// 0x3F: CCF
// Complement Carry Flag: flip the value of the carry flag and reset the N and H flags
func TestCCF(t *testing.T) {
	preconditions()

	// set all flags to 1
	cpu.f = 0x00
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()
	cpu.setCFlag()

	// save the initial state of the cpu
	initialState := getCpuState()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x3F, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the Z flag was left untouched (1)
	if finalState.Z != initialState.Z {
		t.Errorf("Expected Z flag to be 1")
	}

	// check if the N flag was reset
	if finalState.N != false {
		t.Errorf("Expected N flag to be 0")
	}

	// check if the H flag was reset
	if finalState.H != false {
		t.Errorf("Expected H flag to be 0")
	}

	// check if the C flag was set
	if finalState.C == true {
		t.Errorf("Expected C flag to be flipped")
	}
}

// CP: compare 2 memory locations and/or registers by subtracting them without storing the result
// opcodes:
//   - BF = CP A, A
//   - B8 = CP A, B
//   - B9 = CP A, C
//   - BA = CP A, D
//   - BB = CP A, E
//   - BC = CP A, H
//   - BD = CP A, L
//   - BE = CP A, [HL]
//   - FE = CP A, n8
//     flags: Z:Z N:1 H:H C:C
func TestCP(t *testing.T) {
	t.Run("0xBF: CP A, A", test_0xBF_CP_A_A)
	t.Run("0xB8: CP A, B", test_0xB8_CP_A_B)
	t.Run("0xB9: CP A, C", test_0xB9_CP_A_C)
	t.Run("0xBA: CP A, D", test_0xBA_CP_A_D)
	t.Run("0xBB: CP A, E", test_0xBB_CP_A_E)
	t.Run("0xBC: CP A, H", test_0xBC_CP_A_H)
	t.Run("0xBD: CP A, L", test_0xBD_CP_A_L)
	t.Run("0xBE: CP A, [HL]", test_0xBE_CP_A_HL)
	t.Run("0xFE: CP A, n8", test_0xFE_CP_A_n8)
}

type TestData_CP struct {
	minuend    uint8
	subtrahend uint8
	Z          bool
	N          bool
	H          bool
	C          bool
}

var testDataCP = []TestData_CP{
	{0x00, 0x00, true, true, false, false},
	{0x00, 0x01, false, true, true, true},
	{0x01, 0x00, false, true, false, false},
	{0x01, 0x01, true, true, false, false},
	{0x10, 0x00, false, true, false, false},
	{0x10, 0x01, false, true, true, false},
	{0x0F, 0x01, false, true, false, false},
	{0x0F, 0x0F, true, true, false, false},
	{0xF0, 0x01, false, true, true, false},
	{0xF0, 0xF0, true, true, false, false},
	{0xFF, 0x01, false, true, false, false},
	{0xFF, 0xFF, true, true, false, false},
}

func test_0xBF_CP_A_A(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		testProgram := []uint8{0xBF, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB8_CP_A_B] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A is left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xB8_CP_A_B] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		// check flags
		if !cpu.getZFlag() {
			t.Errorf("[test_0xB8_CP_A_B] %d> the Z flag should always be set, got %t \n", idx, cpu.getZFlag())
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0xB8_CP_A_B] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB8_CP_A_B] %d> the H flag should have be reset, got %t \n", idx, cpu.getHFlag())
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB8_CP_A_B] %d> the C flag should have be reset, got %t \n", idx, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xB8_CP_A_B(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.b = data.subtrahend
		testProgram := []uint8{0xB8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB8_CP_A_B] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & B are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xB8_CP_A_B] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.b != data.subtrahend {
			t.Errorf("[test_0xB8_CP_A_B] %d> the value of register B should have been left untouched, got 0x%02X \n", idx, cpu.b)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xB8_CP_A_B] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xB8_CP_A_B] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xB8_CP_A_B] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xB8_CP_A_B] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xB9_CP_A_C(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.c = data.subtrahend
		testProgram := []uint8{0xB9, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB9_CP_A_C] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & C are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xB9_CP_A_C] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.c != data.subtrahend {
			t.Errorf("[test_0xB9_CP_A_C] %d> the value of register C should have been left untouched, got 0x%02X \n", idx, cpu.c)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xB9_CP_A_C] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xB9_CP_A_C] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xB9_CP_A_C] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xB9_CP_A_C] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xBA_CP_A_D(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.d = data.subtrahend
		testProgram := []uint8{0xBA, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xBA_CP_A_D] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & D are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xBA_CP_A_D] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.d != data.subtrahend {
			t.Errorf("[test_0xBA_CP_A_D] %d> the value of register D should have been left untouched, got 0x%02X \n", idx, cpu.d)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xBA_CP_A_D] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xBA_CP_A_D] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xBA_CP_A_D] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xBA_CP_A_D] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xBB_CP_A_E(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.e = data.subtrahend
		testProgram := []uint8{0xBB, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xBB_CP_A_E] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & E are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xBB_CP_A_E] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.e != data.subtrahend {
			t.Errorf("[test_0xBB_CP_A_E] %d> the value of register E should have been left untouched, got 0x%02X \n", idx, cpu.e)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xBB_CP_A_E] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xBB_CP_A_E] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xBB_CP_A_E] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xBB_CP_A_E] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xBC_CP_A_H(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.h = data.subtrahend
		testProgram := []uint8{0xBC, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xBC_CP_A_H] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & H are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xBC_CP_A_H] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.h != data.subtrahend {
			t.Errorf("[test_0xBC_CP_A_H] %d> the value of register H should have been left untouched, got 0x%02X \n", idx, cpu.h)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xBC_CP_A_H] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xBC_CP_A_H] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xBC_CP_A_H] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xBC_CP_A_H] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xBD_CP_A_L(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.l = data.subtrahend
		testProgram := []uint8{0xBD, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xBD_CP_A_L] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & L are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xBD_CP_A_L] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		if cpu.l != data.subtrahend {
			t.Errorf("[test_0xBD_CP_A_L] %d> the value of register L should have been left untouched, got 0x%02X \n", idx, cpu.l)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xBD_CP_A_L] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xBD_CP_A_L] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xBD_CP_A_L] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xBD_CP_A_L] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xBE_CP_A_HL(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		cpu.setHL(0x0002)
		testProgram := []uint8{0xBE, 0x10, data.subtrahend}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the program counter should have stopped at 0x0001, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & [HL] are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		memoryLocationHLValue := cpu.bus.Read(cpu.getHL())
		if memoryLocationHLValue != data.subtrahend {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the value of memory location [HL] should have been left untouched, got 0x%02X \n", idx, memoryLocationHLValue)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xBE_CP_A_HL] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0xFE_CP_A_n8(t *testing.T) {
	for idx, data := range testDataCP {
		preconditions()
		randomizeFlags()
		cpu.a = data.minuend
		testProgram := []uint8{0xFE, data.subtrahend, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the program counter should have stopped at 0x0002, got 0x%04X \n", idx, cpu.pc)
		}
		// check that register A & B are left untouched
		if cpu.a != data.minuend {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the value of register A should have been left untouched, got 0x%02X \n", idx, cpu.a)
		}
		n8OperandValue := cpu.bus.Read(0x0001)
		if n8OperandValue != data.subtrahend {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the value of operand at memory location 0x0001 should have been left untouched, got 0x%02X \n", idx, n8OperandValue)
		}
		// check flags
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the Z flag should have been set to %t, got %t \n", idx, data.Z, cpu.getZFlag())
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the N flag should always be set, got %t \n", idx, cpu.getNFlag())
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the H flag should have been set to %t, got %t \n", idx, data.H, cpu.getHFlag())
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0xFE_CP_A_n8] %d> the C flag should have been set to %t, got %t \n", idx, data.C, cpu.getCFlag())
		}
		postconditions()
	}
}

// CPL: Complement A (flip all bits)
// opcodes: 0x2F=CPL
// flags: Z:- N:1 H:1 C:-
func TestCPL(t *testing.T) {
	var testData_CPL = []uint8{0b00000000, 0b11111111, 0b10101010, 0b01010101, 0b11001100, 0b00110011, 0b11110000, 0b00001111}
	var expectedResults = []uint8{0b11111111, 0b00000000, 0b01010101, 0b10101010, 0b00110011, 0b11001100, 0b00001111, 0b11110000}

	for idx, data := range testData_CPL {
		preconditions()
		randomizeFlags()
		saveZFlag := cpu.getZFlag()
		saveCFlag := cpu.getCFlag()
		cpu.a = data
		testProgram := []uint8{0x2F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check the final state of the cpu
		if cpu.pc != 0x0001 {
			t.Errorf("[TestCPL] the program counter should have stopped at 0x0001, got 0x%04X \n", cpu.pc)
		}
		// check that register A is left untouched
		expectedValue := expectedResults[idx]
		if cpu.a != expectedValue {
			t.Errorf("[TestCPL] the value of register A should have been set to 0x%02X, got 0x%02X \n", expectedValue, cpu.a)
		}
		// check flags
		if cpu.getZFlag() != saveZFlag {
			t.Errorf("[TestCPL] expected Z flag to be unchanged %t, got %t \n", saveZFlag, cpu.getZFlag())
		}
		if !cpu.getNFlag() {
			t.Errorf("[TestCPL] expected N flag to be set, got %t \n", cpu.getNFlag())
		}
		if !cpu.getHFlag() {
			t.Errorf("[TestCPL] expected H flag to be set, got %t \n", cpu.getHFlag())
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("[TestCPL] expected C flag to be unchanged %t, got %t \n", saveCFlag, cpu.getCFlag())
		}
	}
}

// DAA: should adjust the destination to be a valid BCD number
// 0x27: DAA
// Decimal Adjust Accumulator to get a correct BCD representation after an arithmetic instruction.
// The steps performed depend on the type of the previous operation: ADD or SUB
func TestDAA(t *testing.T) {
	// Test Cases
	type params struct {
		A uint8
		Z bool
		N bool
		H bool
		C bool
	}
	type testCase struct {
		initial  params
		expected params
	}
	testCases := []testCase{
		// TC0: A = 0x00, no flags set
		{initial: params{
			A: 0x55,
			Z: false,
			N: false,
			H: false,
			C: false,
		}, expected: params{
			A: 0x55,
			Z: false,
			N: false,
			H: false,
			C: false,
		}},
		// TC1: A = 0x0A, no flags set
		{initial: params{
			A: 0x0A,
			Z: false,
			N: false,
			H: false,
			C: false,
		}, expected: params{
			A: 0x10,
			Z: false,
			N: false,
			H: false,
			C: false,
		}},
		// TC2: A = 0x99, no flags set
		{initial: params{
			A: 0x99,
			Z: false,
			N: false,
			H: false,
			C: false,
		}, expected: params{
			A: 0x99,
			Z: false,
			N: false,
			H: false,
			C: false,
		}},
		// TC3: A = 0x42, no flags set
		{initial: params{
			A: 0x42,
			Z: false,
			N: false,
			H: false,
			C: false,
		}, expected: params{
			A: 0x42,
			Z: false,
			N: false,
			H: false,
			C: false,
		}},
		// TC4: A = 0x00, Z flag set
		{
			initial: params{
				A: 0x00,
				Z: true,
				N: false,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x00,
				Z: true,
				N: false,
				H: false,
				C: false,
			},
		},
		// TC5: A = 0x09, no flags set
		{
			initial: params{
				A: 0x09,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x09,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
		},
		// TC6: A = 0xA0, no flags set
		{
			initial: params{
				A: 0xA0,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x00,
				Z: true,
				N: false,
				H: false,
				C: true,
			},
		},
		// TC7: A = 0x99, no flags set
		{
			initial: params{
				A: 0x99,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x99,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
		},
		// TC8: A=0x99+0x66=0xFF, no flags set
		{
			initial: params{
				A: 0xFF,
				Z: false,
				N: false,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x65,
				Z: false,
				N: false,
				H: false,
				C: true,
			},
		},
		// TC9: A=0x55-0x55=0x00, N flag set
		{
			initial: params{
				A: 0x00,
				Z: true,
				N: true,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x00,
				Z: true,
				N: true,
				H: false,
				C: false,
			},
		},
		// TC10
		{
			initial: params{
				A: 0x10,
				Z: false,
				N: true,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x10,
				Z: false,
				N: true,
				H: false,
				C: false,
			},
		},
		// TC11
		{
			initial: params{
				A: 0x99,
				Z: false,
				N: true,
				H: false,
				C: false,
			},
			expected: params{
				A: 0x99,
				Z: false,
				N: true,
				H: false,
				C: false,
			},
		},
		// TC12: A=0x00-0x01=0xFF, N & C flags set
		{
			initial: params{
				A: 0xFF,
				Z: false,
				N: true,
				H: false,
				C: true,
			},
			expected: params{
				A: 0x05,
				Z: false,
				N: true,
				H: false,
				C: true,
			},
		},
	}

	for tci, tc := range testCases {
		preconditions()

		// set the A register to 10 = 0x0A
		cpu.a = tc.initial.A

		// Set flags
		if tc.initial.Z {
			cpu.setZFlag() // to confirm that Z is reset if the result is not 0
		} else {
			cpu.resetZFlag() // to confirm that Z is reset
		}
		if tc.initial.N {
			cpu.setNFlag() // last operation was an addition
		} else {
			cpu.resetNFlag() // last operation was a subtraction
		}
		if tc.initial.H {
			cpu.setHFlag() // to confirm that H is reset
		} else {
			cpu.resetHFlag() // to confirm that H is reset
		}
		if tc.initial.C {
			cpu.setCFlag() // to confirm that C is reset
		} else {
			cpu.resetCFlag() // to confirm that C is reset
		}

		// Execute the instruction
		cpu.executeInstruction(GetInstruction("0x27", false))

		// check if the program counter was incremented by instruction.Length
		//if cpu.pc != cpuCopy.PC+1 {
		//	t.Errorf("TC%v> Expected PC to be %v, got %v", tci, cpuCopy.PC+1, cpu.pc)
		//}

		// check if the A register now contains the correct BCD representation of 10 (0x10)
		if cpu.a != tc.expected.A {
			t.Errorf("TC%v> Expected A to be 0x%02X, got 0x%02X", tci, tc.expected.A, cpu.a)
		}

		// check flags
		if cpu.getZFlag() != tc.expected.Z {
			t.Errorf("TC%v> Expected Z flag to be %t", tci, tc.expected.Z)
		}
		if cpu.getNFlag() != tc.expected.N {
			t.Errorf("TC%v> Expected N flag to be %t", tci, tc.expected.N)
		}
		if cpu.getHFlag() != tc.expected.H {
			t.Errorf("TC%v> Expected H flag to be %t", tci, tc.expected.H)
		}
		if cpu.getCFlag() != tc.expected.C {
			t.Errorf("TC%v> Expected C flag to be %t", tci, tc.expected.C)
		}
		postconditions()
	}
}

// DEC: should decrement the value of the destination (register or memory)
// opcodes:
//   - 0x3D=DEC A
//   - 0x05=DEC B
//   - 0x0D=DEC C
//   - 0x15=DEC D
//   - 0x1D=DEC E
//   - 0x25=DEC H
//   - 0x2D=DEC L
//   - 0x35=DEC [HL]
//
// flags: Z:Z N:1 H:H C:-
// - 0x0B=DEC BC
// - 0x1B=DEC DE
// - 0x2B=DEC HL
// - 0x3B=DEC SP
// flags: Z:- N:- H:- C:-
func TestDEC(t *testing.T) {
	t.Run("0x3D_DEC_A", test_0x3D_DEC_A)
	t.Run("0x05_DEC_B", test_0x05_DEC_B)
	t.Run("0x0D_DEC_C", test_0x0D_DEC_C)
	t.Run("0x15_DEC_D", test_0x15_DEC_D)
	t.Run("0x1D_DEC_E", test_0x1D_DEC_E)
	t.Run("0x25_DEC_H", test_0x25_DEC_H)
	t.Run("0x2D_DEC_L", test_0x2D_DEC_L)
	t.Run("0x35_DEC_HL", test_0x35_DEC_HL)
	t.Run("0x0B_DEC_BC", test_0x0B_DEC_BC)
	t.Run("0x1B_DEC_DE", test_0x1B_DEC_DE)
	t.Run("0x2B_DEC_HL", test_0x2B_DEC_HL)
	t.Run("0x3B_DEC_SP", test_0x3B_DEC_SP)
}

type TestData_DEC_8bit struct {
	initialValue  uint8
	expectedValue uint8
	Z             bool
	N             bool
	H             bool
	C             bool
}

var testData_DEC_8bit = []TestData_DEC_8bit{
	{
		initialValue:  0x00,
		expectedValue: 0xFF,
		Z:             false,
		N:             true,
		H:             true,
		C:             false,
	},
	{
		initialValue:  0x01,
		expectedValue: 0x00,
		Z:             true,
		N:             true,
		H:             false,
		C:             true,
	},
	{
		initialValue:  0x10,
		expectedValue: 0x0F,
		Z:             false,
		N:             true,
		H:             true,
		C:             false,
	},
	{
		initialValue:  0x0F,
		expectedValue: 0x0E,
		Z:             false,
		N:             true,
		H:             false,
		C:             true,
	},
	{
		initialValue:  0xF0,
		expectedValue: 0xEF,
		Z:             false,
		N:             true,
		H:             true,
		C:             false,
	},
	{
		initialValue:  0xFF,
		expectedValue: 0xFE,
		Z:             false,
		N:             true,
		H:             false,
		C:             true,
	},
}

func test_0x3D_DEC_A(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.a = data.initialValue
		testProgram := []uint8{0x3D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.a != data.expectedValue {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.a)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x3D_DEC_A] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x05_DEC_B(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.b = data.initialValue
		testProgram := []uint8{0x05, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.b != data.expectedValue {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected B to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.b)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x05_DEC_B] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x0D_DEC_C(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.c = data.initialValue
		testProgram := []uint8{0x0D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.c != data.expectedValue {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected C to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.c)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x0D_DEC_C] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x15_DEC_D(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.d = data.initialValue
		testProgram := []uint8{0x15, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.d != data.expectedValue {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected D to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.d)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x15_DEC_D] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x1D_DEC_E(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.e = data.initialValue
		testProgram := []uint8{0x1D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.e != data.expectedValue {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected E to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.e)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x1D_DEC_E] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x25_DEC_H(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.h = data.initialValue
		testProgram := []uint8{0x25, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.h != data.expectedValue {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected H to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.h)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x25_DEC_H] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2D_DEC_L(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.l = data.initialValue
		testProgram := []uint8{0x2D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.l != data.expectedValue {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected L to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.l)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x2D_DEC_L] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x35_DEC_HL(t *testing.T) {
	for idx, data := range testData_DEC_8bit {
		preconditions()
		cpu.f = 0x00
		if data.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.setHL(0x0002)
		testProgram := []uint8{0x35, 0x10, data.initialValue}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		valueAtHL := cpu.bus.Read(cpu.getHL())
		if valueAtHL != data.expectedValue {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, data.expectedValue, valueAtHL)
		}
		if cpu.getZFlag() != data.Z {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected Z flag to be %t", idx, data.Z)
		}
		if cpu.getNFlag() != data.N {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected N flag to be %t", idx, data.N)
		}
		if cpu.getHFlag() != data.H {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected H flag to be %t", idx, data.H)
		}
		if cpu.getCFlag() != data.C {
			t.Errorf("[test_0x35_DEC_HL] TC%v> Expected C flag to be false, got true", idx)
		}
		postconditions()
	}
}

type TestData_DEC_16bit struct {
	initialValue  uint16
	expectedValue uint16
}

var testData_DEC_16bit = []TestData_DEC_16bit{
	{
		initialValue:  0x0000,
		expectedValue: 0xFFFF,
	},
	{
		initialValue:  0x0001,
		expectedValue: 0x0000,
	},
	{
		initialValue:  0x0010,
		expectedValue: 0x000F,
	},
	{
		initialValue:  0x000F,
		expectedValue: 0x000E,
	},
	{
		initialValue:  0x00F0,
		expectedValue: 0x00EF,
	},
	{
		initialValue:  0x00FF,
		expectedValue: 0x00FE,
	},
	{
		initialValue:  0xFFFF,
		expectedValue: 0xFFFE,
	},
	{
		initialValue:  0xFF00,
		expectedValue: 0xFEFF,
	},
	{
		initialValue:  0x1234,
		expectedValue: 0x1233,
	},
}

func test_0x0B_DEC_BC(t *testing.T) {
	for idx, data := range testData_DEC_16bit {
		preconditions()
		// randomize the flags
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setBC(data.initialValue)
		testProgram := []uint8{0x0B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x0B_DEC_BC] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.getBC() != data.expectedValue {
			t.Errorf("[test_0x0B_DEC_BC] TC%v> Expected BC to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.getBC())
		}
		// flags should be untouched
		if cpu.f != saveFlags {
			t.Errorf("[test_0x0B_DEC_BC] TC%v> Expected flags to be untouched 0x%02X, got 0x%02X", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x1B_DEC_DE(t *testing.T) {
	for idx, data := range testData_DEC_16bit {
		preconditions()
		// randomize the flags
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setDE(data.initialValue)
		testProgram := []uint8{0x1B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x1B_DEC_DE] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.getDE() != data.expectedValue {
			t.Errorf("[test_0x1B_DEC_DE] TC%v> Expected DE to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.getDE())
		}
		// flags should be untouched
		if cpu.f != saveFlags {
			t.Errorf("[test_0x1B_DEC_DE] TC%v> Expected flags to be untouched 0x%02X, got 0x%02X", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x2B_DEC_HL(t *testing.T) {
	for idx, data := range testData_DEC_16bit {
		preconditions()
		// randomize the flags
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(data.initialValue)
		testProgram := []uint8{0x2B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x2B_DEC_HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.getHL() != data.expectedValue {
			t.Errorf("[test_0x2B_DEC_HL] TC%v> Expected HL to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.getHL())
		}
		// flags should be untouched
		if cpu.f != saveFlags {
			t.Errorf("[test_0x2B_DEC_HL] TC%v> Expected flags to be untouched 0x%02X, got 0x%02X", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}
func test_0x3B_DEC_SP(t *testing.T) {
	for idx, data := range testData_DEC_16bit {
		preconditions()
		// randomize the flags
		randomizeFlags()
		saveFlags := cpu.f
		cpu.sp = data.initialValue
		testProgram := []uint8{0x3B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x3B_DEC_SP] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.sp != data.expectedValue {
			t.Errorf("[test_0x3B_DEC_SP] TC%v> Expected SP to be 0x%02X, got 0x%02X", idx, data.expectedValue, cpu.sp)
		}
		// flags should be untouched
		if cpu.f != saveFlags {
			t.Errorf("[test_0x3B_DEC_SP] TC%v> Expected flags to be untouched 0x%02X, got 0x%02X", idx, saveFlags, cpu.f)
		}
		postconditions()
	}
}

// SUB: Subtract register/memory 8bit value from A register
// opcodes:
//   - 0xD6 = SUB A, n8
//   - 0x97 = SUB A, A
//   - 0x90 = SUB A, B
//   - 0x91 = SUB A, C
//   - 0x92 = SUB A, D
//   - 0x93 = SUB A, E
//   - 0x94 = SUB A, H
//   - 0x95 = SUB A, L
//   - 0x96 = SUB A, [HL]
//
// flags: Z->Z N->1 H->H C->C (except for 0x97 where Z->1 N->1 H->0 C->0)
func TestSUB(t *testing.T) {
	t.Run("0xD6_SUB_A_n8", test_0xD6_SUB_A_n8)
	t.Run("0x97_SUB_A_A", test_0x97_SUB_A_A)
	t.Run("0x90_SUB_A_B", test_0x90_SUB_A_B)
	t.Run("0x91_SUB_A_C", test_0x91_SUB_A_C)
	t.Run("0x92_SUB_A_D", test_0x92_SUB_A_D)
	t.Run("0x93_SUB_A_E", test_0x93_SUB_A_E)
	t.Run("0x94_SUB_A_H", test_0x94_SUB_A_H)
	t.Run("0x95_SUB_A_L", test_0x95_SUB_A_L)
	t.Run("0x96_SUB_A__HL", test_0x96_SUB_A__HL)
}

var testData_SUB_A = []uint8{0x99, 0x00, 0x0F, 0x10, 0x0A, 0xA0, 0x55, 0xAA, 0x00, 0x01, 0x10, 0x0F, 0xF0, 0xFF, 0x11, 0xBA}
var testData_SUB_8bit_operand = []uint8{0x00, 0x01, 0x10, 0x0F, 0xF0, 0xFF, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}

func test_0xD6_SUB_A_n8(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		n8 := testData_SUB_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xD6, n8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - n8
		if cpu.a != expectedA {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (n8&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (n8&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected H flag to be false", idx)
		}
		if data < n8 && !cpu.getCFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected C flag to be true", idx)
		} else if data >= n8 && cpu.getCFlag() {
			t.Errorf("[test_0xD6_SUB_A_n8] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x97_SUB_A_A(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0x97, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := uint8(0x00)
		if cpu.a != expectedA {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if !cpu.getZFlag() {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected Z flag to be true", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0x97_SUB_A_A] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x90_SUB_A_B(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.b = operand
		cpu.a = data
		testProgram := []uint8{0x90, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x90_SUB_A_B] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x91_SUB_A_C(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.c = operand
		cpu.a = data
		testProgram := []uint8{0x91, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x91_SUB_A_C] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x92_SUB_A_D(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.d = operand
		cpu.a = data
		testProgram := []uint8{0x92, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x92_SUB_A_D] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x93_SUB_A_E(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.e = operand
		cpu.a = data
		testProgram := []uint8{0x93, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x93_SUB_A_E] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x94_SUB_A_H(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.h = operand
		cpu.a = data
		testProgram := []uint8{0x94, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x94_SUB_A_H] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x95_SUB_A_L(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.l = operand
		cpu.a = data
		testProgram := []uint8{0x95, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x95_SUB_A_L] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0x96_SUB_A__HL(t *testing.T) {
	for idx, data := range testData_SUB_A {
		preconditions()
		randomizeFlags()
		operand := testData_SUB_8bit_operand[idx]
		cpu.setHL(0x0002)
		cpu.a = data
		testProgram := []uint8{0x96, 0x10, operand}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value A - operand
		expectedA := data - operand
		if cpu.a != expectedA {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if expectedA == 0 && !cpu.getZFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected Z flag to be true", idx)
		} else if expectedA != 0 && cpu.getZFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected N flag to be true", idx)
		}
		if (data&0x0F) < (operand&0x0F) && !cpu.getHFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected H flag to be true", idx)
		} else if (data&0x0F) >= (operand&0x0F) && cpu.getHFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected H flag to be false", idx)
		}
		if data < operand && !cpu.getCFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected C flag to be true", idx)
		} else if data >= operand && cpu.getCFlag() {
			t.Errorf("[test_0x96_SUB_A__HL] TC%v> Expected C flag to be false", idx)
		}
	}
}

// SBC: Subtract register/memory 8bit plus carry flag from A register
// opcodes:
//   - 0xDE = SBC A, n8
//   - 0x9F = SBC A, A
//   - 0x98 = SBC A, B
//   - 0x99 = SBC A, C
//   - 0x9A = SBC A, D
//   - 0x9B = SBC A, E
//   - 0x9C = SBC A, H
//   - 0x9D = SBC A, L
//   - 0x9E = SBC A, [HL]
func TestSBC(t *testing.T) {
	t.Run("0xDE_SBC_A_n8", test_0xDE_SBC_A_n8)
	t.Run("0x9F_SBC_A_A", test_0x9F_SBC_A_A)
	t.Run("0x98_SBC_A_B", test_0x98_SBC_A_B)
	t.Run("0x99_SBC_A_C", test_0x99_SBC_A_C)
	t.Run("0x9A_SBC_A_D", test_0x9A_SBC_A_D)
	t.Run("0x9B_SBC_A_E", test_0x9B_SBC_A_E)
	t.Run("0x9C_SBC_A_H", test_0x9C_SBC_A_H)
	t.Run("0x9D_SBC_A_L", test_0x9D_SBC_A_L)
	t.Run("0x9E_SBC_A__HL", test_0x9E_SBC_A__HL)
}

var testData_SBC_A = []uint8{0x99, 0x00, 0x0F, 0x10, 0x0A, 0xA0, 0x55, 0xAA, 0x00, 0x01, 0x10, 0x0F, 0xF0, 0xFF, 0x11, 0xBA}
var testData_SBC_8bit_operand = []uint8{0x98, 0x87, 0x1A, 0xDE, 0x35, 0x0A, 0x7F, 0xDD, 0xFF, 0x00, 0x10, 0xF0, 0x0E, 0x01, 0x11, 0x22}
var testData_SBC_carry = []bool{true, true, false, true, false, true, false, true, false, true, false, true, false, true, false, true}

var expected_A_register = []uint8{0x00, 0x78, 0xF5, 0x31, 0xD5, 0x95, 0xD6, 0xCC, 0x01, 0x00, 0x00, 0x1E, 0xE2, 0xFD, 0x00, 0x97}
var expected_Z_flag = []bool{true, false, false, false, false, false, false, false, false, true, true, false, false, false, true, false}
var expected_H_flag = []bool{false, true, false, true, false, true, true, true, true, false, false, false, true, false, false, false}
var expected_C_flag = []bool{false, true, true, true, true, false, true, true, true, false, false, true, false, false, false, false}

func test_0xDE_SBC_A_n8(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		n8 := testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xDE, n8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0xDE_SBC_A_n8] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9F_SBC_A_A(t *testing.T) {
	expected_SBC_H_flag := []bool{true, true, false, true, false, true, false, true, false, true, false, false, false, false, false, true}
	for idx, data := range testData_SBC_A {
		preconditions()
		var expectedA uint8
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
			expectedA = 0xFF
		} else {
			cpu.resetCFlag()
			expectedA = 0x00
		}
		saveCFlag := cpu.getCFlag()
		cpu.a = data
		testProgram := []uint8{0x9F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expectedA {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedA, cpu.a)
		}
		// check flags
		if (cpu.getZFlag() && expectedA != 0) || (!cpu.getZFlag() && expectedA == 0) {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected Z flag to be false", idx)
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_SBC_H_flag[idx] {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected H flag to be %t, got %t", idx, expected_SBC_H_flag[idx], cpu.getHFlag())
		}
		// c flag is not affected
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("[test_0x9F_SBC_A_A] TC%v> Expected C flag to be unaltered", idx)
		}
		postconditions()
	}
}
func test_0x98_SBC_A_B(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.b = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x98, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x98_SBC_A_B] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x99_SBC_A_C(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.c = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x99, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x99_SBC_A_C] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9A_SBC_A_D(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.d = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x9A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9B_SBC_A_E(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.e = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x9B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x9A_SBC_A_D] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9C_SBC_A_H(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.h = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x9C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x9C_SBC_A_H] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9D_SBC_A_L(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.l = testData_SBC_8bit_operand[idx]
		cpu.a = data
		testProgram := []uint8{0x9D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x9D_SBC_A_L] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}
func test_0x9E_SBC_A__HL(t *testing.T) {
	for idx, data := range testData_SBC_A {
		preconditions()
		randomizeFlags()
		if testData_SBC_carry[idx] {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.setHL(0x0002)
		cpu.a = data
		testProgram := []uint8{0x9E, 0x10, testData_SBC_8bit_operand[idx]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		if cpu.a != expected_A_register[idx] {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expected_A_register[idx], cpu.a)
		}
		// check flags
		if cpu.getZFlag() != expected_Z_flag[idx] {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected Z flag to be %v", idx, expected_Z_flag[idx])
		}
		if !cpu.getNFlag() {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected N flag to be true", idx)
		}
		if cpu.getHFlag() != expected_H_flag[idx] {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected H flag to be %v", idx, expected_H_flag[idx])
		}
		if cpu.getCFlag() != expected_C_flag[idx] {
			t.Errorf("[test_0x9E_SBC_A__HL] TC%v> Expected C flag to be %v", idx, expected_C_flag[idx])
		}
	}
}

// SCF: should set the carry flag
// 0x37: SCF
// Set Carry Flag: set the carry flag to 1 and reset the N and H flags
func TestSCF(t *testing.T) {
	preconditions()

	// set all flags to see if they are all correctly reset and have C flag set to be rotated in register A
	cpu.f = 0x00
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// reset the carry flag
	cpu.setCFlag()

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x37, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// check if the Z flag was left untouched (1)
	if finalState.Z != true {
		t.Errorf("Expected Z flag to be 1")
	}

	// check if the N flag was reset
	if finalState.N != false {
		t.Errorf("Expected N flag to be 0")
	}

	// check if the H flag was reset
	if finalState.H != false {
		t.Errorf("Expected H flag to be 0")
	}

	// check if the C flag was set
	if finalState.C != true {
		t.Errorf("Expected C flag to be 1")
	}

	postconditions()
}

// Stores into A register the result of the bitwise OR operation between A and the operand
// opcodes:
//   - 0xB0 = OR A, B
//   - 0xB1 = OR A, C
//   - 0xB2 = OR A, D
//   - 0xB3 = OR A, E
//   - 0xB4 = OR A, H
//   - 0xB5 = OR A, L
//   - 0xB6 = OR A, [HL]
//   - 0xB7 = OR A, A
//   - 0xF6 = OR A, n8
//
// flags: Z:Z N:0 H:0 C:0
func TestOR(t *testing.T) {
	t.Run("0xB0_OR_A_B", test_0xB0_OR_A_B)
	t.Run("0xB1_OR_A_C", test_0xB1_OR_A_C)
	t.Run("0xB2_OR_A_D", test_0xB2_OR_A_D)
	t.Run("0xB3_OR_A_E", test_0xB3_OR_A_E)
	t.Run("0xB4_OR_A_H", test_0xB4_OR_A_H)
	t.Run("0xB5_OR_A_L", test_0xB5_OR_A_L)
	t.Run("0xB6_OR_A_HL", test_0xB6_OR_A_HL)
	t.Run("0xB7_OR_A_A", test_0xB7_OR_A_A)
	t.Run("0xF6_OR_A_n8", test_0xF6_OR_A_n8)
}

var testData_OR_A = []uint8{0x00, 0xFF, 0x0F, 0xF0, 0x0A, 0xA0, 0x55, 0xAA}
var testData_OR_operand = []uint8{0xFF, 0x87, 0x1A, 0xDE, 0x35, 0x0A, 0x7F, 0xDD}

func test_0xF6_OR_A_n8(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		n8 := testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xF6, n8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | n8
		if cpu.a != expectedValue {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xF6_OR_A_n8] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB7_OR_A_A(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.a = data
		testProgram := []uint8{0xB7, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | data
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB7_OR_A_A] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB0_OR_A_B(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.b = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB0, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.b
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the B register was unaffected
		if cpu.b != testData_OR_operand[idx] {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected B register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.b)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB0_OR_A_B] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB1_OR_A_C(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.c = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB1, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.c
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the C register was unaffected
		if cpu.c != testData_OR_operand[idx] {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected C register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.c)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB1_OR_A_C] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB2_OR_A_D(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.d = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.d
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the D register was unaffected
		if cpu.d != testData_OR_operand[idx] {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected D register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.d)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB2_OR_A_D] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB3_OR_A_E(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.e = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB3, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.e
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the E register was unaffected
		if cpu.e != testData_OR_operand[idx] {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected E register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.e)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB3_OR_A_E] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB4_OR_A_H(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.h = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB4, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.h
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the H register was unaffected
		if cpu.h != testData_OR_operand[idx] {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected H register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.h)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB4_OR_A_H] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB5_OR_A_L(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.l = testData_OR_operand[idx]
		cpu.a = data
		testProgram := []uint8{0xB5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | cpu.l
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the L register was unaffected
		if cpu.l != testData_OR_operand[idx] {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected L register to be 0x%02X, got 0x%02X", idx, testData_OR_operand[idx], cpu.l)
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB5_OR_A_L] TC%v> Expected C flag to be false", idx)
		}
	}
}
func test_0xB6_OR_A_HL(t *testing.T) {
	for idx, data := range testData_OR_A {
		preconditions()
		randomizeFlags()
		cpu.setHL(0x0002)
		cpu.a = data
		testProgram := []uint8{0xB6, 0x10, testData_OR_operand[idx]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check if the program stopped at the right place
		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		// check if the A register now contains the correct value
		expectedValue := data | testData_OR_operand[idx]
		if cpu.a != expectedValue {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected A register to be 0x%02X, got 0x%02X", idx, expectedValue, cpu.a)
		}
		// check if the HL register was unaffected
		if cpu.getHL() != 0x0002 {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected HL register to be 0x%02X, got 0x%02X", idx, 0x0002, cpu.getHL())
		}
		// check if the flags were set correctly
		if cpu.getZFlag() != (expectedValue == 0x00) {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected Z flag to be %t", idx, expectedValue == 0x00)
		}
		if cpu.getNFlag() {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected N flag to be false", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected H flag to be false", idx)
		}
		if cpu.getCFlag() {
			t.Errorf("[test_0xB6_OR_A_HL] TC%v> Expected C flag to be false", idx)
		}
	}
}

// XOR: should perform a bitwise XOR between the source and the destination
// opcodes
//   - 0xA8: XOR A, B
//   - 0xA9: XOR A, C
//   - 0xAA: XOR A, D
//   - 0xAB: XOR A, E
//   - 0xAC: XOR A, H
//   - 0xAD: XOR A, L
//   - 0xAE: XOR A, (HL)
//   - 0xAF: XOR A, A
//   - 0xEE: XOR A, n8
//
// flags: Z:Z N:0 H:0 C:0
func TestXOR(t *testing.T) {
	t.Run("0xA8_XOR_A_B", test_0xA8_XOR_A_B)
	t.Run("0xA9_XOR_A_C", test_0xA9_XOR_A_C)
	t.Run("0xAA_XOR_A_D", test_0xAA_XOR_A_D)
	t.Run("0xAB_XOR_A_E", test_0xAB_XOR_A_E)
	t.Run("0xAC_XOR_A_H", test_0xAC_XOR_A_H)
	t.Run("0xAD_XOR_A_L", test_0xAD_XOR_A_L)
	t.Run("0xAE_XOR_A_HL", test_0xAE_XOR_A_HL)
	t.Run("0xAF_XOR_A_A", test_0xAF_XOR_A_A)
	t.Run("0xEE_XOR_A_n8", test_0xEE_XOR_A_n8)
}
func test_0xA8_XOR_A_B(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.b = data[1]
		expected := data[2]
		testProgram := []uint8{0xA8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA8_XOR_A_B] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xA8_XOR_A_B] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xA8_XOR_A_B] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xA8_XOR_A_B] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xA9_XOR_A_C(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.c = data[1]
		expected := data[2]
		testProgram := []uint8{0xA9, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xA9_XOR_A_C] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xA9_XOR_A_C] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xA9_XOR_A_C] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xA9_XOR_A_C] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAA_XOR_A_D(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.d = data[1]
		expected := data[2]
		testProgram := []uint8{0xAA, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAA_XOR_A_D] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xAA_XOR_A_D] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xAA_XOR_A_D] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xAA_XOR_A_D] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAB_XOR_A_E(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.e = data[1]
		expected := data[2]
		testProgram := []uint8{0xAB, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAB_XOR_A_E] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xAB_XOR_A_E] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xAB_XOR_A_E] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xAB_XOR_A_E] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAC_XOR_A_H(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.h = data[1]
		expected := data[2]
		testProgram := []uint8{0xAC, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAC_XOR_A_H] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xAC_XOR_A_H] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xAC_XOR_A_H] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xAC_XOR_A_H] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAD_XOR_A_L(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.l = data[1]
		expected := data[2]
		testProgram := []uint8{0xAD, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAD_XOR_A_L] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xAD_XOR_A_L] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xAD_XOR_A_L] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xAD_XOR_A_L] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAE_XOR_A_HL(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		cpu.setHL(0x0002)
		expected := data[2]
		testProgram := []uint8{0xAE, 0x10, data[1]}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAE_XOR_A_HL] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != expected {
			t.Errorf("[test_0xAE_XOR_A_HL] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xAE_XOR_A_HL] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xAE_XOR_A_HL] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}
func test_0xAF_XOR_A_A(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		expected := data[2]
		testProgram := []uint8{0xAF, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0001 {
			t.Errorf("[test_0xAF_XOR_A_A] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}

		if cpu.a != 0x00 {
			t.Errorf("[test_0xAF_XOR_A_A] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if cpu.f != 0x80 {
			t.Errorf("[test_0xAF_XOR_A_A] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
		}
	}
}
func test_0xEE_XOR_A_n8(t *testing.T) {
	testData := [][3]uint8{
		{0b00000000, 0b00000000, 0b00000000},
		{0b11111111, 0b11111111, 0b00000000},
		{0b01010101, 0b10101010, 0b11111111},
		{0b10101010, 0b01010101, 0b11111111},
		{0b00001111, 0b11110000, 0b11111111},
		{0b00001111, 0b00001111, 0b00000000},
		{0b11110000, 0b00001111, 0b11111111},
		{0b11110000, 0b11110000, 0b00000000},
		{0b00110011, 0b11001100, 0b11111111},
		{0b11001100, 0b00110011, 0b11111111},
		{0b00110011, 0b00110011, 0b00000000},
		{0b11001100, 0b11001100, 0b00000000},
		{0b01101001, 0b11110100, 0b10011101},
	}

	for idx, data := range testData {
		preconditions()
		cpu.f = 0xF0
		cpu.a = data[0]
		expected := data[2]
		testProgram := []uint8{0xEE, data[1], 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		if cpu.pc != 0x0002 {
			t.Errorf("[test_0xEE_XOR_A_n8] TC%v> Expected PC to be 0x0001, got 0x%04X", idx, cpu.pc)
		}
		if cpu.a != expected {
			t.Errorf("[test_0xEE_XOR_A_n8] TC%v> Expected A to be 0x%02X, got 0x%02X", idx, expected, cpu.a)
		}
		if expected == 0x00 {
			if cpu.f != 0x80 {
				t.Errorf("[test_0xEE_XOR_A_n8] TC%v> Expected flags to be reset except Z flag, got 0x%02X", idx, cpu.f)
			}
		} else {
			if cpu.f != 0x00 {
				t.Errorf("[test_0xEE_XOR_A_n8] TC%v> Expected flags to be reset, got 0x%02X", idx, cpu.f)
			}
		}
	}
}

// RLA: should rotate the destination left through the carry
// 0x17: RLA rotates the A register to the left through the carry flag
// The 0th bit of the A register should be set to the value of the carry flag
// The carry flag should be set to the value of the 7th bit of the A register
// All other flags should be reset
func TestRLA(t *testing.T) {
	preconditions()

	// set all flags to see if they are all correctly reset and have C flag set to be rotated in register A
	cpu.f = 0x00

	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()
	cpu.setCFlag()

	// set the A register to b 0101 0101
	cpu.a = 0x55

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x17, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	fmt.Println("Z:", finalState.Z, "N:", finalState.N, "H:", finalState.H, "C:", finalState.C)

	// check if the A register was rotated left giving 0xAB
	if finalState.A != 0xAB {
		t.Errorf("[RLA_TC30_CHK_1] Error> RLA instruction: the A register should have been set to 0xAB, got 0x%02X \n", finalState.A)
	}

	// check if the flags are set correctly
	if finalState.F != 0x00 {
		t.Errorf("[RLA_TC30_CHK_2] Error> RLA instruction: the flags should have been reset, got 0x%02X \n", finalState.F)
	}

	postconditions()
}

// RLCA: should rotate the destination left
// 0x07: RLCA rotates the A register to the left
// The carry flag should be set to the value of the 7th bit of the A register
// All other flags should be reset
func TestRLCA(t *testing.T) {
	preconditions()

	// set all flags to 1 to see if they are reset
	cpu.f = 0x00
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// reset the carry flag
	cpu.resetCFlag()

	// set the A register to b10101010
	cpu.a = 0xAA

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x07, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	if finalState.A != 0x55 {
		t.Errorf("Expected A to be 0x55, got 0x%02X", finalState.A)
	}

	// check if the C flag was set
	if cpu.getCFlag() != true {
		t.Errorf("Expected C flag to be set")
	}

	// check if the Z, N and H flags were reset
	if cpu.getZFlag() != false {
		t.Errorf("Expected Z flag to be reset")
	}
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be reset")
	}
	if cpu.getHFlag() != false {
		t.Errorf("Expected H flag to be reset")
	}

	postconditions()
}

// RRA: should rotate the destination right through the carry
// 0x1F: RRA rotates the A register to the right through the carry flag
// The 7th bit of the A register should be set to the value of the carry flag
// The carry flag should be set to the value of the 0th bit of the A register
// All other flags should be reset
func TestRRA(t *testing.T) {
	preconditions()

	// set all flags to 1 to see if they are reset
	cpu.f = 0x00
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()
	cpu.setCFlag()

	// set the A register to b10101010
	cpu.a = 0xAA

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x1F, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	// 1010 1010 => 1101 0101 (0xD5)
	if finalState.A != 0xD5 {
		t.Errorf("Expected A to be 0xD5, got 0x%02X", finalState.A)
	}

	// check if the carry flag was reset
	if cpu.getCFlag() != false {
		t.Errorf("Expected C flag to be reset")
	}

	// check if the Z, N and H flags were reset
	if cpu.getZFlag() != false {
		t.Errorf("Expected Z flag to be reset")
	}
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be reset")
	}
	if cpu.getHFlag() != false {
		t.Errorf("Expected H flag to be reset")
	}
}

// RRCA: should rotate the destination right
// 0x0F: RRCA rotates the A register to the right
// The carry flag should be set to the value of the 0th bit of the A register
// All other flags should be reset
func TestRRCA(t *testing.T) {
	preconditions()

	// set all flags to 1 to see if they are reset
	cpu.f = 0x00
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// reset the carry flag
	cpu.resetCFlag()

	// set the A register to b10101010
	cpu.a = 0x55

	// load the program into the memory
	testData := []uint8{0x00, 0x00, 0x00, 0x0F, 0x00, 0x10}
	loadProgramIntoMemory(memory1, testData)

	// run the program
	cpu.Run()

	// check the final state of the cpu
	finalState := getCpuState()

	if finalState.A != 0xAA {
		t.Errorf("Expected A to be 0xAA, got 0x%02X", finalState.A)
	}

	// check if the carry flag was set
	if cpu.getCFlag() != true {
		t.Errorf("Expected C flag to be set")
	}

	// check if the Z, N and H flags were reset
	if cpu.getZFlag() != false {
		t.Errorf("Expected Z flag to be reset")
	}
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be reset")
	}
	if cpu.getHFlag() != false {
		t.Errorf("Expected H flag to be reset")
	}

	postconditions()
}
