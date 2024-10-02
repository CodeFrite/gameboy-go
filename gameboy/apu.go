// Audio Processing Unit (APU) for the Gameboy
package gameboy

type APU struct {
	sound bool
}

type ApuState struct {
	sound bool
}

func NewAPU() *APU {
	return &APU{
		sound: false,
	}
}

func (a *APU) getState() *ApuState {
	return &ApuState{
		sound: a.sound,
	}
}

func (a *APU) onTick() {
	a.sound = !a.sound
}
