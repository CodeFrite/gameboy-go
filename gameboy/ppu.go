// Picture Processing Unit (PPU) for the Gameboy
package gameboy

type PPU struct {
	// Memory
	vram [0x2000]byte
	oam  [0xA0]byte
	// Registers
	lcdc byte
	stat byte
	scy  byte
	scx  byte
	ly   byte
	lyc  byte
	dma  byte
	bgp  byte
	obp0 byte
	obp1 byte
	wy   byte
	wx   byte
	// Internal registers
	mode byte
	tick byte
}

func NewPPU() *PPU {
	return &PPU{}
}
