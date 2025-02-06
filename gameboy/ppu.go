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
// 2			Searching for OBJs which overlap this line		1-80 dots			  						VRAM, CGB palettes
// 3			Sending pixels to the LCD											172-289 dots  							None
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
//	|   Tile IDs for...				|	Block 0				|		Block 1				|		Block 2			|
//	|   -------------------		|	----------		|		----------		|		----------	|
//	|													|	$8000–87FF		|		$8800–8FFF		|		$9000–97FF	|
//	|		Objects								|	0–127					|		128–255				|		—						|
//	|		BG/Win, if LCDC.4=1		|	0–127					|		128–255				|		—						|	   using unsigned id (0-255)
//	|		BG/Win, if LCDC.4=0		|	—							|		128–255				|		0–127				|		 using signed id (-128 to 127)
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
	TILE_MAP_0_START_ADDRESS uint16 = 0x9800 // up to 0x9BFF (32x32 tiles of 8x8 pixels with 2 bits color depth
	TILE_MAP_1_START_ADDRESS uint16 = 0x9C00 // up to 0x9FFF (same as above)
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
	FF41_7_UNUSED               uint8 = 7 // unused

	// a dot is a PPU cycle and not a pixel

	// LCD properties
	DOTS_PER_LINE   uint64 = 456
	LINES_PER_FRAME uint64 = 154
	DOTS_PER_FRAME  uint64 = LINES_PER_FRAME * DOTS_PER_LINE // 70,224 dots per frame

	// interrupt parameters
	LCD_X_RESOLUTION uint8  = 160
	LCD_Y_RESOLUTION uint8  = 144
	PIXELS_PER_FRAME uint64 = uint64(LCD_X_RESOLUTION) * uint64(LCD_Y_RESOLUTION) // 23,040 pixels per frame
)

// RenderedImage is a 2D array of 144x160 pixels
// Each pixel can have up to 4 colors which represents 2 bits of information
// if I code every pixel as a uint8 when it only needs 2 bits, I get a state size of 144 * 160 * 2 = 46,080 bits = 5.76 kB
// To optimize the size exchange, I will regroup pixels color information by groups of 4x2 bits = 1 byte
// For a single line, I need 160 * 2 bits = 320 bits, represented as 40 * 8 bits = 40 bytes
// In summary, the rendered 2D image (superposition of background, window & objects) will be arranged as 144 lines of 40 bytes
// In this struct, coord (y, x) will correspond to pixels on line y from [x * 4, (x+1) * 4[
type RenderedImage [LCD_Y_RESOLUTION][LCD_X_RESOLUTION / 4]uint8

type PPU struct {
	bus *Bus

	// Screen
	image      RenderedImage  // rendered image
	background [256][64]uint8 // background layer: 256x256 pixels each coding a color in a byte with optimization to easily extract the rendered image
	mode       uint8          // current mode
	ticks      uint64         // should be able to count up 160 x 144 = 23,040,000
	dotX       uint16         // current scanline dot x position (0-455)
	dotY       uint16         // current scanline dot y position (0-153)

	// simulation parameters: random values for now TODO: remove this after implementing mode3
	mode3Length uint16 // length of mode 3 (sending pixels to the LCD)

	// Memory
	oam *Memory // Object Attribute Memory (0xFE00-0xFE9F) - 40 4-byte entries
}

func NewPPU(bus *Bus) *PPU {
	ppu := &PPU{
		bus:        bus,
		image:      RenderedImage{},
		background: [256][64]uint8{},
		oam:        NewMemory(uint16(OAM_MEMORY_BYTE_SIZE)),
	}
	// attach the OAM memory to the bus
	ppu.bus.AttachMemory("OAM", OAM_MEMORY_START_ADDRESS, ppu.oam)
	return ppu
}

func (p *PPU) reset() {
	p.ticks = 0
	p.dotX = 0
	p.dotY = 0
	p.mode = PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM
	p.mode3Length = uint16(rand.IntN(289-172) + 172) // random value to simulate the real hardware processing
	p.image = RenderedImage{}
	p.background = [256][64]uint8{}
}

// onTick is called at each crystal cycle
// there are 4,194,304 ticks per frame
func (p *PPU) Tick() {
	// check if the PPU & LCD are enabled. If not, set image to blank color and return
	if !p.isEnabled() {
		//fmt.Printf("PPU is disabled (ticks=%d / dotX=%d)\n", p.ticks, p.dotX)
		return
	}

	// new frame: reset ticks, dot x & y, LY register, image
	if p.ticks%DOTS_PER_FRAME == 0 {
		p.dotX = 0
		p.dotY = 0
		p.mode3Length = uint16(172 + rand.IntN(289-172)) // random value to simulate the real hardware processing
	} else if p.ticks%DOTS_PER_LINE == 0 {
		// new scanline
		p.dotX = 0
		p.dotY++
		p.updateLYRegister()

	} else {
		// processing scan line
		p.dotX++
	}

	//fmt.Printf("PPU is processing (ticks:%d / dotX=%d)\n", p.ticks, p.dotX)

	// managing mode switching
	if p.dotY < uint16(LCD_Y_RESOLUTION) && p.dotX == 0 {
		p.mode = PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM
		p.updateSTATRegister_PPUMode()
	} else if p.dotY < uint16(LCD_Y_RESOLUTION) && p.dotX == uint16(MODE2_LENGTH) {
		p.mode = PPU_MODE_3_SEND_PIXEL_LCD
		p.updateSTATRegister_PPUMode()
	} else if p.dotY < uint16(LCD_Y_RESOLUTION) && p.dotX == uint16(MODE2_LENGTH)+p.mode3Length {
		p.mode = PPU_MODE_0_HBLANK
		//fmt.Printf("... SETTING PPU MODE TO 0 (mode3 length=%d, pivot=%d)\n", p.mode3Length, MODE2_LENGTH+uint8(p.mode3Length))
		p.updateSTATRegister_PPUMode()
	} else if p.dotY >= uint16(LCD_Y_RESOLUTION) {
		p.mode = PPU_MODE_1_VBLANK
		p.updateSTATRegister_PPUMode()
	}

	// processing data
	switch p.mode {
	case PPU_MODE_2_SEARCH_OVERLAP_OBJ_OAM:
		// MODE 2: searching for OBJs which overlap this line
		//fmt.Println("... PPU MODE 2 (SEARCH OVERLAP OBJ OAM)")
	case PPU_MODE_3_SEND_PIXEL_LCD:
		// MODE 3: sending pixels to the LCD

		// During the first 160 dots of mode 3 which can be 172-289 dots long, we will compute the pixel color for the screen pixel SCX + dotX - 80
		// After the 160 dots, we will wait until the end of mode 3
		if p.dotX >= uint16(MODE2_LENGTH)+uint16(LCD_X_RESOLUTION) {
			break // and not return because we want to increment the ticks
		}
		//fmt.Printf("... PPU MODE 3 (pixel:%d)\n", p.dotX-uint16(MODE2_LENGTH))

		// DRAW BACKGROUND
		lcdc := p.bus.Read(REG_FF40_LCDC)
		scx := p.bus.Read(REG_FF43_SCX)
		scy := p.bus.Read(REG_FF42_SCY)

		// 1. determine in which block the background tile maps are located by cheking LCDC.3: 0 -> $9800-$9BFF & 1 -> $9C00-$9FFF
		bgTileMapArea := (lcdc >> FF40_3_BG_TILE_MAP_AREA) & 0x01

		// 2. determine in which block the background tile data are located by cheking LCDC.4: 0 -> $8800-$97FF & 1 -> $8000-$8FFF
		bgTileDataArea := (lcdc >> FF40_4_BG_WINDOW_TILE_DATA_AREA) & 0x01

		// 3. compute the position of the pixel to draw relative to the background
		// will overflow at 256 (thanks to go) and that is what we want since the background is 256x256 and we want to loop over it with our viewport of 160x144
		pixelToDrawX := uint16(scx) + p.dotX - uint16(MODE2_LENGTH)
		pixelToDrawY := uint16(scy) + p.dotY

		// 4. compute the tile x, y position in the tile map
		tileX := pixelToDrawX / 8
		tileY := pixelToDrawY / 8

		// 5. get the tile index from the tile map (if LCDC.3 = 0 -> 0x9800-0x9BFF & if LCDC.3 = 1 -> 0x9C00-0x9FFF)
		tileId := uint8(0) // tile id in the tile map (0-1023)
		if bgTileMapArea == 0 {
			tileId = p.bus.Read(TILE_MAP_0_START_ADDRESS + uint16(tileY*32+tileX)) // background has 32x32 tiles
		} else {
			tileId = p.bus.Read(TILE_MAP_1_START_ADDRESS + uint16(tileY*32+tileX))
		}

		// 6. get the tile data from the tile data area (0x8000-0x8FFF or 0x8800-0x97FF depending on the tile data area LCDC.4) given that each tile is 16 bytes long (2 bytes per line)
		tileSubPixelX := pixelToDrawX % 8
		tileSubPixelY := pixelToDrawY % 8

		tileData := [2]uint8{}
		if bgTileDataArea == 0 {
			tileData[0] = p.bus.Read(BLOCK_1_START_ADDRESS + uint16(tileId)*16 + tileSubPixelY*2)
			tileData[1] = p.bus.Read(BLOCK_1_START_ADDRESS + uint16(tileId)*16 + tileSubPixelY*2 + 1)
		} else {
			tileData[0] = p.bus.Read(BLOCK_0_START_ADDRESS + uint16(tileId)*16 + tileSubPixelY*2)
			tileData[1] = p.bus.Read(BLOCK_0_START_ADDRESS + uint16(tileId)*16 + tileSubPixelY*2 + 1)
		}

		// 7. reconstruct the current pixel color
		pixel_low_bit := (tileData[0] >> (7 - tileSubPixelX)) & 0x01
		pixel_high_bit := (tileData[1] >> (7 - tileSubPixelX)) & 0x01

		// 8. Save the pixel color information to the image
		// compute the x, y position of the pixel in the image
		imageX := (p.dotX - uint16(MODE2_LENGTH)) / 4
		imageY := p.dotY

		// due to the fact that we save 4 pixel colors inside the same byte, when writing the current pixel color to the image, we need to:
		// - shift the current pixel color to the left by 2
		// - append the new pixel color information
		currentPixelSlotValue := p.image[imageY][imageX]
		p.image[imageY][imageX] = (currentPixelSlotValue << 2) | (pixel_high_bit << 1) | pixel_low_bit

	case PPU_MODE_0_HBLANK:
		//fmt.Println("... PPU MODE 0 (HBLANK)")
		// waiting until the end of the scanline
		// p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["HBLANK"]}) ??? TODO: check if this is an interrupt that i should trigger or if the CPU pools the PPU mode for this
	case PPU_MODE_1_VBLANK:
		//fmt.Println("... PPU MODE 1 (VBLANK)")
		// waiting until the next frame
	}

	// TODO: DRAW WINDOW

	// increment the ticks at the end of the cycle otherwise tick 0 will be skipped
	p.ticks++
}

// draw the background
func (p *PPU) drawBackground() {}

// draw the window
func (p *PPU) drawWindow() {}

// Registers management

// check whether the PPU is enabled
func (p *PPU) isEnabled() bool {
	lcdc := p.bus.Read(REG_FF40_LCDC)
	return (lcdc>>FF40_7_LCD_PPU_ENABLE)&0x01 == 0x01
}

// update the STAT register FF41 to reflect the current PPU mode
func (p *PPU) updateSTATRegister_PPUMode() {
	// get the register value
	stat := p.bus.Read(REG_FF41_STAT)
	// update the bits 1-0 with the current mode and leave the other bits unchanged
	stat = (stat & 0b11111100) | p.mode
	// write the new value back to the register
	p.bus.Write(REG_FF41_STAT, stat)
}

// update the LY register FF44 to reflect the current scanline and trigger the VBLANK interrupt
func (p *PPU) updateLYRegister() {
	p.bus.Write(REG_FF44_LY, uint8(p.dotY))
	if p.dotY == uint16(LCD_Y_RESOLUTION) {
		// set the LYC == LY bit
	}
}
