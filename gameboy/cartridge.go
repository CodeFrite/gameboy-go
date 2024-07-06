package gameboy

import (
	"fmt"
)

type Cartridge struct {
	cartridgePath     string
	cartridgeName     string
	rom               *Memory
	header            []uint8
	entry_point       []uint8
	nintendo_logo     []uint8
	title             []uint8
	manufacturer_code []uint8
	cgb_flag          []uint8
	new_licensee_code []uint8
	sgb_flag          []uint8
	cartridge_type    []uint8
	rom_size          []uint8
	ram_size          []uint8
	destination_code  []uint8
	old_licensee_code []uint8
	mask_rom_version  []uint8
	header_checksum   []uint8
	global_checksum   []uint8
}

func NewCartridge(uri string, name string) *Cartridge {
	var c Cartridge
	rom, err := LoadRom(uri + "/" + name)
	if err != nil {
		fmt.Println("Error loading ROM:", err)
		return nil
	}
	c.rom = NewMemory(uint16(len(rom[0x0100:])))
	c.rom.setData(rom[0x0100:])
	c.cartridgePath = uri
	c.cartridgeName = name
	c.parseHeader(rom)
	return &c
}

func (c *Cartridge) parseHeader(rom []uint8) {
	c.header = rom[0x0100:0x014F]
	c.entry_point = rom[0x0100:0x0103]
	c.nintendo_logo = rom[0x0104:0x0133]
	c.title = rom[0x0134:0x0143]
	c.manufacturer_code = rom[0x013F:0x0142]
	c.cgb_flag = rom[0x0143:0x0144]
	c.new_licensee_code = rom[0x0144:0x0145]
	c.sgb_flag = rom[0x0146:0x0147]
	c.cartridge_type = rom[0x0147:0x0148]
	c.rom_size = rom[0x0148:0x0149]
	c.ram_size = rom[0x0149:0x014A]
	c.destination_code = rom[0x014A:0x014B]
	c.old_licensee_code = rom[0x014B:0x014C]
	c.mask_rom_version = rom[0x014C:0x014D]
	c.header_checksum = rom[0x014D:0x014E]
	c.global_checksum = rom[0x014E:0x0150]
}

func (c *Cartridge) PrintInfo() {
	// metadata about the cartridge
	fmt.Println("Cartridge Path:", c.cartridgePath)
	fmt.Println("Cartridge Name:", c.cartridgeName)
	fmt.Println("Cartridge Size:", c.rom_size, "bytes")

	// header information
	fmt.Println("Header:", c.header)
	fmt.Println("Entry Point:", c.entry_point)
	fmt.Println("Nintendo Logo:", c.nintendo_logo)
	fmt.Println("Title:", string(c.title))
	fmt.Println("Manufacturer Code:", c.manufacturer_code)
	fmt.Println("CGB Flag:", c.cgb_flag)
	fmt.Println("New Licensee Code:", c.new_licensee_code)
	fmt.Println("SGB Flag:", c.sgb_flag)
	fmt.Println("Cartridge Type:", c.cartridge_type)
	fmt.Println("ROM Size:", c.rom_size)
	fmt.Println("RAM Size:", c.ram_size)
	fmt.Println("Destination Code:", c.destination_code)
	fmt.Println("Old Licensee Code:", c.old_licensee_code)
	fmt.Println("Mask ROM Version:", c.mask_rom_version)
	fmt.Println("Header Checksum:", c.header_checksum)
	fmt.Println("Global Checksum:", c.global_checksum)
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.rom.Read(addr)
}

func (c *Cartridge) Dump(from uint16, to uint16) []uint8 {
	return c.rom.Dump(from, to)
}

func (c *Cartridge) Write(addr uint16, value uint8) {
	c.rom.Write(addr, value)
}

func (c *Cartridge) Size() uint16 {
	return c.rom.Size()
}
