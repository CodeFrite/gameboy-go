// Audio Processing Unit (APU) for the Gameboy
package gameboy

const (
	// Special Sound Registers
	NR10           = 0xFF10 // Sound channel 1 sweep (R/W)
	NR11           = 0xFF11 // Sound channel 1 length timer & duty cycle (Mixed)
	NR12           = 0xFF12 // Sound channel 1 volume & envelope (R/W)
	NR13           = 0xFF13 // Sound channel 1 period low (W)
	NR14           = 0xFF14 // Sound channel 1 period high & control (Mixed)
	NR21           = 0xFF16 // Sound channel 2 length timer & duty cycle (Mixed)
	NR22           = 0xFF17 // Sound channel 2 volume & envelope (R/W)
	NR23           = 0xFF18 // Sound channel 2 period low (W)
	NR24           = 0xFF19 // Sound channel 2 period high & control (Mixed)
	NR30           = 0xFF1A // Sound channel 3 DAC enable (R/W)
	NR31           = 0xFF1B // Sound channel 3 length timer (W)
	NR32           = 0xFF1C // Sound channel 3 output level (R/W)
	NR33           = 0xFF1D // Sound channel 3 period low (W)
	NR34           = 0xFF1E // Sound channel 3 period high & control (Mixed)
	NR41           = 0xFF20 // Sound channel 4 length timer (W)
	NR42           = 0xFF21 // Sound channel 4 volume & envelope (R/W)
	NR43           = 0xFF22 // Sound channel 4 frequency & randomness (R/W)
	NR44           = 0xFF23 // Sound channel 4 control (Mixed)
	NR50           = 0xFF24 // Master volume & VIN panning (R/W)
	NR51           = 0xFF25 // Sound panning (R/W)
	NR52           = 0xFF26 // Sound on/off (Mixed)
	WAVE_RAM_START = 0xFF30 // Wave RAM: Storage for one of the sound channels’ waveform (R/W)
	WAVE_RAM_END   = 0xFF3F
)

type APU struct {
	sound bool
}

type ApuState struct {
	SOUND bool `json:"SOUND"`
}

func NewAPU() *APU {
	return &APU{
		sound: false,
	}
}

func (a *APU) reset() {
	a.sound = false
}

func (a *APU) getState() ApuState {
	return ApuState{
		SOUND: a.sound,
	}
}

func (a *APU) Tick() {
	a.sound = !a.sound
}
