package gameboy

// * 4.1. Vector 0040h – Vertical Blanking Interrupt
// - This interrupt is triggered when the LCD controller enters V-Blank at scanline 144
// - This doesn't happen if the LCD is off (LCDC.7=0: in this implementation, PPU just returns when ticked if this LCD is off)
// * 4.2. Vector 0048h – LCD STAT Interrupt
// - This interrupt can be configured to be triggered when some LCD events happen (like starting to draw a scanline specified in LYC register)
// * 4.3. Vector 0050h – Timer Interrupt
// - This interrupt is requested when TIMA overflows.
// - There is a delay of one CPU cycle between the overflow and the IF flag being set
// * 4.4. Vector 0058h – Serial Interrupt
// - This interrupt is requested when a serial transfer of 1 byte is complete
// * 4.5. Vector 0060h – Joypad Interrupt
// - This interrupt is triggered when there is a transition from '1' to '0' in one of the P1 input lines
// * 4.6. IME – Interrupt Master Enable Flag
// - This flag is not mapped to memory and can't be read by any means.
// - The meaning of the flag is not to enable or disable interrupts.
// - In fact, what it does is enable the jump to the interrupt vectors.
// 0 = Disable jump to interrupt vectors.
// 1 = Enable jump to interrupt vectors.
// DONE> ime flag implemented in cpu.go

// ! IME can only be set to '1' by the instructions EI and RETI, and can only be set to '0' by DI (and the CPU when jumping to an interrupt vector).

const (
	// IE register: controls whether an interrupt is enabled or not
	IE_REGISTER     uint16 = 0xFFFF
	FFFF_0_VBLANK   uint8  = 0
	FFFF_1_LCD_STAT uint8  = 1
	FFFF_2_TIMER    uint8  = 2
	FFFF_3_SERIAL   uint8  = 3
	FFFF_4_JOYPAD   uint8  = 4

	// IF register: flags that are set when an interrupt is requested
	IF_REGISTER     uint16 = 0xFF0F
	FF0F_0_VBLANK   uint8  = 0
	FF0F_1_LCD_STAT uint8  = 1
	FF0F_2_TIMER    uint8  = 2
	FF0F_3_SERIAL   uint8  = 3
	FF0F_4_JOYPAD   uint8  = 4

	// Interrupts jump vectors
	INTERRUPT_VBLANK_JUMP_VECTOR   uint16 = 0x0040
	INTERRUPT_LCD_STAT_JUMP_VECTOR uint16 = 0x0048
	INTERRUPT_TIMER_JUMP_VECTOR    uint16 = 0x0050
	INTERRUPT_SERIAL_JUMP_VECTOR   uint16 = 0x0058
	INTERRUPT_JOYPAD_JUMP_VECTOR   uint16 = 0x0060
)

// Interrupts are used to signal the CPU that an event has occurred and that it should handle it with the appropriate interrupt handler

func (cpu *CPU) handleInterrupts() {
	// check if the interrupt master enable flag is set
	if !cpu.ime {
		return
	}

	// get the IE and IF registers
	ie_register := cpu.GetIEFlag()
	if_register := cpu.bus.Read(IF_REGISTER)

	// check if an interrupt is requested and enabled by priority
	// if it is the case, disable the IME flag and jump to the interrupt handler

	// V-Blank interrupt
	if (ie_register&(1<<FFFF_0_VBLANK))&(if_register&(1<<FF0F_0_VBLANK)) == 1 {
		cpu.ime = false
		cpu.onVBlankInterrupt()
	} else if (ie_register&(1<<FFFF_1_LCD_STAT))&(if_register&(1<<FF0F_1_LCD_STAT)) == 1 {
		// LCD STAT interrupt
		cpu.ime = false
		cpu.onLCDStatInterrupt()
	} else if (ie_register&(1<<FFFF_2_TIMER))&(if_register&(1<<FF0F_2_TIMER)) == 1 {
		// TIMER interrupt
		cpu.ime = false
		cpu.onTimerInterrupt()
	} else if (ie_register&(1<<FFFF_3_SERIAL))&(if_register&(1<<FF0F_3_SERIAL)) == 1 {
		// SERIAL interrupt
		cpu.ime = false
		cpu.onSerialInterrupt()
	} else if (ie_register&(1<<FFFF_4_JOYPAD))&(if_register&(1<<FF0F_4_JOYPAD)) == 1 {
		// JOYPAD interrupt
		cpu.ime = false
		cpu.onJoypadInterrupt()
	}
}

// * The following interrupt service routine is executed when control is being transferred to an interrupt handler (source: https://gbdev.io/pandocs/Interrupts.html)
// - Two wait states are executed (2 M-cycles pass while nothing happens; presumably the CPU is executing nops during this time).
// ! my not on the previous point: the CPU is not executing nops, it is waiting for the interrupt handler to be executed otherwise it would alter the PC which would
// ! impact the push operation that is about to happen ==> I will simply add 5 cycles to the total number of cpu cycles for correct timing
// - The current value of the PC register is pushed onto the stack, consuming 2 more M-cycles.
// - The PC register is set to the address of the handler (one of: $40, $48, $50, $58, $60). This consumes one last M-cycle.
// ==> The entire process lasts 5 M-cycles.

func (cpu *CPU) onVBlankInterrupt() {
	// clear the interrupt request flag before jumping to the interrupt handler @0x0040
	if_register := cpu.bus.Read(IF_REGISTER) & ((1 << FF0F_0_VBLANK) ^ 0xFF)
	cpu.setIEFlag(uint16(if_register))
	// update the program counter to have the address of the next instruction that cpu would have execute if no interrupt was triggered and push it to the stack
	cpu.updatepc()
	cpu.push(cpu.pc)
	// trigger the interrupt handler
	cpu.pc = INTERRUPT_VBLANK_JUMP_VECTOR
	// wait for 5 M-cycles = 20 T-cycles
	cpu.cpuCycles += 5 * 4
}

func (cpu *CPU) onLCDStatInterrupt() {
	// clear the interrupt request flag before jumping to the interrupt handler @0x0048
	if_register := cpu.bus.Read(IF_REGISTER) & ((1 << FF0F_1_LCD_STAT) ^ 0xFF)
	cpu.setIEFlag(uint16(if_register))
	// update the program counter to have the address of the next instruction that cpu would have execute if no interrupt was triggered and push it to the stack}
	cpu.updatepc()
	cpu.push(cpu.pc)
	// trigger the interrupt handler
	cpu.pc = INTERRUPT_LCD_STAT_JUMP_VECTOR
	// wait for 5 M-cycles = 20 T-cycles
	cpu.cpuCycles += 5 * 4
}

func (cpu *CPU) onTimerInterrupt() {
	// clear the interrupt request flag before jumping to the interrupt handler @0x0050
	if_register := cpu.bus.Read(IF_REGISTER) & ((1 << FF0F_2_TIMER) ^ 0xFF)
	cpu.setIEFlag(uint16(if_register))
	// update the program counter to have the address of the next instruction that cpu would have execute if no interrupt was triggered and push it to the stack}
	cpu.updatepc()
	cpu.push(cpu.pc)
	// trigger the interrupt handler
	cpu.pc = INTERRUPT_TIMER_JUMP_VECTOR
	// wait for 5 M-cycles = 20 T-cycles
	cpu.cpuCycles += 5 * 4
}

func (cpu *CPU) onSerialInterrupt() {
	// clear the interrupt request flag before jumping to the interrupt handler @0x0050
	if_register := cpu.bus.Read(IF_REGISTER) & ((1 << FF0F_3_SERIAL) ^ 0xFF)
	cpu.setIEFlag(uint16(if_register))
	// update the program counter to have the address of the next instruction that cpu would have execute if no interrupt was triggered and push it to the stack}
	cpu.updatepc()
	cpu.push(cpu.pc)
	// trigger the interrupt handler
	cpu.pc = INTERRUPT_SERIAL_JUMP_VECTOR
	// wait for 5 M-cycles = 20 T-cycles
	cpu.cpuCycles += 5 * 4
}

func (cpu *CPU) onJoypadInterrupt() {
	// clear the interrupt request flag before jumping to the interrupt handler @0x0050
	if_register := cpu.bus.Read(IF_REGISTER) & ((1 << FF0F_4_JOYPAD) ^ 0xFF)
	cpu.setIEFlag(uint16(if_register))
	// update the program counter to have the address of the next instruction that cpu would have execute if no interrupt was triggered and push it to the stack}
	cpu.updatepc()
	cpu.push(cpu.pc)
	// trigger the interrupt handler
	cpu.pc = INTERRUPT_JOYPAD_JUMP_VECTOR
	// wait for 5 M-cycles = 20 T-cycles
	cpu.cpuCycles += 5 * 4
}
