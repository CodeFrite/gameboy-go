# Instructions

In this document we will go through the different instructions and see what they do. We will understand how a game is executed from the ROM and discuss how to approach their implementation in Go.

## Ressources

In order to make sense of the gameboy instructions set, I am basing myself on the following resources:

- [optables](https://gbdev.io/gb-opcodes/optables/): It presents the instructions set in a concise and easy to read table format, listing all the instructions in an octal format. However, apart from explaining the notation used, it does not provide any further explanation on how to interpret the instructions.
- [RGBDS](https://rgbds.gbdev.io/docs/v0.7.0/gbz80.7): This ressource list the instructions and regroup the related ones in categories that reflect their purpose or the structure they operate on. It also provides a brief explanation of what each instruction does and how to interpret them. It is a very good starting point to understand the instructions set and the CPU logic.

More over, when it comes to the implementation of the instructions, I will refer to the following ressource json file which list all the opcodes, mnemonic, instruction length, cycles, flags, access mode and operands for the 512 instructions: [opcodes.json](https://gbdev.io/gb-opcodes/Opcodes.json).

## Introduction

The gameboy has 512 instructions regrouped in two main categories: the non-prefixed instructions and the prefixed instructions. The non-prefixed instructions are the most common and are used to perform basic operations like loading data from one memory location to another, performing arithmetic operations, or branching. The prefixed instructions are used to perform more complex operations like bit manipulation, stack operations, or jumps.

To distinguish between these two categories, a prefixed instruction will begin with a `0xBC prefix byte` encountered by the CPU during the opcode fetch phase when it reads the next opcode after finishing the execution of the previous instruction.

## Approach

A naïve and straightforward approach would be to code the instructions one by one and define for each a handler function that would be called when the corresponding opcode is encountered:

```go
func LD_A_A() {// do stuff}
func LD_A_B() {// do stuff}
func LD_A_C() {// do stuff}
func LD_A_D() {// do stuff}
func LD_A_E() {// do stuff}
func LD_A_H() {// do stuff}
func LD_A_L() {// do stuff}
...
func LB_B_A() {// do stuff}
func LB_B_B() {// do stuff}
func LB_B_C() {// do stuff}
func LB_B_D() {// do stuff}
func LB_B_E() {// do stuff}
func LB_B_H() {// do stuff}
func LB_B_L() {// do stuff}
...
up to 90 functions
```

While this approach will certainly work, it is cumbersome and not very efficient. Instead, I propose to first regroup the instructions by their mnemonic and their operands and try to recognize patterns in order to regroup them in a more efficient way. This will lead to a more maintainable and scalable code that will make it possible to reuse the same utility functions to fetch operands and also take advantage of some Go features like reflection, interfaces and function pointers. These pattern will emerge as we go through the instructions set.

## Execution cycle

Let's now discuss the general structure of our code. The CPU architecture is a CISC and follows the fetch-decode-execute cycle. The CPU will:

- Fetch the next opcode from the bus @PC
- Decode the instruction based on the opcodes.json file by routing the CPU execution to the corresponding handler function. The corresponding handler function will contain the appropriate logic to decode the operands and
- Execute the instruction
  - Fetch the operands
  - Update the program counter
  - Update the flags
  - Update the clock

### Fetching the next opcode

When begining a new cycle, the CPU will fetch the next opcode from the bus located at the address pointed by the program counter:

```go
// Fetch the opcode from bus at address PC and store it in the instruction register
func (c *CPU) fetchOpcode() {
	// Fetch the opcode from memory at the address in the program counter
	opcode := c.bus.Read(c.PC)

	// Store the opcode in the instruction register
	c.IR = opcode
}
```

### Decoding the instruction

Based on the opcode fetched, we can now decode the instruction. This step involves finding the corresponding instruction in the `opcodes.json` file and figuring out which operands from memory will be needed by the execution step.

We will distinguish between instructions that require 0, 1 or 2 operands. To get a better understanding of what this will imply, let's first extract the opcodes and instructions for each of these categories.

But before going any further, please note that this categorisation is based on the data present in the `opcodes.json` file as found in the [gb-opcodes](https://gbdev.io/gb-opcodes/Opcodes.json). Once again, I chose to use this file to have a comprehensive, community accepted and up-to-date list of the instructions set. After my initial implementation, if I find a more efficient way to categorize the instructions, I will propose an update to this document or propose a new one.

#### Instructions with 0 operands

There are a total of 15 instructions that do not require any operands:

| #   | opcode | mnemonic | description                    | length (bytes) | cycles | flags            |
| --- | ------ | -------- | ------------------------------ | -------------- | ------ | ---------------- |
| 1   | 0x00   | NOP      | No operation                   | 1              | 4      | Z:- N:- H:- C:-  |
| 2   | 0x07   | RLCA     | Rotate A left                  | 1              | 4      | Z:0 N:0 H:0 C:C  |
| 3   | 0x0F   | RRCA     | Rotate A right                 | 1              | 4      | Z:0 N:0 H:0 C:C  |
| 4   | 0x17   | RLA      | Rotate A left through carry    | 1              | 4      | Z:0 N:0 H:0 C:C  |
| 5   | 0x1F   | RRA      | Rotate A right through carry   | 1              | 4      | Z:0 N:0 H:0 C:C  |
| 6   | 0x27   | DAA      | Decimal adjust A               | 1              | 4      | Z:Z N:- H:0 C:C  |
| 7   | 0x2F   | CPL      | Complement A                   | 1              | 4      | Z:- N:1 H:1 C:-  |
| 8   | 0x37   | SCF      | Set carry flag                 | 1              | 4      | Z:- N:0 H:0 C:1  |
| 9   | 0x3F   | CCF      | Complement carry flag          | 1              | 4      | Z:- N:0 H:0 C:^C |
| 10  | 0x76   | HALT     | Halt until an interrupt occurs | 1              | 4      | Z:- N:- H:- C:-  |
| 11  | 0xC9   | RET      | Return                         | 1              | 16     | Z:- N:- H:- C:-  |
| 12  | 0xCB   | PREFIX   | Prefix byte                    | 1              | 4      | Z:- N:- H:- C:-  |
| 13  | 0xD9   | RETI     | Return from interrupt          | 1              | 16     | Z:- N:- H:- C:-  |
| 14  | 0xF3   | DI       | Disable interrupts             | 1              | 4      | Z:- N:- H:- C:-  |
| 15  | 0xFB   | EI       | Enable interrupts              | 1              | 4      | Z:- N:- H:- C:-  |

Aside from these instructions, there are also 11 illegal opcodes that are not used by the gameboy and which have an empty array of operands in the `opcodes.json` file. Calling these opcodes will result in an illegal opcode exception, and in our implementation will simply trigger a Go panic:

| #   | opcode | mnemonic   |
| --- | ------ | ---------- |
| 1   | 0xDB   | ILLEGAL_DB |
| 2   | 0xDD   | ILLEGAL_DD |
| 13  | 0xD3   | ILLEGAL_D3 |
| 3   | 0xE3   | ILLEGAL_E3 |
| 4   | 0xE4   | ILLEGAL_E4 |
| 5   | 0xEB   | ILLEGAL_EB |
| 6   | 0xEC   | ILLEGAL_EC |
| 7   | 0xED   | ILLEGAL_ED |
| 8   | 0xF4   | ILLEGAL_F4 |
| 9   | 0xFC   | ILLEGAL_FC |
| 10  | 0xFD   | ILLEGAL_FD |

#### Instructions with 1 operand

Due to the large number of instructions, I will used a condensed table format to present the instructions with 1 operand without explaning what each instruction does. This table on the following one for 2 operands were generated using a script using `opcodes.json` file as input. The script is available in the `scripts` folder.

| #   | opcode | mnemonic | operand | description                                         | length (bytes) | cycles | flags           |
| --- | ------ | -------- | ------- | --------------------------------------------------- | -------------- | ------ | --------------- |
| 1   | 0x03   | INC      | BC      | Increment BC                                        | 1              | 8      | Z:- N:- H:- C:- |
| 2   | 0x04   | INC      | B       | Increment B                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 3   | 0x05   | DEC      | B       | Decrement B                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 4   | 0x0B   | DEC      | BC      | Decrement BC                                        | 1              | 8      | Z:- N:- H:- C:- |
| 5   | 0x0C   | INC      | C       | Increment C                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 6   | 0x0D   | DEC      | C       | Decrement C                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 7   | 0x10   | STOP     | n8      | Stop (n8 does not matter but generally set to 0x00) | 2              | 4      | Z:- N:- H:- C:- |
| 8   | 0x13   | INC      | DE      | Increment DE                                        | 1              | 8      | Z:- N:- H:- C:- |
| 9   | 0x14   | INC      | D       | Increment D                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 10  | 0x15   | DEC      | D       | Decrement D                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 11  | 0x18   | JR       | e8      | Jump relative from PC to signed e8                  | 2              | 12     | Z:- N:- H:- C:- |
| 12  | 0x1B   | DEC      | DE      | Decrement DE                                        | 1              | 8      | Z:- N:- H:- C:- |
| 13  | 0x1C   | INC      | E       | Increment E                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 14  | 0x1D   | DEC      | E       | Decrement E                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 15  | 0x23   | INC      | HL      | Increment HL                                        | 1              | 8      | Z:- N:- H:- C:- |
| 16  | 0x24   | INC      | H       | Increment H                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 17  | 0x25   | DEC      | H       | Decrement H                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 18  | 0x2B   | DEC      | HL      | Decrement HL                                        | 1              | 8      | Z:- N:- H:- C:- |
| 19  | 0x2C   | INC      | L       | Increment L                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 20  | 0x2D   | DEC      | L       | Decrement L                                         | 1              | 4      | Z:Z N:1 H:H C:- |
| 21  | 0x33   | INC      | SP      | Increment SP                                        | 1              | 8      | Z:- N:- H:- C:- |
| 22  | 0x34   | INC      | [HL]    | Increment the byte located @HL                      | 1              | 12     | Z:Z N:0 H:H C:- |
| 23  | 0x35   | DEC      | [HL]    | Decrement the byte located @HL                      | 1              | 12     | Z:Z N:1 H:H C:- |
| 24  | 0x3B   | DEC      | SP      | Decrement SP                                        | 1              | 8      | Z:- N:- H:- C:- |
| 25  | 0x3C   | INC      | A       | Increment A                                         | 1              | 4      | Z:Z N:0 H:H C:- |
| 26  | 0x3D   | DEC      | A       | Decrement A                                         | 1              | 4      | Z:Z N:1 H:H C:- |

#### Instructions with 2 operands

| #   | opcode | mnemonic | operand1 | operand2 | description                                           | length (bytes) | cycles | flags           |
| --- | ------ | -------- | -------- | -------- | ----------------------------------------------------- | -------------- | ------ | --------------- |
| 1   | 0x01   | LD       | BC       | n16      | Load 16-bits immediate into BC                        | 3              | 12     | Z:- N:- H:- C:- |
| 2   | 0x02   | LD       | [BC]     | A        | Load A into address BC                                | 1              | 8      | Z:- N:- H:- C:- |
| 3   | 0x06   | LD       | B        | n8       | Load 8-bits immediate into B                          | 2              | 8      | Z:- N:- H:- C:- |
| 4   | 0x08   | LD       | [a16]    | SP       | Load SP into address a16                              | 3              | 20     | Z:- N:- H:- C:- |
| 5   | 0x09   | ADD      | HL       | BC       | Add BC to HL                                          | 1              | 8      | Z:- N:0 H:H C:C |
| 6   | 0x0A   | LD       | A        | [BC]     | Load value at address BC into A                       | 1              | 8      | Z:- N:- H:- C:- |
| 7   | 0x0E   | LD       | C        | n8       | Load 8-bits immediate into C                          | 2              | 8      | Z:- N:- H:- C:- |
| 8   | 0x11   | LD       | DE       | n16      | Load 16-bits immediate into DE                        | 3              | 12     | Z:- N:- H:- C:- |
| 9   | 0x12   | LD       | [DE]     | A        | Load A into address DE                                | 1              | 8      | Z:- N:- H:- C:- |
| 10  | 0x16   | LD       | D        | n8       | Load 8-bits immediate into D                          | 2              | 8      | Z:- N:- H:- C:- |
| 11  | 0x19   | ADD      | HL       | DE       | Add DE to HL                                          | 1              | 8      | Z:- N:0 H:H C:C |
| 12  | 0x1A   | LD       | A        | [DE]     | Load value at address DE into A                       | 1              | 8      | Z:- N:- H:- C:- |
| 13  | 0x1E   | LD       | E        | n8       | Load 8-bits immediate into E                          | 2              | 8      | Z:- N:- H:- C:- |
| 14  | 0x20   | JR       | NZ       | e8       | Jump relative if Z is not set                         | 2              | 12/8   | Z:- N:- H:- C:- |
| 15  | 0x21   | LD       | HL       | n16      | Load 16-bits immediate into HL                        | 3              | 12     | Z:- N:- H:- C:- |
| 16  | 0x22   | LD       | [HL+]    | A        | Load A into address HL and then increment HL          | 1              | 8      | Z:- N:- H:- C:- |
| 17  | 0x26   | LD       | H        | n8       | Load 8-bits immediate into H                          | 2              | 8      | Z:- N:- H:- C:- |
| 18  | 0x28   | JR       | Z        | e8       | Jump relative if Z is set                             | 2              | 12/8   | Z:- N:- H:- C:- |
| 19  | 0x29   | ADD      | HL       | HL       | Add HL to HL                                          | 1              | 8      | Z:- N:0 H:H C:C |
| 20  | 0x2A   | LD       | A        | [HL+]    | Load value at address HL into A and then increment HL | 1              | 8      | Z:- N:- H:- C:- |
| 21  | 0x2E   | LD       | L        | n8       | Load 8-bits immediate into L                          | 2              | 8      | Z:- N:- H:- C:- |
| 22  | 0x30   | JR       | NC       | e8       | Jump relative if C is not set                         | 2              | 12/8   | Z:- N:- H:- C:- |
| 23  | 0x31   | LD       | SP       | n16      | Load 16-bits immediate into SP                        | 3              | 12     | Z:- N:- H:- C:- |
| 24  | 0x32   | LD       | [HL-]    | A        | Load A into address HL and then decrement HL          | 1              | 8      | Z:- N:- H:- C:- |
| 25  | 0x36   | LD       | [HL]     | n8       | Load 8-bits immediate into address HL                 | 2              | 12     | Z:- N:- H:- C:- |
| 26  | 0x38   | JR       | C        | e8       | Jump relative if C is set                             | 2              | 12/8   | Z:- N:- H:- C:- |
| 27  | 0x39   | ADD      | HL       | SP       | Add SP to HL                                          | 1              | 8      | Z:- N:0 H:H C:C |
| 28  | 0x3A   | LD       | A        | [HL-]    | Load value at address HL into A and then decrement HL | 1              | 8      | Z:- N:- H:- C:- |
| 29  | 0x3E   | LD       | A        | n8       | Load 8-bits immediate into A                          | 2              | 8      | Z:- N:- H:- C:- |
| 30  | 0x40   | LD       | B        | B        | Load B into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 31  | 0x41   | LD       | B        | C        | Load C into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 32  | 0x42   | LD       | B        | D        | Load D into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 33  | 0x43   | LD       | B        | E        | Load E into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 34  | 0x44   | LD       | B        | H        | Load H into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 35  | 0x45   | LD       | B        | L        | Load L into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 36  | 0x46   | LD       | B        | [HL]     | Load value at address HL into B                       | 1              | 8      | Z:- N:- H:- C:- |
| 37  | 0x47   | LD       | B        | A        | Load A into B                                         | 1              | 4      | Z:- N:- H:- C:- |
| 38  | 0x48   | LD       | C        | B        | Load B into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 39  | 0x49   | LD       | C        | C        | Load C into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 40  | 0x4A   | LD       | C        | D        | Load D into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 41  | 0x4B   | LD       | C        | E        | Load E into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 42  | 0x4C   | LD       | C        | H        | Load H into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 43  | 0x4D   | LD       | C        | L        | Load L into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 44  | 0x4E   | LD       | C        | [HL]     | Load value at address HL into C                       | 1              | 8      | Z:- N:- H:- C:- |
| 45  | 0x4F   | LD       | C        | A        | Load A into C                                         | 1              | 4      | Z:- N:- H:- C:- |
| 46  | 0x50   | LD       | D        | B        | Load B into D                                         | 1              | 4      | Z:- N:- H:- C:- |
| 47  | 0x51   | LD       | D        | C        | Load C into D                                         | 1              | 4      | Z:- N:- H:- C:- |
| 48  | 0x52   | LD       | D        | D        | Load D into D                                         | 1              | 4      | Z:- N:- H:- C:- |
| 49  | 0x53   | LD       | D        | E        | Load E into D                                         | 1              | 4      | Z:- N:- H:- C:- |
| 50  | 0x54   | LD       | D        | H        | Load H into D                                         | 1              | 4      | Z:- N:- H:- C:- |

<table>
  <tr>
    <th>Initial A</th>
    <th>Initial C</th>
    <th>Result A</th>
    <th>Result C</th>
  </tr>
  <tr>
    <td style="color: blue;">01011010</td>
    <td style="color: green;">0</td>
    <td style="color: red;">00101101</td>
    <td style="color: purple;">0</td>
  </tr>
  <tr>
    <td style="color: blue;">01011010</td>
    <td style="color: green;">1</td>
    <td style="color: red;">10101101</td>
    <td style="color: purple;">0</td>
  </tr>
</table>

### Routing the execution

Once we have the opcode, we can route it to the corresponding handler function using the `opcodes.json` file that list the corresponding mnemonic. So basically, after fetching the opcode and saving it to the `instruction register IR`, we will call the `LD` handler function that will take care of the rest:

```go

// Route the execution to the corresponding instruction handler
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
  ...
  }
```

### Executing the instruction

## NOP

This instruction is designed to wait for 1 cycle (4 clock ticks). Making the emulator time accurate is a quite complex task. Even if by design, taking care of this instruction will for sure make the CPU 'waste some time', please keep in mind that we are not pursuing this noble goal. In this first implementation, we will just increment the program counter by 1 and move on to the next instruction.

```go
func (c *CPU) NOP() {
  c.PC++
}
```

## LD

The LD instruction is used to load data from one memory location to another. With a total of 90 opcodes, `LD` is the most common instruction. The different underlying opcodes can be split into the following categories:

| Mnemonic             | operand 1             | operand 2                  | description                                                                                              | count |
| -------------------- | --------------------- | -------------------------- | -------------------------------------------------------------------------------------------------------- | ----- |
| LD r8, r8            | A/B/C/D/E/H/L         | A/B/C/D/E/H/L              | Load an 8-bits register into another 8-bits register                                                     | 49    |
| LD r8, n8            | A/B/C/D/E/H/L         | n8                         | Load a byte located @PC+1 into an 8-bits register                                                        | 7     |
| LD r8, [r8] + 0xFF00 | A                     | [C]                        | Load the byte located @address pointed by C + 0xFF00 into the accumulator                                | 1     |
| LD r8, [a16]         | A                     | [a16]                      | Load the 2 bytes located @address pointed by (PC+2, PC+1) into the accumulator                           | 1     |
| LD [r8], r8          | [C]                   | A                          | Load accumulator value into the address pointed by register C                                            | 1     |
| LD r8, [r16]         | A                     | [BC]/[DE]/[HL]/[HL+]/[HL-] | Load the byte located @address pointed by a 16-bits register into the accumulator                        | 5     |
|                      | B/C/D/E/L/H           | [HL]                       | Load the byte located @address pointed by HL register into an 8-bits register                            | 6     |
| LD r16, n16          | BC/DE/HL              | n16                        | Load the value of the 2 next bytes (PC+2, PC+1) into a 16-bits register                                  | 3     |
| LD [r16], n8         | [HL]                  | n8                         | Load the next byte @PC+1 into the address pointed by the HL register                                     | 1     |
| LD [a16], r8         | [a16]                 | A                          | Load the accumulator value into memory location pointed by the next 2 bytes @(PC+2,PC+1)                 | 1     |
| LD [r16], r8         | [BC]/[DE]/[HL+]/[HL-] | A                          | Load the accumulator value into the memory location pointed by a 16-bits register                        | 4     |
|                      | [HL]                  | A/B/C/D/E/H/L              | Load the value of an 8-bits register into the memory location pointed by the HL register                 | 7     |
| LD SP, n16           | SP                    | n16                        | Push to the stack pointer the 2 next bytes @(PC+2, PC+1)                                                 | 1     |
| LD SP, r16           | SP                    | HL                         | Push to the stack pointer the value of register HL                                                       | 1     |
| LD [a16], SP         | [a16]                 | SP                         | Pop the stack pointer value and load it to the memory location pointed by the 2 next bytes @(PC+2, PC+1) | 1     |
| LD r16, r16 + e8     | HL                    | SP+e8                      | Pop the stack pointer value, add a signed 8-bits integer and load it to the HL register                  | 1     |

### Implementation

The `LD`instruction basically copies a value from one memory location to another. The goal here is to use a common handler function that will take care of all the 90 opcodes. In order to achieve this goal, we need a standard way to refer to the destination address and the source value using Go language features.
