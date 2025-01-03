package gameboy

type PpuState struct {
	MODE  uint8 `json:"MODE"`
	DOT_X uint8 `json:"DOT_X"`
	DOT_Y uint8 `json:"DOT_Y"`
	IMAGE Image `json:"IMAGE"`
}

func (p *PPU) getState() PpuState {
	return PpuState{
		MODE:  p.mode,
		DOT_X: p.dotX,
		DOT_Y: p.dotY,
		IMAGE: p.image,
	}
}
