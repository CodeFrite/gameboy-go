package gameboy

import "fmt"

// the interrupt struct is used to pass information about the interrupt to the subscribers
type Interrupt struct {
	_type   uint8
	payload interface{}
}

var interrupt_types = map[string]uint8{
	"VBLANK": 0,
	"LCD":    1,
	"TIMER":  2,
	"SERIAL": 3,
	"JOYPAD": 4,
}

// Interrupts are used to signal the CPU that an event has occurred and that it should handle it with the appropriate interrupt handler

// * 4.1. Vector 0040h – Vertical Blanking Interrupt
// - This interrupt is triggered when the LCD controller enters V-Blank at scanline 144
// - This doesn't happen if the LCD is off
type vblank_interrupt_payload struct{}

// * 4.2. Vector 0048h – LCD STAT Interrupt
// - This interrupt can be configured to be triggered when some LCD events happen (like starting to draw a scanline specified in LYC register)
type lcd_stat_interrupt_payload struct{}

// * 4.3. Vector 0050h – Timer Interrupt
// - This interrupt is requested when TIMA overflows.
// - There is a delay of one CPU cycle between the overflow and the IF flag being set
type timer_interrupt_payload struct{}

// * 4.4. Vector 0058h – Serial Interrupt
// - This interrupt is requested when a serial transfer of 1 byte is complete
type serial_interrupt_payload struct{}

// * 4.5. Vector 0060h – Joypad Interrupt
// - This interrupt is triggered when there is a transition from '1' to '0' in one of the P1 input lines
type joypad_interrupt_payload struct{}

// * 4.6. IME – Interrupt Master Enable Flag
// - This flag is not mapped to memory and can't be read by any means.
// - The meaning of the flag is not to enable or disable interrupts.
// - In fact, what it does is enable the jump to the interrupt vectors.
// 0 = Disable jump to interrupt vectors.
// 1 = Enable jump to interrupt vectors.
// DONE> ime flag implemented in cpu.go
// ! IME can only be set to '1' by the instructions EI and RETI, and can only be set to '0' by DI (and the CPU when jumping to an interrupt vector).

func (cpu *CPU) TriggerInterrupt(interrupt Interrupt) {
	// trigger the correct interrupt handler
	switch interrupt._type {
	case interrupt_types["VBLANK"]:
		cpu.onVBlankInterrupt(interrupt)
	case interrupt_types["LCD"]:
		cpu.onLCDInterrupt(interrupt)
	/*case interrupt_types["TIMER"]:
		cpu.onTimerInterrupt(interrupt)
	case interrupt_types["SERIAL"]:
		cpu.onSerialInterrupt(interrupt)
	case interrupt_types["JOYPAD"]:
		cpu.onJoypadInterrupt(interrupt)*/
	default:
		panic("CPU> unknown interrupt type")
	}
}

func (cpu *CPU) onVBlankInterrupt(interrupt Interrupt) {
	// set the interrupt flag
	cpu.setIEFlag(0x01)
	// trigger the interrupt handler
	//cpu.onInterruptHandler(interrupt)
}

// Synchronizable interface implementation
func (cpu *CPU) onTick() {
	// return if CPU is locked, otherwise lock CPU and run
	if cpu.state == CPU_EXECUTION_STATE_LOCKED {
		fmt.Println("CPU is locked")
		return
	} else {
		cpu.state = CPU_EXECUTION_STATE_LOCKED
	}

	// check if the CPU is halted
	if cpu.halted {
		// check if the interrupt master enable flag is set
		if cpu.ime {
			// wake up the CPU
			cpu.halted = false
		}
	} else {
		cpu.Step()
	}

	// unlock the CPU
	cpu.state = CPU_EXECUTION_STATE_FREE
}

func (cpu *CPU) onLCDInterrupt(interrupt Interrupt) {
	// before handling the interrupt, we need to disable the interrupt master enable flag
	cpu.ime = false
	// trigger the interrupt handler
	//cpu.onInterruptHandler(interrupt)
}
