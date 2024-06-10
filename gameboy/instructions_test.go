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
	"testing"
)

// instantiate a new gameboy
func createNewGameboy() (*CPU, *Bus, *Cartridge) {
	// 1. Init VRAM
	vram := NewVRAM()

	// 2. Init WRAM
	wram := NewWRAM()

	// 3. Init Cartridge
	cartridge := Cartridge{}

	// 4. init BUS
	bus := NewBus(&cartridge, vram, wram)

	// 4. instantiate a new CPU
	cpu := NewCPU(bus)

	return cpu, bus, &cartridge
}

// Deep copy of the CPU struct
func copyCPU(cpu *CPU) *CPU {
	return &CPU{
		PC: cpu.PC,
		SP: cpu.SP,
		A: cpu.A,
		F: cpu.F,
		BC: cpu.BC,
		DE: cpu.DE,
		HL: cpu.HL,
		HRAM: cpu.HRAM,
		IR: cpu.IR,
		bus: cpu.bus,
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
	if cpu.BC != ref.BC {
		differences["BC"] = [2]interface{}{ref.BC, cpu.BC}
	}
	if cpu.DE != ref.DE {
		differences["DE"] = [2]interface{}{ref.DE, cpu.DE}
	}
	if cpu.HL != ref.HL {
		differences["HL"] = [2]interface{}{ref.HL, cpu.HL}
	}

	return differences
}

// Write to WRAM and read from WRAM with relative address
func writeWRAM(bus *Bus, addr [2]byte, value byte) byte {
	var WRAMStart [2]byte = [2]byte{0xC0, 0x00}
	relativeAddr := uint16ToBytes(bytesToUint16(addr) + bytesToUint16(WRAMStart))
	bus.Write(relativeAddr, value)
	return bus.Read(relativeAddr)
}

// 0x00: Nothing should happen apart from the program counter incrementing
func TestNOP(t *testing.T) {
	cpu, _, _ := createNewGameboy()

	// save cpu state
	cpuCopy := copyCPU(cpu)

	// NOP
	cpu.executeInstruction(Instructions[0x00])

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

// 0x01: Load the next two bytes into the BC register
func TestLD_BC_n16(t *testing.T) {
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000
	
	// write the instruction and operand to vram
	writeWRAM(bus, uint16ToBytes(0x00), 0x01) // instruction
	writeWRAM(bus, uint16ToBytes(0x01), 0x1A) // operand 1
	writeWRAM(bus, uint16ToBytes(0x02), 0x2B) // operand 2
	
	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Step()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := Instructions[0x01]
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+instruction.Length {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+instruction.Length, cpu.PC)
	}

	// ... the BC register was loaded with the operand
	if _, ok := differences["BC"]; !ok {
		t.Errorf("Expected BC to be modified")
	}
	if cpu.BC != 0x2B1A {
		t.Errorf("Expected BC to be 0x2B1A, got 0x%X", cpu.BC)
	}
}

// 0x11: Load the next two bytes into the DE register
func TestLD_DE_n16(t *testing.T) {
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000
	
	// write the instruction and operand to vram
	writeWRAM(bus, uint16ToBytes(0x00), 0x11) // instruction
	writeWRAM(bus, uint16ToBytes(0x01), 0x1A) // operand 1
	writeWRAM(bus, uint16ToBytes(0x02), 0x2B) // operand 2
	
	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Step()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 2 differences
	if len(differences) != 2 {
		t.Errorf("Expected 2 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := Instructions[0x01]
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+instruction.Length {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+instruction.Length, cpu.PC)
	}

	// ... the BC register was loaded with the operand
	if _, ok := differences["DE"]; !ok {
		t.Errorf("Expected BC to be modified")
	}
	if cpu.DE != 0x2B1A {
		t.Errorf("Expected BC to be 0x2B1A, got 0x%X", cpu.DE)
	}
}


// 0x2C: Increment the value of register L
func TestINC_L_NO_FLAG_SET(t *testing.T) {
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the Z, N and H flags to 1 to see if they are reset
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.setHFlag()
	
	// write the instruction and operand to vram
	writeWRAM(bus, uint16ToBytes(0x00), 0x2C) // instruction
	
	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Step()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := Instructions[0x2C]
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+instruction.Length {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+instruction.Length, cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.HL != cpuCopy.HL+1 {
		t.Errorf("Expected HL to be 0x%X, got 0x%X", cpuCopy.HL+1, cpu.HL)
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
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the Z, N flags and reset H flag to see if they are computed correctly
	cpu.setZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()

	// set the L register to 15 (0x0F)
	cpu.HL = 0x000F
	
	// write the instruction and operand to vram
	writeWRAM(bus, uint16ToBytes(0x00), 0x2C) // instruction
	
	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Step()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := Instructions[0x2C]
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+instruction.Length {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+instruction.Length, cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.HL != cpuCopy.HL+1 {
		t.Errorf("Expected HL to be 0x%X, got 0x%X", cpuCopy.HL+1, cpu.HL)
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
	cpu, bus, _ := createNewGameboy()

	// set the program counter to VRAM 0x00 address (0xC000)
	cpu.PC = 0xC000

	// set the N flags and reset Z & H flags to see if they are computed correctly
	cpu.resetZFlag()
	cpu.setNFlag()
	cpu.resetHFlag()

	// set the L register to 0xFF
	cpu.HL = 0x00FF
	
	// write the instruction and operand to vram
	writeWRAM(bus, uint16ToBytes(0x00), 0x2C) // instruction
	
	// save cpu state
	cpuCopy := copyCPU(cpu)

	// Execute the instruction
	cpu.Step()

	// get the differences
	differences := compareCPU(cpu, cpuCopy)

	// check that there are 3 differences
	if len(differences) != 3 {
		t.Errorf("Expected 3 differences, got %v", len(differences))
	}

	// ... the program counter was incremented by instruction.Length
	instruction := Instructions[0x2C]
	if _, ok := differences["PC"]; !ok {
		t.Errorf("Expected PC to be modified")
	}
	if cpu.PC != cpuCopy.PC+instruction.Length {
		t.Errorf("Expected PC to be %v, got %v", cpuCopy.PC+instruction.Length, cpu.PC)
	}

	// ... the HL register was incremented by 1
	if _, ok := differences["HL"]; !ok {
		t.Errorf("Expected HL to be modified")
	}
	if cpu.HL != 0x0000 {
		t.Errorf("Expected HL to be 0x%X, got 0x%X", cpuCopy.HL+1, cpu.HL)
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
