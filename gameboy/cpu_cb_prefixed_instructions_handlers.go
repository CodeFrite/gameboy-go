package gameboy

// > instructions handlers (PREFIX CB)
func (c *CPU) RLC(instruction *Instruction) {
	panic("RLC not implemented")
}
func (c *CPU) RRC(instruction *Instruction) {
	panic("RRC not implemented")
}
func (c *CPU) RL(instruction *Instruction) {
	panic("RL not implemented")
}
func (c *CPU) RR(instruction *Instruction) {
	panic("RR not implemented")
}
func (c *CPU) SLA(instruction *Instruction) {
	panic("SLA not implemented")
}
func (c *CPU) SRA(instruction *Instruction) {
	panic("SRA not implemented")
}
func (c *CPU) SWAP(instruction *Instruction) {
	panic("SWAP not implemented")
}
func (c *CPU) SRL(instruction *Instruction) {
	panic("SRL not implemented")
}

/*
 BIT b, r8 / [HL]
 Test bit b in register r8 or [HL]: If bit b is 0, Z is set.
 opcodes: 0x40-0x7F
 flags: Z=Z N=0 H=1 C=-
*/
func (c *CPU) BIT(instruction *Instruction) {
	// get the bit position to test
	b := c.bus.Read(c.PC + 1)
	// check if bit b of operand is 0
	if c.Operand&(1<<b) == 0 {
		c.setZFlag()
	} else {
		c.resetZFlag()
	}
	c.resetNFlag()
	c.setHFlag()
	// increment the program counter
	c.incrementPC(uint16(instruction.Bytes))
}
func (c *CPU) RES(instruction *Instruction) {
	panic("RES not implemented")
}
func (c *CPU) SET(instruction *Instruction) {
	panic("SET not implemented")
}
