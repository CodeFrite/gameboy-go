package gameboy

type Bus struct {
	cartridge *Cartridge
	vram *VRAM
	wram *WRAM
}

func NewBus(c *Cartridge, vram *VRAM, wram *WRAM) *Bus {
	return &Bus{
		cartridge: c,
		vram: vram,
		wram: wram,
	}
}

func (b *Bus) Read(addr [2]byte) byte {
	addr16 := bytesToUint16(addr)
	if addr16 <= 0x7FFF {
		return b.cartridge.Read(addr)
	} else if addr16 >= 0x8000 && addr16 <= 0x9FFF {
		return b.vram.Read(addr16 - 0x8000)
	} else if addr16 >= 0xC000 && addr16 <= 0xDFFF {
		return b.wram.Read(addr16 - 0xC000)
	}
	return 0
}

func (b *Bus) Write(addr [2]byte, value byte) {
	addr16 := bytesToUint16(addr)
	if addr16 <= 0x7FFF {
		b.cartridge.Write(addr, value)
	} else if addr16 >= 0x8000 && addr16 <= 0x9FFF {
		b.vram.Write(addr16 - 0x8000, value)
	} else if addr16 >= 0xC000 && addr16 <= 0xDFFF {
		b.wram.Write(addr16 - 0xC000, value)
	}
}