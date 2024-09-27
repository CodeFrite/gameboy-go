// Picture Processing Unit (PPU)
// -----------------------------
// + responsible for rendering the graphics in memory to the LCD screen
// + runs @ the crystal clock frequency to produce 59.727 images per second
// + each image counts 144 * 160 dots
// + a PPU rendering cycle is divided into 154 lines of 456 dots (dot = 1 PPU cycle during which a pixel is drawn to compare to 3 CPU stages to execute 1 instruction)
// + during these 4,194,304 ticks, the PPU will render 144 lines of 546 dots out of which 160 are used to draw the 160 pixels of the current line on the LCD screen
// + after drawing the current line, 87-204 crystal cycles will be spent in the HBlank period during which the PPU is idle and the CPU can access all the memory
// + scanlines 145-153 are spent in the VBlank period during which the PPU is idle and the CPU can access all the memory
package gameboy

const (
	// memory
	OAM_START uint16 = 0xFE00
	OAM_SIZE  uint8  = 0xA0

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

	// LCD properties
	LINES_PER_FRAME uint16 = 154
	DOTS_PER_LINE   uint16 = 456

	// interrupt parameters
	YRES uint8 = 144
	XRES uint8 = 160
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
	ticks uint64 // should be able to count up 160 x 144 =
	x     uint8  // x coordinate of the current pixel being drawn 0-159
	y     uint8  // y coordinate of the current line being drawn 0-143
	tile  uint8  // tile number
	mode  uint8  // current mode

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
	ppu.oam = NewMemory(uint16(OAM_SIZE))
	ppu.bus.AttachMemory("OAM", OAM_START, ppu.oam)
	return ppu
}

func (p *PPU) updateLy() {
	// increment y line counter
	p.y++
	// update LY register
	p.cpu.bus.Write(LY_ADDR, p.y)
	// trigger LYC=LY interrupt
	if p.y == p.lyc.uint8 {
		p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["LCD"]})
	}
}

func (p *PPU) onTick() {
	p.ticks++

	// reset ticks if we reach the end of the frame
	if p.ticks%uint64(DOTS_PER_LINE) == 0 {
		p.updateLy()
	} else if p.ticks%uint64(DOTS_PER_LINE)*uint64(LINES_PER_FRAME) == 0 {
		p.ticks = 0
		p.y = 0
		p.x = 0
	}

	// manage interrupts and draw pixels
	if p.ticks%uint64(DOTS_PER_LINE) < uint64(OAM_SIZE) {
		// 0-80 dots: OAM scan (mode 2)
	} else if p.ticks%uint64(DOTS_PER_LINE) > 80 && p.x < XRES {
		// 81- 81+160 dots (mode 3) TODO: this is not always true depending on sprites/background priority calculation. Could be more realist but good enough for now

		// draw pixel
		// increment x
		p.x++

	} else if p.y >= YRES {
		// trigger VBLANK interrupt on the CPU
		p.cpu.TriggerInterrupt(Interrupt{_type: interrupt_types["VBLANK"]})
	}
}
