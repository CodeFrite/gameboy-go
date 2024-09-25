// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

const (
	// start and length of the cpu internal memory regions
	BOOT_ROM_START     uint16 = 0x0000
	BOOT_ROM_LEN       uint16 = 0x0100
	IO_REGISTERS_START uint16 = 0xFF00
	IO_REGISTERS_LEN   uint16 = 0x0080
	HRAM_START         uint16 = 0xFF80
	HRAM_LEN           uint16 = 0x007F
	IE_FLAG_START      uint16 = 0xFFFF
	IE_FLAG_LEN        uint16 = 0x0001
)

/*
 * CPU: executes instructions fetched from memory, reads and writes to memory (internal registers, flags & bus)
 */
type CPU struct {
	// Registers
	PC               Register16 // Program Counter
	SP               Register16 // Stack Pointer
	A                Register8  // Accumulator
	F                Register8  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	B, C, D, E, H, L Register8  // 16-bit general purpose registers

	// Instruction
	IR        Register8 // Instruction Register
	Prefixed  bool      // Is the current instruction prefixed with 0xCB
	Operand   uint16    // Current operand fetched from memory (this register doesn't physically exist in the CPU)
	Offset    uint16    // offset used in some instructions
	CpuCycles uint16    // number of cycles the CPU has executed TODO: change to the correct type and implement the interrupt (overflow) handling

	// Interrupts
	IME                    bool // interrupt master enable
	IME_ENABLE_NEXT_CYCLE  bool // enable the IME on the next cycle
	IME_DISABLE_NEXT_CYCLE bool // disable the IME on the next cycle
	Halted                 bool // is the CPU halted (waiting for an interrupt to wake up)
	Stopped                bool // is the CPU stopped (waiting for an interrupt from the joypad)

	// CPU SoC Internal Memories
	bus          *Bus    // reference to the bus
	bootrom      *Memory // 0x0000-0x00FF: (256 bytes) - Boot ROM
	io_registers *Memory // 0xFF00-0xFF7F: (128 bytes) - I/O Registers
	hram         *Memory // 0xFF80-0xFFFE: (127 bytes) - High RAM
	IE           *Memory // 0xFFFF: Interrupt Enable
}

// Create a new CPU
func NewCPU(bus *Bus) *CPU {

	randValue := func(base int, exponent int) int {
		return rand.Intn(int(math.Pow(float64(base), float64(exponent))))
	}

	cpu := &CPU{
		bus: bus,
		// on startup, simulate the CPU registers being in an unknown state
		PC: 0x0000, // only value set by the cpu on startup
		SP: Register16(randValue(2, 16)),
		A:  Register8(randValue(2, 8)),
		F:  Register8(randValue(2, 8)),
		B:  Register8(randValue(2, 8)),
		C:  Register8(randValue(2, 8)),
		D:  Register8(randValue(2, 8)),
		E:  Register8(randValue(2, 8)),
		H:  Register8(randValue(2, 8)),
		L:  Register8(randValue(2, 8)),
	}

	return cpu
}

// initialize the bootrom & memories
func (c *CPU) Init() {
	c.initBootRom()
	c.initMemory()
}

/**
 * initializes the bootrom @ 0x0000
 */
func (c *CPU) initBootRom() {
	bootRomData, err := LoadRom("/Users/codefrite/Desktop/CODE/codefrite-emulator/gameboy/gameboy-go/roms/dmg_boot.bin")
	if err != nil {
		log.Fatal(err)
	}
	c.bootrom = NewMemory(BOOT_ROM_LEN)
	c.bus.AttachMemory("Boot ROM", BOOT_ROM_START, c.bootrom)
	c.bus.WriteBlob(BOOT_ROM_START, bootRomData)
}

/**
 * initializes the memories and attaches them to the bus
 * HRAM: 127 bytes @ 0xFF80
 * VRAM: 8KB bytes @ 0x8000
 * WRAM: 8KB @ 0xC000
 * I/O Registers: 128 bytes @ 0xFF00
 */
func (c *CPU) initMemory() {
	// initialize memories
	c.hram = NewMemory(HRAM_LEN)                 // High RAM (127 bytes)
	c.io_registers = NewMemory(IO_REGISTERS_LEN) // I/O Registers (128 bytes)
	c.IE = NewMemory(IE_FLAG_LEN)                // Interrupt Enable Register (1 byte)

	// attach memories to the bus
	c.bus.AttachMemory("High RAM (HRAM)", HRAM_START, c.hram)
	c.bus.AttachMemory("I/O Registers", IO_REGISTERS_START, c.io_registers)
	c.bus.AttachMemory("Interrupt Enable Register", IE_FLAG_START, c.IE)
}

// Increment the Program Counter by the given offset
func (c *CPU) updatePC() {
	c.PC.Set(c.Offset)
}

// Stack operations

// Push a value to the stack
func (c *CPU) push(value uint16) {
	// decrement the stack pointer
	c.SP -= 1
	// write the high byte to the stack
	c.bus.Write(c.SP.Get(), byte(value>>8))
	// decrement the stack pointer
	c.SP -= 1
	// write the low byte to the stack
	c.bus.Write(c.SP.Get(), byte(value))
}

// Pop a value from the stack
func (c *CPU) pop() uint16 {
	// read the low byte from the stack
	low := c.bus.Read(c.SP.Get())
	// increment the stack pointer
	c.SP += 1
	// read the high byte from the stack
	high := c.bus.Read(c.SP.Get())
	// increment the stack pointer
	c.SP += 1
	return uint16(high)<<8 | uint16(low)
}

// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() (opcode uint8, prefixed bool) {
	// Fetch the opcode from memory at the address in the program counter
	opcode = c.bus.Read(c.PC.Get())

	// is it a prefixed instruction?
	if opcode == 0xCB {
		prefixed = true
		// fetch the next opcode
		opcode = c.bus.Read(c.PC.Get() + 1)
	}
	return opcode, prefixed
}

/*
 * Fetch the value of an operand
 * Save the result in cpu.Operand as an uint16 (must be casted to the correct type inside the different instruction handlers)
 */
func (c *CPU) fetchOperandValue(operand Operand) uint16 {
	var value, addr uint16
	switch operand.Name {

	// n8: immediate 8-bit data
	case "n8":
		value = uint16(c.bus.Read(c.PC.Get() + 1))

	// n16: immediate little-endian 16-bit data
	case "n16":
		value = c.bus.Read16(c.PC.Get() + 1)

	// a8: 8-bit unsigned data, which is added to $FF00 in certain instructions to create a 16-bit address in HRAM (High RAM)
	case "a8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.PC.Get() + 1))
		} else {
			//addr = 0xFF00 + c.bus.Read16(c.PC+1)
			addr = 0xFF00 + uint16(c.bus.Read(c.PC.Get()+1))
			value = uint16(c.bus.Read(addr))
		}
	// a16: little-endian 16-bit address
	case "a16": // not always immediate
		if operand.Immediate {
			value = c.bus.Read16(c.PC.Get() + 1)
		} else {
			addr := c.bus.Read16(c.PC.Get() + 1)
			value = c.bus.Read16(addr)
		}
	// e8 means 8-bit signed data
	case "e8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.PC.Get() + 1))
		} else {
			panic("e8 non immediate operand not implemented yet")
		}
	case "A":
		if operand.Immediate {
			value = uint16(c.A)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "B":
		if operand.Immediate {
			value = uint16(c.B)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "C":
		if operand.Immediate {
			value = uint16(c.C)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "D":
		if operand.Immediate {
			value = uint16(c.D)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "E":
		if operand.Immediate {
			value = uint16(c.E)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "H":
		if operand.Immediate {
			value = uint16(c.H)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "L":
		if operand.Immediate {
			value = uint16(c.L)
		} else {
			panic("Non immediate operand not implemented yet")
		}

	case "BC":
		if operand.Immediate {
			value = c.getBC()
		} else {
			value = c.bus.Read16(c.getBC())
		}
	case "DE":
		if operand.Immediate {
			value = c.getDE()
		} else {
			value = c.bus.Read16(c.getDE())
		}
	case "HL":
		if operand.Immediate {
			value = c.getHL()
		} else {
			value = c.bus.Read16(c.getHL())
		}
	case "SP": // always immediate
		value = c.SP.Get()
	case "$00": // RST $00
		value = 0x00
	case "$08": // RST $08
		value = 0x08
	case "$10": // RST $10
		value = 0x10
	case "$18": // RST $18
		value = 0x18
	case "$20": // RST $20
		value = 0x20
	case "$28": // RST $28
		value = 0x28
	case "$30": // RST $30
		value = 0x30
	case "$38": // RST $38
		value = 0x38
	case "Z":
		value = uint16(c.F & 0x80)
	case "NZ":
		value = uint16(c.F & 0x80)
	case "NC":
		value = uint16(c.F & 0x10)
	default:
		err := fmt.Sprintf("cpu.fetchOperandValue> Unknown operand name: %s (0x%02X)", operand.Name, c.IR)
		panic(err)
	}
	return value
}

// Execute one cycle of the CPU: fetch, decode and execute the next instruction
// TODO: i am supposed to return an error but i am always returning nil. Chose an error handling strategy and implement it
func (c *CPU) Step() error {
	// return if the CPU is halted or stopped
	if c.Halted || c.Stopped {
		return nil
	}

	// update the pc and reset the offset
	c.updatePC()
	c.Offset = 0

	// reset the prefixed flag
	c.Prefixed = false

	// Store the opcode in the instruction register & prefix flag
	opCode, prefixed := c.fetchOpcode()
	c.Prefixed = prefixed
	c.IR.Set(opCode)

	// Decode the instruction
	// get instruction from opcodes.json file with IR used as key
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.IR)), c.Prefixed)
	// get the operands of the instruction
	operands := instruction.Operands
	// fetch the operand value
	if len(operands) == 1 {
		c.Operand = c.fetchOperandValue(operands[0])
	} else if len(operands) == 2 {
		// decode operand 2
		c.Operand = c.fetchOperandValue(operands[1])
	}

	// Handle the IME
	if c.IME_ENABLE_NEXT_CYCLE {
		// Execute the instruction
		if !c.Prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
		// enable the IME
		c.IME = true
		c.IME_ENABLE_NEXT_CYCLE = false
	} else if c.IME_DISABLE_NEXT_CYCLE {
		// Execute the instruction
		if !c.Prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
		// disable the IME
		c.IME = false
		c.IME_DISABLE_NEXT_CYCLE = false
	} else {
		// Execute the instruction
		if !c.Prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
	}

	return nil
}

// Run the CPU
func (c *CPU) Run() {
	for {
		if c.Halted {
			// if the CPU is halted, wiat for an interrupt to wake it up
			// TODO: implement the interrupt handling
			// ! for the moment, we will break the loop to avoid an infinite loop
			break
		}
		if c.Stopped {
			// if the CPU is stopped, wait for an interrupt from the joypad
			break
		}
		// Execute the next instruction
		if err := c.Step(); err != nil {
			panic(err)
		}
	}
}

// Boot the CPU and returns when the boot process is done (PC=0x0100)
func (c *CPU) Boot() {
	c.PC = 0x0000
	for c.PC != 0x0100 {
		err := c.Step()
		if err != nil {
			panic(err)
		}
	}
}
