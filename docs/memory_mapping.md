# Memory Mapping

The gameboy cpu bus is 16-bit wide and can thus address 2^16 = 65536 bytes of memory. The memory is divided into several regions, physically or logically, each with a specific purpose. The following table lists the memory regions, their purpose and the location where they reside:

| Address Range | Name                      | Description                                             | Chip Where it Resides |
| ------------- | ------------------------- | ------------------------------------------------------- | --------------------- |
| **0000–00FF** | Boot ROM                  | Game Boy startup code (disabled after boot)             | CPU SoC internal      |
| **0000–3FFF** | ROM Bank 0                | Fixed portion of the game ROM                           | cartridge board       |
| **4000–7FFF** | Switchable ROM Bank       | Switchable portion of the game ROM                      | cartridge board       |
| **8000–9FFF** | VRAM                      | Video RAM (used by GPU for graphics)                    | PPU SoC internal      |
| **A000–BFFF** | External RAM              | Optional RAM in some cartridges for save data           | cartridge board       |
| **C000–CFFF** | WRAM Bank 0               | Internal Work RAM                                       | mainboard             |
| **D000–DFFF** | WRAM Bank 1               | Switchable Work RAM bank                                | mainboard             |
| **E000–FDFF** | Echo RAM                  | Mirrors of C000–DDFF (Not recommended for use)          | mainboard             |
| **FE00–FE9F** | OAM                       | Object Attribute Memory (stores sprite information)     | PPU SoC internal      |
| **FEA0–FEFF** | Unusable Memory           | Prohibited area (no physical memory attached)           | N/A                   |
| **FF00–FF7F** | I/O Registers             | Memory-mapped registers for input/output (joypad, etc.) | CPU SoC internal      |
| **FF80–FFFE** | High RAM (HRAM)           | High-speed internal Work RAM                            | CPU SoC internal      |
| **FFFF**      | Interrupt Enable Register | Controls which interrupts are enabled                   | CPU SoC internal      |

When it comes to implementing the memory mapping, I see two main options:

- initialize every single memory area in the gameboy_init.go file like if the gameboy was responsible to build
- leave the responsibility to the different components to connect the memories to the appropriate bus and memory areas

I'll go with the second option, as it reflects better the real hardware.

## Implementation

Based on the table above, we will have to load the memory areas in the following structs:

- Gameboy: responsible for instantiating the 2 buses/mmu's, the CPU, the PPU, the APU and the Cartridge as well as:
  - **0x8000-0x9FFF** : VRAM
  - **0xC000-0xCFFF** : WRAM Bank 0
  - **0xD000-0xDFFF** : WRAM Bank 1
  - **0xE000-0xFDFF** : Echo RAM
- CPU: responsible for connecting the CPU SoC memory areas. Please note that other registers PC, SP, AF, BC, DE, HL, etc. are part of the CPU struct and are not accessible through the buses. Therefore, they are not listed here below.
  - **0x0000-0x00FF** : Boot ROM
  - **0xFF00-0xFF7F** : I/O Registers
  - **0xFF80-0xFFFE** : High RAM (HRAM)
  - **0xFFFF** : Interrupt Enable Register
- PPU: responsible for connecting the PPU SoC memory areas
  - **0xFE00-0xFE9F** : OAM
- Cartridge: responsible for connecting the Cartridge memory areas
  - **0x0000-0x3FFF** : ROM Bank 0
  - **0x4000-0x7FFF** : Switchable ROM Bank
  - **0xA000-0xBFFF** : External RAM
