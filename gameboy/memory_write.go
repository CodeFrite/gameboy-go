package gameboy

import (
	"fmt"
	"strings"
)

/**
 * JSONableSlice is a type alias for a slice of uint8 used to make it JSONable by defining a custom JSON marshalling method.
 */
type JSONableSlice []uint8

/**
* represents a memory write operation used by the mmu to keep track of all memory changes between 2 states
* memory writes can be either reset:
* - before every step from the debugger
* - before a new run from a breakpoint to another breakpoint
 */
type MemoryWrite struct {
	Name    string        `json:"name"`
	Address uint16        `json:"address"`
	Data    JSONableSlice `json:"data"`
}

/**
* MarshalJSON is a custom JSON marshalling method for the JSONableSlice type
* it converts the slice of uint8 to a string of comma-separated values
 */
func (j JSONableSlice) MarshalJSON() ([]byte, error) {
	result := strings.Join(strings.Fields(fmt.Sprintf("%d", j)), ",")
	return []byte(result), nil
}
