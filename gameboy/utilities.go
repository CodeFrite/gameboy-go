package gameboy

import (
	"errors"
	"fmt"
	"os"
)

func LoadRom(uri string) ([]byte, error) {
	rom, err := os.ReadFile(uri)
	if err != nil {
		errText := fmt.Sprint("Error loading ROM:", err)
		return nil, errors.New(errText)
	}
	fmt.Println("ROM loaded successfully:")
	for i, word := range rom {
		if i == 0 {
			fmt.Println()
			fmt.Print("      ")
			for j := 0; j < 16; j++ {
				fmt.Printf("%02X ", j)
			}
		}
		if i%16 == 0 {
			fmt.Println()
			fmt.Printf("%04X: ", i)
		}
		fmt.Printf("%02X ", word)
	}
	fmt.Print("\n\n")
	return rom, nil
}
