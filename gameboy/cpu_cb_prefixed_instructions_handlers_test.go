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
