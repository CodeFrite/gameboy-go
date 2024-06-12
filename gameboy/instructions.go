package gameboy

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type GameboyInstructionsMap map[string]InstructionsMap // map of instructions "unprefixed" and "cbprefixed" (512 in total)
type InstructionsMap map[Opcode]Instruction // map of instructions by opcode
type Opcode string // instruction opcode in string format ["0x00"- "0xFF"]
type Instruction struct {
	Mnemonic string 		`json:"mnemonic"`		// instruction mnemonic
	Bytes int 					`json:"bytes"`			// number of bytes the instruction takes
	Cycles []int 				`json:"cycles"`			// number of cycles the instruction takes to execute. The first element is the number of cycles the instruction takes when the condition is met, the second element is the number of cycles the instruction takes when the condition is not met (see RETZ for example)
	Operands []Operand 	`json:"operands"`		// instruction operands used as function arguments
	Immediate bool 			`json:"immediate"`	// is the operand an immediate value or should it be fetched from memory
	Flags Flags 				`json:"flags"`			// cpu flags affected by the instruction
}
type Operand struct {
	Name string `json:"name"`						// operand name: register, n8/n16 (immediate unsigned value), e8 (immediate signed value), a8/a16 (memory location)
	Bytes int `json:"bytes,omitempty"`	// number of bytes the operand takes (optional)
	Immediate bool `json:"immediate"`		// is the operand an immediate value or should it be fetched from memory
}
type Flags struct {
	Z string `json:"Z"`	// Zero flag: set if the result is zero (all bits are 0)
	N string `json:"N"` // Subtract flag: set if the instruction is a subtraction
	H string `json:"H"` // Half carry flag: set if there was a carry from bit 3 (result is 0x0F)
	C string `json:"C"` // Carry flag: set if there was a carry from bit 7 (result is 0xFF)
}

// Load the gameboy instructions set from the JSON file
func loadJSONOpcodeTable() GameboyInstructionsMap {
	// CREDITS: Please note that I am using the https://gbdev.io/gb-opcodes/Opcodes.json file
	// It is a reliable community accepted opcode table for the Gameboy CPU that has been used in many projects and was updated many times
	
	// Open the json file containing the opcode table
	content, err := ioutil.ReadFile("./gameboy/opcodes.json")
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
var instructions GameboyInstructionsMap = loadJSONOpcodeTable()

func GetInstruction(opcode Opcode, prefixed bool) Instruction {
	if !prefixed {
		return instructions["unprefixed"][opcode]
	}
	return instructions["cbprefixed"][opcode]
}