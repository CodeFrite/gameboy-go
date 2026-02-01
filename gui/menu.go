package gui

//===> MENU
// Menu system for the emulator:
// + Recursive list of options to chose from
// + Each menu item is rendered as a Label (string)
// + Each option can either
// 		- trigger an action named Callback
// 		- open a new sub-menu allowing selecting from multiple options
// 		- opening a dialog for inputting a value validated by a Validator func which returns a boolean

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type MenuType int

const (
	MTAction  MenuType = iota // triggers the execution of the associated callback
	MTSubMenu                 // opens a sub-menu for further selection
	MTInt                     // opens a dialog for the user to input an integer value
	MTFloat                   // opens a dialog for the user to input a float value
	MTString                  // opens a dialog for the user to input a string value
	MTVector                  // opens a dialog for the user to input a vector value
)

type MenuItem struct {
	Type      MenuType                    // Type of menu item
	Label     string                      // Display label for the menu item
	Action    func()                      // Callback action to be executed when the item is selected
	Validator func(string) (bool, string) // Validation function for user input
}

type Menu struct {
	items        []MenuItem
	selectedItem []int  // keeps track of the selected item and potentially sub-menu selections: [1] or [1, 2] (max 2 elements for now)
	subMenuData  string // holds the current input for sub-menu items
}

func NewMenu(items []MenuItem) Menu {
	return Menu{
		items:        items,
		selectedItem: []int{0},
	}
}

func (m *Menu) Draw(screen *ebiten.Image) {
	// Draw menu background
	menuX, menuY := 50, 50
	menuWidth, menuHeight := 150, 100
	vector.DrawFilledRect(screen, float32(menuX), float32(menuY), float32(menuWidth), float32(menuHeight), color.RGBA{50, 50, 50, 200}, false)

	// Draw menu border
	vector.DrawFilledRect(screen, float32(menuX), float32(menuY), float32(menuWidth), 2, color.White, false)
	vector.DrawFilledRect(screen, float32(menuX), float32(menuY+menuHeight-2), float32(menuWidth), 2, color.White, false)
	vector.DrawFilledRect(screen, float32(menuX), float32(menuY), 2, float32(menuHeight), color.White, false)
	vector.DrawFilledRect(screen, float32(menuX+menuWidth-2), float32(menuY), 2, float32(menuHeight), color.White, false)

	// Draw menu title
	ebitenutil.DebugPrintAt(screen, "File Menu", menuX+10, menuY+10)

	// Draw menu items
	if len(m.selectedItem) == 0 {
		for i, item := range m.items {
			y := menuY + 40 + i*20
			prefix := "  "
			if i == m.selectedItem[0] {
				prefix = "> "
			}
			ebitenutil.DebugPrintAt(screen, prefix+item.Label, menuX+10, y)
		}
	} else {
		// Draw sub-menu

	}
}

func (m *Menu) Next() {
	// can only be called on Level 0 menu, not on sub-menus
	if len(m.selectedItem) == 1 {
		m.selectedItem[0] = (m.selectedItem[0] + 1) % len(m.items)

	}
}

func (m *Menu) Previous() {
	// can only be called on Level 0 menu, not on sub-menus
	if len(m.selectedItem) == 1 {
		m.selectedItem[0] = (m.selectedItem[0] - 1 + len(m.items)) % len(m.items)
	}
}

func (m *Menu) Validate() {
	if len(m.selectedItem) == 1 && m.items[m.selectedItem[0]].Type == MTAction {
		m.items[m.selectedItem[0]].Action()
	} else {

	}
}

func (m *Menu) OnKeyPress(k ebiten.Key) {
	if len(m.selectedItem) == 2 {

		// any alpha or numerical key: append to sub-menu data
		if k >= ebiten.KeyA && k <= ebiten.KeyZ || k >= ebiten.Key0 && k <= ebiten.Key9 {
			m.subMenuData += k.String()
		}

		// backspace: erase the last character if any
		if k == ebiten.KeyBackspace && len(m.subMenuData) > 0 {
			m.subMenuData = m.subMenuData[:len(m.subMenuData)-1]
			return
		}

		// Enter
	}
}

/*


func (e *Emulator) drawScaleInputDialog(screen *ebiten.Image) {
	// Draw dialog background
	x, y := 40, 60
	w, h := 240, 100
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{60, 60, 60, 230}, false)
	// Border
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), 2, color.White, false)
	vector.DrawFilledRect(screen, float32(x), float32(y+h-2), float32(w), 2, color.White, false)
	vector.DrawFilledRect(screen, float32(x), float32(y), 2, float32(h), color.White, false)
	vector.DrawFilledRect(screen, float32(x+w-2), float32(y), 2, float32(h), color.White, false)
	// Title
	ebitenutil.DebugPrintAt(screen, "Set Scale Factor (1-10)", x+10, y+10)
	// Input field
	ebitenutil.DebugPrintAt(screen, "> "+e.scaleInputValue, x+10, y+40)
	// Error message
	if e.scaleInputError != "" {
		ebitenutil.DebugPrintAt(screen, e.scaleInputError, x+10, y+70)
	}
}


func (e *Emulator) handleMenuSelection() {
	switch e.selectedItem {
	case 0: // Open ROM
		e.openRomDialog()
	case 1: // Change scale factor
		e.changeScaleFactor()
	case 2: // Exit
		os.Exit(0)
	}
	e.showMenu = false
}

func (e *Emulator) openRomDialog() {
	// Set flag to show dialog is opening
	e.dialogOpening = true

	// Run the file dialog in a goroutine to avoid blocking the main thread
	go func() {
		// Set initial directory to the roms folder
		romsPath := "/Users/codefrite/Documents/CODE/codefrite-emulator/gameboy/gameboy-go/roms"
		filename, err := dialog.File().Filter("Game Boy ROMs", "gb").SetStartDir(romsPath).Load()

		// Reset dialog flag
		e.dialogOpening = false

		if err != nil {
			if err != dialog.ErrCancelled {
				log.Printf("Error opening file dialog: %v", err)
			}
			return
		}
		e.selectedRom = filename
		log.Printf("Selected ROM: %s", filename)
	}()
}

func (e *Emulator) changeScaleFactor() {
	e.showScaleInput = true
	e.scaleInputValue = ""
	e.scaleInputError = ""
}

*/
