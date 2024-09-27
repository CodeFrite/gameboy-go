# Comparison

## Solution 1: Struct based approach

```go
type MemoryRegister8 struct {
	uint8
}

func (r8 MemoryRegister8) Get() uint8 {
	return r8.uint8
}

func (r8 MemoryRegister8) Set(value uint8) {
	r8.uint8 = value
}

func (r8 MemoryRegister8) GetBit(bit uint8) bool {
	if bit > 7 {
		panic(fmt.Sprintf("MemoryRegister8> getBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	return (r8.uint8 & op) == op
}

func (r8 MemoryRegister8) SetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("MemoryRegister8> setBit: bit out of range: %v", bit))
	}
	op := uint8(1 << bit)
	r8.uint8 |= op
}

func (r8 MemoryRegister8) ResetBit(bit uint8) {
	if bit > 7 {
		panic(fmt.Sprintf("MemoryRegister8> resetBit: bit out of range: %v", bit))
	}
	r8.uint8 ^= 1 << bit
}
```

And here is how it can be used:

```go
var r8 MemoryRegister8
r8.Set(0x00)
fmt.Printf("Register value: %v\n", r8.Get()) // Register value: 0
```

## Solution 2:
