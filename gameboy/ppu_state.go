package gameboy

import "fmt"

type PpuState struct {
	MODE  uint8 `json:"MODE"`
	DOT_X uint8 `json:"DOT_X"`
	DOT_Y uint8 `json:"DOT_Y"`
}

func (p *PPU) getState() PpuState {
	return PpuState{
		MODE:  p.mode,
		DOT_X: p.dotX,
		DOT_Y: p.dotY,
	}
}

func (p *PpuState) print() {
	fmt.Println()
	fmt.Println("PPU> state:")
	fmt.Println("-----------")
	fmt.Printf("mode: %d\n", p.MODE)
	fmt.Printf("x: %d\n", p.DOT_X)
	fmt.Printf("y: %d\n", p.DOT_Y)
	fmt.Println()
}
