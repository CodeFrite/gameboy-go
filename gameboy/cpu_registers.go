package gameboy

const (
	Z_FLAG_POSITION uint8 = 0x07
	N_FLAG_POSITION uint8 = 0x06
	H_FLAG_POSITION uint8 = 0x05
	C_FLAG_POSITION uint8 = 0x04
)

/*
 * F flags register
 * 7 6 5 4 3 2 1 0 (position)
 * Z N H C 0 0 0 0 (flag)
 */

func getFlag(value uint8, position uint8) bool {
	return value&(0x01<<position) != 0
}

func setFlag(value uint8, position uint8) uint8 {
	return value | (0x01 << position)
}

func resetFlag(value uint8, position uint8) uint8 {
	return value & ^(0x01 << position)
}

// Zero Flag operations
// Get the Z flag from the F register
func (c *CPU) getZFlag() bool {
	return getFlag(c.f, Z_FLAG_POSITION)
}

// Set the Z flag in the F register
func (c *CPU) setZFlag() {
	c.f = setFlag(c.f, Z_FLAG_POSITION)
}

// Reset the Z flag in the F register
func (c *CPU) resetZFlag() {
	c.f = resetFlag(c.f, Z_FLAG_POSITION)
}

// Carry Flag operations
// Get the N flag from the F register
func (c *CPU) getNFlag() bool {
	return getFlag(c.f, N_FLAG_POSITION)
}

// Set the N flag in the F register
func (c *CPU) setNFlag() {
	c.f = setFlag(c.f, N_FLAG_POSITION)
}

// Reset the N flag in the F register
func (c *CPU) resetNFlag() {
	c.f = resetFlag(c.f, N_FLAG_POSITION)
}

// Half Carry Flag operations
// Get the H flag from the F register
func (c *CPU) getHFlag() bool {
	return getFlag(c.f, H_FLAG_POSITION)
}

// Set the H flag in the F register
func (c *CPU) setHFlag() {
	c.f = setFlag(c.f, H_FLAG_POSITION)
}

// Reset the H flag in the F register
func (c *CPU) resetHFlag() {
	c.f = resetFlag(c.f, H_FLAG_POSITION)
}

// Carry Flag operations
// Get the C flag from the F register
func (c *CPU) getCFlag() bool {
	return getFlag(c.f, C_FLAG_POSITION)
}

// Set the C flag in the F register
func (c *CPU) setCFlag() {
	c.f = setFlag(c.f, C_FLAG_POSITION)
}

// Reset the C flag in the F register
func (c *CPU) resetCFlag() {
	c.f = resetFlag(c.f, C_FLAG_POSITION)
}

/*
 * 16-bit registers accessors
 */
func (c *CPU) getBC() uint16 {
	return uint16(c.b)<<8 | uint16(c.c)
}

func (c *CPU) setBC(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.b = high
	c.c = low
}

func (c *CPU) getDE() uint16 {
	return uint16(c.d)<<8 | uint16(c.e)
}

func (c *CPU) setDE(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.d = high
	c.e = low
}

func (c *CPU) getHL() uint16 {
	return uint16(c.h)<<8 | uint16(c.l)
}

func (c *CPU) setHL(value uint16) {
	low := uint8(value)
	high := uint8(value >> 8)
	c.h = high
	c.l = low
}

func (c *CPU) GetIEFlag() uint8 {
	return c.bus.Read(IE_FLAG_START)
}

func (c *CPU) setIEFlag(value uint16) {
	c.bus.Write(IE_FLAG_START, byte(value))
}
