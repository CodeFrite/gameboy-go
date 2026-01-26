# Gameboy emulator in Go

Gameboy emulator written in Go. It is a work in progress and is not yet complete. It relies on the documentation found at [GBDEV/Pan](https://gbdev.io/pandocs/About.html)

# Usage

# Architecture

The gameboy package contains the implementation of the different components of the Gameboy. The main components are:

- **CPU:** runs the instructions read from the bootrom and the cartridge. It contains the different registers and flags and accesses the memory through the bus
- **Memory:** contain the different memory regions of the Gameboy (ROM, RAM, VRAM, etc). It is accessed through the bus
- **Bus/MMU:** the bus allows the CPU, GPU and APU to access the memory. It also contains the interrupt controller. The MMU allows to select the right memory region from which to read/write data
- **Cartridge:** contains the program aka 'game' to run. It might also contain additional circuitry to provide additional functionalities (e.g. MBC, save states, sound modules, ...) which are currently out of scope
- **GPU:** the processing made by this unit are comparable to what the CPU does with instructions: it transforms data present in the VRAM so that they are ready to be rendered on screen. Rendering the image to the screen is left to the consumer application to allow for different implementations (e.g. terminal, SDL, webapp, ...). What we are basically e
- **APU:** the audio processing unit is responsible for generating sound. Here again, the actual sound generation is left to the consumer application to allow for different implementations (e.g. terminal, SDL, webapp, ...). What is taken care of here is the generation of the sound data
- **Timer:** the timer is responsible for keeping track of the time and generating interrupts when the timer overflows

To glue everything together, the Gameboy struct is used. It contains the different components and provides the necessary methods to run the emulator.

Finally, the Debugger struct is used to provide debugging functionalities to the emulator. It contains the Gameboy struct and provides methods to run the emulator in debug mode and inspect the state of the different components.
