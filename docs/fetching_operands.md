# Fetching operands easily and grouping related opcodes

This file contains all the 512 opcodes and their respective information (cpu cycles, bytes, addressing mode, etc). for the Gameboy CPU (LR35902).

This file comes handy when you want to implement the CPU of the Gameboy in your emulator. If used wisely, it can help you emulate the CPU behavior without having to write a function for each opcode. Indeed, by pooling the addressing modes and operands, you can regroup the implementation of several opcodes in a single function (105 LD instructions, 16 ADD instructions, etc).

## The issue with the file

The Gameboy has the following registers: AB, BC, DE, HL, SP, PC. The F register is the 'flags' register used to keep track of the result of the last operation. Unfortunately, the `opcodes.json` file does not make the distinction between the registers & flags `C` and `H`. Therefore, we have to find a way to differentiate them.

## The `H`operand

Let's list all the opcodes that use the `H` operand as `operand 1` or `operand 2` and see if it is sometimes used to reference the `H` register and sometimes the `C` flag register.

### `INC` & `DEC` instructions

They both use the `H` register as `operand 1` and do not have any `operand 2`. The increase or decrease the value of the `H` register by 1 and affect the `H`flag register.

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                         |
| --- | ------ | -------- | --------- | --------- | ----------------------------------- |
| 1   | 0x24   | INC H    | H         | -         | Increment value in register r8 by 1 |
| 2   | 0x25   | DEC H    | H         | -         | Decrement value in register r8 by 1 |

These 2 instructions do not bring any ambiguity as they both refer to the `H` register.

### `LD` instructions

There are a total of 16 `LD`instructions using the `H`register either as `operand 1` or `operand 2.

| #   | Opcode | Mnemonic   | operand 1 | operand 2 | description                                    |
| --- | ------ | ---------- | --------- | --------- | ---------------------------------------------- |
| 3   | 0x26   | LD H, n8   | H         | n8        | Load immediate 8-bit data in register H        |
| 4   | 0x67   | LD H, A    | H         | A         | Load 8-bit register value into register H      |
| 5   | 0x60   | LD H, B    | H         | B         | -                                              |
| 6   | 0x61   | LD H, C    | H         | C         | -                                              |
| 7   | 0x62   | LD H, D    | H         | D         | -                                              |
| 8   | 0x63   | LD H, E    | H         | E         | -                                              |
| 9   | 0x64   | LD H, H    | H         | H         | -                                              |
| 10  | 0x65   | LD H, L    | H         | L         | -                                              |
| 11  | 0x66   | LD H, [HL] | H         | [HL]      | Load 8-bit value located @(HL) into H register |

| #   | Opcode | Mnemonic   | operand 1 | operand 2 | description                                           |
| --- | ------ | ---------- | --------- | --------- | ----------------------------------------------------- |
| 12  | 0x7C   | LD A, H    | A         | H         | Load H register value into another 8-bit register     |
| 13  | 0x44   | LD B, H    | B         | H         | -                                                     |
| 14  | 0x4C   | LD C, H    | C         | H         | -                                                     |
| 15  | 0x54   | LD D, H    | D         | H         | -                                                     |
| 16  | 0x5C   | LD E, H    | E         | H         | -                                                     |
| 17  | 0x6C   | LD L, H    | L         | H         | -                                                     |
| 18  | 0x74   | LD [HL], H | [HL]      | H         | Load H register value into the address in HL register |

Here again, there is no ambiguity. All instructions above refer to the `H` register.

### `ADD`, `SUB` instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                                       |
| --- | ------ | -------- | --------- | --------- | ------------------------------------------------- |
| 19  | 0x84   | ADD A, H | A         | H         | Add the value of register H to register A         |
| 20  | 0x94   | SUB A, H | A         | H         | Substract the value of register H from register A |

No ambiguity here.

### `AND`, `XOR`, `OR`, `CP` instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                                                                                                              |
| --- | ------ | -------- | --------- | --------- | ------------------------------------------------------------------------------------------------------------------------ |
| 21  | 0xA4   | AND A, H | A         | H         | Bitwise AND operation between register A and H                                                                           |
| 22  | 0xAC   | XOR A, H | A         | H         | Bitwise XOR operation between register A and H                                                                           |
| 23  | 0xB4   | OR A, H  | A         | H         | Bitwise OR operation between register A and H                                                                            |
| 24  | 0xBC   | CP A, H  | A         | H         | Compare values of A and H register by substracting H from A and setting flags accordingly without saving the result in A |

No ambiguity here.

### `ADC`, `SBC` instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                                                          |
| --- | ------ | -------- | --------- | --------- | -------------------------------------------------------------------- |
| 25  | 0x8C   | ADC A, H | A         | H         | Add the value of register H and the Carry Flag to register A         |
| 26  | 0x9C   | SBC A, H | A         | H         | Substract the value of register H and the Carry Flag from register A |

No ambiguity here.

### `RLC`, `RRC`, `RL`, `RR`, `SLA`, `SRA`, `SWAP`, `SRL` instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                                                     |
| --- | ------ | -------- | --------- | --------- | --------------------------------------------------------------- |
| 27  | 0x04   | RLC H    | H         | -         | Rotate H register left                                          |
| 28  | 0x0C   | RRC H    | H         | -         | Rotate H register right                                         |
| 29  | 0x14   | RL H     | H         | -         | Rotate H register left through the Carry flag                   |
| 30  | 0x1C   | RR H     | H         | -         | Rotate H register right through the Carry flag                  |
| 31  | 0x24   | SLA H    | H         | -         | Shift Arithmetically H register left (bit 0 is replaced by 0)   |
| 32  | 0x2C   | SRA H    | H         | -         | Shift Arithmetically H register right (bit 7 is left unchanged) |
| 33  | 0x3C   | SRL H    | H         | -         | Shift Logically H register left (bit 7 is replaced by 0)        |
| 34  | 0x34   | SWAP H   | H         | -         | Swap register H upper and lower nibbles                         |

No ambiguity here.

### `BIT`instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                                         |
| --- | ------ | -------- | --------- | --------- | --------------------------------------------------- |
| 35  | 0x44   | BIT 0, H | 0         | H         | Test bit x of register H and set Z flag accordingly |
| 36  | 0x4C   | BIT 1, H | 1         | H         | -                                                   |
| 37  | 0x54   | BIT 2, H | 2         | H         | -                                                   |
| 38  | 0x5C   | BIT 3, H | 3         | H         | -                                                   |
| 39  | 0x64   | BIT 4, H | 4         | H         | -                                                   |
| 40  | 0x6C   | BIT 5, H | 5         | H         | -                                                   |
| 41  | 0x74   | BIT 6, H | 6         | H         | -                                                   |
| 42  | 0x7C   | BIT 7, H | 7         | H         | -                                                   |

No ambiguity here.

### `RES`instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                    |
| --- | ------ | -------- | --------- | --------- | ------------------------------ |
| 43  | 0x84   | RES 0, H | 0         | H         | Reset bit x of register H (=0) |
| 44  | 0x8C   | RES 1, H | 1         | H         | -                              |
| 45  | 0x94   | RES 2, H | 2         | H         | -                              |
| 46  | 0x9C   | RES 3, H | 3         | H         | -                              |
| 47  | 0xA4   | RES 4, H | 4         | H         | -                              |
| 48  | 0xAC   | RES 5, H | 5         | H         | -                              |
| 49  | 0xB4   | RES 6, H | 6         | H         | -                              |
| 50  | 0xBC   | RES 7, H | 7         | H         | -                              |

No ambiguity here.

### `SET`instructions

| #   | Opcode | Mnemonic | operand 1 | operand 2 | description                  |
| --- | ------ | -------- | --------- | --------- | ---------------------------- |
| 51  | 0xC4   | SET 0, H | 0         | H         | Set bit x of register H (=1) |
| 52  | 0xCC   | SET 1, H | 1         | H         | -                            |
| 53  | 0xD4   | SET 2, H | 2         | H         | -                            |
| 54  | 0xDC   | SET 3, H | 3         | H         | -                            |
| 55  | 0xE4   | SET 4, H | 4         | H         | -                            |
| 56  | 0xEC   | SET 5, H | 5         | H         | -                            |
| 57  | 0xF4   | SET 6, H | 6         | H         | -                            |
| 58  | 0xFC   | SET 7, H | 7         | H         | -                            |

No ambiguity here.

### Conclusions on the `H` operand

The `H` operand is always used to reference the `H` register and never the `C` flag register. Therefore, we can safely assume that the `H` operand always refers to the `H` register. No changes are needed in the `opcodes.json` file.

## The `C` operand

By analogy, looking at the instructions using the `C`operand, I come to the conclusion that all opcodes using the `C`operand actually refer to the `C`register except for the conditional jumps instructions `CALL`/`JP`/ `JR`/ `RET`, the `CCF` (complement carry flag) and finally the `SCF` (set carry flag) instructions.

For these instructions, special care should be taken when implementing them in the emulator. Indeed, instead of relying on the operand value fetched during the `fetch` state, we will directly access the value of the `C` flag register.

For example, here is the implementation of the `JP` instruction where the opcode `0xDA = JP C, a16` is treated differently from the other `JP`opcodes:

## Conclusions

When fetching the operand value, we will always return the value of register `C` and `H`and NOT the value of the corresponding flag register. If for some instruction, the value of the flag register is needed, we will directly access the flag register from the instruction implementation.

Here is the updated `fetchOperandValue` function:

```go


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
			panic("Non immediate operand not implemented yet")
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
	default:
		err := fmt.Sprintf("cpu.fetchOperandValue> Unknown operand name: %s (0x%02X)", operand.Name, c.ir)
		panic(err)
	}
	return value
}
```
