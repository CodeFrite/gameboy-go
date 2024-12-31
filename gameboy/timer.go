package gameboy

// DIV register is incremented at 16384Hz (always)
// TIMA register is incremented at a rate specified by TAC only if enabled (timer control register)
// When TIMA overflows, it is reset to TMA and an interrupt is requested
// TAC register is used to enable/disable the timer and set its frequency
// TAC.2: timer enable
// TAC.1-0: timer clock select = Frequency at which TIMA is incremented. If TMA = 0xFF, TIMA will overflow at the following frequencies:
// 00:   4,096 Hz = 4.194,304 MHz / 1024 T-cycles (256 M-cycles)
// 01: 262,144 Hz = 4.194,304 MHz /   16 T-cycles (  4 M-cycles)
// 10:  65,536 Hz = 4.194,304 MHz /   64 T-cycles ( 16 M-cycles)
// 11:  16,384 Hz	= 4.194,304 MHz /  256 T-cycles ( 64 M-cycles)

const (
	// Timer Special Registers
	REG_FF04_DIV  = 0xFF04 // divider register: incremented at a rate of 16384 Hz. Always incremented, regardless of the TAC
	REG_FF05_TIMA = 0xFF05 // timer counter: incremented at the rate specified by TAC. Overflows when it reaches 0xFF
	REG_FF06_TMA  = 0xFF06 // timer modulo: when TIMA overflows, it is reset to TMA
	REG_FF07_TAC  = 0xFF07 // timer control: enables/disables the timer TIMA and sets its frequency
)

// generates a tick at a given frequency
type Timer struct {
	bus           *Bus
	internalClock uint16
}

func NewTimer(bus *Bus) *Timer {
	return &Timer{
		bus: bus,
	}
}

// on tick, increment DIV and TIMA registers if enabled
func (t *Timer) Tick() {
	// increment the internal M-cycle clock
	t.internalClock++

	// always increment DIV register at 16384Hz independently of the TAC
	if t.internalClock%256 == 0 {
		div := t.bus.Read(REG_FF04_DIV)
		div++
		t.bus.mmu.timerWrite(REG_FF04_DIV, div)
	}

	// increment TIMA register if enabled based on the TAC clock select
	tac := t.bus.Read(REG_FF07_TAC)
	tima_enabled := tac&0x04 == 0x04
	tima_clock_select := tac & 0x03

	// when TIMA overflows, reset it to TMA and request an interrupt
	tima := t.bus.Read(REG_FF05_TIMA)

	// increment TIMA at the rate specified by TAC
	if tima_enabled {
		switch tima_clock_select {
		// 00:   4,096 Hz = 4.194,304 MHz / 1024 T-cycles (256 M-cycles)
		case 0x00:
			if t.internalClock%1024 == 0 {
				if tima == 0xFF {
					tima = t.bus.Read(REG_FF06_TMA)
					if_register := t.bus.Read(IF_REGISTER)
					t.bus.Write(IF_REGISTER, if_register|0x04)
				} else {
					tima++
				}
			}
		// 01: 262,144 Hz = 4.194,304 MHz /   16 T-cycles (  4 M-cycles)
		case 0x01:
			if t.internalClock%16 == 0 {
				if tima == 0xFF {
					tima = t.bus.Read(REG_FF06_TMA)
					if_register := t.bus.Read(IF_REGISTER)
					t.bus.Write(IF_REGISTER, if_register|0x04)
				} else {
					tima++
				}
			}
		// 10:  65,536 Hz = 4.194,304 MHz /   64 T-cycles ( 16 M-cycles)
		case 0x02:
			if t.internalClock%64 == 0 {
				if tima == 0xFF {
					tima = t.bus.Read(REG_FF06_TMA)
					if_register := t.bus.Read(IF_REGISTER)
					t.bus.Write(IF_REGISTER, if_register|0x04)
				} else {
					tima++
				}
			}
		// 11:  16,384 Hz	= 4.194,304 MHz /  256 T-cycles ( 64 M-cycles)
		case 0x03:
			if t.internalClock%256 == 0 {
				if tima == 0xFF {
					tima = t.bus.Read(REG_FF06_TMA)
					if_register := t.bus.Read(IF_REGISTER)
					t.bus.Write(IF_REGISTER, if_register|0x04)
				} else {
					tima++
				}
			}
		}
	}

	// update the TIMA register
	t.bus.Write(REG_FF05_TIMA, tima)
}
