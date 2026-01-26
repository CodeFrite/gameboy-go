# PPU

# Layers

## Background Layer

- The background layer is a 256x256 pixel map that is displayed behind the window layer.
- The background layer is made up of 32x32 tiles, each 8x8 pixels in size.
- Each tile is made up of 2-bit pixels, allowing for 4 colors per tile.
- The background layer can be scrolled horizontally and vertically.

## Window Layer

- The window layer has no transparency
- Only the top left corner and the bottom right corner of the window layer can be parameterized

## Objects/Sprite Layer

- Allows objects to move separately from the background
- Made of 1 or 2 stacked tiles that can be 8x8 or 8x16 pixels in size
- Multiple objects can be grouped together to create a meta-sprite

# Registers

## DIV@FF04 - Divider Register (R/W)

- This register is incremented at rate of 16384Hz (~16779Hz on SGB). Note that DIV is always counting, regardless of the TAC register.
- Writing any value to this register resets it to 0x00
- This register is reset when executing the stop instruction, and only begins ticking again once stop mode ends.

## TIMA@FF05 - Timer Counter (R/W)

- This timer is incremented at the clock frequency specified by the TAC register ($FF07).
- When the value overflows (exceeds $FF) it is reset to the value specified in TMA (FF06) and an interrupt is requested.

## TMA@FF06 - Timer Modulo (R/W)

- When TIMA overflows, it is reset to the value in this register and an interrupt is requested.
- Example of use: if TMA is set to $FF, an interrupt is requested at the clock frequency selected in TAC (because every increment is an overflow)
- However, if TMA is set to $FE, an interrupt is only requested every two increments, which effectively divides the selected clock by two. Setting TMA to $FD would divide the clock by three, and so on.
- If a TMA write is executed on the same M-cycle as the content of TMA is transferred to TIMA due to a timer overflow,
- the old value is transferred to TIMA.

## TAC@FF07 - Timer Control (R/W)

- This register controls the timer, determining the clock source and frequency.

| | 7 | 6| 5| 4| 3| 2| 1| 0|
|TAC| | | | Œ Enable Clock select
// ? Enable: Controls whether TIMA is incremented.
// ? Clock select: Controls the frequency at which TIMA is incremented, as follows:
// ! Clock select Increment every DMG, SGB2, CGB in single-speed mode SGB1 Frequency (Hz) CGB in double-speed mode
// ! 00 256 M-cycles 4096 ~4194 8192
// ! 01 4 M-cycles 262144 ~268400 524288
// ! 10 16 M-cycles 65536 ~67110 131072
// ! 11 64 M-cycles 16384 ~16780 32768
