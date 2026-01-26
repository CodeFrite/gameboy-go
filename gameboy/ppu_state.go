package gameboy

import "errors"

type PpuState struct {
	ImageNumber uint64        `json:"imageNumber"`
	MODE        uint8         `json:"MODE"`
	DOT_X       uint16        `json:"DOT_X"`
	DOT_Y       uint16        `json:"DOT_Y"`
	IMAGE       RenderedImage `json:"IMAGE"`
}

// Returns the current state of the PPU if the PPU is enabled
func (p *PPU) getState() (*PpuState, error) {
	if p.isEnabled() {
		return &PpuState{
			ImageNumber: p.ticks / DOTS_PER_FRAME,
			MODE:        p.mode,
			DOT_X:       p.dotX,
			DOT_Y:       p.dotY,
			IMAGE:       p.image,
		}, nil
	} else {
		return nil, errors.New("PPU is disabled")
	}
}
