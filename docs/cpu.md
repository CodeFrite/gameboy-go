# CPU

The gameboy CPU is an 8-bit processor with a 16-bit address bus. This means that the CPU can access up to 2^16 = 65536 memory location of one byte each: this means that programs can be up to 64KB in size.
But the gameboy make a clever use of banks to allow for larger programs (typically up to 2MB) by using cartridge internal ROM and RAM banks.

## Registers

Since the gameboy CPU has an 8bit architecture, I decided to represent the CPU registers as follows:

```go
type CPU struct {
  PC uint16 					// Program Counter
	SP uint16 					// Stack Pointer
	IR byte 						// Instruction Register
	A, F byte 					// Accumulator & Flags: Zero (position 7), Subtraction (position 6), Half Carry (position 5), Carry (position 4)
  // 8-bits general purpose registers
	B, C byte
  D, E byte
  H, L byte
	IE byte 						// Interrupt Enable
  ...
}
```

As you can see, all registers are declared as byte except for the Program Counter (PC) and Stack Pointer (SP) which are declared as uint16. Why this choice? Why not declaring the registers AF, BC, DE, HL as uint16 as well? Why not declaring them as [2]byte?

The short answer would be that it doesn't really matter. As long as we are consistent in our implementation and make sure to handle the registers correctly, we can choose any representation we want.

The long answer is that after testing a few different versions of the code, I found it more convenient to have `memory` as well as registers as byte. Moreover, the gameboy instruction set has instructions that acts on the 8-bit registers as well as the 16-bit registers. So in the end, we will have to handle conversions between the two representations anyway.

The reason why i chose to declare PC and SP as uint16 and not as [2]byte is that it makes it easier to handle the incrementation and decrementation of the registers. Finally, setting the value of these registers as [2]byte in Go is a bit more cumbersome than setting them as uint16:

```go
cpu.PC = [2]byte{0x00, 0x00} // Cumbersome
cpu.PC = 0x0000 // Easier
```

This is for example less of an issue in a programming language like C++ where we can make use of struct unions to define at the same time the 16-bit and 8-bit representation of the registers, which grants us direct access to both representations:

```cpp
union {
  struct {
    uint8_t A;
    uint8_t F;
  };
  uint16_t AF;
} pc;

// access the full AF register
pc.AF = 0x00FF;

// access the A register
pc.A = 0x00;

// access the F register
pc.F = 0xFF;
```

### Flags

The flags register is a special register that contains 4 flags: Zero, Subtraction, Half Carry, Carry. These flags are used to store the result of the last operation and are used by the CPU to make decisions when executing instructions. They are located in the F register as follows:

```
7 6 5 4 3 2 1 0
Z N H C 0 0 0 0
```

In order to manipulate these flags, we will use bitwise operations. Indeed, if we consider the atomic case of a single bit operand and see what the &, |, ^ operators, we can see their effect on the bit:

| A   | B   | A & B | A \| B | A ^ B |
| --- | --- | ----- | ------ | ----- |
| 0   | 0   | 0     | 0      | 0     |
| 1   | 0   | 0     | 1      | 1     |
| 0   | 1   | 0     | 1      | 1     |
| 1   | 1   | 1     | 1      | 0     |

If we consider the two first lines and two last lines of the table, we can deduce for each operator the function they perform:

| Operator | Value | Function |
| -------- | ----- | -------- |
| &        | 0     | Reset    |
| &        | 1     | Keep     |
| \|       | 0     | Keep     |
| \|       | 1     | Set      |
| ^        | 0     | Keep     |
| ^        | 1     | Toggle   |

Therefore, we will use the `&` operator to reset a flag, the `|` operator to set a flag and the `^` operator to toggle a flag.
When it comes to checking the value of a flag (check wether it is set or not), we can use any of these operators. In my implementation, I will be using the `&`.

`Setting a flag`

In order to set a flag, we will be using the bitwise OR operator `|`:

```go
cpu.F |= 0x80 // Set the Zero flag
cpu.F |= 0x40 // Set the Subtraction flag
cpu.F |= 0x20 // Set the Half Carry flag
cpu.F |= 0x10 // Set the Carry flag
```

`Clearing a flag`

In order to clear a flag, we will be using the bitwise AND operator `&`:

```go
cpu.F &= 0x80 // Clear the Zero flag
cpu.F &= 0x40 // Clear the Subtraction flag
cpu.F &= 0x20 // Clear the Half Carry flag
cpu.F &= 0x10 // Clear the Carry flag
```

`Checking a flag`

In order to check a flag, we will be using the bitwise AND operator `&`:

```go
cpu.F & 0x80 == 0x80 // Check the Zero flag
cpu.F & 0x40 == 0x40 // Check the Subtraction flag
cpu.F & 0x20 == 0x20 // Check the Half Carry flag
cpu.F & 0x10 == 0x10 // Check the Carry flag
```

`Toggling a flag`

In order to toggle a flag, we will be using the bitwise XOR operator `^`:

```go
cpu.F ^= 0x80 // Toggle the Zero flag
cpu.F ^= 0x40 // Toggle the Subtraction flag
cpu.F ^= 0x20 // Toggle the Half Carry flag
cpu.F ^= 0x10 // Toggle the Carry flag
```
