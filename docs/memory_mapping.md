| Address Range  | Name                      | Description                                           | Chip Where it Resides              |
|----------------|---------------------------|-------------------------------------------------------|------------------------------------|
| **0000–00FF**  | **Boot ROM**               | Game Boy startup code (disabled after boot)            | **CPU SoC (internal)**            |
| **0000–3FFF**  | ROM Bank 0                 | Fixed portion of the game ROM                          | Cartridge (ROM)                   |
| **4000–7FFF**  | Switchable ROM Bank        | Switchable portion of the game ROM                     | Cartridge (ROM)                   |
| **8000–9FFF**  | VRAM                       | Video RAM (used by GPU for graphics)                   | **PPU SoC (internal)**        |
| **A000–BFFF**  | External RAM               | Optional RAM in some cartridges for save data          | Cartridge (External RAM, if present) |
| **C000–CFFF**  | WRAM Bank 0                | Internal Work RAM                                      | **WRAM (CPU SoC)**                |
| **D000–DFFF**  | WRAM Bank 1                | Switchable Work RAM bank                               | **WRAM (CPU SoC)**                |
| **E000–FDFF**  | Echo RAM                   | Mirrors of C000–DDFF (Not recommended for use)         | **WRAM (CPU SoC)**                |
| **FE00–FE9F**  | OAM                        | Object Attribute Memory (stores sprite information)    | **PPU SoC (internal)**        |
| **FEA0–FEFF**  | Unusable Memory            | Prohibited area (no physical memory attached)          | N/A                               |
| **FF00–FF7F**  | I/O Registers              | Memory-mapped registers for input/output (joypad, etc.)| **CPU SoC (internal)**            |
| **FF80–FFFE**  | High RAM (HRAM)            | High-speed internal Work RAM                           | **CPU SoC (internal)**            |
| **FFFF**       | Interrupt Enable Register  | Controls which interrupts are enabled                  | **CPU SoC (internal)**            |