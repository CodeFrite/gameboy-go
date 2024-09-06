package gameboy

/* Represents the GameBoy bus which connects the CPU, PPU, APU, RAM, ROM and other components together.
 * TODO: Apparently in my implementation, I am just calling the MMU through the bus without adding any extra functionality.
 * I should either remember why I made this choice and introduce some extra functionality or remove either the bus or mmu.
 */
type Bus struct {
	mmu *MMU
}

/* constructor for the Bus struct */
func NewBus() *Bus {
	return &Bus{
		mmu: NewMMU(),
	}
}

/* attach a memory struct to the bus
 * @param name: of the memory used in the frontend
 * @param address: starting address of the memory
 * @param memory: to attach to the bus
 * @return nothing
 */
func (b *Bus) AttachMemory(name string, address uint16, memory Accessible) {
	b.mmu.AttachMemory(name, address, memory)
}

/* reads a uint8 data from the bus at a given address
 * @param addr: uint16 address to read from
 * @return uint8 value found at the address
 */
func (b *Bus) Read(addr uint16) uint8 {
	return b.mmu.Read(addr)
}

/* reads a uint16 data from the bus at a given address
 * @param addr: uint16 address to read from
 * @return uint16 value found at the address
 */
func (b *Bus) Read16(addr uint16) uint16 {
	return b.mmu.Read16(addr)
}

/* dumps the memory content from a starting address to an ending address
 * @param from: uint16 starting address
 * @param to: uint16 ending address
 * @return []uint8 slice of memory content
 */
func (b *Bus) Dump(from uint16, to uint16) []uint8 {
	return b.mmu.Dump(from, to)
}

/* writes a uint8 value to the bus at a given address
 * @param addr: uint16 address to write to
 * @param value: uint8 value to write
 * @return nothing
 */
func (b *Bus) Write(addr uint16, value uint8) {
	b.mmu.Write(addr, value)
}

/* writes a blob of uint8 values to the bus at a given address
 * @param addr: uint16 address to write to
 * @param blob: []uint8 slice of values to write
 * @return nothing
 */
func (b *Bus) WriteBlob(addr uint16, blob []uint8) {
	b.mmu.WriteBlob(addr, blob)
}
