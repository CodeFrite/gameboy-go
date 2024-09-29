package gameboy

import "fmt"

type PpuState struct {
	Mode uint8
	x, y uint8
}

func (p *PPU) getState() *PpuState {
	return &PpuState{
		Mode: p.mode,
		x:    p.dotX,
		y:    p.dotY,
	}
}

func (p *PpuState) print() {
	fmt.Println()
	fmt.Println("PPU> state:")
	fmt.Println("-----------")
	fmt.Printf("mode: %d\n", p.Mode)
	fmt.Printf("x: %d\n", p.x)
	fmt.Printf("y: %d\n", p.y)
	fmt.Println()
}
