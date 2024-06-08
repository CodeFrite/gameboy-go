package gameboy

// Convert [2]byte to uint16
func bytesToUint16(b [2]byte) uint16 {
	return uint16(b[0])<<8 | uint16(b[1])
}

// Convert uint16 to [2]byte
func uint16ToBytes(u uint16) [2]byte {
	return [2]byte{byte(u >> 8), byte(u)}
}