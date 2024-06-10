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

