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
	OAM_START uint16 = 0xFE00
	OAM_SIZE  uint8  = 0xA0

	LINES_PER_FRAME uint16 = 154
	DOTS_PER_LINE   uint16 = 456

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
	ticks uint64 // should be able to cound up to 4,194,304
	x     uint8  // x coordinate of the current pixel being drawn 0-159
	y     uint8  // y coordinate of the current line being drawn 0-143
	tile  uint8  // tile number

	// Memory
	oam *Memory // Object Attribute Memory (0xFE00-0xFE9F) - 40 4-byte entries

	// Registers
	lcdc uint8 // lcd control
	stat uint8 // lcd status
	scy  uint8 // scroll y
	scx  uint8 // scroll x
	ly   uint8 // lcd y
	lyc  uint8 // lcd y compare
	dma  uint8 // dma transfer
	bgp  uint8 // bg palette
	obp0 uint8 // obj palette 0
	obp1 uint8 // obj palette 1
	wy   uint8 // window y
	wx   uint8 // window x

	// Internal registers
	mode uint8 // current mode
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
	p.cpu.bus.Write(0xFF44, p.y)
	// trigger LYC=LY interrupt
	if p.y == p.lyc {
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
