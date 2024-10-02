package gameboy

import "fmt"

type JoypadState struct {
	Up, Down, Left, Right, A, B, Start, Select bool
}

type Joypad struct {
	eventchannel <-chan *JoypadState
}

func NewJoypad(eventchannel <-chan *JoypadState) *Joypad {
	j := &Joypad{
		eventchannel: eventchannel,
	}
	j.initInterruptHandler()
	return j
}

func (j *Joypad) initInterruptHandler() {
	go func() {
		for {
			select {
			case event := <-j.eventchannel:
				j.handleEvent(event)
			}
		}
	}()
}

func (j *Joypad) handleEvent(event *JoypadState) {
	fmt.Println("Joypad event received: ", event)
}
