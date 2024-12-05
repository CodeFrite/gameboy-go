package gameboy

import (
	"testing"
)

// TEST CASES

// Rotate r8 register left and save bit 7 to the Carry flag
// opcodes:
//   - 0x00 =	RLC B
//   - 0x01 =	RLC C
//   - 0x02 =	RLC D
//   - 0x03 =	RLC E
//   - 0x04 =	RLC H
//   - 0x05 =	RLC L
//   - 0x06 =	RLC [HL]
//   - 0x07 =	RLC A
//
// flags: Z=Z N=0 H=0 C=C
func TestRLC(t *testing.T) {
	t.Run("0x00 RLC B", test_0x00_RLC_B)
	t.Run("0x01 RLC C", test_0x01_RLC_C)
	t.Run("0x02 RLC D", test_0x02_RLC_D)
	t.Run("0x03 RLC E", test_0x03_RLC_E)
	t.Run("0x04 RLC H", test_0x04_RLC_H)
	t.Run("0x05 RLC L", test_0x05_RLC_L)
	t.Run("0x06 RLC [HL]", test_0x06_RLC__HL)
	t.Run("0x07 RLC A", test_0x07_RLC_A)
}

type TestCaseRLC struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesRLC = []TestCaseRLC{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000000, false},  // 1
	{0b11111111, false, 0b11111111, true},  // 2
	{0b11111111, true, 0b11111111, true},   // 3
	{0b10101010, false, 0b01010101, true},  // 4
	{0b10101010, true, 0b01010101, true},   // 5
	{0b01010101, true, 0b10101010, false},  // 6
	{0b00110011, false, 0b01100110, false}, // 7
	{0b00110011, true, 0b01100110, false},  // 8
	{0b11001100, false, 0b10011001, true},  // 9
	{0b11001100, true, 0b10011001, true},   // 10
	{0b00010001, false, 0b00100010, false}, // 11
	{0b00010001, true, 0b00100010, false},  // 12
	{0b10001000, false, 0b00010001, true},  // 13
	{0b10001000, true, 0b00010001, true},   // 14
}

func test_0x00_RLC_B(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x00, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x01_RLC_C(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x01, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x02_RLC_D(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x02, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x03_RLC_E(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x03, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x04_RLC_H(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x04, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x05_RLC_L(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x05, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x06_RLC__HL(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x06, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value @[HL] to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x07_RLC_A(t *testing.T) {
	for idx, tc := range testCasesRLC {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x07, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}

// Rotate r8 register right and save bit 0 to the Carry flag
// opcodes:
//   - 0x08 =	RRC B
//   - 0x09 =	RRC C
//   - 0x0A =	RRC D
//   - 0x0B =	RRC E
//   - 0x0C =	RRC H
//   - 0x0D =	RRC L
//   - 0x0E =	RRC [HL]
//   - 0x0F =	RRC A
//
// flags: Z=Z N=0 H=0 C=C
func TestRRC(t *testing.T) {
	t.Run("0x08 RRC B", test_0x08_RRC_B)
	t.Run("0x09 RRC C", test_0x09_RRC_C)
	t.Run("0x0A RRC D", test_0x0A_RRC_D)
	t.Run("0x0B RRC E", test_0x0B_RRC_E)
	t.Run("0x0C RRC H", test_0x0C_RRC_H)
	t.Run("0x0D RRC L", test_0x0D_RRC_L)
	t.Run("0x0E RRC [HL]", test_0x0E_RRC__HL)
	t.Run("0x0F RRC A", test_0x0F_RRC_A)
}

type TestCaseRRC struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesRRC = []TestCaseRRC{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000000, false},  // 1
	{0b11111111, false, 0b11111111, true},  // 2
	{0b11111111, true, 0b11111111, true},   // 3
	{0b10101010, false, 0b01010101, false}, // 4
	{0b10101010, true, 0b01010101, false},  // 5
	{0b01010101, true, 0b10101010, true},   // 6
	{0b00110011, false, 0b10011001, true},  // 7
	{0b00110011, true, 0b10011001, true},   // 8
	{0b11001100, false, 0b01100110, false}, // 9
	{0b11001100, true, 0b01100110, false},  // 10
	{0b00010001, false, 0b10001000, true},  // 11
	{0b00010001, true, 0b10001000, true},   // 12
	{0b10001000, false, 0b01000100, false}, // 13
	{0b10001000, true, 0b01000100, false},  // 14
}

func test_0x08_RRC_B(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x08, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x09_RRC_C(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x09, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0A_RRC_D(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0B_RRC_E(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0C_RRC_H(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0D_RRC_L(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0E_RRC__HL(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0E, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value @[HL] to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x0F_RRC_A(t *testing.T) {
	for idx, tc := range testCasesRRC {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x0F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}

// RL r8 / [HL]
// Rotate r8 or [HL] left through carry: old bit 7 to Carry flag, new bit 0 to bit 7.
// opcodes:
//   - 0x10:	RL B
//   - 0x11:	RL C
//   - 0x12:	RL D
//   - 0x13:	RL E
//   - 0x14:	RL H
//   - 0x15:	RL L
//   - 0x16:	RL [HL]
//   - 0x17:	RL A
//
// flags: Z=Z N=0 H=0 C=C
func TestRL(t *testing.T) {
	t.Run("0x10 RL B", test_0x10_RL_B)
	t.Run("0x11 RL C", test_0x11_RL_C)
	t.Run("0x12 RL D", test_0x12_RL_D)
	t.Run("0x13 RL E", test_0x13_RL_E)
	t.Run("0x14 RL H", test_0x14_RL_H)
	t.Run("0x15 RL L", test_0x15_RL_L)
	t.Run("0x16 RL [HL]", test_0x16_RL__HL)
	t.Run("0x17 RL A", test_0x17_RL_A)
}

type TestCaseRL struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesRL = []TestCaseRL{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000001, false},  // 1
	{0b11111111, false, 0b11111110, true},  // 2
	{0b11111111, true, 0b11111111, true},   // 3
	{0b10101010, false, 0b01010100, true},  // 4
	{0b10101010, true, 0b01010101, true},   // 5
	{0b01010101, true, 0b10101011, false},  // 6
	{0b00110011, false, 0b01100110, false}, // 7
	{0b00110011, true, 0b01100111, false},  // 8
	{0b11001100, false, 0b10011000, true},  // 9
	{0b11001100, true, 0b10011001, true},   // 10
	{0b00010001, false, 0b00100010, false}, // 11
	{0b00010001, true, 0b00100011, false},  // 12
	{0b10001000, false, 0b00010000, true},  // 13
	{0b10001000, true, 0b00010001, true},   // 14
}

func test_0x10_RL_B(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x10, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x11_RL_C(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x11, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.c)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x12_RL_D(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x12, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.d)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x13_RL_E(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x13, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.e)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x14_RL_H(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x14, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.h)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x15_RL_L(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x15, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.l)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x16_RL__HL(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x16, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value @[HL] to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x17_RL_A(t *testing.T) {
	for idx, tc := range testCasesRL {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x17, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}

// RR r8 / [HL]
// Rotate r8 or [HL] right through carry: old bit 0 to Carry flag, new bit 7 from carry flag
// opcodes:
//   - 0x18:	RR B
//   - 0x19:	RR C
//   - 0x1A:	RR D
//   - 0x1B:	RR E
//   - 0x1C:	RR H
//   - 0x1D:	RR L
//   - 0x1E:	RR [HL]
//   - 0x1F:	RR A
//
// flags: Z=Z N=0 H=0 C=C
func TestRR(t *testing.T) {
	t.Run("0x18 RR B", test_0x18_RR_B)
	t.Run("0x19 RR C", test_0x19_RR_C)
	t.Run("0x1A RR D", test_0x1A_RR_D)
	t.Run("0x1B RR E", test_0x1B_RR_E)
	t.Run("0x1C RR H", test_0x1C_RR_H)
	t.Run("0x1D RR L", test_0x1D_RR_L)
	t.Run("0x1E RR [HL]", test_0x1E_RR__HL)
	t.Run("0x1F RR A", test_0x1F_RR_A)
}

type TestCaseRR struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesRR = []TestCaseRR{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b10000000, false},  // 1
	{0b11111111, false, 0b01111111, true},  // 2
	{0b11111111, true, 0b11111111, true},   // 3
	{0b10101010, false, 0b01010101, false}, // 4
	{0b10101010, true, 0b11010101, false},  // 5
	{0b01010101, true, 0b10101010, true},   // 6
	{0b00110011, false, 0b00011001, true},  // 7
	{0b00110011, true, 0b10011001, true},   // 8
	{0b11001100, false, 0b01100110, false}, // 9
	{0b11001100, true, 0b11100110, false},  // 10
	{0b00010001, false, 0b00001000, true},  // 11
	{0b00010001, true, 0b10001000, true},   // 12
	{0b10001000, false, 0b01000100, false}, // 13
	{0b10001000, true, 0b11000100, false},  // 14
}

func test_0x18_RR_B(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x18, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x19_RR_C(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x19, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1A_RR_D(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1B_RR_E(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1C_RR_H(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1D_RR_L(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1E_RR__HL(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1E, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value @[HL] to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x1F_RR_A(t *testing.T) {
	for idx, tc := range testCasesRR {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x1F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was rotated correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		postconditions()
	}
}

// Shift Left Arithmetically register r8.
// Carry ←╂─ b7 ← ... ← b0 ←╂─ 0
// opcodes:
//   - 0x20:	SLA B
//   - 0x21:	SLA C
//   - 0x22:	SLA D
//   - 0x23:	SLA E
//   - 0x24:	SLA H
//   - 0x25:	SLA L
//   - 0x26:	SLA [HL]
//   - 0x27:	SLA A
//
// flags: Z=Z N=0 H=0 C=b7
func TestSLA(t *testing.T) {
	t.Run("0x20 SLA B", test_0x20_SLA_B)
	t.Run("0x21 SLA C", test_0x21_SLA_C)
	t.Run("0x22 SLA D", test_0x22_SLA_D)
	t.Run("0x23 SLA E", test_0x23_SLA_E)
	t.Run("0x24 SLA H", test_0x24_SLA_H)
	t.Run("0x25 SLA L", test_0x25_SLA_L)
	t.Run("0x26 SLA [HL]", test_0x26_SLA__HL)
	t.Run("0x27 SLA A", test_0x27_SLA_A)
}

type TestCaseSLA struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesSLA = []TestCaseSLA{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000000, false},  // 1
	{0b11111111, false, 0b11111110, true},  // 2
	{0b11111111, true, 0b11111110, true},   // 3
	{0b10101010, false, 0b01010100, true},  // 4
	{0b10101010, true, 0b01010100, true},   // 5
	{0b01010101, true, 0b10101010, false},  // 6
	{0b00110011, false, 0b01100110, false}, // 7
	{0b00110011, true, 0b01100110, false},  // 8
	{0b11001100, false, 0b10011000, true},  // 9
	{0b11001100, true, 0b10011000, true},   // 10
	{0b00010001, false, 0b00100010, false}, // 11
	{0b00010001, true, 0b00100010, false},  // 12
	{0b10001000, false, 0b00010000, true},  // 13
	{0b10001000, true, 0b00010000, true},   // 14
}

func test_0x20_SLA_B(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x20, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x21_SLA_C(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x21, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.c)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x22_SLA_D(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x22, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.d)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x23_SLA_E(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x23, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected E register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.e)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x24_SLA_H(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x24, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.h)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x25_SLA_L(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x25, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.l)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x26_SLA__HL(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x26, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value pointed by [HL] register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x27_SLA_A(t *testing.T) {
	for idx, tc := range testCasesSLA {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x27, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}

// Shift Right Arithmetically register r8 (bit 7 of r8 is unchanged).
// ┃ b7 → ... → b0 ─╂→ Carry
// opcodes:
//   - 0x28:	SRA B
//   - 0x29:	SRA C
//   - 0x2A:	SRA D
//   - 0x2B:	SRA E
//   - 0x2C:	SRA H
//   - 0x2D:	SRA L
//   - 0x2E:	SRA [HL]
//   - 0x2F:	SRA A
//
// flags: Z=Z N=0 H=0 C=b0
func TestSRA(t *testing.T) {
	t.Run("0x28 SRA B", test_0x28_SRA_B)
	t.Run("0x29 SRA C", test_0x29_SRA_C)
	t.Run("0x2A SRA D", test_0x2A_SRA_D)
	t.Run("0x2B SRA E", test_0x2B_SRA_E)
	t.Run("0x2C SRA H", test_0x2C_SRA_H)
	t.Run("0x2D SRA L", test_0x2D_SRA_L)
	t.Run("0x2E SRA [HL]", test_0x2E_SRA__HL)
	t.Run("0x2F SRA A", test_0x2F_SRA_A)
}

type TestCaseSRA struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCaseSRA = []TestCaseSRA{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000000, false},  // 1
	{0b11111111, false, 0b11111111, true},  // 2
	{0b11111111, true, 0b11111111, true},   // 3
	{0b10101010, false, 0b11010101, false}, // 4
	{0b10101010, true, 0b11010101, false},  // 5
	{0b01010101, true, 0b00101010, true},   // 6
	{0b00110011, false, 0b00011001, true},  // 7
	{0b00110011, true, 0b00011001, true},   // 8
	{0b11001100, false, 0b11100110, false}, // 9
	{0b11001100, true, 0b11100110, false},  // 10
	{0b00010001, false, 0b00001000, true},  // 11
	{0b00010001, true, 0b00001000, true},   // 12
	{0b10001000, false, 0b11000100, false}, // 13
	{0b10001000, true, 0b11000100, false},  // 14
}

func test_0x28_SRA_B(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.b = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x28, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x29_SRA_C(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.c = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x29, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.c)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2A_SRA_D(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.d = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.d)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2B_SRA_E(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.e = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected E register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.e)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2C_SRA_H(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.h = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.h)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2D_SRA_L(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.l = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.l)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2E_SRA__HL(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.setHL(0x0003)
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2E, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value pointed by [HL] register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x2F_SRA_A(t *testing.T) {
	for idx, tc := range testCaseSRA {
		preconditions()
		cpu.a = tc.value
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		testProgram := []uint8{0xCB, 0x2F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()

		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if cpu.getZFlag() != (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be %t, got %t", idx, tc.expectedValue == 0, cpu.getZFlag())
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}

// Swap upper & lower nibles of r8 register
// opcodes:
//   - 0x30:	SWAP B
//   - 0x31:	SWAP C
//   - 0x32:	SWAP D
//   - 0x33:	SWAP E
//   - 0x34:	SWAP H
//   - 0x35:	SWAP L
//   - 0x36:	SWAP [HL]
//   - 0x37:	SWAP A
//
// flags: Z=Z N=0 H=0 C=0
func TestSWAP(t *testing.T) {
	t.Run("0x30 SWAP B", test_0x30_SWAP_B)
	t.Run("0x31 SWAP C", test_0x31_SWAP_C)
	t.Run("0x32 SWAP D", test_0x32_SWAP_D)
	t.Run("0x33 SWAP E", test_0x33_SWAP_E)
	t.Run("0x34 SWAP H", test_0x34_SWAP_H)
	t.Run("0x35 SWAP L", test_0x35_SWAP_L)
	t.Run("0x36 SWAP [HL]", test_0x36_SWAP__HL)
	t.Run("0x37 SWAP A", test_0x37_SWAP_A)
}

type TestCaseSWAP struct {
	value         uint8
	expectedValue uint8
}

var testCaseSWAP = []TestCaseSWAP{
	{0b00000000, 0b00000000}, // 0
	{0b00000001, 0b00010000}, // 1
	{0b00000010, 0b00100000}, // 2
	{0b00000011, 0b00110000}, // 3
	{0b00000100, 0b01000000}, // 4
	{0b00000101, 0b01010000}, // 5
	{0b00000110, 0b01100000}, // 6
	{0b00000111, 0b01110000}, // 7
	{0b00001000, 0b10000000}, // 8
	{0b00001001, 0b10010000}, // 9
	{0b00001010, 0b10100000}, // 10
	{0b00001011, 0b10110000}, // 11
	{0b00001100, 0b11000000}, // 12
	{0b00001101, 0b11010000}, // 13
	{0b00001110, 0b11100000}, // 14
	{0b00001111, 0b11110000}, // 15
	{0b11110000, 0b00001111}, // 16
	{0b11110001, 0b00011111}, // 17
	{0b11110010, 0b00101111}, // 18
	{0b11110011, 0b00111111}, // 19
	{0b11110100, 0b01001111}, // 20
	{0b11110101, 0b01011111}, // 21
	{0b11110110, 0b01101111}, // 22
	{0b11110111, 0b01111111}, // 23
	{0b11111000, 0b10001111}, // 24
	{0b11111001, 0b10011111}, // 25
	{0b11111010, 0b10101111}, // 26
	{0b11111011, 0b10111111}, // 27
	{0b11111100, 0b11001111}, // 28
	{0b11111101, 0b11011111}, // 29
	{0b11111110, 0b11101111}, // 30
	{0b11111111, 0b11111111}, // 31
}

func test_0x30_SWAP_B(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.b = tc.value
		testProgram := []uint8{0xCB, 0x30, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x31_SWAP_C(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.c = tc.value
		testProgram := []uint8{0xCB, 0x31, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.c)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x32_SWAP_D(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.d = tc.value
		testProgram := []uint8{0xCB, 0x32, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.d)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x33_SWAP_E(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.e = tc.value
		testProgram := []uint8{0xCB, 0x33, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected E register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.e)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x34_SWAP_H(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.h = tc.value
		testProgram := []uint8{0xCB, 0x34, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.h)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x35_SWAP_L(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.l = tc.value
		testProgram := []uint8{0xCB, 0x35, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.l)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x36_SWAP__HL(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x36, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected the value pointed by [HL] register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x37_SWAP_A(t *testing.T) {
	for idx, tc := range testCaseSWAP {
		preconditions()
		cpu.a = tc.value
		testProgram := []uint8{0xCB, 0x37, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was swapped correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if tc.expectedValue == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Test case %d: Expected Z flag to be false, got true", idx)
			}
		}
		if cpu.getCFlag() {
			t.Errorf("Test case %d: Expected C flag to be false, got true", idx)
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}

// Shift Right Logically register r8.
// 0 ─╂→ b7 → ... → b0 ─╂→ Carry
// opcodes:
//   - 0x38:	SRL B
//   - 0x39:	SRL C
//   - 0x3A:	SRL D
//   - 0x3B:	SRL E
//   - 0x3C:	SRL H
//   - 0x3D:	SRL L
//   - 0x3E:	SRL [HL]
//   - 0x3F:	SRL A
//
// flags: Z=Z N=0 H=0 C=b0
func TestSRL(t *testing.T) {
	t.Run("0x38 SRL B", test_0x38_SRL_B)
	t.Run("0x39 SRL C", test_0x39_SRL_C)
	t.Run("0x3A SRL D", test_0x3A_SRL_D)
	t.Run("0x3B SRL E", test_0x3B_SRL_E)
	t.Run("0x3C SRL H", test_0x3C_SRL_H)
	t.Run("0x3D SRL L", test_0x3D_SRL_L)
	t.Run("0x3E SRL [HL]", test_0x3E_SRL__HL)
	t.Run("0x3F SRL A", test_0x3F_SRL_A)
}

type TestCaseSRL struct {
	value         uint8
	C             bool
	expectedValue uint8
	expectedCFlag bool
}

var testCasesSRL = []TestCaseSRL{
	{0b00000000, false, 0b00000000, false}, // 0
	{0b00000000, true, 0b00000000, false},  // 1
	{0b11111111, false, 0b01111111, true},  // 2
	{0b11111111, true, 0b01111111, true},   // 3
	{0b10101010, false, 0b01010101, false}, // 4
	{0b10101010, true, 0b01010101, false},  // 5
	{0b01010101, true, 0b00101010, true},   // 6
	{0b00110011, false, 0b00011001, true},  // 7
	{0b00110011, true, 0b00011001, true},   // 8
	{0b11001100, false, 0b01100110, false}, // 9
	{0b11001100, true, 0b01100110, false},  // 10
	{0b00010001, false, 0b00001000, true},  // 11
	{0b00010001, true, 0b00001000, true},   // 12
	{0b10001000, false, 0b01000100, false}, // 13
	{0b10001000, true, 0b01000100, false},  // 14
}

func test_0x38_SRL_B(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.b = tc.value
		testProgram := []uint8{0xCB, 0x38, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.b != tc.expectedValue {
			t.Errorf("Test case %d: Expected B register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.b)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x39_SRL_C(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.c = tc.value
		testProgram := []uint8{0xCB, 0x39, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.c != tc.expectedValue {
			t.Errorf("Test case %d: Expected C register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.c)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3A_SRL_D(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.d = tc.value
		testProgram := []uint8{0xCB, 0x3A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.d != tc.expectedValue {
			t.Errorf("Test case %d: Expected D register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.d)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3B_SRL_E(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.e = tc.value
		testProgram := []uint8{0xCB, 0x3B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.e != tc.expectedValue {
			t.Errorf("Test case %d: Expected E register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.e)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3C_SRL_H(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.h = tc.value
		testProgram := []uint8{0xCB, 0x3C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.h != tc.expectedValue {
			t.Errorf("Test case %d: Expected H register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.h)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3D_SRL_L(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.l = tc.value
		testProgram := []uint8{0xCB, 0x3D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.l != tc.expectedValue {
			t.Errorf("Test case %d: Expected L register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.l)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3E_SRL__HL(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x3E, 0x10, tc.value}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		valueAtHL := cpu.bus.Read(0x0003)
		if valueAtHL != tc.expectedValue {
			t.Errorf("Test case %d: Expected value pointed by [HL] register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, valueAtHL)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}
func test_0x3F_SRL_A(t *testing.T) {
	for idx, tc := range testCasesSRL {
		preconditions()
		randomizeFlags()
		if tc.C {
			cpu.setCFlag()
		} else {
			cpu.resetCFlag()
		}
		cpu.a = tc.value
		testProgram := []uint8{0xCB, 0x3F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Test case %d: Expected PC to be 0x0002, got 0x%04X", idx, cpu.pc)
		}
		// check that the value was shifted correctly
		if cpu.a != tc.expectedValue {
			t.Errorf("Test case %d: Expected A register to be 0x%02X, got 0x%02X", idx, tc.expectedValue, cpu.a)
		}
		// check that the flags were set correctly
		if !cpu.getZFlag() && (tc.expectedValue == 0) {
			t.Errorf("Test case %d: Expected Z flag to be true, got false", idx)
		}
		if cpu.getCFlag() != tc.expectedCFlag {
			t.Errorf("Test case %d: Expected C flag to be %t, got %t", idx, tc.expectedCFlag, cpu.getCFlag())
		}
		if cpu.getNFlag() {
			t.Errorf("Test case %d: Expected N flag to be false, got true", idx)
		}
		if cpu.getHFlag() {
			t.Errorf("Test case %d: Expected H flag to be false, got true", idx)
		}
		postconditions()
	}
}

// Test bit b in register r8 or [HL]: If bit b is 0, Z is set.
// opcodes:
//   - 0x40:	BIT 0, B
//   - 0x41:	BIT 0, C
//   - 0x42:	BIT 0, D
//   - 0x43:	BIT 0, E
//   - 0x44:	BIT 0, H
//   - 0x45:	BIT 0, L
//   - 0x46:	BIT 0, [HL]
//   - 0x47:	BIT 0, A
//   - 0x48:	BIT 1, B
//   - 0x49:	BIT 1, C
//   - 0x4A:	BIT 1, D
//   - 0x4B:	BIT 1, E
//   - 0x4C:	BIT 1, H
//   - 0x4D:	BIT 1, L
//   - 0x4E:	BIT 1, [HL]
//   - 0x4F:	BIT 1, A
//   - 0x50:	BIT 2, B
//   - 0x51:	BIT 2, C
//   - 0x52:	BIT 2, D
//   - 0x53:	BIT 2, E
//   - 0x54:	BIT 2, H
//   - 0x55:	BIT 2, L
//   - 0x56:	BIT 2, [HL]
//   - 0x57:	BIT 2, A
//   - 0x58:	BIT 3, B
//   - 0x59:	BIT 3, C
//   - 0x5A:	BIT 3, D
//   - 0x5B:	BIT 3, E
//   - 0x5C:	BIT 3, H
//   - 0x5D:	BIT 3, L
//   - 0x5E:	BIT 3, [HL]
//   - 0x5F:	BIT 3, A
//   - 0x60:	BIT 4, B
//   - 0x61:	BIT 4, C
//   - 0x62:	BIT 4, D
//   - 0x63:	BIT 4, E
//   - 0x64:	BIT 4, H
//   - 0x65:	BIT 4, L
//   - 0x66:	BIT 4, [HL]
//   - 0x67:	BIT 4, A
//   - 0x68:	BIT 5, B
//   - 0x69:	BIT 5, C
//   - 0x6A:	BIT 5, D
//   - 0x6B:	BIT 5, E
//   - 0x6C:	BIT 5, H
//   - 0x6D:	BIT 5, L
//   - 0x6E:	BIT 5, [HL]
//   - 0x6F:	BIT 5, A
//   - 0x70:	BIT 6, B
//   - 0x71:	BIT 6, C
//   - 0x72:	BIT 6, D
//   - 0x73:	BIT 6, E
//   - 0x74:	BIT 6, H
//   - 0x75:	BIT 6, L
//   - 0x76:	BIT 6, [HL]
//   - 0x77:	BIT 6, A
//   - 0x78:	BIT 7, B
//   - 0x79:	BIT 7, C
//   - 0x7A:	BIT 7, D
//   - 0x7B:	BIT 7, E
//   - 0x7C:	BIT 7, H
//   - 0x7D:	BIT 7, L
//   - 0x7E:	BIT 7, [HL]
//   - 0x7F:	BIT 7, A
//
// flags: Z=Z N=0 H=1 C=-
func TestBIT(t *testing.T) {
	// BIT 0
	t.Run("0x40 BIT 0, B", test_0x40_BIT_0_B)
	t.Run("0x41 BIT 0, C", test_0x41_BIT_0_C)
	t.Run("0x42 BIT 0, D", test_0x42_BIT_0_D)
	t.Run("0x43 BIT 0, E", test_0x43_BIT_0_E)
	t.Run("0x44 BIT 0, H", test_0x44_BIT_0_H)
	t.Run("0x45 BIT 0, L", test_0x45_BIT_0_L)
	t.Run("0x46 BIT 0, [HL]", test_0x46_BIT_0__HL)
	t.Run("0x47 BIT 0, A", test_0x47_BIT_0_A)

	// BIT 1
	t.Run("0x48 BIT 1, B", test_0x48_BIT_1_B)
	t.Run("0x49 BIT 1, C", test_0x49_BIT_1_C)
	t.Run("0x4A BIT 1, D", test_0x4A_BIT_1_D)
	t.Run("0x4B BIT 1, E", test_0x4B_BIT_1_E)
	t.Run("0x4C BIT 1, H", test_0x4C_BIT_1_H)
	t.Run("0x4D BIT 1, L", test_0x4D_BIT_1_L)
	t.Run("0x4E BIT 1, [HL]", test_0x4E_BIT_1__HL)
	t.Run("0x4F BIT 1, A", test_0x4F_BIT_1_A)

	// BIT 2
	t.Run("0x50 BIT 2, B", test_0x50_BIT_2_B)
	t.Run("0x51 BIT 2, C", test_0x51_BIT_2_C)
	t.Run("0x52 BIT 2, D", test_0x52_BIT_2_D)
	t.Run("0x53 BIT 2, E", test_0x53_BIT_2_E)
	t.Run("0x54 BIT 2, H", test_0x54_BIT_2_H)
	t.Run("0x55 BIT 2, L", test_0x55_BIT_2_L)
	t.Run("0x56 BIT 2, [HL]", test_0x56_BIT_2__HL)
	t.Run("0x57 BIT 2, A", test_0x57_BIT_2_A)

	// BIT 3
	t.Run("0x58 BIT 3, B", test_0x58_BIT_3_B)
	t.Run("0x59 BIT 3, C", test_0x59_BIT_3_C)
	t.Run("0x5A BIT 3, D", test_0x5A_BIT_3_D)
	t.Run("0x5B BIT 3, E", test_0x5B_BIT_3_E)
	t.Run("0x5C BIT 3, H", test_0x5C_BIT_3_H)
	t.Run("0x5D BIT 3, L", test_0x5D_BIT_3_L)
	t.Run("0x5E BIT 3, [HL]", test_0x5E_BIT_3__HL)
	t.Run("0x5F BIT 3, A", test_0x5F_BIT_3_A)

	// BIT 4
	t.Run("0x60 BIT 4, B", test_0x60_BIT_4_B)
	t.Run("0x61 BIT 4, C", test_0x61_BIT_4_C)
	t.Run("0x62 BIT 4, D", test_0x62_BIT_4_D)
	t.Run("0x63 BIT 4, E", test_0x63_BIT_4_E)
	t.Run("0x64 BIT 4, H", test_0x64_BIT_4_H)
	t.Run("0x65 BIT 4, L", test_0x65_BIT_4_L)
	t.Run("0x66 BIT 4, [HL]", test_0x66_BIT_4__HL)
	t.Run("0x67 BIT 4, A", test_0x67_BIT_4_A)

	// BIT 5
	t.Run("0x68 BIT 5, B", test_0x68_BIT_5_B)
	t.Run("0x69 BIT 5, C", test_0x69_BIT_5_C)
	t.Run("0x6A BIT 5, D", test_0x6A_BIT_5_D)
	t.Run("0x6B BIT 5, E", test_0x6B_BIT_5_E)
	t.Run("0x6C BIT 5, H", test_0x6C_BIT_5_H)
	t.Run("0x6D BIT 5, L", test_0x6D_BIT_5_L)
	t.Run("0x6E BIT 5, [HL]", test_0x6E_BIT_5__HL)
	t.Run("0x6F BIT 5, A", test_0x6F_BIT_5_A)

	// BIT 6
	t.Run("0x70 BIT 6, B", test_0x70_BIT_6_B)
	t.Run("0x71 BIT 6, C", test_0x71_BIT_6_C)
	t.Run("0x72 BIT 6, D", test_0x72_BIT_6_D)
	t.Run("0x73 BIT 6, E", test_0x73_BIT_6_E)
	t.Run("0x74 BIT 6, H", test_0x74_BIT_6_H)
	t.Run("0x75 BIT 6, L", test_0x75_BIT_6_L)
	t.Run("0x76 BIT 6, [HL]", test_0x76_BIT_6__HL)
	t.Run("0x77 BIT 6, A", test_0x77_BIT_6_A)

	// BIT 7
	t.Run("0x78 BIT 7, B", test_0x78_BIT_7_B)
	t.Run("0x79 BIT 7, C", test_0x79_BIT_7_C)
	t.Run("0x7A BIT 7, D", test_0x7A_BIT_7_D)
	t.Run("0x7B BIT 7, E", test_0x7B_BIT_7_E)
	t.Run("0x7C BIT 7, H", test_0x7C_BIT_7_H)
	t.Run("0x7D BIT 7, L", test_0x7D_BIT_7_L)
	t.Run("0x7E BIT 7, [HL]", test_0x7E_BIT_7__HL)
	t.Run("0x7F BIT 7, A", test_0x7F_BIT_7_A)
}

// BIT 0
func test_0x40_BIT_0_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x40, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x41_BIT_0_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x41, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x42_BIT_0_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x42, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x43_BIT_0_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x43, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x44_BIT_0_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x44, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x45_BIT_0_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x45, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x46_BIT_0__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x46, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x47_BIT_0_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x47, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&0x01 == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 1
func test_0x48_BIT_1_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x48, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x49_BIT_1_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x49, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4A_BIT_1_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x4A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4B_BIT_1_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x4B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4C_BIT_1_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x4C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4D_BIT_1_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x4D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4E_BIT_1__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x4E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x4F_BIT_1_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x4F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<1) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 2
func test_0x50_BIT_2_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x50, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x51_BIT_2_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x51, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x52_BIT_2_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x52, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x53_BIT_2_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x53, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x54_BIT_2_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x54, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x55_BIT_2_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x55, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x56_BIT_2__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x56, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x57_BIT_2_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x57, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<2) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 3
func test_0x58_BIT_3_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x58, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x59_BIT_3_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x59, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5A_BIT_3_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x5A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5B_BIT_3_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x5B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5C_BIT_3_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x5C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5D_BIT_3_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x5D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5E_BIT_3__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x5E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x5F_BIT_3_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x5F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<3) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 4
func test_0x60_BIT_4_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x60, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x61_BIT_4_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x61, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x62_BIT_4_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x62, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x63_BIT_4_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x63, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x64_BIT_4_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x64, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x65_BIT_4_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x65, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x66_BIT_4__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x66, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x67_BIT_4_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x67, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<4) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 5
func test_0x68_BIT_5_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x68, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x69_BIT_5_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x69, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6A_BIT_5_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x6A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6B_BIT_5_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x6B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6C_BIT_5_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x6C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6D_BIT_5_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x6D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6E_BIT_5__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x6E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x6F_BIT_5_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x6F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<5) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 6
func test_0x70_BIT_6_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x70, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x71_BIT_6_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x71, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x72_BIT_6_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x72, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x73_BIT_6_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x73, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x74_BIT_6_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x74, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x75_BIT_6_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x75, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x76_BIT_6__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x76, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x77_BIT_6_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x77, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<6) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// BIT 7
func test_0x78_BIT_7_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x78, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x79_BIT_7_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x79, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7A_BIT_7_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x7A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7B_BIT_7_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x7B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7C_BIT_7_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x7C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7D_BIT_7_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x7D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7E_BIT_7__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x7E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}
func test_0x7F_BIT_7_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveCFlag := cpu.getCFlag()
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x7F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the flags were set correctly
		if uint8(i)&(1<<7) == 0 {
			if !cpu.getZFlag() {
				t.Errorf("Expected Z flag to be true, got false")
			}
		} else {
			if cpu.getZFlag() {
				t.Errorf("Expected Z flag to be false, got true")
			}
		}
		if cpu.getNFlag() {
			t.Errorf("Expected N flag to be false, got true")
		}
		if !cpu.getHFlag() {
			t.Errorf("Expected H flag to be true, got false")
		}
		if cpu.getCFlag() != saveCFlag {
			t.Errorf("Expected C flag to be unchanged, got %t", cpu.getCFlag())
		}
		postconditions()
	}
}

// Reset bit b in register r8 or [HL].
// opcodes:
//   - 0x80:	RES 0, B
//   - 0x81:	RES 0, C
//   - 0x82:	RES 0, D
//   - 0x83:	RES 0, E
//   - 0x84:	RES 0, H
//   - 0x85:	RES 0, L
//   - 0x86:	RES 0, [HL]
//   - 0x87:	RES 0, A
//   - 0x88:	RES 1, B
//   - 0x89:	RES 1, C
//   - 0x8A:	RES 1, D
//   - 0x8B:	RES 1, E
//   - 0x8C:	RES 1, H
//   - 0x8D:	RES 1, L
//   - 0x8E:	RES 1, [HL]
//   - 0x8F:	RES 1, A
//   - 0x90:	RES 2, B
//   - 0x91:	RES 2, C
//   - 0x92:	RES 2, D
//   - 0x93:	RES 2, E
//   - 0x94:	RES 2, H
//   - 0x95:	RES 2, L
//   - 0x96:	RES 2, [HL]
//   - 0x97:	RES 2, A
//   - 0x98:	RES 3, B
//   - 0x99:	RES 3, C
//   - 0x9A:	RES 3, D
//   - 0x9B:	RES 3, E
//   - 0x9C:	RES 3, H
//   - 0x9D:	RES 3, L
//   - 0x9E:	RES 3, [HL]
//   - 0x9F:	RES 3, A
//   - 0xA0:	RES 4, B
//   - 0xA1:	RES 4, C
//   - 0xA2:	RES 4, D
//   - 0xA3:	RES 4, E
//   - 0xA4:	RES 4, H
//   - 0xA5:	RES 4, L
//   - 0xA6:	RES 4, [HL]
//   - 0xA7:	RES 4, A
//   - 0xA8:	RES 5, B
//   - 0xA9:	RES 5, C
//   - 0xAA:	RES 5, D
//   - 0xAB:	RES 5, E
//   - 0xAC:	RES 5, H
//   - 0xAD:	RES 5, L
//   - 0xAE:	RES 5, [HL]
//   - 0xAF:	RES 5, A
//   - 0xB0:	RES 6, B
//   - 0xB1:	RES 6, C
//   - 0xB2:	RES 6, D
//   - 0xB3:	RES 6, E
//   - 0xB4:	RES 6, H
//   - 0xB5:	RES 6, L
//   - 0xB6:	RES 6, [HL]
//   - 0xB7:	RES 6, A
//   - 0xB8:	RES 7, B
//   - 0xB9:	RES 7, C
//   - 0xBA:	RES 7, D
//   - 0xBB:	RES 7, E
//   - 0xBC:	RES 7, H
//   - 0xBD:	RES 7, L
//   - 0xBE:	RES 7, [HL]
//   - 0xBF:	RES 7, A
//
// flags: None affected
func TestRES(t *testing.T) {
	// BIT 0
	t.Run("0x80 RES 0, B", test_0x80_RES_0_B)
	t.Run("0x81 RES 0, C", test_0x81_RES_0_C)
	t.Run("0x82 RES 0, D", test_0x82_RES_0_D)
	t.Run("0x83 RES 0, E", test_0x83_RES_0_E)
	t.Run("0x84 RES 0, H", test_0x84_RES_0_H)
	t.Run("0x85 RES 0, L", test_0x85_RES_0_L)
	t.Run("0x86 RES 0, [HL]", test_0x86_RES_0__HL)
	t.Run("0x87 RES 0, A", test_0x87_RES_0_A)

	// BIT 1
	t.Run("0x88 RES 1, B", test_0x88_RES_1_B)
	t.Run("0x89 RES 1, C", test_0x89_RES_1_C)
	t.Run("0x8A RES 1, D", test_0x8A_RES_1_D)
	t.Run("0x8B RES 1, E", test_0x8B_RES_1_E)
	t.Run("0x8C RES 1, H", test_0x8C_RES_1_H)
	t.Run("0x8D RES 1, L", test_0x8D_RES_1_L)
	t.Run("0x8E RES 1, [HL]", test_0x8E_RES_1__HL)
	t.Run("0x8F RES 1, A", test_0x8F_RES_1_A)

	// BIT 2
	t.Run("0x90 RES 2, B", test_0x90_RES_2_B)
	t.Run("0x91 RES 2, C", test_0x91_RES_2_C)
	t.Run("0x92 RES 2, D", test_0x92_RES_2_D)
	t.Run("0x93 RES 2, E", test_0x93_RES_2_E)
	t.Run("0x94 RES 2, H", test_0x94_RES_2_H)
	t.Run("0x95 RES 2, L", test_0x95_RES_2_L)
	t.Run("0x96 RES 2, [HL]", test_0x96_RES_2__HL)
	t.Run("0x97 RES 2, A", test_0x97_RES_2_A)

	// BIT 3
	t.Run("0x98 RES 3, B", test_0x98_RES_3_B)
	t.Run("0x99 RES 3, C", test_0x99_RES_3_C)
	t.Run("0x9A RES 3, D", test_0x9A_RES_3_D)
	t.Run("0x9B RES 3, E", test_0x9B_RES_3_E)
	t.Run("0x9C RES 3, H", test_0x9C_RES_3_H)
	t.Run("0x9D RES 3, L", test_0x9D_RES_3_L)
	t.Run("0x9E RES 3, [HL]", test_0x9E_RES_3__HL)
	t.Run("0x9F RES 3, A", test_0x9F_RES_3_A)

	// BIT 4
	t.Run("0xA0 RES 4, B", test_0xA0_RES_4_B)
	t.Run("0xA1 RES 4, C", test_0xA1_RES_4_C)
	t.Run("0xA2 RES 4, D", test_0xA2_RES_4_D)
	t.Run("0xA3 RES 4, E", test_0xA3_RES_4_E)
	t.Run("0xA4 RES 4, H", test_0xA4_RES_4_H)
	t.Run("0xA5 RES 4, L", test_0xA5_RES_4_L)
	t.Run("0xA6 RES 4, [HL]", test_0xA6_RES_4__HL)
	t.Run("0xA7 RES 4, A", test_0xA7_RES_4_A)

	// BIT 5
	t.Run("0xA8 RES 5, B", test_0xA8_RES_5_B)
	t.Run("0xA9 RES 5, C", test_0xA9_RES_5_C)
	t.Run("0xAA RES 5, D", test_0xAA_RES_5_D)
	t.Run("0xAB RES 5, E", test_0xAB_RES_5_E)
	t.Run("0xAC RES 5, H", test_0xAC_RES_5_H)
	t.Run("0xAD RES 5, L", test_0xAD_RES_5_L)
	t.Run("0xAE RES 5, [HL]", test_0xAE_RES_5__HL)
	t.Run("0xAF RES 5, A", test_0xAF_RES_5_A)

	// BIT 6
	t.Run("0xB0 RES 6, B", test_0xB0_RES_6_B)
	t.Run("0xB1 RES 6, C", test_0xB1_RES_6_C)
	t.Run("0xB2 RES 6, D", test_0xB2_RES_6_D)
	t.Run("0xB3 RES 6, E", test_0xB3_RES_6_E)
	t.Run("0xB4 RES 6, H", test_0xB4_RES_6_H)
	t.Run("0xB5 RES 6, L", test_0xB5_RES_6_L)
	t.Run("0xB6 RES 6, [HL]", test_0xB6_RES_6__HL)
	t.Run("0xB7 RES 6, A", test_0xB7_RES_6_A)

	// BIT 7
	t.Run("0xB8 RES 7, B", test_0xB8_RES_7_B)
	t.Run("0xB9 RES 7, C", test_0xB9_RES_7_C)
	t.Run("0xBA RES 7, D", test_0xBA_RES_7_D)
	t.Run("0xBB RES 7, E", test_0xBB_RES_7_E)
	t.Run("0xBC RES 7, H", test_0xBC_RES_7_H)
	t.Run("0xBD RES 7, L", test_0xBD_RES_7_L)
	t.Run("0xBE RES 7, [HL]", test_0xBE_RES_7__HL)
	t.Run("0xBF RES 7, A", test_0xBF_RES_7_A)
}

// BIT 0
func test_0x80_RES_0_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x80, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x81_RES_0_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x81, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x82_RES_0_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x82, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x83_RES_0_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x83, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x84_RES_0_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x84, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x85_RES_0_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x85, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x86_RES_0__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x86, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x87_RES_0_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x87, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 0)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<0), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 1
func test_0x88_RES_1_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x88, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x89_RES_1_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x89, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8A_RES_1_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x8A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8B_RES_1_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x8B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8C_RES_1_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x8C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8D_RES_1_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x8D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8E_RES_1__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x8E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x8F_RES_1_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x8F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 1)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<1), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 2
func test_0x90_RES_2_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x90, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x91_RES_2_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x91, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x92_RES_2_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x92, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x93_RES_2_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x93, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x94_RES_2_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x94, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x95_RES_2_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x95, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x96_RES_2__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x96, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x97_RES_2_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x97, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 2)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<2), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 3
func test_0x98_RES_3_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0x98, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x99_RES_3_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0x99, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9A_RES_3_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0x9A, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9B_RES_3_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0x9B, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9C_RES_3_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0x9C, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9D_RES_3_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0x9D, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9E_RES_3__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0x9E, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0x9F_RES_3_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0x9F, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 3)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<3), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 4
func test_0xA0_RES_4_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0xA0, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA1_RES_4_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0xA1, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA2_RES_4_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0xA2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA3_RES_4_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0xA3, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA4_RES_4_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0xA4, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA5_RES_4_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0xA5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA6_RES_4__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0xA6, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA7_RES_4_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0xA7, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 4)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<4), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 5
func test_0xA8_RES_5_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0xA8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xA9_RES_5_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0xA9, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAA_RES_5_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0xAA, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAB_RES_5_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0xAB, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAC_RES_5_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0xAC, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAD_RES_5_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0xAD, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAE_RES_5__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0xAE, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xAF_RES_5_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0xAF, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 5)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<5), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 6
func test_0xB0_RES_6_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0xB0, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB1_RES_6_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0xB1, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB2_RES_6_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0xB2, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB3_RES_6_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0xB3, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB4_RES_6_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0xB4, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB5_RES_6_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0xB5, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB6_RES_6__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0xB6, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB7_RES_6_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0xB7, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 6)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<6), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}

// BIT 7
func test_0xB8_RES_7_B(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.b = uint8(i)
		testProgram := []uint8{0xCB, 0xB8, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.b != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.b)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xB9_RES_7_C(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.c = uint8(i)
		testProgram := []uint8{0xCB, 0xB9, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.c != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.c)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBA_RES_7_D(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.d = uint8(i)
		testProgram := []uint8{0xCB, 0xBA, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.d != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.d)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBB_RES_7_E(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.e = uint8(i)
		testProgram := []uint8{0xCB, 0xBB, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.e != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.e)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBC_RES_7_H(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.h = uint8(i)
		testProgram := []uint8{0xCB, 0xBC, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.h != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.h)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBD_RES_7_L(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.l = uint8(i)
		testProgram := []uint8{0xCB, 0xBD, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.l != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.l)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBE_RES_7__HL(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.setHL(0x0003)
		testProgram := []uint8{0xCB, 0xBE, 0x10, uint8(i)}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		valueAtHL := memory1.Read(cpu.getHL())
		if valueAtHL != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), valueAtHL)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
func test_0xBF_RES_7_A(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		preconditions()
		randomizeFlags()
		saveFlags := cpu.f
		cpu.a = uint8(i)
		testProgram := []uint8{0xCB, 0xBF, 0x10}
		loadProgramIntoMemory(memory1, testProgram)
		cpu.Run()
		// check that the program stopped at the right place
		if cpu.pc != 0x0002 {
			t.Errorf("Expected PC to be 0x0002, got 0x%04X", cpu.pc)
		}
		// check that the bit was reset
		if cpu.a != (uint8(i) &^ (1 << 7)) {
			t.Errorf("Expected B to be 0x%02X, got 0x%02X", uint8(i)&^(1<<7), cpu.a)
		}
		// check that the flags were left unchanged
		if cpu.f != saveFlags {
			t.Errorf("Expected flags to be unchanged, got 0x%02X", cpu.f)
		}
		postconditions()
	}
}
