package gameboy

type Memory struct {
	data []uint8
}

func NewMemory(size uint16) *Memory {
	return &Memory{data: make([]uint8, size)}
}

func (m *Memory) Size() uint16 {
	return uint16(len(m.data))
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.data[addr]
}

func (m *Memory) Write(addr uint16, value uint8) {
	m.data[addr] = value
}