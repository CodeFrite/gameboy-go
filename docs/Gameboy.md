```mermaid
classDiagram
    %% Core Gameboy Components
    class Gameboy {
        -ticks: uint64
        -state: GameBoyState
        -timer: Timer*
        -bus: Bus*
        -cpu: CPU*
        -ppu: PPU*
        -apu: APU*
        -bootrom: Memory*
        -cartridge: Cartridge*
        -vram: Memory*
        -wram: Memory*
        -joypad: Joypad*
        -gameboyActionChannel: chan GameboyActionMessage
        -cpuStateChannel: chan CpuState
        -ppuStateChannel: chan PpuState
        -apuStateChannel: chan ApuState
        -memoryStateChannel: chan []MemoryWrite
        +NewGameboy() *Gameboy
        +Tick()
        +GetCpuState() CpuState
        +GetMemoryWrites() []MemoryWrite
        +LoadRom(romName string)
        +Run() / Pause() / Stop()
    }

    %% CPU Component
    class CPU {
        -clock: uint64
        -state: CPU_EXECUTION_STATE
        -pc: uint16
        -sp: uint16
        -a, f: uint8
        -b, c, d, e, h, l: uint8
        -ir: uint8
        -instruction: Instruction
        -prefixed: bool
        -operand: uint16
        -cpuCycles: uint64
        -ime: bool
        -halted: bool
        -stopped: bool
        -bus: Bus*
        -io_registers: Memory*
        -hram: Memory*
        -ie: Memory*
        +NewCPU(bus *Bus) *CPU
        +Tick()
        +fetch() / decode() / execute()
        +getState() CpuState
        +handleInterrupts()
        +executeInstruction(instruction Instruction)
    }

    %% PPU Component
    class PPU {
        -bus: Bus*
        -image: RenderedImage
        -background: [256][64]uint8
        -mode: uint8
        -ticks: uint64
        -dotX: uint16
        -dotY: uint16
        -mode3Length: uint16
        -oam: Memory*
        +NewPPU(bus *Bus) *PPU
        +Tick()
        +drawBackground()
        +drawWindow()
        +isEnabled() bool
        +updateSTATRegister_PPUMode()
        +updateLYRegister()
    }

    %% Bus/Memory Management
    class Bus {
        -memoryMaps: []MemoryMap
        -MemoryWrites: []MemoryWrite
        -writeHandlers: map[uint16]func(uint8) uint8
        +NewBus() *Bus
        +AttachMemory(name string, address uint16, memory Accessible)
        +Read(address uint16) uint8
        +Write(address uint16, value uint8)
        +Read16(address uint16) uint16
        +Dump(address uint16, size uint16) []uint8
        +DisableBootrom()
    }

    class Memory {
        -data: []uint8
        +NewMemory(size uint16) *Memory
        +NewMemoryWithData(size uint16, data []uint8) *Memory
        +NewMemoryWithRandomData(size uint16) *Memory
        +Read(addr uint16) uint8
        +Write(addr uint16, value uint8)
        +Dump(from uint16, to uint16) []uint8
        +Size() uint16
    }

    class MemoryMap {
        +Name: string
        +Address: uint16
        +Memory: Accessible
    }

    %% Interfaces
    class Accessible {
        <<interface>>
        +Read(uint16) uint8
        +Write(uint16, uint8)
        +Dump(uint16, uint16) []uint8
        +Size() uint16
    }

    %% Instructions and State
    class Instruction {
        +Mnemonic: string
        +Bytes: int
        +Cycles: JSONableSlice
        +Operands: []Operand
        +Immediate: bool
        +Flags: Flags
    }

    class Operand {
        +Name: string
        +Bytes: int
        +Immediate: bool
        +Increment: bool
        +Decrement: bool
    }

    class Flags {
        +Z: string
        +N: string
        +H: string
        +C: string
    }

    %% State Objects
    class CpuState {
        +CPU_CYCLES: uint64
        +PC: uint16
        +SP: uint16
        +A, F: uint8
        +BC, DE, HL: uint16
        +INSTRUCTION: Instruction
        +PREFIXED: bool
        +IR: uint8
        +OPERAND_VALUE: uint16
        +IE: uint8
        +IME: bool
        +HALTED: bool
        +STOPPED: bool
    }

    class PpuState {
        +MODE: uint8
        +DOT_X: uint8
        +DOT_Y: uint8
        +IMAGE: RenderedImage
    }

    class ApuState {
        +SOUND: bool
    }

    class MemoryWrite {
        +Name: string
        +Address: uint16
        +Value: uint8
    }

    %% Other Components
    class APU {
        -sound: bool
        +NewAPU() *APU
        +Tick()
        +getState() ApuState
    }

    class Timer {
        -bus: Bus*
        -internalClock: uint16
        +NewTimer(bus *Bus) *Timer
        +Tick()
    }

    class Cartridge {
        -cartridgePath: string
        -cartridgeName: string
        -rom: Memory*
        -header: []uint8
        +NewCartridge(uri string, name string) *Cartridge
        +Read(addr uint16) uint8
        +Write(addr uint16, value uint8)
        +Size() uint16
        +PrintInfo()
    }

    class Joypad {
        -state: uint8
        +NewJoypad() *Joypad
        +Write(value uint8)
    }

    %% Debugger
    class Debugger {
        -gameboy: Gameboy*
        -programFlow: Fifo[uint16]*
        -breakpoints: []uint16
        -cpuStateQueue: Fifo[CpuState]*
        -memoryStateQueue: Fifo[[]MemoryWrite]*
        +NewDebugger() *Debugger
        +LoadRom(romName string)
        +Tick() / Run() / Pause() / Stop()
        +AddBreakPoint(addr uint16)
        +RemoveBreakPoint(addr uint16)
    }

    %% Data Structures
    class Fifo~T~ {
        -capacity: uint64
        -count: uint64
        -head: Node[T]*
        +NewFifo(capacity uint64) *Fifo[T]
        +Push(value *T) uint64
        +Pop() *T
        +GetHead() *Node[T]
    }

    class Node~T~ {
        -value: T*
        -next: Node[T]*
        +NewNode(value *T, next *Node[T]) *Node[T]
        +GetValue() *T
        +GetNext() *Node[T]
        +SetNext(next *Node[T])
    }

    class Iterable~T~ {
        <<interface>>
        +GetHead() *Node[T]
    }

    class UpdatableList~T~ {
        <<interface>>
        +Push(*T) uint64
        +Pop() *T
    }

    class UpdatableIterator~T~ {
        <<interface>>
        +GetHead() *Node[T]
        +Push(*T) uint64
        +Pop() *T
    }

    %% Relationships
    Gameboy *-- Timer
    Gameboy *-- Bus
    Gameboy *-- CPU
    Gameboy *-- PPU
    Gameboy *-- APU
    Gameboy *-- Memory : "bootrom, vram, wram"
    Gameboy *-- Cartridge
    Gameboy *-- Joypad

    CPU *-- Bus
    CPU *-- Memory : "io_registers, hram, ie"
    CPU *-- Instruction
    CPU *-- CpuState

    PPU *-- Bus
    PPU *-- Memory : "oam"
    PPU *-- PpuState

    Bus *-- MemoryMap
    Bus *-- MemoryWrite
    MemoryMap *-- Accessible

    Instruction *-- Operand
    Instruction *-- Flags

    CpuState *-- Instruction

    Debugger *-- Gameboy
    Debugger *-- Fifo

    Fifo *-- Node
    Node --|> Iterable
    Fifo --|> UpdatableIterator
    UpdatableIterator --|> Iterable
    UpdatableIterator --|> UpdatableList

    Memory ..|> Accessible
    Bus ..|> Accessible
    Cartridge ..|> Accessible

    Timer *-- Bus
    APU *-- ApuState

    %% Color coding for better visualization
    classDef coreComponent fill:#e1f5fe
    classDef memoryComponent fill:#f3e5f5
    classDef stateComponent fill:#e8f5e8
    classDef interfaceComponent fill:#fff3e0
    classDef dataStructure fill:#fce4ec

    class Gameboy,CPU,PPU,APU,Timer coreComponent
    class Memory,Bus,MemoryMap,Cartridge memoryComponent
    class CpuState,PpuState,ApuState,MemoryWrite stateComponent
    class Accessible,Iterable,UpdatableList,UpdatableIterator interfaceComponent
    class Fifo,Node dataStructure
```
