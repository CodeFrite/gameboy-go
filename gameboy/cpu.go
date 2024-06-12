// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"fmt"
)

type CPU struct {
	IR uint8 // Instruction Register
	IE uint8 // Interrupt Enable
	SP uint16 // Stack Pointer
	PC uint16 // Program Counter

	// 8-bit registers
	A uint8 // Accumulator
	F uint8 // Flags (Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4))

	// 16-bit general purpose registers
	BC, DE, HL uint16

	// 127bytes of High-Speed RAM
	HRAM [127]byte

	// reference to the bus
	bus *Bus
}

func NewCPU(bus *Bus) *CPU {
	// initialize all registers to 0 except the program counter which starts at 0x100 (in cartridge ROM)
	return &CPU{
		PC:0x0100, // Start at 0x100 and ignore the boot ROM for now
		bus: bus,
	}
}

func (c *CPU) PrintRegisters() {
	// Print the registers
	fmt.Println("> CPU Registers:")
	fmt.Printf("+ IR: 0x%X\n", c.IR)
	fmt.Printf("+ IE: 0x%X\n", c.IE)
	fmt.Printf("+ SP: 0x%X\n", c.SP)
	fmt.Printf("+ PC: 0x%X\n", c.PC)
	fmt.Printf("+ A: 0x%X\n", c.A)
	fmt.Printf("+ F: 0x%X\n", c.F)
	fmt.Printf("+ BC: 0x%X\n", c.BC)
	fmt.Printf("+ DE: 0x%X\n", c.DE)
	fmt.Printf("+ HL: 0x%X\n", c.HL)
	fmt.Printf("+ HRAM: 0x%X\n", c.HRAM)
}

func (c *CPU) PrintIR() {
	// Print the instruction register
	fmt.Printf("IR: 0x%X\n", c.IR)
}

func (c *CPU) PrintPC() {
	// Print the program counter
	fmt.Printf("PC: 0x%X\n", c.PC)
}

// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return c.F & 0x80 == 0x80
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.F = c.F | 0x80
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.F = c.F & 0x7F
}

// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return c.F & 0x40 == 0x40
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.F = c.F | 0x40
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.F = c.F & 0xBF
}

// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return c.F & 0x20 == 0x20
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.F = c.F | 0x20
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.F = c.F & 0xDF
}

func (c *CPU) getB() byte {
	return byte(c.BC >> 8)
}

func (c *CPU) getC() byte {
	return byte(c.BC)
}

func (c *CPU) getD() byte {
	return byte(c.DE >> 8)
}

func (c *CPU) getE() byte {
	return byte(c.DE)
}

func (c *CPU) getH() byte {
	return byte(c.HL >> 8)
}

func (c *CPU) getL() byte {
	return byte(c.HL)
}

// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() {
	// Fetch the opcode from memory at the address in the program counter
	opcode := c.bus.Read(uint16ToBytes(c.PC))
	// Store the opcode in the instruction register
	c.IR = opcode
}

// Fetch the instruction corresponding to the opcode stored in the instruction register
func (c *CPU) fetchInstruction() Instruction {
	return GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), false)
}

// 
func (c *CPU) executeInstruction(instruction Instruction) {
	// Execute the corresponding instruction
	switch instruction.Mnemonic {
	case "NOP":
		c.NOP(&instruction)
	case "STOP":
		c.STOP(&instruction)
	case "HALT":
		c.HALT(&instruction)
	case "DI":
		c.DI(&instruction)
	case "EI":
		c.EI(&instruction)
	case "PREFIX":
		c.PREFIX(&instruction)
	case "JP":
		c.JP(&instruction)
	case "JR":
		c.JR(&instruction)
	case "CALL":
		c.CALL(&instruction)
	case "RET":
		c.RET(&instruction)
	case "RETI":
		c.RETI(&instruction)
	case "RST":
		c.RST(&instruction)
	case "LD":
		c.LD(&instruction)
	case "LDH":
		c.LDH(&instruction)
	case "PUSH":
		c.PUSH(&instruction)
	case "POP":
		c.POP(&instruction)
	case "ADD":
		c.ADD(&instruction)
	case "ADC":
		c.ADC(&instruction)
	case "AND":
		c.AND(&instruction)
	case "INC":
		c.INC(&instruction)
	case "CCF":
		c.CCF(&instruction)
	case "CP":
		c.CP(&instruction)
	case "CPL":
		c.CPL(&instruction)
	case "DAA":
		c.DAA(&instruction)
	case "DEC":
		c.DEC(&instruction)
	case "SUB":
		c.SUB(&instruction)
	case "SBC":
		c.SBC(&instruction)
	case "SCF":
		c.SCF(&instruction)
	case "OR":
		c.OR(&instruction)
	case "XOR":
		c.XOR(&instruction)
	case "RLA":
		c.RLA(&instruction)
	case "RLCA":
		c.RLCA(&instruction)
	case "RRA":
		c.RRA(&instruction)
	case "RRCA":
		c.RRCA(&instruction)
	default:
		fmt.Println("Unknown instruction")
	}
}

func (c *CPU) incrementPC(offset uint16) {
	c.PC += uint16(offset)
}

func (c *CPU) Step() error {
	// Fetch the opcode from memory
	c.fetchOpcode()
	c.PrintPC()
	c.PrintIR()
	
	// Fetch the instruction corresponding to the opcode contained in the instruction register
	instruction := c.fetchInstruction()

	// Execute the instruction
	c.executeInstruction(instruction)

	return nil
}

// > instructions handlers (NO PREFIX)

// Misc / Control instructions
func (c *CPU) DI(instruction *Instruction) {
	panic("DI not implemented")
}
func (c *CPU) EI(instruction *Instruction) {
	panic("EI not implemented")
}
func (c *CPU) HALT(instruction *Instruction) {
	panic("HALT not implemented")
}

// Fetch the value of an operand
func (c *CPU) fetchOperandValue(operand Operand) interface{} {
	var value interface{}
	switch operand.Name {
		case "n8": // always immediate
			value = c.bus.Read(uint16ToBytes(c.PC + 1))
		case "n16": // always immediate
			// little-endian
			value = [2]byte{
				c.bus.Read(uint16ToBytes(c.PC + 2)),
				c.bus.Read(uint16ToBytes(c.PC + 1)),
			}
		case "a8": // not always immediate
			if operand.Immediate {
				value = c.bus.Read(uint16ToBytes(c.PC + 1))
			} else {
				addr := c.bus.Read(uint16ToBytes(c.PC + 1))
				value = c.bus.Read(uint16ToBytes(uint16(addr)))
			}
		case "a16": // not always immediate
			if operand.Immediate {
				value = [2]byte{
					c.bus.Read(uint16ToBytes(c.PC + 2)),
					c.bus.Read(uint16ToBytes(c.PC + 1)),
				}
			} else {
				addr := [2]byte{
					c.bus.Read(uint16ToBytes(c.PC + 2)),
					c.bus.Read(uint16ToBytes(c.PC + 1)),
				}
				value = c.bus.Read(addr)
			}
		case "A":
			if operand.Immediate {
				value = c.A
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.A)))
			}
		case "B":
			if operand.Immediate {
				value = c.getB()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getB())))
			}
		case "C":
			if operand.Immediate {
				value = c.getC()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getC())))
			}
		case "D":
			if operand.Immediate {
				value = c.getD()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getD())))
			}
		case "E":
			if operand.Immediate {
				value = c.getE()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getE())))
			}
		case "H":
			if operand.Immediate {
				value = c.getH()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getH())))
			}
		case "L":
			if operand.Immediate {
				value = c.getL()
			} else {
				value = c.bus.Read(uint16ToBytes(uint16(c.getL())))
			}
		default:
			panic("Unknown operand type")
	}
	return value
}

// Fetch the address of an operand
func (c *CPU) fetchOperandAddress(operand Operand) uint16 {
	var address uint16
	switch operand.Name {
	case "A":
		address = *c.A
	case "B":
		address = *c.B
	case "C":
		address = *c.C
	case "D":
		address = *c.D
	case "E":
		address = *c.E
	case "H":
		address = *c.H
	case "L":
		address = *c.L
	case "BC":
		if operand.Immediate {
			address = c.BC
		} else {
			address = c.bus.Read(uint16ToBytes(c.BC))
		}
	case "DE":
		if operand.Immediate {
			address = c.DE
		} else {
			address = c.bus.Read(uint16ToBytes(c.DE))
		}
	case "HL":
		if operand.Immediate {
			address = c.HL
		} else {
			address = c.bus.Read(uint16ToBytes(c.HL))
		}
	case "SP":
		if operand.Immediate {
			address = c.SP
		} else {
			address = c.bus.Read(uint16ToBytes(c.SP))
		}
	default:
		panic("Unknown operand type")
	}
	return address
}


/* 
	NOP: No operation, does nothing apart from incrementing the program counter
	opcodes: 0x00=NOP
	flags impacted: -
*/
func (c *CPU) NOP(instruction *Instruction) {
	c.incrementPC(uint16(instruction.Bytes))
}
func (c *CPU) PREFIX(instruction *Instruction) {
	panic("PREFIX not implemented")
}
func (c *CPU) STOP(instruction *Instruction) {
	panic("STOP not implemented")
}

// Jump / Call instructions
func (c *CPU) CALL(instruction *Instruction) {
	panic("CALL not implemented")
}

/*
	JP: Jumps to an address
	opcodes: 
		- 0xC2 = JP NZ, a16
		- 0xC3 = JP a16
		- 0xCA = JP Z, a16
		- 0xD2 = JP NC, a16
		- 0xDA = JP C, a16
		- 0xE9 = JP HL
	flags: -
*/
func (c *CPU) JP(instruction *Instruction) {
	switch c.IR {
		case 0xC2:
			// JP NZ, a16
			if !c.getZFlag() {
				operand := c.fetchOperandValue(instruction.Operands[1])
				c.PC = bytesToUint16(operand.([2]byte))
			}
		case 0xC3:
			// JP a16
			operand := c.fetchOperandValue(instruction.Operands[0])
			c.PC = bytesToUint16(operand.([2]byte))
		case 0xCA:
			// JP Z, a16
			if c.getZFlag() {
				operand := c.fetchOperandValue(instruction.Operands[1])
				c.PC = bytesToUint16(operand.([2]byte))
			}
		case 0xD2:
			// JP NC, a16
			if !c.getNFlag() {
				operand := c.fetchOperandValue(instruction.Operands[1])
				c.PC = bytesToUint16(operand.([2]byte))
			}
		case 0xDA:
			// JP C, a16
			if c.getNFlag() {
				operand := c.fetchOperandValue(instruction.Operands[1])
				c.PC = bytesToUint16(operand.([2]byte))
			}
		case 0xE9:
			// JP HL
			c.PC = c.HL
		default:
			panic("JP not implemented")
	}
}
func (c *CPU) JR(instruction *Instruction) {
	panic("JR not implemented")
}
func (c *CPU) RET(instruction *Instruction) {
	panic("RET not implemented")
}
func (c *CPU) RETI(instruction *Instruction) {
	panic("RETI not implemented")
}
func (c *CPU) RST(instruction *Instruction) {
	panic("RST not implemented")
}

// Load / Store instructions

/*
	LD: Load data from one location to another
	opcodes: 
		- 0x01 = LD BC, n16
		- 0x02 = LD [BC], A
		- 0x06 = LD B, n8
		- 0x08 = LD [a16], SP
		- 0x0A = LD A, [BC]
		- 0x0E = LD C, n8
		- 0x11 = LD DE, n16
		- 0x12 = LD [DE], A
		- 0x16 = LD D, n8
		- 0x1A = LD A, [DE]
		- 0x1E = LD E, n8
		- 0x21 = LD HL, n16
		- 0x22 = LD [HL+], A	=> increment HL after loading A into [HL]
		- 0x26 = LD H, n8
		- 0x2A = LD A, [HL+] 	=> increment HL after loading [HL] into A
		- 0x2E = LD L, n8
		- 0x31 = LD SP, n16
		- 0x32 = LD [HL-], A	=> decrement HL after loading A into [HL] ([HL]=A--)
		- 0x36 = LD [HL], n8
		- 0x3A = LD A, [HL-]	=> decrement HL after loading [HL] into A (A=[HL] & [HL]=[HL]--)
		- 0x3E = LD A, n8
		- 0x40 = LD B, B
		- 0x41 = LD B, C
		- 0x42 = LD B, D
		- 0x43 = LD B, E
		- 0x44 = LD B, H
		- 0x45 = LD B, L
		- 0x46 = LD B, [HL]
		- 0x47 = LD B, A
		- 0x48 = LD C, B
		- 0x49 = LD C, C
		- 0x4A = LD C, D
		- 0x4B = LD C, E
		- 0x4C = LD C, H
		- 0x4D = LD C, L
		- 0x4E = LD C, [HL]
		- 0x4F = LD C, A
		- 0x50 = LD D, B
		- 0x51 = LD D, C
		- 0x52 = LD D, D
		- 0x53 = LD D, E
		- 0x54 = LD D, H
		- 0x55 = LD D, L
		- 0x56 = LD D, [HL]
		- 0x57 = LD D, A
		- 0x58 = LD E, B
		- 0x59 = LD E, C
		- 0x5A = LD E, D
		- 0x5B = LD E, E
		- 0x5C = LD E, H
		- 0x5D = LD E, L
		- 0x5E = LD E, [HL]
		- 0x5F = LD E, A
		- 0x60 = LD H, B
		- 0x61 = LD H, C
		- 0x62 = LD H, D
		- 0x63 = LD H, E
		- 0x64 = LD H, H
		- 0x65 = LD H, L
		- 0x66 = LD H, [HL]
		- 0x67 = LD H, A
		- 0x68 = LD L, B
		- 0x69 = LD L, C
		- 0x6A = LD L, D
		- 0x6B = LD L, E
		- 0x6C = LD L, H
		- 0x6D = LD L, L
		- 0x6E = LD L, [HL]
		- 0x6F = LD L, A
		- 0x70 = LD [HL], B
		- 0x71 = LD [HL], C
		- 0x72 = LD [HL], D
		- 0x73 = LD [HL], E
		- 0x74 = LD [HL], H
		- 0x75 = LD [HL], L
		- 0x77 = LD [HL], A
		- 0x78 = LD A, B
		- 0x79 = LD A, C
		- 0x7A = LD A, D
		- 0x7B = LD A, E
		- 0x7C = LD A, H
		- 0x7D = LD A, L
		- 0x7E = LD A, [HL]
		- 0x7F = LD A, A
		- 0xE2 = LD [C], A
		- 0xEA = LD [a16], A
		- 0xF2 = LD A, [C]
		- 0xF8 = LD HL, SP+e8	=> add signed immediate to SP and store the result in HL
		- 0xF9 = LD SP, HL
		- 0xFA = LD A, [a16]
	flags: - except for 0xF8 where Z->0 N->0 H->C C->C

	NOTE: all LD instructions have 2 operands, the first one is always the destination and the second one is always the source (except for LD HL, SP+r8)
	=> we will 'automate' the process of fetching the operands expect for LD HL, SP+r8 that will be handled manually
*/
func (c *CPU) LD(instruction *Instruction) {
	if c.IR != 0xF8 {
		// fetch the first operand as an address
		address := c.fetchOperandAddress(instruction.Operands[0])
		// fetch the second operand as a value
		value := c.fetchOperandValue(instruction.Operands[1])
		
		
}
func (c *CPU) LDH(instruction *Instruction) {
	panic("LDH not implemented")
}
func (c *CPU) PUSH(instruction *Instruction) {
	panic("PUSH not implemented")
}
func (c *CPU) POP(instruction *Instruction) {
	panic("POP not implemented")
}

// Arithmetic / Logical instructions
func (c *CPU) ADC(instruction *Instruction) {
	panic("ADC not implemented")
}
func (c *CPU) ADD(instruction *Instruction) {
	panic("ADD not implemented")
}
func (c *CPU) AND(instruction *Instruction) {
	panic("AND not implemented")
}
func (c *CPU) CCF(instruction *Instruction) {
	panic("CCF not implemented")
}
func (c *CPU) CP(instruction *Instruction) {
	panic("CP not implemented")
}
func (c *CPU) CPL(instruction *Instruction) {
	panic("CPL not implemented")
}
func (c *CPU) DAA(instruction *Instruction) {
	panic("DAA not implemented")
}
func (c *CPU) DEC(instruction *Instruction) {
	panic("DEC not implemented")
}
func (c *CPU) INC(instruction *Instruction) {
	panic("INC not implemented")
}
func (c *CPU) SUB(instruction *Instruction) {
	panic("SUB not implemented")
}
func (c *CPU) SBC(instruction *Instruction) {
	panic("SBC not implemented")
}
func (c *CPU) SCF(instruction *Instruction) {
	panic("SCF not implemented")
}
func (c *CPU) OR(instruction *Instruction) {
	panic("OR not implemented")
}
/*
	XOR: Bitwise XOR
	opcodes: 
		- 0xA8 = XOR A, B
		- 0xA9 = XOR A, C
		- 0xAA = XOR A, D
		- 0xAB = XOR A, E
		- 0xAC = XOR A, H
		- 0xAD = XOR A, L
		- 0xAE = XOR A, [HL]
		- 0xAF = XOR A, A
		- 0xEE = XOR A, n8
	flags: Z->Z N->0 H->0 C->0
	note: 0xAF XOR 0xAF = 0x00 (Z flag is always set)
*/
func (c *CPU) XOR(instruction *Instruction) {
	c.A = c.A ^ c.fetchOperandValue(instruction.Operands[0]).(byte)
	c.incrementPC(uint16(instruction.Bytes))
}

// Shift / Rotate and Bit instructions
func (c *CPU) RLA(instruction *Instruction) {
	panic("RLA not implemented")
}
func (c *CPU) RLCA(instruction *Instruction) {
	panic("RLCA not implemented")
}
func (c *CPU) RRA(instruction *Instruction) {
	panic("RRA not implemented")
}
func (c *CPU) RRCA(instruction *Instruction) {
	panic("RRCA not implemented")
}

// > instructions handlers (PREFIX CB)
func (c *CPU) BIT(instruction *Instruction) {
	panic("BIT not implemented")
}
func (c *CPU) RES(instruction *Instruction) {
	panic("RES not implemented")
}
func (c *CPU) RL(instruction *Instruction) {
	panic("RL not implemented")
}
func (c *CPU) RLC(instruction *Instruction) {
	panic("RLC not implemented")
}
func (c *CPU) RR(instruction *Instruction) {
	panic("RR not implemented")
}
func (c *CPU) RRC(instruction *Instruction) {
	panic("RRC not implemented")
}
func (c *CPU) SET(instruction *Instruction) {
	panic("SET not implemented")
}
func (c *CPU) SLA(instruction *Instruction) {
	panic("SLA not implemented")
}
func (c *CPU) SRA(instruction *Instruction) {
	panic("SRA not implemented")
}
func (c *CPU) SRL(instruction *Instruction) {
	panic("SRL not implemented")
}
func (c *CPU) SWAP(instruction *Instruction) {
	panic("SWAP not implemented")
}
