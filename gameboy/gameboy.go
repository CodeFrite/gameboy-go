package gameboy

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Gameboy struct {
	cpu          *CPU
	bootrom      *Memory
	cartridge    *Cartridge
	vram         *Memory
	wram         *Memory
	io_registers *Memory
	hram         *Memory
	bus          *Bus
}

type CPURegistersState struct {
	PC uint16 `json:"PC"` // Program Counter
	SP uint16 `json:"SP"` // Stack Pointer
	A  uint8  `json:"A"`  // Accumulator
	F  uint8  `json:"F"`  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	// 16-bits general purpose registers
	BC uint16 `json:"BC"`
	DE uint16 `json:"DE"`
	HL uint16 `json:"HL"`

	IE uint8 `json:"IE"` // Interrupt Enable

	IR       uint8  `json:"IR"`       // Instruction Register
	PREFIXED bool   `json:"PREFIXED"` // Is the current instruction prefixed with 0xCB
	OPERAND  uint16 `json:"OPERAND"`  // Current operand fetched from memory (this register doesn't physically exist in the CPU)

	IME    bool `json:"IME"`    // interrupt master enable
	HALTED bool `json:"HALTED"` // is the CPU halted
}

type currentInstruction struct {
	PC       uint16   `json:"PC"`       // Program Counter
	IR       uint16   `json:"IR"`       // Stack Pointer
	BIN      []byte   `json:"BIN"`      // Binary instruction
	MNE      string   `json:"MNE"`      // Mnemonic
	OPERANDS []string `json:"OPERANDS"` // Operands
}

func NewGameboy() *Gameboy {
	g := &Gameboy{}
	g.initMemory()
	g.initBus()
	g.initCPU()
	return g
}

// Initializers

func (g *Gameboy) initMemory() {
	// Bootrom
	g.bootrom = NewMemory(0x0100)
	// VRAM
	g.vram = NewMemory(0x2000)
	// WRAM
	g.wram = NewMemory(0x2000)
	// I/O Registers
	g.io_registers = NewMemory(0x007F)
	// high ram
	g.hram = NewMemory(0x007F)
}

func (g *Gameboy) initBus() {
	g.bus = NewBus()
	g.bus.AttachMemory(0x8000, g.vram)
	g.bus.AttachMemory(0xC000, g.wram)
	g.bus.AttachMemory(0xFF00, g.io_registers)
	g.bus.AttachMemory(0xFF80, g.hram)
}

func (g *Gameboy) initCPU() {
	g.cpu = NewCPU(g.bus)
}

// Utility functions
func getBootRomData() []uint8 {
	rom, err := LoadRom("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	return rom
}

func (g *Gameboy) loadCartridge(uri string, name string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	g.cartridge = NewCartridge(currentDir+uri, name)
	g.bus.AttachMemory(0x0100, g.cartridge)
}

func printCurrentInstruction(c *CPU) {
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), c.prefixed)
	// get the bytes following the opcode and corresponding to the operands, if any
	getBytes := func() []byte {
		var bytes []byte
		for i := 0; i < instruction.Bytes; i++ {
			bytes = append(bytes, c.bus.Read(c.PC+uint16(i)))
		}
		return bytes
	}
	//
	getOperands := func() string {
		var operands []string
		for _, operand := range instruction.Operands {
			var value string
			if operand.Name == "n8" {
				value = fmt.Sprintf("$%02X", c.bus.Read(c.PC+1))
			} else if operand.Name == "n16" {
				value = fmt.Sprintf("$%04X", c.bus.Read16(c.PC+1))
			} else {
				value = operand.Name
			}

			if operand.Increment {
				value += "+"
			} else if operand.Decrement {
				value += "-"
			}

			if !operand.Immediate {
				value = "[" + value + "]"
			}

			operands = append(operands, value)
		}
		return strings.Join(operands, ", ")
	}

	fmt.Printf("PC: 0x%04X, IR: 0x%02X", c.PC, c.IR)
	fmt.Printf(", memory: %-6X", getBytes())
	fmt.Printf(", asm: %s %s\n", instruction.Mnemonic, getOperands())
}

func printRegisters(c *CPU) {
	fmt.Printf("A: 0x%02X, B: 0x%02X, C: 0x%02X, D: 0x%02X, E: 0x%02X, H: 0x%02X, L: 0x%02X; SP: 0x%04X \n", c.A, c.B, c.C, c.D, c.E, c.H, c.L, c.SP)
	fmt.Printf("F: 0x%02X, Z: %t, N: %t, H: %t, C: %t\n\n", c.F, c.getZFlag(), c.getNFlag(), c.getHFlag(), c.getCFlag())
}

func (g *Gameboy) Step() *CPURegistersState {
	PC := g.cpu.PC // since the CPU will increment the PC after the instruction is executed, when save it before
	printCurrentInstruction(g.cpu)
	g.cpu.step()
	printRegisters(g.cpu)
	return &CPURegistersState{
		PC: PC,
		SP: g.cpu.SP,
		A:  g.cpu.A,
		F:  g.cpu.F,
		BC: uint16(g.cpu.B)<<8 | uint16(g.cpu.C),
		DE: uint16(g.cpu.D)<<8 | uint16(g.cpu.E),
		HL: uint16(g.cpu.H)<<8 | uint16(g.cpu.L),
		IE: g.cpu.IE,
		IR: g.cpu.IR,
		//PREFIXED: g.cpu.PREFIXED,
		OPERAND: g.cpu.Operand,
		IME:     g.cpu.IME,
		//HALTED:   g.cpu.HALTED,
	}
}

/*
 * Run the bootrom and then the game
 */
func (g *Gameboy) Run(uri string, name string) {
	g.loadCartridge(uri, name)
	g.cpu.Boot()
}

func (g *Gameboy) Init() {
	bootRomData := getBootRomData()
	g.bus.AttachMemory(0x0000, g.bootrom)
	g.bus.WriteBlob(0x0000, bootRomData)
}
