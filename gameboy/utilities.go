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
	fmt.Println("ROM", uri, "loaded successfully:")
	PrintByteTable(rom, 16)

	return rom, nil
}

func PrintByteTable(data []byte, maxLines int) {
	for i, word := range data {
		if i == 0 {
			fmt.Println()
			fmt.Print("      ")
			for j := 0; j < 16; j++ {
				colHeader := AnsiColor(fmt.Sprintf("%02X ", j), "yellow")
				fmt.Print(colHeader)
			}
		}
		// print only the first maxLines*16 bytes
		if i >= maxLines*16 {
			break
		}
		if i%16 == 0 {
			fmt.Println()
			lineHeader := AnsiColor(fmt.Sprintf("%04X: ", i), "yellow")
			fmt.Print(lineHeader)
		}
		fmt.Printf("%02X ", word)
	}
	fmt.Print("\n\n")
}

func AnsiColor(text string, color string) string {
	switch color {
	case "red":
		return fmt.Sprintf("\033[31m%v\033[0m", text)
	case "green":
		return fmt.Sprintf("\033[32m%v\033[0m", text)
	case "yellow":
		return fmt.Sprintf("\033[33m%v\033[0m", text)
	case "blue":
		return fmt.Sprintf("\033[34m%v\033[0m", text)
	case "magenta":
		return fmt.Sprintf("\033[35m%v\033[0m", text)
	case "cyan":
		return fmt.Sprintf("\033[36m%v\033[0m", text)
	case "white":
		return fmt.Sprintf("\033[37m%v\033[0m", text)
	default:
		fmt.Println("Unknown color:", color)
		return text
	}
}
