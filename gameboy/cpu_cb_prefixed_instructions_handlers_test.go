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
