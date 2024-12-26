package gameboy

import (
	"fmt"
	"testing"
)

const CLOCK_FREQUENCY int = 4_194_304

// DIV register is incremented at 16384Hz (always)
func Test_DIV_Incremented_at_16384Hz(t *testing.T) {
	preconditions()
	timer := NewTimer(bus)
	for i := 0; i < 256; i++ {
		for j := 0; j <= 0xFF; j++ {
			timer.Tick()
		}
		div := bus.Read(REG_FF04_DIV)
		if div != uint8(i+1) {
			t.Errorf("Expected DIV register to increment every 256 T-cycles. Expected DIV to be %d, got %d", i+1, div)
		}
	}
}

// DIV register is reset to 0 when written to from outside the Timer
func Test_DIV_Reset(t *testing.T) {
	preconditions()
	timer := NewTimer(bus)

	// tick the timer a few times up to value 0x42
	for i := 0; i < 0x42; i++ {
		for j := 0; j <= 0xFF; j++ {
			timer.Tick()
		}
	}

	// Expect the DIV register to be 0x42
	div := bus.Read(REG_FF04_DIV)
	if div != 0x42 {
		t.Errorf("Expected DIV register to be 0x42, got %d", div)
	}

	// Write 0xCA to the DIV register and observe the reset
	bus.Write(REG_FF04_DIV, 0xCA)
	div = bus.Read(REG_FF04_DIV)
	if div != 0x00 {
		t.Errorf("Expected DIV register to be reset to 0x00 when written to, got %d", div)
	}
}

// TIMA register increments only if TAC.2 is set
func Test_TIMA_Disabled(t *testing.T) {
	preconditions()
	timer := NewTimer(bus)

	// set TAC to 0x00 (TAC.2 is reset)
	bus.Write(REG_FF07_TAC, 0x00)

	// tick the timer 10 * 4096 * 256 (DIV increments every 256 T-cycles) and observe TIMA not incrementing
	for i := 0; i < 10*4096*256; i++ {
		timer.Tick()
	}

	// expect TIMA to be 0x00
	tima := bus.Read(REG_FF05_TIMA)
	if tima != 0x00 {
		t.Errorf("Expected TIMA to be disabled and still have a value of 0x00, got %d", tima)
	}
}

// TIMA register is incremented at the rate specified by TAC
// We run the clock at 4.194304 MHz and observe the TIMA increment at the rate specified by TAC
func Test_TIMA_Increment_Frequency(t *testing.T) {
	t.Run("TIMA Rate 4096", test_TIMA_Increment_Frequency_4096)
	t.Run("TIMA Rate 262144", test_TIMA_Increment_Frequency_262144)
	t.Run("TIMA Rate 65536", test_TIMA_Increment_Frequency_65536)
	t.Run("TIMA Rate 16384", test_TIMA_Increment_Frequency_16384)
}

// Clock select = 00
func test_TIMA_Increment_Frequency_4096(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 4096 (0x00)
	TARGET_FREQUENCY := 4_096
	TICKS_PER_INCREMENT := CLOCK_FREQUENCY / TARGET_FREQUENCY
	bus.Write(REG_FF07_TAC, 0x04)

	// TIMA should increment every 1024 T-cycles
	expectedTIMAValue := 0
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		timer.Tick()
		if i%TICKS_PER_INCREMENT == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(expectedTIMAValue) {
				t.Fatalf("Expected TIMA to increment every %d T-cycles. Expected TIMA to be %d, got %d", TICKS_PER_INCREMENT, uint8(expectedTIMAValue), tima)
			}
			expectedTIMAValue++
		}
	}
}

// Clock select = 01
func test_TIMA_Increment_Frequency_262144(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 262144 (0x01)
	TARGET_FREQUENCY := 262_144
	TICKS_PER_INCREMENT := CLOCK_FREQUENCY / TARGET_FREQUENCY
	bus.Write(REG_FF07_TAC, 0x05)

	// TIMA should increment every 1024 T-cycles
	expectedTIMAValue := 0
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		timer.Tick()
		if i%TICKS_PER_INCREMENT == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(expectedTIMAValue) {
				t.Fatalf("Expected TIMA to increment every %d T-cycles. Expected TIMA to be %d, got %d", TICKS_PER_INCREMENT, uint8(expectedTIMAValue), tima)
			}
			expectedTIMAValue++
		}
	}
}

// Clock select = 10
func test_TIMA_Increment_Frequency_65536(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 65536 (0x02)
	TARGET_FREQUENCY := 65_536
	TICKS_PER_INCREMENT := CLOCK_FREQUENCY / TARGET_FREQUENCY
	bus.Write(REG_FF07_TAC, 0x06)

	// TIMA should increment every 1024 T-cycles
	expectedTIMAValue := 0
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		timer.Tick()
		if i%TICKS_PER_INCREMENT == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(expectedTIMAValue) {
				t.Fatalf("Expected TIMA to increment every %d T-cycles. Expected TIMA to be %d, got %d", TICKS_PER_INCREMENT, uint8(expectedTIMAValue), tima)
			}
			expectedTIMAValue++
		}
	}
}

// Clock select = 11
func test_TIMA_Increment_Frequency_16384(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 4096 (0x03)
	TARGET_FREQUENCY := 16_384
	TICKS_PER_INCREMENT := CLOCK_FREQUENCY / TARGET_FREQUENCY
	bus.Write(REG_FF07_TAC, 0x07)

	// TIMA should increment every 1024 T-cycles
	expectedTIMAValue := 0
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		timer.Tick()
		if i%TICKS_PER_INCREMENT == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(expectedTIMAValue) {
				t.Fatalf("Expected TIMA to increment every %d T-cycles. Expected TIMA to be %d, got %d", TICKS_PER_INCREMENT, uint8(expectedTIMAValue), tima)
			}
			expectedTIMAValue++
		}
	}
}

// TIMA resets to TMA when it overflows
func Test_TIMA_Overflow_Frequency(t *testing.T) {
	t.Run("Overflow Frequency 2,048Hz (TAC clock frequency = 0x00 (4096) - TMA=0x7F)", test_TIMA_Overflow_Frequency_4096_0x7F)
	t.Run("Overflow Frequency 16,384Hz (TAC clock frequency = 0x01 (262,144) - TMA=0xF0)", test_TIMA_Overflow_Frequency_262144_0xF0)
	t.Run("Overflow Frequency 65,536Hz (TAC clock frequency = 0x02 (65,536) - TMA=0xFF)", test_TIMA_Overflow_Frequency_65536_0xFF)
	t.Run("Overflow Frequency 15,424Hz (TAC clock frequency = 0x03 (16,384) - TMA=0x0F)", test_TIMA_Overflow_Frequency_16384_0x0F)
}

// TAC clock frequency = 0x00 (4,096 Hz) & TMA = 128 => overflow freq = 4096 / ((0xFF-0x7F)+1) = 16 Hz (every 4194304/32 = 131,077 T-cycles)
func test_TIMA_Overflow_Frequency_4096_0x7F(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 4096 (0x00)
	// set TMA & TIMA to 0x7F (TIMA is normally only reset to TMA when it overflows not when TMA changes during runtime)
	BASE_FREQUENCY := 4_096
	TMA_RESET_VALUE := 128
	TARGET_OVERFLOW_FREQUENCY := BASE_FREQUENCY / ((0xFF - TMA_RESET_VALUE) + 1)
	TICKS_PER_OVERFLOW := CLOCK_FREQUENCY / TARGET_OVERFLOW_FREQUENCY
	bus.Write(REG_FF05_TIMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF06_TMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF07_TAC, 0x04)

	// TIMA should increment every 2048 T-cycles
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		if i%TICKS_PER_OVERFLOW == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(TMA_RESET_VALUE) {
				t.Fatalf("Expected TIMA to reset to TMA when it overflows. Expected TIMA to be 0x%02X, got 0x%02X", TMA_RESET_VALUE, tima)
			}
		}
		timer.Tick()
	}
}

// TAC clock frequency = 0x01 (262,144 Hz) & TMA = 0xF0 => overflow freq = 262144 / (0xFF-0xF0+1) = 16384 Hz (every 4194304/16384 = 256 T-cycles)
func test_TIMA_Overflow_Frequency_262144_0xF0(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 262144 (0x01)
	// set TMA & TIMA to 0x7F (TIMA is normally only reset to TMA when it overflows not when TMA changes during runtime)
	BASE_FREQUENCY := 262_144
	TMA_RESET_VALUE := 0xF0
	TARGET_OVERFLOW_FREQUENCY := BASE_FREQUENCY / ((0xFF - TMA_RESET_VALUE) + 1)
	TICKS_PER_OVERFLOW := CLOCK_FREQUENCY / TARGET_OVERFLOW_FREQUENCY
	bus.Write(REG_FF05_TIMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF06_TMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF07_TAC, 0x05)

	// TIMA should increment every 2048 T-cycles
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		if i%TICKS_PER_OVERFLOW == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(TMA_RESET_VALUE) {
				t.Fatalf("Expected TIMA to reset to TMA when it overflows. Expected TIMA to be 0x%02X, got 0x%02X", TMA_RESET_VALUE, tima)
			}
		}
		timer.Tick()
	}
}

// TAC clock frequency = 0x02 (65,536 Hz) & TMA = 0xFF => overflow freq = 65536 / ((0xFF-0xFF)+1) = 65,536 Hz (every 4194304/65536 = 64 T-cycles)
func test_TIMA_Overflow_Frequency_65536_0xFF(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 65536 (0x02)
	// set TMA & TIMA to 0x7F (TIMA is normally only reset to TMA when it overflows not when TMA changes during runtime)
	BASE_FREQUENCY := 65_536
	TMA_RESET_VALUE := 0xFF
	TARGET_OVERFLOW_FREQUENCY := BASE_FREQUENCY / ((0xFF - TMA_RESET_VALUE) + 1)
	TICKS_PER_OVERFLOW := CLOCK_FREQUENCY / TARGET_OVERFLOW_FREQUENCY
	bus.Write(REG_FF05_TIMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF06_TMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF07_TAC, 0x06)

	// TIMA should increment every 2048 T-cycles
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		if i%TICKS_PER_OVERFLOW == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(TMA_RESET_VALUE) {
				t.Fatalf("Expected TIMA to reset to TMA when it overflows. Expected TIMA to be 0x%02X, got 0x%02X", TMA_RESET_VALUE, tima)
			}
		}
		timer.Tick()
	}
}

// TAC clock frequency = 0x03 (16,384 Hz) & TMA = 0xF0 => overflow freq = 16384 / ((0xFF-0xF0)+1) = 1024 Hz (every 4194304/1024 = 4096 T-cycles)
func test_TIMA_Overflow_Frequency_16384_0x0F(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 16384 (0x03)
	// set TMA & TIMA to 0xF0 (TIMA is normally only reset to TMA when it overflows not when TMA changes during runtime)
	BASE_FREQUENCY := 16_384
	TMA_RESET_VALUE := 0xF0
	TARGET_OVERFLOW_FREQUENCY := BASE_FREQUENCY / ((0xFF - TMA_RESET_VALUE) + 1)
	TICKS_PER_OVERFLOW := CLOCK_FREQUENCY / TARGET_OVERFLOW_FREQUENCY
	bus.Write(REG_FF05_TIMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF06_TMA, uint8(TMA_RESET_VALUE))
	bus.Write(REG_FF07_TAC, 0x07)

	fmt.Println("TICKS_PER_OVERFLOW", TICKS_PER_OVERFLOW)

	// TIMA should increment every 2048 T-cycles
	for i := 0; i < CLOCK_FREQUENCY; i++ {
		if i%TICKS_PER_OVERFLOW == 0 {
			tima := bus.Read(REG_FF05_TIMA)
			if tima != uint8(TMA_RESET_VALUE) {
				t.Fatalf("Expected TIMA to reset to TMA when it overflows. Expected TIMA to be 0x%02X, got 0x%02X", TMA_RESET_VALUE, tima)
			}
		}
		timer.Tick()
	}
}

// Timer interrupt occurs when TIMA overflows
func Test_Timer_Interrupt(t *testing.T) {
	// set up
	preconditions()
	timer := NewTimer(bus)

	// enable TIMA (0x04) and set TAC to 4096 (0x00)
	bus.Write(REG_FF05_TIMA, 0x00)
	bus.Write(REG_FF06_TMA, 0x00)
	bus.Write(REG_FF07_TAC, 0x04)

	// make sure the IME, IE (0x02) and IF registers (0x00) are set up correctly
	cpu.ime = true
	bus.Write(IE_REGISTER, 0x02)
	bus.Write(IF_REGISTER, 0x00)

	// TIMA should overflow after 1024 ticks
	for i := 0; i <= 1024*256; i++ {
		if i < 1024*256 {
			if bus.Read(IF_REGISTER)&0x04 == 0x04 {
				t.Fatalf("%d>Expected no timer interrupt to be requested before TIMA overflows. Got interrupt request", i)
			}
		} else if i == 1024*256 {
			if bus.Read(IF_REGISTER)&0x04 != 0x04 {
				t.Fatalf("%d>Expected timer interrupt to be requested after TIMA overflows. Got no interrupt request", i)
			}
		}
		timer.Tick()
	}
}
