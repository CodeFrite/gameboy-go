package gameboy

type VRAM struct {
	// Video RAM 8 KiB [0x8000-0x9FFF]
	data [0x2000]byte
}

func NewVRAM() *VRAM {
	return &VRAM{}
}

func (v *VRAM) Read(addr uint16) byte {
	return v.data[addr]
}

func (v *VRAM) Write(addr uint16, value byte) {
	v.data[addr] = value
}