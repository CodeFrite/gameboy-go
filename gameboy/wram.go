package gameboy

type WRAM struct {
	// Work RAM 8 KiB [0xC000-0xDFFF]
	data [0x2000]byte
}

func NewWRAM() *WRAM {
	return &WRAM{}
}

func (v *WRAM) Read(addr uint16) byte {
	return v.data[addr]
}

func (v *WRAM) Write(addr uint16, value byte) {
	v.data[addr] = value
}