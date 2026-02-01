package gui

import (
	"fmt"
	"os"

	"github.com/codefrite/gameboy-go/gameboy"
	"github.com/veandco/go-sdl2/sdl"
)

// DMG (original Game Boy) color palette - 4 shades of gray
var palette = [4]sdl.Color{
	{R: 224, G: 248, B: 208, A: 255}, // 0: lightest (white)
	{R: 136, G: 192, B: 112, A: 255}, // 1: light gray
	{R: 52, G: 104, B: 86, A: 255},   // 2: dark gray
	{R: 8, G: 24, B: 32, A: 255},     // 3: darkest (black)
}

type GUI struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	areas    map[string]sdl.Rect
}

func NewGUI() (*GUI, error) {
	gui := &GUI{}
	error := gui.initialize()
	return gui, error
}

// initializes the SDL window, renderer and define the screen regions (menu, LCD, debugger, ...). Returns any error encountered during the process
func (g *GUI) initialize() error {
	err := g.initSDLWindow()
	if err != nil {
		return err
	}
	err = g.createSDLRenderer()
	if err != nil {
		return err
	}
	g.createScreenAreas()
	return nil
}

func (g *GUI) initSDLWindow() error {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init SDL: %s\n", err)

	}

	// Create window
	window, err := sdl.CreateWindow(
		"Gameboy-Go",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(gameboy.LCD_X_RESOLUTION), int32(gameboy.LCD_Y_RESOLUTION),
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		return fmt.Errorf("failed to create SDL window: %s", err)
	}

	g.window = window
	return nil
}

func (g *GUI) createSDLRenderer() error {
	renderer, err := sdl.CreateRenderer(g.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return fmt.Errorf("failed to create SDL renderer: %s", err)
	}

	g.renderer = renderer
	return nil
}

func (g *GUI) createScreenAreas() {
	g.areas = make(map[string]sdl.Rect)
	g.areas["lcd"] = sdl.Rect{X: 0, Y: 0, W: 160, H: 144}
}

func (g *GUI) LCDDrawPixel(x int, y int, color sdl.Color) {
	if area, ok := g.areas["lcd"]; ok {
		if (x >= 0 && x < int(area.W)) && (y >= 0 && y < int(area.H)) {
			g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			g.renderer.DrawPoint(int32(x), int32(y))
		}
	}
}

func (g *GUI) LCDDrawRect(x int, y int, w int, h int, color sdl.Color) {
	if area, ok := g.areas["lcd"]; ok {
		if (x+w) >= 0 && (x+w) < int(area.W) && (y+h) >= 0 && (y+h) < int(area.H) {
			g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
			rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
			g.renderer.FillRect(&rect)
		}
	}
}

func (g *GUI) LCDDrawImage(image gameboy.RenderedImage) {

	// Iterate through each row (y)
	for y := 0; y < int(gameboy.LCD_Y_RESOLUTION); y++ {
		// Iterate through each packed byte in the row
		for packedX := 0; packedX < int(gameboy.LCD_X_RESOLUTION)/4; packedX++ {
			packedByte := image[y][packedX]

			// Extract 4 pixels from this byte (each pixel is 2 bits)
			// Pixels are packed from most significant bits to least significant bits
			for pixelInByte := 0; pixelInByte < 4; pixelInByte++ {
				shift := 6 - (pixelInByte * 2)
				colorIndex := (packedByte >> shift) & 0x03
				color := palette[colorIndex]
				g.LCDDrawPixel(packedX*4+pixelInByte, y, color)
			}
		}
	}
}

func (g *GUI) LCDClear() {
	g.renderer.SetDrawColor(0, 0, 0, 255)
	g.renderer.Clear()
}

func (g *GUI) LCDPresent() {
	g.renderer.Present()
}
