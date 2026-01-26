# Halt & Stop modes

The Halt and Stop modes are low-power modes that can be used to reduce the power consumption of the device when it is not needed to be fully operational. Let's check the differences between these two modes:

The STOP command halts the GameBoy processor and screen until any button is pressed. The GB and GBP screen goes white with a single dark horizontal line. The GBC screen goes black.

The HALT command stops the system clock reducing the power consumption of both the CPU and ROM. The screen remains on and the sound continues to play. The HALT command is used to wait for an interrupt to occur.

## Comparison table

| Mode | CPU     | Screen | Sound | Trigger               | Wake up condition     |
| ---- | ------- | ------ | ----- | --------------------- | --------------------- |
| STOP | Halted  | Off    | Off   | stop instruction 0x10 | Any button is pressed |
| HALT | Stopped | On     | On    | halt instruction 0x76 | Any interrupt occurs  |

## CPU execution loop during low power modes

The fact that the CPU watches for interrupts during the HALT and STOP modes shows that even if the CPU does not execute any more instruction from the ROM, it is still waiting for an interrupt to occur.

Basically, the CPU is in a loop that (or more specifically, the CPU is stepped at regular intervals by the clock):

- check if the cpu is in HALT mode and if there is an interrupt requested and wakes up
- check if the cpu is in STOP mode and if any button is pressed and wakes up
- then it always execute the next instruction

## Sources

1. Game BoyTM CPU Manual @v1.01 (Pan of Anthrox, GABY, Marat Fayzullin, Pascal Felber, Paul Robson, Martin Korth, kOOPa, Bowser)
