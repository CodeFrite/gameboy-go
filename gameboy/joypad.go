package gameboy

// Gameboy Joypad
// --------------
// Only one register FF00 is used to read the joypad state since the gameboy has only one joypad controller.
// FF00 lower nibble is read-only and upper nibble is read-write.
// In order to read the joypad state, the cpu must reset bit:
// - FF00.4 to read the direction pad state from FF00.0-3: Right, Left, Up, Down (low = pressed)
// - FF00.5 to read the buttons state from FF00.0-3: A, B, Start, Select (low = pressed)
// When a button (D-pad or A, B, Start, Select) is pressed, an interrupt is requested (IF).
// If IME & IE are set for the joypad, the cpu will jump to vector 0x60.

const (
	REG_FF00_JOYP            = 0xFF00
	FF00_0_A_RIGHT_BUTTON    = 0 // if 0, A or RIGHT is being pressed
	FF00_1_B_LEFT_BUTTON     = 1 // if 0, B or LEFT is being pressed
	FF00_2_SELECT_UP_BUTTON  = 2 // if 0, SELECT or UP is being pressed
	FF00_3_START_DOWN_BUTTON = 3 // if 0, START or DOWN is being pressed
	FF00_4_SELECT_DPAD       = 4 // if 0, then direction pad can be read from lower nibble
	FF00_5_SELECT_BUTTONS    = 5 // if 0, then buttons can be read from lower nibble

	FF00_INITIAL_STATE = 0x3F // all buttons are released and nor the buttons neither the direction pad are selected for reading
)

type JoypadEvent struct {
	Up, Down, Left, Right, A, B, Start, Select bool
}

// Whenever a button is pressed or released in the frontend, the Joypad state is updated.
// It is only when the FF00.4/5 bits are set that the Joypad state is effectively stored in the FF00 register.
type Joypad struct {
	// state of the buttons in the same order as the FF00 register: Right, Left, Up, Down, A, B, Select, Start
	state uint8
}

// returns a new joypad that listens to joypad input events on the eventChannel
func NewJoypad() *Joypad {
	return &Joypad{
		state: FF00_INITIAL_STATE,
	}
}

// MMU redirects the write to the joypad register FF00 to the joypad
func (j *Joypad) Write(value uint8) {
	// if FF00.4 is reset, the state of the direction pad is written to FF00.0-3
	if value&(1<<FF00_4_SELECT_DPAD) == 0 && value&(1<<FF00_5_SELECT_BUTTONS) == 1 {

	}
}
