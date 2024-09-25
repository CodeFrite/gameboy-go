package gameboy

const (
	Z_FLAG_POSITION = 7
	N_FLAG_POSITION = 6
	H_FLAG_POSITION = 5
	C_FLAG_POSITION = 4
)

/*
 * F flags register
 * 7 6 5 4 3 2 1 0 (position)
 * Z N H C 0 0 0 0 (flag)
 */

// Zero Flag operations
// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return c.F.GetBit(Z_FLAG_POSITION)
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.F.SetBit(Z_FLAG_POSITION)
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.F.ResetBit(Z_FLAG_POSITION)
}

// Carry Flag operations
// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return c.F.GetBit(N_FLAG_POSITION)
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.F.SetBit(N_FLAG_POSITION)
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.F.ResetBit(N_FLAG_POSITION)
}

// Half Carry Flag operations
// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return c.F.GetBit(H_FLAG_POSITION)
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.F.SetBit(H_FLAG_POSITION)
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.F.ResetBit(H_FLAG_POSITION)
}

// Carry Flag operations
// Get the C flag from the F register
func (c *CPU) getCFlag() bool {
	return c.F.GetBit(C_FLAG_POSITION)
}

// Set the C flag in the F register
func (c *CPU) setCFlag() {
	c.F.SetBit(C_FLAG_POSITION)
}

// Reset the C flag in the F register
func (c *CPU) resetCFlag() {
	c.F.ResetBit(C_FLAG_POSITION)
}

/*
 * 16-bit registers accessors
 */
func (c *CPU) getBC() uint16 {
	return uint16(c.B)<<8 | uint16(c.C)
}

func (c *CPU) setBC(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.B.Set(high)
	c.C.Set(low)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

func (c *CPU) setDE(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.D.Set(high)
	c.E.Set(low)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

func (c *CPU) setHL(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.H.Set(high)
	c.L.Set(low)
}

func (c *CPU) GetIEFlag() uint8 {
	return c.bus.mmu.Read(IE_FLAG_START)
}

func (c *CPU) setIEFlag(value uint16) {
	c.bus.Write(IE_FLAG_START, byte(value))
}
