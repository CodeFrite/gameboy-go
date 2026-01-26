# Ebiten GUI Example for Gameboy Emulator

This is a toy example showing how to create a window with Ebiten for the Gameboy emulator.

## Installation Instructions

1. First, install the required Go modules:

```bash
# Install Ebiten (2D game library)
go get github.com/hajimehoshi/ebiten/v2

# Install dialog library for file selection
go get github.com/sqweek/dialog
```

2. Run the example:

```bash
go run gui_example.go
```

## Features

- **Window**: 400x600 pixels
- **Graphics**: Draws a 100x100 blue square at coordinates (200, 200)
- **Menu**: Press ESC to toggle File menu
  - "Open ROM..." - Opens a file dialog to select Game Boy ROM files
  - "Exit" - Closes the application
- **Navigation**: Use arrow keys to navigate menu, Enter to select

## Controls

- **ESC**: Toggle menu visibility
- **Arrow Keys**: Navigate menu items when menu is open
- **Enter**: Select menu item

## Requirements

- Go 1.16 or later
- Works on Windows, macOS, and Linux
- On Linux, you may need to install additional dependencies:

  ```bash
  # Ubuntu/Debian
  sudo apt-get install libc6-dev libglu1-mesa-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config

  # Fedora
  sudo dnf install mesa-libGL-devel mesa-libGLU-devel libXcursor-devel libXi-devel libXinerama-devel libXrandr-devel libXxf86vm-devel alsa-lib-devel pkg-config
  ```

## Integration with Gameboy Emulator

This example can be integrated with your existing Gameboy emulator by:

1. Importing your gameboy package
2. Creating gameboy instance in the Game struct
3. Rendering the gameboy screen buffer in the Draw method
4. Handling gameboy controls through Ebiten input events

Example integration:

```go
import "github.com/codefrite/gameboy-go/gameboy"

type Game struct {
    gb *gameboy.Gameboy
    // ... other fields
}
```
