package gameboy

import "errors"

type PpuState struct {
	MODE  uint8         `json:"MODE"`
	DOT_X uint8         `json:"DOT_X"`
	DOT_Y uint8         `json:"DOT_Y"`
	IMAGE RenderedImage `json:"IMAGE"`
}

// Returns the current state of the PPU if the PPU is enabled
func (p *PPU) getState() (*PpuState, error) {
	if p.isEnabled() {
		return &PpuState{
			MODE:  p.mode,
			DOT_X: p.dotX,
			DOT_Y: p.dotY,
			IMAGE: p.image,
		}, nil
	} else {
		return nil, errors.New("PPU is disabled")
	}
}
