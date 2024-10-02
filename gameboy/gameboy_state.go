package gameboy

import (
	"fmt"
)

/*
func (d *Debugger) printCPUState() {
	fmt.Println("")
	fmt.Println("\n> CPU State:")
	fmt.Println("------------")
	// if previous and current states are nil, there is nothing to print
	if (gbs.PREV_CPU_STATE == nil) && (gbs.CURR_CPU_STATE == nil) {
		fmt.Println("> No CPU state to print")
		// if only the current state is available, print it
	} else if gbs.PREV_CPU_STATE == nil {
		curr := reflect.Indirect(reflect.ValueOf(gbs.CURR_CPU_STATE))
		typeOfCpu := curr.Type()
		for i := 0; i < curr.NumField(); i++ {
			if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
				fmt.Printf("- %s: %t\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
				fmt.Printf("- %s: 0x%02X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
				fmt.Printf("- %s: 0x%04X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
				fmt.Printf("- %s: %s\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
			}
		}
	} else {
		prev := reflect.Indirect(reflect.ValueOf(gbs.PREV_CPU_STATE))
		curr := reflect.Indirect(reflect.ValueOf(gbs.CURR_CPU_STATE))
		typeOfCpu := prev.Type()

		for i := 0; i < prev.NumField(); i++ {
			if prev.Field(i).Interface() != curr.Field(i).Interface() {
				if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
					fmt.Printf("- %s: %t -> %t \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
					fmt.Printf("- %s: 0x%02X -> 0x%02X \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
					fmt.Printf("- %s: 0x%04X -> 0x%04X \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
					fmt.Printf("- %s: %s -> %s \n", typeOfCpu.Field(i).Name, prev.Field(i).Interface(), curr.Field(i).Interface())
				}
			} else {
				if typeOfCpu.Field(i).Type.Kind() == reflect.Bool {
					fmt.Printf("- %s: %t \n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint8 {
					fmt.Printf("- %s: 0x%02X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.Uint16 {
					fmt.Printf("- %s: 0x%04X\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				} else if typeOfCpu.Field(i).Type.Kind() == reflect.String {
					fmt.Printf("- %s: %s\n", typeOfCpu.Field(i).Name, curr.Field(i).Interface())
				}
			}
		}
	}
}
*/

// currInstruction: returns the current instruction being processed based on cpu IR and prefix values
func (gb *Gameboy) currInstruction() *Instruction {
	instruction := GetInstruction(Opcode(fmt.Sprintf("0x%02X", gb.cpu.ir)), gb.cpu.prefixed)
	return &instruction
}

// clear memory writes
func (gb *Gameboy) clearMemoryWrites() {
	gb.cpuBus.mmu.clearMemoryWrites()
}

// returns the current memory writes
func (gb *Gameboy) currMemoryWrites() []MemoryWrite {
	return gb.cpuBus.mmu.memoryWrites
}

/**
 * print the properties of the memories attached to the bus
 */
func (gb *Gameboy) PrintMemoryProperties() {
	memoryMaps := gb.cpuBus.mmu.GetMemoryMaps()
	fmt.Println("")
	fmt.Println("\n> Memory Mapping:")
	fmt.Println("-----------------")
	for _, memoryMap := range memoryMaps {
		fmt.Printf("> Memory %s: %d bytes @ 0x%04X->0x%04X\n", memoryMap.Name, len(memoryMap.Data), memoryMap.Address, memoryMap.Address+uint16(len(memoryMap.Data))-1)
	}
}
