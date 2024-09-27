// A Program is a annotated list of instructions corresponding to the translated code of a gameboy ROM.
package gameboy

type ExecutionContext struct {
	reads  []MemoryWrite // list of memory reads
	writes []MemoryWrite // list of memory writes
}

type Step struct {
	prefixed bool   // is the instruction prefixed by 0xCB
	opcode   string // instruction meta data (mnemonic, cycles, operands, flags, ...) in .json format
	comment  string // links an address to a comment .md file name
}

type Program struct {
	address uint16          // address of the program in the program memory
	romName string          // name of the ROM file from which the program is extracted
	rawCode []uint8         // raw code extracted from the ROM file
	code    map[uint16]Step // list of instructions with their address as key
}

func (p *Program) NewProgram(romName string, address uint16, rawCode []uint8) *Program {
	return &Program{
		address: address,
		romName: romName,
		rawCode: rawCode,
		code:    make(map[uint16]Step),
	}
}

func (p *Program) Save() {
	panic("Program> 'save' modifier not implemented")
}

func (p *Program) IsAddressInMemory(addr uint16) bool {
	return addr >= p.address && addr < p.address+uint16(len(p.rawCode))
}

func (p *Program) LoadComment(addr uint16, comment string) {
	temp := p.code[addr]
	temp.comment = comment
	p.code[addr] = temp
}
