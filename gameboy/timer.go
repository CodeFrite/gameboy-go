package gameboy

// Registers
// * DIV - Divider Register (R/W)
// This register is incremented at rate of 16384Hz (~16779Hz on SGB). Writing any value to this register resets it to 0x00.
// Additionally, this register is reset when executing the stop instruction, and only begins ticking again once stop mode ends.
// TODO: check that the stop instruction is reseting the DIV register correctly
// * TIMA - Timer Counter (R/W)
// This timer is incremented at the clock frequency specified by the TAC register ($FF07).
// When the value overflows (exceeds $FF) it is reset to the value specified in TMA (FF06) and an interrupt is requested.
// * TMA - Timer Modulo (R/W)
// When TIMA overflows, it is reset to the value in this register and an interrupt is requested.
// Example of use: if TMA is set to $FF, an interrupt is requested at the clock frequency selected in TAC
// (because every increment is an overflow).
// However, if TMA is set to $FE, an interrupt is only requested every two increments,
// which effectively divides the selected clock by two. Setting TMA to $FD would divide the clock by three, and so on.
// If a TMA write is executed on the same M-cycle as the content of TMA is transferred to TIMA due to a timer overflow,
// the old value is transferred to TIMA.
//   - TAC - Timer Control (R/W)
//     |	7	6	5	4	3	2				1			0
//
// TAC	|						Enable	Clock select
// + Enable: Controls whether TIMA is incremented. Note that DIV is always counting, regardless of this bit.
// + Clock select: Controls the frequency at which TIMA is incremented, as follows:
//
//	Clock select		Increment every		DMG, SGB2, CGB in single-speed mode		SGB1 Frequency (Hz) 	CGB in double-speed mode
//						00			 256 M-cycles			 														 4096		  						~4194												8192
//						01				 4 M-cycles																 262144								~268400											524288
//						10				16 M-cycles																	65536		 						 ~67110											131072
//						11				64 M-cycles																	16384								 ~16780											 32768
var TIMER_REGISTERS map[string]uint16 = map[string]uint16{
	"DIV":  0xFF04,
	"TIMA": 0xFF05,
	"TMA":  0xFF06,
	"TAC":  0xFF07,
}

type Timer struct {
	// state
	Enabled bool   `json:"enabled"`
	count   uint16 `json:"count"`

	// registers
	Div   uint16 `json:"div"`   // Divider Register
	Tima  uint8  `json:"tima"`  // Timer Counter
	Tma   uint8  `json:"tma"`   // Timer Modulo
	Tac   uint8  `json:"tac"`   // Timer Control
	Timer uint16 `json:"timer"` // Timer
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Increment() {
	t.count++
}

func (t *Timer) Run() {
	for t.Enabled {
		t.count++
	}
}

func (t *Timer) Reset() {
	t.count = 0
}

func (t *Timer) Enable() {
	t.Enabled = true
	t.Run()
}

func (t *Timer) Disable() {
	t.Enabled = false
}
