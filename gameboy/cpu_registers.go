package gameboy

/*
 * F flags register
 * 7 6 5 4 3 2 1 0 (position)
 * Z N H C 0 0 0 0 (flag)
 */

// Zero Flag operations
// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return c.F&0x80 == 0x80
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.F = c.F | 0x80
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.F = c.F & 0x7F
}

// Toggle the Z flag in the F register
func (c *CPU) toggleZFlag() {
	c.F = c.F ^ 0x80
}

// Carry Flag operations
// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return c.F&0x40 == 0x40
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.F = c.F | 0x40
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.F = c.F & 0xBF
}

// Toggle the N flag in the F register
func (c *CPU) toggleNFlag() {
	c.F = c.F ^ 0x40
}

// Half Carry Flag operations
// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return c.F&0x20 == 0x20
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.F = c.F | 0x20
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.F = c.F & 0xDF
}

// Toggle the H flag in the F register
func (c *CPU) toggleHFlag() {
	c.F = c.F ^ 0x20
}

// Carry Flag operations
// Get the C flag from the F register
func (c *CPU) getCFlag() bool {
	return c.F&0x10 == 0x10
}

// Set the C flag in the F register
func (c *CPU) setCFlag() {
	c.F = c.F | 0x10
}

// Reset the C flag in the F register
func (c *CPU) resetCFlag() {
	c.F = c.F & 0xEF
}

// Toggle the C flag in the F register
func (c *CPU) toggleCFlag() {
	c.F = c.F ^ 0x10
}

/*
 * 16-bit registers accessors
 */
func (c *CPU) getBC() uint16 {
	return uint16(c.B)<<8 | uint16(c.C)
}

func (c *CPU) setBC(value uint16) {
	c.B = byte(value >> 8)
	c.C = byte(value)
}

func (c *CPU) getDE() uint16 {
	return uint16(c.D)<<8 | uint16(c.E)
}

func (c *CPU) setDE(value uint16) {
	c.D = byte(value >> 8)
	c.E = byte(value)
}

func (c *CPU) setHL(value uint16) {
	c.H = byte(value >> 8)
	c.L = byte(value)
}

func (c *CPU) getHL() uint16 {
	return uint16(c.H)<<8 | uint16(c.L)
}

// Interrupt Enable
func (c *CPU) getIE() uint8 {
	return c.bus.Read(0xFFFF)
}

func (c *CPU) setIE(value uint8) {
	c.bus.Write(0xFFFF, value)
}
