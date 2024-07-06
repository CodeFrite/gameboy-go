package gameboy

type Memory struct {
	data []uint8
}

func NewMemory(size uint16) *Memory {
	return &Memory{data: make([]uint8, size)}
}

func (m *Memory) setData(data []uint8) {
	m.data = data
}

func (m *Memory) Size() uint16 {
	return uint16(len(m.data))
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.data[addr]
}

/*
 * Dump memory from address 'from' to address 'to'
 */
func (m *Memory) Dump(from uint16, to uint16) []uint8 {
	return m.data[from:to]
}

func (m *Memory) Write(addr uint16, value uint8) {
	m.data[addr] = value
}
