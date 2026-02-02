package gameboy

// update the STAT register FF41 to reflect the current PPU mode
func (p *PPU) updateSTATRegister_PPUMode() {
	// get the register value
	stat := p.bus.Read(REG_FF41_STAT)
	// update the bits 1-0 with the current mode and leave the other bits unchanged
	stat = (stat & 0b11111100) | p.mode
	// write the new value back to the register
	p.bus.Write(REG_FF41_STAT, stat)
}

// update the LY register FF44 to reflect the current scanline and trigger the VBLANK interrupt
func (p *PPU) updateLYRegister() {
	p.bus.Write(REG_FF44_LY, uint8(p.dotY))
}

// evaluate if a STAT interrupt should be triggered based on the current PPU mode and LY=LYC condition
func (p *PPU) evaluateSTATInterrupt() {
	stat := p.bus.Read(REG_FF41_STAT)

	// Check for mode-based STAT interrupts
	mode := p.mode
	var modeInterruptEnabled bool
	switch mode {
	case PPU_MODE_0_HBLANK:
		modeInterruptEnabled = (stat>>FF41_3_MODE_0_HBLANK_SELECT)&0x01 == 1
	case PPU_MODE_1_VBLANK:
		modeInterruptEnabled = (stat>>FF41_4_MODE_1_VBLANK_SELECT)&0x01 == 1
	case PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM:
		modeInterruptEnabled = (stat>>FF41_5_MODE_2_OAM_SELECT)&0x01 == 1
	}

	if modeInterruptEnabled {
		p.requestSTATInterrupt()
		return
	}

	// Check for LY=LYC STAT interrupt
	ly := p.bus.Read(REG_FF44_LY)
	lyc := p.bus.Read(REG_FF45_LYC)
	if ly == lyc {
		lycInterruptEnabled := (stat>>FF41_6_MODE_3_LYC_SELECT)&0x01 == 1
		if lycInterruptEnabled {
			p.requestSTATInterrupt()
		}
	}
}

// request VBLANK interrupt
func (p *PPU) requestVBLANKInterrupt() {
	if_register := p.bus.Read(IF_REGISTER)
	p.bus.Write(IF_REGISTER, if_register|FF0F_0_VBLANK)
}

// request STAT interrupt
func (p *PPU) requestSTATInterrupt() {
	if_register := p.bus.Read(IF_REGISTER)
	p.bus.Write(IF_REGISTER, if_register|FF0F_1_LCD_STAT)
}
