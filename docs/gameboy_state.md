# Gameboy State

The gameboy should be able to share its internal state with the user. In real life, it does it through the LCD screen for the image and with the speaker for the sound. But as it was designed, the gameboy-go core doesn't provide any of these features for displaying the image and playing the sound. It was done intentionally to keep allow the user to choose and implement the way to render these signals. Therefore, the gameboy-go core should provide a way to access its internal state.

## State @v0.3.2

Up to the current version 0.3.2, the gameboy public interface provides the `Debugger.Step` and `Debugger.Run`functions to control the execution flow. They both return a `State` object that contains the gameboy internal state at the end of their execution, which in the case of the Run function means when the gameboy is halted, stopped, a breakpoint is reached or a panic occurs.

## @v0.4.0

I want to communicate the gameboy state with the user through channels. Channels for cpu, ppu and apu are hardcoded in the `Debugger` and `Gameboy`struct and initialized in the `NewGameboy`func. Here is the list of channels:

| Channel   | Data Type  | Initializer        | Trigger                                  | Description            |
| --------- | ---------- | ------------------ | ---------------------------------------- | ---------------------- |
| cpu_state | `CpuState` | Gameboy.NewGameboy | `Step`& `Run`& `Tick`                    | The complete CPU state |
| ppu_state | `PpuState` | Gameboy.NewGameboy | only if changes in VRAM or lcd registers | The complete PPU state |
| apu_state | `ApuState` | Gameboy.NewGameboy | `Tick`???                                | The complete APU state |

For version 0.4.0, I will only implement the `cpu_state` channel:

| Channel   | Data Type  | Initializer        | Trigger               | Description            |
| --------- | ---------- | ------------------ | --------------------- | ---------------------- |
| cpu_state | `CpuState` | Gameboy.NewGameboy | `Step`& `Run`& `Tick` | The complete CPU state |

## @v0.4.1

I want to implement the `ppu_state` channel. The `apu_state` channel will be implemented in a future version.

| Channel   | Data Type  | Initializer        | Trigger                                  | Description            |
| --------- | ---------- | ------------------ | ---------------------------------------- | ---------------------- |
| ppu_state | `PpuState` | Gameboy.NewGameboy | only if changes in VRAM or lcd registers | The complete PPU state |

## @v0.4.2

I want to implement the `apu_state` channel.

| Channel   | Data Type  | Initializer        | Trigger   | Description            |
| --------- | ---------- | ------------------ | --------- | ---------------------- |
| apu_state | `ApuState` | Gameboy.NewGameboy | `Tick`??? | The complete APU state |

## @v0.4.3

I want to use a channel to communicate the joypad inputs to the gameboy. This channel should be defined on the server side where the gameboy-go core is running and passed to the gameboy through the `NewGameboy` func. The gameboy should read the joypad inputs from this channel and update the joypad registers accordingly. The CPU will then take over from there.

| Channel | Data Type | Initializer                         | Trigger             | Description       |
| ------- | --------- | ----------------------------------- | ------------------- | ----------------- |
| joypad  | `Joypad`  | Server/any program using gameboy-go | user button presses | The joypad inputs |

Since gameboy-go allows the user to choose a way to interact with the gameboy (keyboard, joypad, mouse, web interface buttons, ...), this is the responsability of the user to define a channel that will provide the joypad inputs in a very basic format to ease the user development:

```go
type Joypad struct {
  Up, Down, Left, Right, A, B, Start, Select bool
}
```
