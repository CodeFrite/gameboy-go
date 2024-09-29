// Picture Processing Unit (PPU)
// -----------------------------
// + responsible for rendering the graphics in memory to the LCD screen
// + runs @ the crystal clock frequency to produce 59.727 images per second
// + each image counts 144 * 160 pixels
// + a PPU rendering cycle is divided into 154 lines of 456 dots (dot = 1 PPU cycle during which a pixel can be drawed to compare to 3 CPU stages to execute 1 instruction)
// + during these 4,194,304 ticks, the PPU will render 144 lines of 546 dots out of which 160 are used to draw the 160 pixels of the current line on the LCD screen
// + after drawing the current line, 87-204 crystal cycles will be spent in the HBlank period during which the PPU is idle and the CPU can access all the memory
// + scanlines 145-153 are spent in the VBlank period during which the PPU is idle and the CPU can access all the memory
//
// During its operations, the PPU cycles between 4 modes during a scanline:
//
// Mode		Action																				Duration										Accessible video memory
// ----   ------																				--------										-----------------------
// 2			Searching for OBJs which overlap this line		1 - 80 dots									VRAM, CGB palettes
// 3			Sending pixels to the LCD											172 - 289 dots							None
// 0			Waiting until the end of the scanline					376 - mode 3’s duration			VRAM, OAM, CGB palettes
// 1			Waiting until the next frame									4560 dots (10 scanlines)		VRAM, OAM, CGB palettes

package gameboy

import (
	"fmt"
	"math/rand/v2"
)

const (
	// modes
	MODE_2_SEARCH_OVERLAP_OBJ_OAM uint8 = 2 // Searching for OBJs which overlap this line
	MODE_3_SEND_PIXEL_LCD         uint8 = 3 // Sending pixels to the LCD
	MODE_0_HBLANK                 uint8 = 0 // Waiting until the end of the scanline
	MODE_1_VBLANK                 uint8 = 1 // Waiting until the next frame

	MODE2_LENGTH uint8 = 80

	// memory
	OAM_MEMORY_START_ADDRESS uint16 = 0xFE00
	OAM_MEMORY_BYTE_SIZE     uint8  = 0xA0

	// registers
	LCDC_ADDR uint16 = 0xFF40 // LCD Control
	STAT_ADDR uint16 = 0xFF41 // LCD Status
	SCY_ADDR  uint16 = 0xFF42 // LCD Scroll Y
	SCX_ADDR  uint16 = 0xFF43 // LCD Scroll X
	LY_ADDR   uint16 = 0xFF44 // LCD Y
	LYC_ADDR  uint16 = 0xFF45 // LCD Y Compare
	DMA_ADDR  uint16 = 0xFF46 // DMA Transfer
	BGP_ADDR  uint16 = 0xFF47 // BG Palette
	OBP0_ADDR uint16 = 0xFF48 // OBP0 Palette
	OBP1_ADDR uint16 = 0xFF49 // OBP1 Palette
	WY_ADDR   uint16 = 0xFF4A // Window Y
	WX_ADDR   uint16 = 0xFF4B // Window X

	// a dot is a PPU cycle and not a pixel

	// LCD properties
	DOTS_PER_LINE   uint16 = 456
	LINES_PER_FRAME uint16 = 154
	DOTS_PER_FRAME  uint64 = uint64(LINES_PER_FRAME) * uint64(DOTS_PER_LINE) // 70,224 dots per frame

	// interrupt parameters
	LCD_X_RESOLUTION uint8  = 160
	LCD_Y_RESOLUTION uint8  = 144
	PIXELS_PER_FRAME uint64 = uint64(LCD_X_RESOLUTION) * uint64(LCD_Y_RESOLUTION) // 23,040 pixels per frame
)

var (
	STATE = map[byte]string{
		0: "TILE",
		1: "DATA0",
		2: "DATA1",
		3: "IDLE",
		4: "PUSH",
	}
)

type PPU struct {
	// reference to the cpu to trigger interrupts
	cpu *CPU
	bus *Bus

	// Screen
	mode  uint8  // current mode
	ticks uint64 // should be able to count up 160 x 144 = 23,040,000
	dotX  uint8  // current scanline dot x position (0-455)
	dotY  uint8  // current scanline dot y position (0-153)

	// simulation parameters: random values for now TODO: remove this after implementing mode3
	mode3Length uint // length of mode 3 (sending pixels to the LCD)

	// Registers
	lcdc Register8 // lcd control
	stat Register8 // lcd status
	scy  Register8 // scroll y
	scx  Register8 // scroll x
	ly   Register8 // lcd y
	lyc  Register8 // lcd y compare
	dma  Register8 // dma transfer
	bgp  Register8 // bg palette
	obp0 Register8 // obj palette 0
	obp1 Register8 // obj palette 1
	wx   Register8 // window x
	wy   Register8 // window y

	// Memory
	oam *Memory // Object Attribute Memory (0xFE00-0xFE9F) - 40 4-byte entries
}

func NewPPU(cpu *CPU, bus *Bus) *PPU {
	ppu := &PPU{cpu: cpu, bus: bus}
	// initialize memory
	ppu.oam = NewMemory(uint16(OAM_MEMORY_BYTE_SIZE))
	ppu.bus.AttachMemory("OAM", OAM_MEMORY_START_ADDRESS, ppu.oam)
	return ppu
}

func (p *PPU) updateLy() {
	// increment y line counter
	line := uint8(p.ticks / uint64(DOTS_PER_LINE))
	// update LY register
	p.cpu.bus.Write(LY_ADDR, line)
	// trigger LYC=LY interrupt TODO: check that there are no other conditions to trigger this interrupt (like on STAT register)
	if line == p.lyc.uint8 {
		p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["LCD"]})
	}
}

// onTick is called at each crystal cycle
// there are 4,194,304 ticks per
func (p *PPU) onTick() {
	p.ticks++

	// new frame: reset ticks, dot x & y, LY register
	if p.ticks%uint64(DOTS_PER_FRAME) == 0 {
		p.ticks = 0
		p.dotX = 0
		p.dotY = 0
		p.mode3Length = uint(rand.IntN(289-172) + 172) // random value for now TODO: remove this after processing mode3

		// new scanline
	} else if p.ticks%uint64(DOTS_PER_LINE) == 0 {
		p.updateLy()
		p.dotX = 0
		p.dotY++

		// processing scan line
	} else {
		p.dotX++
	}

	// managing mode switching
	if p.ticks == 0 {
		p.mode = MODE_2_SEARCH_OVERLAP_OBJ_OAM

	} else if p.dotY < LCD_Y_RESOLUTION && p.dotX == MODE2_LENGTH {
		p.mode = MODE_3_SEND_PIXEL_LCD

	} else if p.dotY < LCD_Y_RESOLUTION && p.dotX == MODE2_LENGTH+uint8(p.mode3Length) {
		p.mode = MODE_0_HBLANK

	} else if p.dotY >= LCD_Y_RESOLUTION {
		p.mode = MODE_1_VBLANK
	}

	// processing data
	switch p.mode {
	case MODE_2_SEARCH_OVERLAP_OBJ_OAM:
		// searching for OBJs which overlap this line
		fmt.Println("PPU> mode 2: searching for OBJs which overlap this line")
	case MODE_3_SEND_PIXEL_LCD:
		// sending pixels to the LCD
		fmt.Println("PPU> mode 3: sending pixels to the LCD")
	case MODE_0_HBLANK:
		// waiting until the end of the scanline
		fmt.Println("PPU> mode 0: waiting until the end of the scanline")
		// p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["HBLANK"]}) ??? TODO: check if this is an interrupt that i should trigger or if the CPU pools the PPU mode for this
	case MODE_1_VBLANK:
		// waiting until the next frame
		fmt.Println("PPU> mode 1: waiting until the next frame")
		p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["VBLANK"]})
	}
}
