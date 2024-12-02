package gameboy

import "testing"

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

var testCasesRRC = []TestCaseRLC{
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
