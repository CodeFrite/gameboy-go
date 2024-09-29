package gameboy

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

/*
 *	Simulate a web server using the debugger and waiting for gameboy state updates on a channel
 */

func TestCpuStateChannel(t *testing.T) {
	// create a state channel
	cpu_state_channel := make(chan *CpuState)

	// create a new gameboy
	//gb := NewGameboy(cpu_state_channel)

	// stepping the gameboy 10 times
	i := 1
	//go gb.Step()

	// send cpu state
	done_channel := make(chan bool)
	go func() {
	cpu_state_loop:
		for {
			select {
			case <-done_channel:
				close(done_channel)
				break cpu_state_loop
			default:
				sleepTime := rand.IntN(3000)
				time.Sleep(time.Duration(sleepTime) * time.Millisecond)
				cpu_state_channel <- &CpuState{
					PC: uint16(i),
				}
			}
		}
	}()

	// receive cpu state

cpu_receive_loop:
	for {
		i++
		cpu_state := <-cpu_state_channel
		t.Logf("	<- received cpu State: %v\n", cpu_state)

		if i <= 10 {
			//go gb.Step()
			fmt.Println("	stepping ...")
		} else if i > 10 {
			close(cpu_state_channel)
			break cpu_receive_loop
		}
	}
	t.Logf("TestCpuStateChannel done ... steps=%v\n", i)
}
