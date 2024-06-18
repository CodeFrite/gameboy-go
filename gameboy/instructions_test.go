/* Test the instructions of the Gameboy CPU
 *
 * Note: in order to test the instructions, we need to bring the CPU to execute some OpCodes along with operands.
 * To ease our work, we will use the WRAM to write the instructions and operands and then read them from the CPU.
 * Indeed, unlike VRAM and WRAM being initialized with a known size (2KiB), the ROM is loaded from a rom file that will
 * dictate the final size of the ROM and we don't have direct access to it.
 * This means that we can't write to it before specifying a size which we will not be doing to preserve the struct code.
 */

package gameboy

import (
	"fmt"
	"testing"
)

// instantiate a new gameboy
func createNewGameboy() (*CPU, *Bus, *Cartridge) {
	// 1. Init VRAM
	vram := NewMemory(0x2000)

	// 2. Init WRAM
	wram := NewMemory(0x2000)

	// 3. Init Cartridge
	cartridge := Cartridge{}

	// 4. init BUS
	bus := NewBus()
	bus.AttachMemory(&cartridge, 0x0000)
	bus.AttachMemory(vram, 0x8000)
	bus.AttachMemory(wram, 0xC000)

	// 4. instantiate a new CPU
	cpu := NewCPU(bus)

	return cpu, bus, &cartridge
}

// Deep copy of the CPU struct
func copyCPU(cpu *CPU) *CPU {
	return &CPU{
		PC:   cpu.PC,
		IR:   cpu.IR,
		SP:   cpu.SP,
		A:    cpu.A,
		F:    cpu.F,
		B:    cpu.B,
		C:    cpu.C,
		D:    cpu.D,
		E:    cpu.E,
		H:    cpu.H,
		L:    cpu.L,
		HRAM: cpu.HRAM,
		bus:  cpu.bus,
	}
}

// compare the CPU structs and return all fields that were modified
func compareCPU(cpu *CPU, ref *CPU) map[string][2]interface{} {
	differences := make(map[string][2]interface{})
	// compare all fields and store the differences
	if cpu.PC != ref.PC {
		differences["PC"] = [2]interface{}{ref.PC, cpu.PC}
	}
	if cpu.SP != ref.SP {
		differences["SP"] = [2]interface{}{ref.SP, cpu.SP}
	}
	if cpu.A != ref.A {
		differences["A"] = [2]interface{}{ref.A, cpu.A}
	}
	if cpu.F != ref.F {
		differences["F"] = [2]interface{}{ref.F, cpu.F}
	}
	if cpu.B != ref.B {
		differences["B"] = [2]interface{}{ref.B, cpu.B}
	}
	if cpu.C != ref.C {
		differences["C"] = [2]interface{}{ref.C, cpu.C}
	}
	if cpu.D != ref.D {
		differences["D"] = [2]interface{}{ref.D, cpu.D}
	}
	if cpu.E != ref.E {
		differences["E"] = [2]interface{}{ref.E, cpu.E}
	}
	if cpu.H != ref.H {
		differences["H"] = [2]interface{}{ref.H, cpu.H}
	}
	if cpu.L != ref.L {
		differences["L"] = [2]interface{}{ref.L, cpu.L}
	}
	return differences
}

// Write to WRAM and read from WRAM with relative address
func writeWRAM(bus *Bus, addr uint16, value uint8) uint8 {
	var WRAMStart uint16 = 0xC000
	relativeAddr := addr + WRAMStart
	bus.Write(relativeAddr, value)
	return bus.Read(relativeAddr)
}

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

//=========================//
// no operand instructions //
//=========================//

/*
 * 0x00: NOP
 * PC should be incremented by 1
 * TODO: update the CPU clock
 */
func TestNOP(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// NOP
	cpu.executeInstruction(GetInstruction("0x00", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if only the program counter was modified
	if len(differences) != 1 {
		t.Errorf("Expected 1 difference, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	} else if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}
}

/*
 * 0x07: RLCA rotates the A register to the left
 * The carry flag should be set to the value of the 7th bit of the A register
 * All other flags should be reset
 */
func TestRLCA(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set all flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// reset the carry flag
	cpu.resetCFlag()

	// set the A register to b10101010
	cpu.A = 0xAA

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x07", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 5 differences
	if len(differences) != 3 {
		fmt.Println(differences)
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the A register was rotated to the left
	if _, ok := differences["A"]; !ok {
		t.Errorf("Expected A to be modified")
	}

	if cpu.A != 0x55 {
		t.Errorf("Expected A to be 0x55, got 0x%02X", cpu.A)
	}

	// check if the carry flag was set
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

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
}

/*
 * 0x0F: RRCA rotates the A register to the right
 * The carry flag should be set to the value of the 0th bit of the A register
 * All other flags should be reset
 */
func TestRRCA(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set all flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// reset the carry flag
	cpu.resetCFlag()

	// set the A register to 0101 0101
	cpu.A = 0x55

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x0F", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 5 differences
	if len(differences) != 3 {
		fmt.Println(differences)
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the A register was rotated to the left
	if _, ok := differences["A"]; !ok {
		t.Errorf("Expected A to be modified")
	}

	if cpu.A != 0xAA {
		t.Errorf("Expected A to be 0xAA, got 0x%02X", cpu.A)
	}

	// check if the carry flag was set
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

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
}

/*
 * 0x17: RLA rotates the A register to the left through the carry flag
 * The 0th bit of the A register should be set to the value of the carry flag
 * The carry flag should be set to the value of the 7th bit of the A register
 * All other flags should be reset
 */
func TestRLA(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set all flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// set the carry flag
	cpu.setCFlag()

	// set the A register to b 0101 0101
	cpu.A = 0x55

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x17", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 5 differences
	if len(differences) != 3 {
		fmt.Println(differences)
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the A register was rotated to the left
	if _, ok := differences["A"]; !ok {
		t.Errorf("Expected A to be modified")
	}

	if cpu.A != 0xAB {
		t.Errorf("Expected A to be 0xAB, got 0x%02X", cpu.A)
	}

	// check if the carry flag was reset
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

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

/*
 * 0x1F: RRA rotates the A register to the right through the carry flag
 * The 7th bit of the A register should be set to the value of the carry flag
 * The carry flag should be set to the value of the 0th bit of the A register
 * All other flags should be reset
 */
func TestRRA(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set all flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// set the carry flag
	cpu.setCFlag()

	// set the A register to b 1010 1010
	cpu.A = 0xAA

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x1F", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 5 differences
	if len(differences) != 3 {
		fmt.Println(differences)
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the A register was rotated to the left
	if _, ok := differences["A"]; !ok {
		t.Errorf("Expected A to be modified")
	}
	// 1010 1010 => 1101 0101 (0xD5)
	if cpu.A != 0xD5 {
		t.Errorf("Expected A to be 0xD5, got 0x%02X", cpu.A)
	}

	// check if the carry flag was reset
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

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

/*
 * 0x27: DAA
 * Decimal Adjust Accumulator to get a correct BCD representation after an arithmetic instruction.
 * The steps performed depend on the type of the previous operation: ADD or SUB
 */
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
		cpu, _, _ := createNewGameboy()

		// set the A register to 10 = 0x0A
		cpu.A = tc.initial.A

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

		// save cpu state
		cpuCopy := copyCPU(cpu)

		// Execute the instruction
		cpu.executeInstruction(GetInstruction("0x27", false), nil, nil)

		// check if the program counter was incremented by instruction.Length
		if cpu.PC != cpuCopy.PC+1 {
			t.Errorf("TC%v> Expected PC to be %v, got %v", tci, cpuCopy.PC+1, cpu.PC)
		}

		// check if the A register now contains the correct BCD representation of 10 (0x10)
		if cpu.A != tc.expected.A {
			t.Errorf("TC%v> Expected A to be 0x%02X, got 0x%02X", tci, tc.expected.A, cpu.A)
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
	}
}

/*
 * 0x2F: CPL
 * ComPLement accumulator: flips all bits of the accumulator (A = ~A).
 */
func TestCPL(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// reset all F register to 0110
	cpu.resetZFlag()
	cpu.resetNFlag()
	cpu.resetHFlag()
	cpu.resetCFlag()

	// set the A register to b 1010 1010
	cpu.A = 0xAA

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x2F", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 5 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the A register was rotated to the left
	if _, ok := differences["A"]; !ok {
		t.Errorf("Expected A to be modified")
	}
	// 1010 1010 => 1101 0101 (0xD5)
	if cpu.A != 0x55 {
		t.Errorf("Expected A to be 0x55, got 0x%02X", cpu.A)
	}

	// check if the carry flag was reset
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

	// check if the Z, N and H flags were reset
	if cpu.getZFlag() != false {
		t.Errorf("Expected Z flag to be reset")
	}
	if cpu.getNFlag() != true {
		t.Errorf("Expected N flag to be set")
	}
	if cpu.getHFlag() != true {
		t.Errorf("Expected H flag to be set")
	}
	if cpu.getCFlag() != false {
		t.Errorf("Expected C flag to be reset")
	}
}

/*
 * 0x37: SCF
 * Set Carry Flag: set the carry flag to 1 and reset the N and H flags
 */
func TestSCF(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set flags
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()
	cpu.resetCFlag()

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x37", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// check if the Z flag was left untouched (1)
	if cpu.getZFlag() != true {
		t.Errorf("Expected Z flag to be 1")
	}

	// check if the N flag was reset
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be 0")
	}

	// check if the H flag was reset
	if cpu.getHFlag() != false {
		t.Errorf("Expected H flag to be 0")
	}

	// check if the C flag was set
	if cpu.getCFlag() != true {
		t.Errorf("Expected C flag to be 1")
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}
}

/*
 * 0x3F: CCF
 * Complement Carry Flag: flip the value of the carry flag and reset the N and H flags
 */
func TestCCF(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// set C flag and expect it to be reset
	cpu.setCFlag()

	// set flags
	cpu.setZFlag() // should be left untouched
	cpu.setNFlag() // should be reset
	cpu.setHFlag() // should be reset

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x3F", false), nil, nil)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// check if the Z flag was left untouched (1)
	if cpuCopy.getZFlag() != cpu.getZFlag() {
		t.Errorf("Expected Z flag to be 1")
	}

	// check if the N flag was reset
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be 0")
	}

	// check if the H flag was reset
	if cpu.getHFlag() != false {
		t.Errorf("Expected H flag to be 0")
	}

	// check if the C flag was set
	if cpu.getCFlag() == cpuCopy.getCFlag() {
		t.Errorf("Expected C flag to be flipped")
	}

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}
}

/*
 * 0x76: HALT
 * Halt the CPU until an interrupt occurs
 */
func TestHALT(t *testing.T) {
	// ! Should be tested in the context of the CPU.Run() method when interrupts are implemented
	cpu, _, _ := createNewGameboy()

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x76", false), nil, nil)

	// check if the program counter was not modified
	if cpu.PC != cpuCopy.PC {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC, cpu.PC)
	}

	// check if the CPU is halted
	if !cpu.halted {
		t.Errorf("Expected CPU to be halted")
	}
}

/*
 * 0xC9: RET
 * Pop two bytes from the stack and jump to that address
 */
func TestRET(t *testing.T) {
	t.Skip("Skipping until PUSH instruction is implemented")
}

/*
 * 0xCB: PREFIX CB
 * Execute the next instruction from the CB instruction set
 */
func TestPREFIX_CB(t *testing.T) {
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0xCB) // PREFIX CB
	writeWRAM(bus, 0x01, 0xFC) // just need to make sure the next byte is loaded into the IR (any value different than CB is ok)

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction (may panic)
	mayPanic(cpu.Run)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check if the program counter was incremented by instruction.Length
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+1 {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+1, cpu.PC)
	}

	// check if the IR register was loaded with the next byte
	if cpu.IR != 0xFC {
		t.Errorf("Expected IR to be 0xFC, got 0x%02X", cpu.IR)
	}
}

//=========================//
// 2 operands instructions //
//=========================//

// 0x01: Load the next two bytes into the BC register
func TestLD_BC_n16(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 2 operands")
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0x01) // instruction
	writeWRAM(bus, 0x01, 0x1A) // operand 1
	writeWRAM(bus, 0x02, 0x2B) // operand 2

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Run()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := GetInstruction("0x01", false)
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+uint16(instruction.Bytes) {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+uint16(instruction.Bytes), cpu.PC)
	}

	// ... the BC register was loaded with the operand
	if _, ok := differences["BC"]; !ok {
		t.Errorf("Expected BC to be modified")
	}
	if cpu.getBC() != 0x2B1A {
		t.Errorf("Expected BC to be 0x2B1A, got 0x%04X", cpu.getBC())
	}
}

// 0x11: Load the next two bytes into the DE register
func TestLD_DE_n16(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 2 operands")
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0x11) // instruction
	writeWRAM(bus, 0x01, 0x1A) // operand 1
	writeWRAM(bus, 0x02, 0x2B) // operand 2

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Run()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := GetInstruction("0x01", false)
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+uint16(instruction.Bytes) {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+uint16(instruction.Bytes), cpu.PC)
	}

	// ... the BC register was loaded with the operand
	if _, ok := differences["DE"]; !ok {
		t.Errorf("Expected BC to be modified")
	}
	if cpu.getDE() != 0x2B1A {
		t.Errorf("Expected BC to be 0x2B1A, got 0x%X", cpu.getDE())
	}
}

// 0x2C: Increment the value of register L
func TestINC_L_NO_FLAG_SET(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 1 operand")
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the Z, N and H flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0x2C) // instruction

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Run()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := GetInstruction("0x2C", false)
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+uint16(instruction.Bytes) {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+uint16(instruction.Bytes), cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.getHL() != cpuCopy.getHL()+1 {
		t.Errorf("Expected HL to be 0x%04X, got 0x%04X", cpuCopy.getHL()+1, cpu.getHL())
	}

	// ... the Z, N and H flags were reset
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

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

func TestINC_L_FLAG_H_SET(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 1 operand")
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the Z, N flags and reset H flag to see if they are computed correctly
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()

	// set the L register to 15 (0x0F)
	cpu.L = 0x000F

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0x2C) // instruction

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Run()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := GetInstruction("0x2C", false)
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+uint16(instruction.Bytes) {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+uint16(instruction.Bytes), cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.getHL() != cpuCopy.getHL()+1 {
		t.Errorf("Expected HL to be 0x%04X, got 0x%04X", cpuCopy.getHL()+1, cpu.getHL())
	}

	// ... the H flag was set
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

	// check if the Z flag was reset
	if cpu.getZFlag() != false {
		t.Errorf("Expected Z flag to be reset")
	}

	// check if the N flag was reset
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be reset")
	}

	// check if the H flag was set
	if cpu.getHFlag() != true {
		t.Errorf("Expected H flag to be set")
	}
}

func TestINC_L_FLAGS_Z_H_SET(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 1 operand")
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the N flags and reset Z & H flags to see if they are computed correctly
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()

	// set the L register to 0xFF
	cpu.L = 0xFF

	// write the instruction and operand to vram
	writeWRAM(bus, 0x00, 0x2C) // instruction

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Run()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := GetInstruction("0x2C", false)
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+uint16(instruction.Bytes) {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+uint16(instruction.Bytes), cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.getHL() != 0x0000 {
		t.Errorf("Expected HL to be 0x%04X, got 0x%04X", cpuCopy.getHL()+1, cpu.getHL())
	}

	// ... the H flag was updated
	if _, ok := differences["F"]; !ok {
		t.Errorf("Expected F to be modified")
	}

	// check if the Z flag was set
	if cpu.getZFlag() != true {
		t.Errorf("Expected Z flag to be set")
	}

	// check if the N flag was reset
	if cpu.getNFlag() != false {
		t.Errorf("Expected N flag to be reset")
	}

	// check if the H flag was set
	if cpu.getHFlag() != true {
		t.Errorf("Expected H flag to be set")
	}
}

// 0x4A: Load the value of register D into register C
func TestLD_C_D(t *testing.T) {
	t.Skip("Skipping until CPU supports instructions with 2 operands")
	// instantiate a new gameboy
	cpu, _, _ := createNewGameboy()

	// set the C register to 0xAB and the D register to 0xCD
	cpu.C = 0xAB
	cpu.D = 0xCD

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.executeInstruction(GetInstruction("0x4A", false), cpu.D, cpu.C)

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// ... the C register was loaded with the value of the D register
	if _, ok := differences["C"]; !ok {
		t.Errorf("Expected C register to be modified")
	}

	if byte(cpu.getBC()&0x00FF) != byte((cpu.getDE()&0xFF00)>>8) {
		t.Errorf("Expected C register to be 0xCD, got 0x%04X", byte(cpu.getBC()&0x00FF))
	}
}
