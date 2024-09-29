package gameboy

import (
	"fmt"
	"time"
)

type Synchronizable interface {
	onTick()
}

// timer registers
var TIMER_REGISTERS map[string]uint16 = map[string]uint16{
	"DIV":  0xFF04, // divider register: incremented at a rate of 16384 Hz
	"TIMA": 0xFF05, // timer counter: incremented at a rate of 16384 Hz
	"TMA":  0xFF06, // timer modulo: when TIMA overflows, it is reset to TMA
	"TAC":  0xFF07, // timer control register: used to start/stop the timer
}

// generates a tick at a given frequency
type Timer struct {
	// state
	Frequency uint32 `json:"frequency"` // freq in Hz

	// communication
	DoneChan    chan bool        `json:"channel"`     // boolean channel used to stop the timer
	TickChan    <-chan time.Time `json:"tickChannel"` // boolean channel used to signal a tick
	Subscribers []Synchronizable `json:"subscribers"` // subscribers to the timer
}

func NewTimer(frequency uint32) *Timer {
	if frequency == 0 {
		frequency = 1
		fmt.Println("Timer> invalid frequency, reverting to default value of 1000 ms (1 Hz)")
	}
	return &Timer{
		Frequency:   frequency,
		DoneChan:    make(chan bool),
		TickChan:    make(<-chan time.Time),
		Subscribers: make([]Synchronizable, 0),
	}
}

// add a subscriber
func (t *Timer) Subscribe(subscriber Synchronizable) {
	t.Subscribers = append(t.Subscribers, subscriber)
}

// start the timer
func (t *Timer) Start() chan bool {
	// compute and convert the period to a time.Duration
	period := 1000000000 / float64(t.Frequency)
	tickRate := time.Duration(period) * time.Nanosecond
	t.TickChan = time.NewTicker(tickRate).C

	go func() {
	loop:
		for {
			select {
			case <-t.DoneChan:
				break loop
			case <-t.TickChan:
				t.Tick()
			}
		}
	}()

	// mark the timer as running
	return t.DoneChan
}

// stop the timer
func (t *Timer) Stop() {
	// send a signal to the timer to stop
	t.DoneChan <- true
}

// on tick, increment the count and notify all subscribers
func (t *Timer) Tick() {
	for _, subscriber := range t.Subscribers {
		go subscriber.onTick()
	}

}
