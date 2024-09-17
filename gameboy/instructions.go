package gameboy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type GameboyInstructionsMap map[string]InstructionsMap // map of instructions "unprefixed" and "cbprefixed" (512 in total)
type InstructionsMap map[Opcode]Instruction            // map of instructions by opcode
type Opcode string                                     // instruction opcode in string format ["0x00"- "0xFF"]
type Instruction struct {
	Mnemonic  string    `json:"mnemonic"`  // instruction mnemonic
	Bytes     int       `json:"bytes"`     // number of bytes the instruction takes
	Cycles    []uint16  `json:"cycles"`    // number of cycles the instruction takes to execute. The first element is the number of cycles the instruction takes when the condition is met, the second element is the number of cycles the instruction takes when the condition is not met (see RETZ for example)
	Operands  []Operand `json:"operands"`  // instruction operands used as function arguments
	Immediate bool      `json:"immediate"` // is the operand an immediate value or should it be fetched from memory
	Flags     Flags     `json:"flags"`     // cpu flags affected by the instruction
}
type Operand struct {
	Name      string `json:"name"`                // operand name: register, n8/n16 (immediate unsigned value), e8 (immediate signed value), a8/a16 (memory location)
	Bytes     int    `json:"bytes,omitempty"`     // number of bytes the operand takes (optional)
	Immediate bool   `json:"immediate"`           // is the operand an immediate value or should it be fetched from memory
	Increment bool   `json:"increment,omitempty"` // should the program counter be incremented after fetching the operand
	Decrement bool   `json:"decrement,omitempty"` // should the program counter be decreased after fetching the operand
}
type Flags struct {
	Z string `json:"Z"` // Zero flag: set if the result is zero (all bits are 0)
	N string `json:"N"` // Subtract flag: set if the instruction is a subtraction
	H string `json:"H"` // Half carry flag: set if there was a carry from bit 3 (result is 0x0F)
	C string `json:"C"` // Carry flag: set if there was a carry from bit 7 (result is 0xFF)
}

// getBasePath returns the directory of the current file
func getBasePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}
	return filepath.Dir(file), nil
}

// Load the gameboy instructions set from the JSON file
func LoadJSONOpcodeTable() GameboyInstructionsMap {
	// CREDITS: Please note that I am using the https://gbdev.io/gb-opcodes/Opcodes.json file
	// It is a reliable community accepted opcode table for the Gameboy CPU that has been used in many projects and was updated many times

	// Open the json file containing the opcode table but first, get the execution directory
	basePath, err := getBasePath()
	if err != nil {
		panic("Error when getting the base path: " + err.Error())
	}

	content, err := ioutil.ReadFile(basePath + "/opcodes.json")
	if err != nil {
		log.Fatal("Error when opening opcodes.json file: ", err)
	}

	// format the json file content to gameboyInstructionsMap
	var payload GameboyInstructionsMap
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return payload
}

// load the gameboy instructions set from the JSON file
var instructions GameboyInstructionsMap = LoadJSONOpcodeTable()

func GetInstruction(opcode Opcode, prefixed bool) Instruction {
	if !prefixed {
		return instructions["unprefixed"][opcode]
	}
	return instructions["cbprefixed"][opcode]
}
