# Registers

Registers are quick access memory locations that are used to control some parameters of an electronic device. We find them in various locations inside the Gameboy. Some are accessible through the bus, others are internal to the component, some are 8-bit, others are 16-bit, ... . And what about flags like Z/C/N/H? What about their relations to the F register?

As we can see, in order to find a good representation of the registers, we need to take a step back, list all the different use cases and find a way to code them efficiently.

# Use Cases

In this section, we will list all the different use cases of registers in the Gameboy and find an appropriate set of types to represent them.

## Case 1: Work Variables

The CPU uses a series of registers to store temporary values. These are the A, B, C, D, E, H, L registers. They are 8-bit registers and can be used in pairs to form 16-bit registers.

These are actual physical electronic "registers" characterized by a fast-access time and located very near the ALU on the CPU SoC. They are not mapped to the bus. Instead, their addresses are directly 'hardcoded' by the cpu architecture. When an instruction like LD A, n8 is encountered, the right signal logic will activate the physical register write signal and store the n8 value.

Therefore, unless we want to represent more accurately the gameboy physical architecture by defining a Memory component and map it to a CPU internal bus, we can simply represent it as a Go variable.

These can be represented as in-memory variables.

# Summary

## CPU Registers

| Register  | Size   | Address | Description                        | Type                               |
| --------- | ------ | ------- | ---------------------------------- | ---------------------------------- |
| A         | 8-bit  | -       | Accumulator                        | Register8(uint8)                   |
| B         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| C         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| D         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| E         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| F         | 8-bit  | -       | Flags                              | Register8(uint8)                   |
| F.Z       | 1-bit  | 7       | Zero Flag                          | Flag(Register8/16, position uint8) |
| F.N       | 1-bit  | 6       | Substraction Flag                  | Flag(Register8/16, position uint8) |
| F.H       | 1-bit  | 5       | Half-Carry Flag                    | Flag(Register8/16, position uint8) |
| F.C       | 1-bit  | 4       | Carry Flag                         | Flag(Register8/16, position uint8) |
| H         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| L         | 8-bit  | -       | General Purpose                    | Register8(uint8)                   |
| AF        | 16-bit | -       | Accumulator and Flags              | Register16(Register8, Register8)   |
| BC        | 16-bit | -       | General Purpose                    | Register16(Register8, Register8)   |
| DE        | 16-bit | -       | General Purpose                    | Register16(Register8, Register8)   |
| HL        | 16-bit | -       | General Purpose                    | Register16(Register8, Register8)   |
| SP        | 16-bit | -       | Stack Pointer                      | Register16(Register8, Register8)   |
| PC        | 16-bit | -       | Program Counter                    | Register16(Register8, Register8)   |
| IME       | 1-bit  | -       | Interrupt Master Enable            | Flag(Register8/16, position uint8) |
| HALTED    | -      | -       | Halt Flag                          | boolean                            |
| STOPPED - | -      | -       | Flag(Register8/16, position uint8) | boolean                            |
| EI        | 1-bit  | -       | Enable Interrupts                  | Flag(Register8/16, position uint8) |
| DI        | 1-bit  | -       | Disable Interrupts                 | Flag(Register8/16, position uint8) |

## PPU Registers

| Location | Register | Size  | Address | Description            | Type      |
| -------- | -------- | ----- | ------- | ---------------------- | --------- |
| PPU      | LCDC     | 8-bit | 0xFF40  | LCD Control            | Register8 |
| PPU      | STAT     | 8-bit | 0xFF41  | LCD Status             | Register8 |
| PPU      | SCY      | 8-bit | 0xFF42  | Scroll Y               | Register8 |
| PPU      | SCX      | 8-bit | 0xFF43  | Scroll X               | Register8 |
| PPU      | LY       | 8-bit | 0xFF44  | LCDC Y-Coordinate      | Register8 |
| PPU      | LYC      | 8-bit | 0xFF45  | LY Compare             | Register8 |
| PPU      | DMA      | 8-bit | 0xFF46  | DMA Transfer and Start | Register8 |
| PPU      | BGP      | 8-bit | 0xFF47  | BG Palette Data        | Register8 |
| PPU      | OBP0     | 8-bit | 0xFF48  | Object Palette 0 Data  | Register8 |
| PPU      | OBP1     | 8-bit | 0xFF49  | Object Palette 1 Data  | Register8 |
| PPU      | WY       | 8-bit | 0xFF4A  | Window Y Position      | Register8 |
| PPU      | WX       | 8-bit | 0xFF4B  | Window X Position      | Register8 |
| PPU      | VBK      | 8-bit | 0xFF4F  | VRAM Bank              | Register8 |
| PPU      | HDMA1    | 8-bit | 0xFF51  | HDMA1                  | Register8 |
| PPU      | HDMA2    | 8-bit | 0xFF52  | HDMA2                  | Register8 |
| PPU      | HDMA3    | 8-bit | 0xFF53  | HDMA3                  | Register8 |
| PPU      | HDMA4    | 8-bit | 0xFF54  | HDMA4                  | Register8 |
| PPU      | HDMA5    | 8-bit | 0xFF55  | HDMA5                  | Register8 |
| PPU      | RP       | 8-bit | 0xFF56  | IR Port                | Register8 |
| PPU      | BCPS     | 8-bit | 0xFF68  | BG Color Palette Spec  | Register8 |
| PPU      | BCPD     | 8-bit | 0xFF69  | BG Color Palette Data  | Register8 |
| PPU      | OCPS     | 8-bit | 0xFF6A  | OBJ Color Palette Spec | Register8 |
| PPU      | OCPD     | 8-bit | 0xFF6B  | OBJ Color Palette Data | Register8 |
| PPU      | SVBK     | 8-bit | 0xFF70  | WRAM Bank              | Register8 |

## Timer Registers

| Location | Register | Size  | Address | Description      | Type      |
| -------- | -------- | ----- | ------- | ---------------- | --------- |
| Timer    | DIV      | 8-bit | 0xFF04  | Divider Register | Register8 |
| Timer    | TIMA     | 8-bit | 0xFF05  | Timer Counter    | Register8 |
| Timer    | TMA      | 8-bit | 0xFF06  | Timer Modulo     | Register8 |
| Timer    | TAC      | 8-bit | 0xFF07  | Timer Control    | Register8 |

## APU Registers

| Location | Register | Size  | Address | Description             | Type      |
| -------- | -------- | ----- | ------- | ----------------------- | --------- |
| APU      | NR10     | 8-bit | 0xFF10  | Sound Mode 1 Sweep      | Register8 |
| APU      | NR11     | 8-bit | 0xFF11  | Sound Mode 1 Length     | Register8 |
| APU      | NR12     | 8-bit | 0xFF12  | Sound Mode 1 Envelope   | Register8 |
| APU      | NR13     | 8-bit | 0xFF13  | Sound Mode 1 Frequency  | Register8 |
| APU      | NR14     | 8-bit | 0xFF14  | Sound Mode 1 Control    | Register8 |
| APU      | NR21     | 8-bit | 0xFF16  | Sound Mode 2 Length     | Register8 |
| APU      | NR22     | 8-bit | 0xFF17  | Sound Mode 2 Envelope   | Register8 |
| APU      | NR23     | 8-bit | 0xFF18  | Sound Mode 2 Frequency  | Register8 |
| APU      | NR24     | 8-bit | 0xFF19  | Sound Mode 2 Control    | Register8 |
| APU      | NR30     | 8-bit | 0xFF1A  | Sound Mode 3 Control    | Register8 |
| APU      | NR31     | 8-bit | 0xFF1B  | Sound Mode 3 Length     | Register8 |
| APU      | NR32     | 8-bit | 0xFF1C  | Sound Mode 3 Output     | Register8 |
| APU      | NR33     | 8-bit | 0xFF1D  | Sound Mode 3 Frequency  | Register8 |
| APU      | NR34     | 8-bit | 0xFF1E  | Sound Mode 3 Control    | Register8 |
| APU      | NR41     | 8-bit | 0xFF20  | Sound Mode 4 Length     | Register8 |
| APU      | NR42     | 8-bit | 0xFF21  | Sound Mode 4 Envelope   | Register8 |
| APU      | NR43     | 8-bit | 0xFF22  | Sound Mode 4 Polynomial | Register8 |
| APU      | NR44     | 8-bit | 0xFF23  | Sound Mode 4 Control    | Register8 |
| APU      | NR50     | 8-bit | 0xFF24  | Sound Control           | Register8 |
| APU      | NR51     | 8-bit | 0xFF25  | Sound Output Selection  | Register8 |
| APU      | NR52     | 8-bit | 0xFF26  | Sound On/Off            | Register8 |
| APU      | WAVE     | 8-bit | 0xFF30  | Wave Pattern RAM        | Register8 |
