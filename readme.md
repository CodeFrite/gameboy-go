# Gameboy emulator in Go

This is a Gameboy emulator written in Go. It is a work in progress and is not yet complete. It relies on the documentation found at [GBDEV/Pan](https://gbdev.io/pandocs/About.html)

## Memory access

The gameboy has 2 distincts memory access modes: 8-bit and 16-bit. Therefore, in order to represent memory locations I'll be using GO `byte` and `[2]byte` types. This will allow us to access registers in one of these 2 modes. For example, to access the general purpose register `BC` we can use the following code:

```go
type CPU struct {
  BC [2]byte
}
```

This will allow us to access the `B` and `C` registers as follows:

```go
cpu := CPU{}
cpu.BC[0] = 0x12 // B register
cpu.BC[1] = 0x34 // C register
```
