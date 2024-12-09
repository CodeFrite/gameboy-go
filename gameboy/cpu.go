// Computing Processing Unit (CPU) for the Gameboy
package gameboy

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	// start and length of the memory regions inside memory map
	IO_REGISTERS_START uint16 = 0xFF00
	IO_REGISTERS_LEN   uint16 = 0x0080
	HRAM_START         uint16 = 0xFF80
	HRAM_LEN           uint16 = 0x007F
	IE_FLAG_START      uint16 = 0xFFFF
	IE_FLAG_LEN        uint16 = 0x0001

	CPU_EXECUTION_STATE_FREE   CPU_EXECUTION_STATE = "free"
	CPU_EXECUTION_STATE_LOCKED CPU_EXECUTION_STATE = "locked"
)

type CPU_EXECUTION_STATE = string

/*
 * CPU: executes instructions fetched from memory, reads and writes to memory (internal registers, flags & bus)
 */
type CPU struct {
	// lock the CPU to prevent concurrent access
	busyChannel chan bool

	state CPU_EXECUTION_STATE // CPU state (locked, free)

	// Work Registers (not mapped to memory)
	pc               uint16 // Program Counter
	sp               uint16 // Stack Pointer
	a                uint8  // Accumulator
	f                uint8  // Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
	b, c, d, e, h, l uint8  // 16-bit general purpose registers
	ir               uint8  // Instruction Register

	// Work variables
	instruction Instruction // Current instruction
	prefixed    bool        // Is the current instruction prefixed with 0xCB
	operand     uint16      // Current operand fetched from memory (this register doesn't physically exist in the CPU)
	offset      uint16      // offset used in some instructions
	cpuCycles   uint64      // number of cycles the CPU has executed since the last reset up to uint64 max value (18,446,744,073,709,551,615 =

	// Interrupts
	ime                    bool // interrupt master enable
	ime_enable_next_cycle  bool // enable the IME on the next cycle
	ime_disable_next_cycle bool // disable the IME on the next cycle
	halted                 bool // is the CPU halted (waiting for an interrupt to wake up)
	stopped                bool // is the CPU & LCD stopped (waiting for an interrupt from the joypad)

	// CPU SoC Internal Memories (not exported in json)
	bus          *Bus    // reference to the bus
	io_registers *Memory // 0xFF00-0xFF7F: (128 bytes) - I/O Registers
	hram         *Memory // 0xFF80-0xFFFE: (127 bytes) - High RAM
	ie           *Memory // 0xFFFF: Interrupt Enable
}

// Create a new CPU
func NewCPU(bus *Bus) *CPU {

	randValue := func(base int, exponent int) int {
		return rand.Intn(int(math.Pow(float64(base), float64(exponent))))
	}

	cpu := &CPU{
		busyChannel: make(chan bool, 1),
		state:       CPU_EXECUTION_STATE_FREE,
		bus:         bus,
		// on startup, simulate the CPU registers being in an unknown state
		cpuCycles: 0,
		pc:        0x0000, // only value set by the cpu on startup, others are randomized
		sp:        uint16(randValue(2, 16)),
		a:         uint8(randValue(2, 8)),
		f:         uint8(randValue(2, 8)),
		b:         uint8(randValue(2, 8)),
		c:         uint8(randValue(2, 8)),
		d:         uint8(randValue(2, 8)),
		e:         uint8(randValue(2, 8)),
		h:         uint8(randValue(2, 8)),
		l:         uint8(randValue(2, 8)),
	}

	return cpu
}

// initializes the memories and attaches them to the bus
//   - HRAM: 127 bytes @ 0xFF80
//   - VRAM: 8KB bytes @ 0x8000
//   - WRAM: 8KB @ 0xC000
//   - I/O Registers: 128 bytes @ 0xFF00
func (c *CPU) init() {
	// initialize memories
	c.hram = NewMemory(HRAM_LEN)                 // High RAM (127 bytes)
	c.io_registers = NewMemory(IO_REGISTERS_LEN) // I/O Registers (128 bytes)
	c.ie = NewMemory(IE_FLAG_LEN)                // Interrupt Enable Register (1 byte)

	// attach memories to the bus
	c.bus.AttachMemory("High RAM (HRAM)", HRAM_START, c.hram)
	c.bus.AttachMemory("I/O Registers", IO_REGISTERS_START, c.io_registers)
	c.bus.AttachMemory("Interrupt Enable Register", IE_FLAG_START, c.ie)
}

// Increment the Program Counter by the given offset
func (c *CPU) updatepc() {
	c.pc = c.offset
}

// Stack operations

// Push a value to the stack
func (c *CPU) push(value uint16) {
	// decrement the stack pointer
	c.sp = c.sp - 1
	// write the high byte to the stack
	c.bus.Write(c.sp, byte(value>>8))
	// decrement the stack pointer
	c.sp = c.sp - 1
	// write the low byte to the stack
	c.bus.Write(c.sp, byte(value))
}

// Pop a value from the stack
func (c *CPU) pop() uint16 {
	// read the low byte from the stack
	low := c.bus.Read(c.sp)
	// increment the stack pointer
	c.sp += 1
	// read the high byte from the stack
	high := c.bus.Read(c.sp)
	// increment the stack pointer
	c.sp += 1
	return uint16(high)<<8 | uint16(low)
}

// Fetch the opcode from bus at address pc and store it in the instruction register
func (c *CPU) fetchOpcode() (opcode uint8, prefixed bool) {
	// Fetch the opcode from memory at the address in the program counter
	opcode = c.bus.Read(c.pc)

	// is it a prefixed instruction?
	if opcode == 0xCB {
		prefixed = true
		// fetch the next opcode
		opcode = c.bus.Read(c.pc + 1)
	}
	return opcode, prefixed
}

/*
 * Fetch the value of an operand
 * Save the result in cpu.operand as an uint16 (must be casted to the correct type inside the different instruction handlers)
 */
func (c *CPU) fetchOperandValue(operand Operand) uint16 {
	var value, addr uint16
	switch operand.Name {

	// n8: immediate 8-bit data
	case "n8":
		value = uint16(c.bus.Read(c.pc + 1))

	// n16: immediate little-endian 16-bit data
	case "n16":
		value = c.bus.Read16(c.pc + 1)

	// a8: 8-bit unsigned data, which is added to $FF00 in certain instructions to create a 16-bit address in HRAM (High RAM)
	case "a8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.pc + 1))
		} else {
			addr = 0xFF00 + uint16(c.bus.Read(c.pc+1))
			value = uint16(c.bus.Read(addr))
		}
	// a16: little-endian 16-bit address
	case "a16": // not always immediate
		if operand.Immediate {
			value = c.bus.Read16(c.pc + 1)
		} else {
			addr := c.bus.Read16(c.pc + 1)
			value = c.bus.Read16(addr)
		}
	// e8 means 8-bit signed data
	case "e8": // not always immediate
		if operand.Immediate {
			value = uint16(c.bus.Read(c.pc + 1))
		} else {
			panic("e8 non immediate operand not implemented yet")
		}
	case "A":
		if operand.Immediate {
			value = uint16(c.a)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "B":
		if operand.Immediate {
			value = uint16(c.b)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "C":
		if operand.Immediate {
			value = uint16(c.c)
		} else {
			addr = 0xFF00 + uint16(c.c)
			value = uint16(c.bus.Read(addr))
		}
	case "D":
		if operand.Immediate {
			value = uint16(c.d)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "E":
		if operand.Immediate {
			value = uint16(c.e)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "H":
		if operand.Immediate {
			value = uint16(c.h)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "L":
		if operand.Immediate {
			value = uint16(c.l)
		} else {
			panic("Non immediate operand not implemented yet")
		}
	case "AF":
		if operand.Immediate {
			value = uint16(c.a)<<8 | uint16(c.f)
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
		// increment or decrement the value of HL
		if operand.Increment {
			c.setHL(c.getHL() + 1)
		} else if operand.Decrement {
			c.setHL(c.getHL() - 1)
		}
	case "SP": // always immediate
		value = c.sp
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

	// flags
	case "flag_Z":
		if c.getZFlag() {
			value = uint16(1)
		} else {
			value = uint16(0)
		}
	case "flag_NZ":
		if c.getZFlag() {
			value = uint16(1)
		} else {
			value = uint16(0)
		}
	case "flag_C":
		if c.getCFlag() {
			value = uint16(1)
		} else {
			value = uint16(0)
		}
	case "flag_NC":
		if c.getCFlag() {
			value = uint16(1)
		} else {
			value = uint16(0)
		}
	default:
		err := fmt.Sprintf("cpu.fetchOperandValue> Unknown operand name: %s (0x%02X)", operand.Name, c.ir)
		panic(err)
	}
	return value
}

// Execute one cycle of the CPU: fetch, decode and execute the next instruction
// TODO: i am supposed to return an error but i am always returning nil. Chose an error handling strategy and implement it
func (c *CPU) Step() error {

	// return if the CPU is halted or stopped
	if c.halted || c.stopped {
		return nil
	}

	// update the pc and reset the offset
	c.updatepc()
	c.offset = 0

	// reset the prefixed flag
	c.prefixed = false

	// Store the opcode in the instruction register & prefix flag
	opCode, prefixed := c.fetchOpcode()
	c.prefixed = prefixed
	c.ir = opCode

	// Decode the instruction
	// get instruction from opcodes.json file with IR used as key
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", c.ir)), c.prefixed)
	c.instruction = instruction
	// get the operands of the instruction
	operands := instruction.Operands
	// fetch the last operand value
	idx := len(operands) - 1
	if idx >= 0 {
		c.operand = c.fetchOperandValue(operands[idx])
	}

	// Handle the IME
	if c.ime_enable_next_cycle {
		// Execute the instruction
		if !c.prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
		// enable the IME
		c.ime = true
		c.ime_enable_next_cycle = false
	} else if c.ime_disable_next_cycle {
		// Execute the instruction
		if !c.prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
		// disable the IME
		c.ime = false
		c.ime_disable_next_cycle = false
	} else {
		// Execute the instruction
		if !c.prefixed {
			c.executeInstruction(instruction)
		} else {
			c.executeCBInstruction(instruction)
		}
	}

	return nil
}

// Run the CPU
func (c *CPU) Run() {
	// return if CPU is locked, otherwise lock CPU and run
	if c.state == CPU_EXECUTION_STATE_LOCKED {
		fmt.Println("CPU is locked")
		return
	} else {
		c.state = CPU_EXECUTION_STATE_LOCKED
	}
	for {
		if c.halted {
			// if the CPU is halted, wait for an interrupt to wake it up
			// TODO: implement the interrupt handling
			// ! for the moment, we will break the loop to avoid an infinite loop
			break
		}
		if c.stopped {
			// if the CPU is stopped, wait for an interrupt from the joypad
			break
		}
		// Execute the next instruction
		if err := c.Step(); err != nil {
			panic(err)
		}
	}

	// unlock the CPU
	c.state = CPU_EXECUTION_STATE_FREE
}

// Boot the CPU and returns when the boot process is done (pc=0x0100)
func (c *CPU) Boot() {
	c.pc = 0x0000
	for c.pc != 0x0100 {
		err := c.Step()
		if err != nil {
			panic(err)
		}
	}
}
