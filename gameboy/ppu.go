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
//
// The image produced by the PPU is composed of 3 layers:
// 1. Background layer (BG) - 256x256 pixels
// 2. Window layer (W) - 160x144 pixels
// 3. Object layer (OBJ) - 8x8 or 8x16 pixels
//
// Images are composed of tiles (8x8 pixels) which have a depth of 2 bits per pixel (4 colors) (8x8x2 = 128 bits = 16 bytes per tile)
// Tiles are stored in VRAM $8000-$97FF with each tile taking 16 bytes (2 bytes per line) meaning that there is a maximum of 384 tiles in VRAM
//
// Gameboy images are composed of tiles which are referenced by the tile map (BG Map) which is stored in VRAM $9800-$9BFF
// In order to recompose the image, we need a way to reference the tiles in the tile map:
//
//   - Tile IDs for...				+	Block 0				+		Block 1				+		Block 2			+
//     |   -------------------		|	----------		|		----------		|		----------	|
//     |													|	$8000–87FF		|		$8800–8FFF		|		$9000–97FF	|
//     |		Objects								|	0–127					|		128–255				|		—						|
//     |		BG/Win, if LCDC.4=1		|	0–127					|		128–255				|		—						|	   using unsigned id (0-255)
//     |		BG/Win, if LCDC.4=0		|	—							|		128–255				|		0–127				|		 using signed id (-128 to 127)
//
// There are two 32x32 (1024 bytes) tile maps in VRAM @$9800-$9BFF and @$9C00-$9FFF used to display the background or window layers.
// Each tile map reference 1 byte indexes of the tiles in VRAM @$8000-$97FF using an unsigned id (0-255) or a signed id (-128 to 127).
// Since each tile is 8x8 pixels, the BG map is 256x256 pixels. Only 160x144 of these pixels are displayed on the LCD screen (using scroll registers).
package gameboy

import (
	"math/rand/v2"
)

const (
	// modes
	PPU_MODE_0_HBLANK                 uint8 = 0 // Waiting until the end of the scanline
	PPU_MODE_1_VBLANK                 uint8 = 1 // Waiting until the next frame
	PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM uint8 = 2 // Searching for OBJs which overlap this line
	PPU_MODE_3_SEND_PIXEL_LCD         uint8 = 3 // Sending pixels to the LCD

	MODE2_LENGTH uint8 = 80

	// memory
	BLOCK_0_START_ADDRESS    uint16 = 0x8000
	BLOCK_1_START_ADDRESS    uint16 = 0x8800
	BLOCK_2_START_ADDRESS    uint16 = 0x9000
	TILE_MAP_0_START_ADDRESS uint16 = 0x9800
	TILE_MAP_1_START_ADDRESS uint16 = 0x9C00
	OAM_MEMORY_START_ADDRESS uint16 = 0xFE00
	OAM_MEMORY_BYTE_SIZE     uint8  = 0xA0

	// registers
	REG_FF40_LCDC uint16 = 0xFF40 // LCD Control
	REG_FF41_STAT uint16 = 0xFF41 // LCD Status
	REG_FF42_SCY  uint16 = 0xFF42 // LCD Scroll Y
	REG_FF43_SCX  uint16 = 0xFF43 // LCD Scroll X
	REG_FF44_LY   uint16 = 0xFF44 // LCD Y (set by the PPU)
	REG_FF45_LYC  uint16 = 0xFF45 // LCD Y Compare (set by the programmer)
	REG_FF46_DMA  uint16 = 0xFF46 // DMA Transfer
	REG_FF47_BGP  uint16 = 0xFF47 // BG Palette
	REG_FF48_OBP0 uint16 = 0xFF48 // OBP0 Palette
	REG_FF49_OBP1 uint16 = 0xFF49 // OBP1 Palette
	REG_FF4A_WY   uint16 = 0xFF4A // Window Y
	REG_FF4B_WX   uint16 = 0xFF4B // Window X

	// FF40 LCDC register bits
	FF40_0_BG_WINDOW_DISPLAY_ENABLE uint8 = 0 // if 0, BG and Window are white (overrides FF40_5)
	FF40_1_OBJ_DISPLAY_ENABLE       uint8 = 1 // if 0, OBJs are not displayed (can be changed mid-frame to prevent obj parts to be displayed over the window)
	FF40_2_OBJ_SIZE                 uint8 = 2 // if 0, 8x8, if 1, 8x16
	FF40_3_BG_TILE_MAP_AREA         uint8 = 3 // if 0, $9800-$9BFF, if 1, $9C00-$9FFF
	FF40_4_BG_WINDOW_TILE_DATA_AREA uint8 = 4 // if 0, unsigned addressing mode accross Block 0 & 1 (0-255), if 1, signed addressing mode accross Block 1 & 2 (-128 to 127)
	FF40_5_WINDOW_DISPLAY_ENABLE    uint8 = 5 // if 0, window is not displayed (overriden by FF40_0)
	FF40_6_WINDOW_TILE_MAP_AREA     uint8 = 6 // if 0, $9800-$9BFF, if 1, $9C00-$9FFF
	FF40_7_LCD_PPU_ENABLE           uint8 = 7 // if 0, LCD and PPU are turned off

	// FF41 STAT register bits
	FF41_01_PPU_MODE            uint8 = 1 // 2 bits to indicate the current PPU mode (0-3)
	FF41_2_LYC_EQ_LY            uint8 = 2 // 1 if LYC == LY
	FF41_3_MODE_0_HBLANK_SELECT uint8 = 3 // if set, trigger HBlank interrupt for STAT interrupt
	FF41_4_MODE_1_VBLANK_SELECT uint8 = 4 // if set, trigger VBlank interrupt for STAT interrupt
	FF41_5_MODE_2_OAM_SELECT    uint8 = 5 // if set, trigger OAM interrupt for STAT interrupt
	FF41_6_MODE_3_LYC_SELECT    uint8 = 6 // if set, trigger LCD interrupt for STAT interrupt
	// FF41_7_UNUSED               uint8 = 7 // unused

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

// 256x256 pixels, 4 colors
// if I code every pixel as a uint8 when it only needs 2 bits, I get a state size of 256x256x8 > 524 kB which exceeds the 64kB limit for channel message communication
// Using a bool instead of uint8 doesn't change the size of the state since a bool is represented as a byte in Go
// The size of the whole image (including the part not visible in the viewport) is 256*256*2 bits ∼ 16kB which is acceptable
// I will therefore regroup pixels color information by groups of 8 x 2bits = 1 byte
// I need 256 * 2 bits = 512 bits, represented as 64 * 8 bits = 64 bytes
// In summary, the 2D image (superposition of background, window & objects) will be arranged as 256 lines of 64 bytes
// In this struct, coord (x,y) will correspond to pixels on line y from [y*8, y*9[
type Image [256][32]uint8

type PPU struct {
	// reference to the cpu to trigger interrupts
	cpu *CPU
	bus *Bus

	// Screen
	image Image  // current image
	mode  uint8  // current mode
	ticks uint64 // should be able to count up 160 x 144 = 23,040,000
	dotX  uint8  // current scanline dot x position (0-455)
	dotY  uint8  // current scanline dot y position (0-153)

	// simulation parameters: random values for now TODO: remove this after implementing mode3
	mode3Length uint // length of mode 3 (sending pixels to the LCD)

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

func (p *PPU) reset() {
	p.ticks = 0
	p.dotX = 0
	p.dotY = 0
	p.mode = PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM
	p.mode3Length = uint(rand.IntN(289-172) + 172) // random value for now TODO: remove this after processing mode3
}

// onTick is called at each crystal cycle
// there are 4,194,304 ticks per
func (p *PPU) Tick() {
	// check if the PPU & LCD are enabled
	// TODO! change this position to return blank background and window if the LCD is off since this is the normal behavior. Good enough for now
	stat := p.cpu.bus.Read(REG_FF41_STAT)
	if (stat & 0b10000000) == 0 {
		return
	}

	// new frame: reset ticks, dot x & y, LY register, image
	if p.ticks%uint64(DOTS_PER_FRAME) == 0 {
		p.ticks = 0
		p.dotX = 0
		p.dotY = 0
		p.mode3Length = uint(rand.IntN(289-172) + 172) // random value for now TODO: remove this after processing mode3
	} else if p.ticks%uint64(DOTS_PER_LINE) == 0 {
		// new scanline
		p.dotX = 0
		p.dotY++
		p.updateLYRegister()

	} else {
		// processing scan line
		p.dotX++
	}

	// managing mode switching
	if p.ticks == 0 {
		p.mode = PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM
		p.updateSTATRegister_PPUMode()
	} else if p.dotY < LCD_Y_RESOLUTION && p.dotX == MODE2_LENGTH {
		p.mode = PPU_MODE_3_SEND_PIXEL_LCD
		p.updateSTATRegister_PPUMode()
	} else if p.dotY < LCD_Y_RESOLUTION && p.dotX == MODE2_LENGTH+uint8(p.mode3Length) {
		p.mode = PPU_MODE_0_HBLANK
		p.updateSTATRegister_PPUMode()
	} else if p.dotY >= LCD_Y_RESOLUTION {
		p.mode = PPU_MODE_1_VBLANK
		p.updateSTATRegister_PPUMode()
	}

	// processing data
	switch p.mode {
	case PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM:
		// searching for OBJs which overlap this line
	case PPU_MODE_3_SEND_PIXEL_LCD:
		// sending pixels to the LCD

		// first, determine from which block the background tiles are coming from
		lcdc := p.bus.Read(REG_FF40_LCDC)
		tileMapArea := (lcdc >> FF40_3_BG_TILE_MAP_AREA) & 0x01
		tileDataArea := (lcdc >> FF40_4_BG_WINDOW_TILE_DATA_AREA) & 0x01
		tileX := p.dotX / 8 // x position of the tile among the 32 tiles of the line
		tileY := p.dotY / 8 // y position of the tile among the 32 tiles of the column
		tileId := uint8(0)

		// get the tile index from the tile map (0x9800-0x9BFF or 0x9C00-0x9FFF depending on the tile map area LCDC.3)
		if tileMapArea == 0 {
			tileId = p.bus.Read(TILE_MAP_0_START_ADDRESS + uint16(tileY*32+tileX))
		} else {
			tileId = p.bus.Read(TILE_MAP_1_START_ADDRESS + uint16(tileY*32+tileX))
		}

		// get the tile data from the tile data area (0x8000-0x8FFF or 0x8800-0x97FF depending on the tile data area LCDC.4)
		pixelX := p.dotX % 8
		pixelY := p.dotY % 8

		tileData := [2]uint8{}
		if tileDataArea == 0 {
			tileData[0] = p.bus.Read(BLOCK_0_START_ADDRESS + uint16(tileId*16+pixelY*2+pixelX))
			tileData[1] = p.bus.Read(BLOCK_0_START_ADDRESS + uint16(tileId*16+pixelY*2+pixelX+1))
		} else {
			tileData[0] = p.bus.Read(BLOCK_1_START_ADDRESS + uint16(tileId*16+pixelY*2+pixelX))
			tileData[1] = p.bus.Read(BLOCK_1_START_ADDRESS + uint16(tileId*16+pixelY*2+pixelX+1))
		}

		// reconstruct the current pixel color and save it to the image
		pixel_low_bit := (tileData[0] >> (7 - pixelX)) & 0x01
		pixel_high_bit := (tileData[1] >> (7 - pixelX)) & 0x01

		// since pixel color information are regrouped by 4 in a single byte, we need to shift the bits to the left before adding the new pixel color information

		// compute the x, y position of the pixel in the image
		imageX := uint8(p.dotX / 4)
		imageY := uint8(p.dotY)

		// write the pixel color to the image
		currentPixelSlotValue := p.image[imageY][imageX]
		p.image[imageX][imageY] = (currentPixelSlotValue << 2) & (pixel_high_bit << 1) & pixel_low_bit

	case PPU_MODE_0_HBLANK:
		// waiting until the end of the scanline
		// p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["HBLANK"]}) ??? TODO: check if this is an interrupt that i should trigger or if the CPU pools the PPU mode for this
	case PPU_MODE_1_VBLANK:
		// waiting until the next frame
	}

	// increment the ticks at the end of the cycle otherwise tick 0 will be skipped
	p.ticks++
}

// draw the background
func (p *PPU) drawBackground() {}

// draw the window
func (p *PPU) drawWindow() {}

// Registers management

// update the STAT register FF41 to reflect the current PPU mode
func (p *PPU) updateSTATRegister_PPUMode() {
	// get the register value
	stat := p.cpu.bus.Read(REG_FF41_STAT)
	// update the bits 1-0 with the current mode and leave the other bits unchanged
	stat = (stat & 0b11111100) | p.mode
	// write the new value back to the register
	p.cpu.bus.Write(REG_FF41_STAT, stat)
}

// update the LY register FF44 to reflect the current scanline and trigger the VBLANK interrupt
func (p *PPU) updateLYRegister() {
	p.cpu.bus.Write(REG_FF44_LY, p.dotY)
	if p.dotY == LCD_Y_RESOLUTION {
		// set the LYC == LY bit
	}
}
