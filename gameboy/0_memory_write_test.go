package gameboy

import (
	"encoding/json"
	"testing"
)

func TestMarshalAndUnmarshalMemoryWrite(t *testing.T) {
	p := MemoryWrite{
		Name:    "ROM",
		Address: 0x01,
		Data:    JSONableSlice{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
	}

	// marshal
	bytes, error := json.Marshal(p)
	if error != nil {
		t.Error("ERROR WHILE MARSHALING")
		return
	}
	t.Log("json marshalled", bytes)

	// unmarshal
	var mw MemoryWrite
	error = json.Unmarshal(bytes, &mw)
	if error != nil {
		t.Error("ERROR WHILE UNMARSHALING", error)
		return
	}
	t.Log("json unmarshalled", &mw)
}

func TestMarshalAndUnmarshalMemoryWriteArray(t *testing.T) {
	p := []MemoryWrite{
		{
			Name:    "ROM0",
			Address: 0x00,
			Data:    JSONableSlice{0x00, 0x01, 0x02, 0x03},
		},
		{
			Name:    "ROM1",
			Address: 0x11,
			Data:    JSONableSlice{0x04, 0x05, 0x06, 0x07},
		},
		{
			Name:    "ROM2",
			Address: 0x22,
			Data:    JSONableSlice{0x08, 0x09, 0x0A, 0x0B},
		},
		{
			Name:    "ROM3",
			Address: 0x33,
			Data:    JSONableSlice{0x0C, 0x0D, 0x0E, 0x0F},
		},
		{
			Name:    "Interrupt Enable Register",
			Address: 0,
			Data:    JSONableSlice{1},
		},
		{
			Name:    "Interrupt Enable Register",
			Address: 0,
			Data:    JSONableSlice{1},
		},
		{
			Name:    "Interrupt Enable Register",
			Address: 0,
			Data:    JSONableSlice{1},
		},
	}

	// marshal
	bytes, error := json.Marshal(p)
	if error != nil {
		t.Error("ERROR WHILE MARSHALING")
		return
	}
	t.Log("json marshalled", bytes)

	// unmarshal
	var mw []MemoryWrite
	error = json.Unmarshal(bytes, &mw)
	if error != nil {
		t.Error("ERROR WHILE UNMARSHALING", error)
		return
	}
	t.Log("json unmarshalled", &mw)
}
