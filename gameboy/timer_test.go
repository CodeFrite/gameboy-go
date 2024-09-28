package gameboy

import (
	"fmt"
	"testing"
	"time"
)

func TestInstantiation(t *testing.T) {
	timer := NewTimer(1)
	if timer.Frequency != 1 {
		t.Errorf("Expected frequency to be 1, got %v", timer.Frequency)
	}
	// on instantiation, the subscribers list should be empty
	if len(timer.Subscribers) != 0 {
		t.Errorf("Expected subscribers list to be empty, got %v", len(timer.Subscribers))
	}
	// on instantiation, the done channel should be nil
	if timer.DoneChan == nil {
		t.Errorf("Expected done channel to be non-nil, got nil")
	}
	// on instantiation, the tick channel should be nil
	if timer.TickChan != nil {
		t.Errorf("Expected tick channel to be nil, got %v", timer.TickChan)
	}
}

// / helper func that implements the Synchronizable interface

var tickCount int

type synch struct {
}

func NewSynch() *synch {
	return &synch{}
}

func (s *synch) onTick() {
	tickCount++
}

func TestSubscription(t *testing.T) {

	timer := NewTimer(1)
	subscriber := NewSynch()
	timer.Subscribe(subscriber)
	if len(timer.Subscribers) != 1 {
		t.Errorf("Expected subscribers list to have 1 subscriber, got %v", len(timer.Subscribers))
	}

	// add another subscriber
	subscriber2 := NewSynch()
	timer.Subscribe(subscriber2)
	if len(timer.Subscribers) != 2 {
		t.Errorf("Expected subscribers list to have 2 subscribers, got %v", len(timer.Subscribers))
	}
}

func TestTick(t *testing.T) {

	timer := NewTimer(1)
	subscriber := NewSynch()
	timer.Subscribe(subscriber)

	tickCount = 0
	for i := 1; i < 11; i++ {
		timer.Tick()
		if tickCount != i {
			t.Errorf("Expected to have received 1 tick, got %v", tickCount)
		}
	}
}

func TestStartStop(t *testing.T) {
	timer := NewTimer(1)
	subscriber := NewSynch()
	timer.Subscribe(subscriber)

	tickCount = 0
	fmt.Println("Starting timer ...")
	doneChan := timer.Start()
	if timer.TickChan == nil {
		t.Errorf("Expected tick channel to be non-nil, got nil")
	}
	fmt.Println("Waiting for 10 seconds ...")
	time.Sleep(10 * time.Second)
	fmt.Println("Stopping timer ... tick count is", tickCount)
	timer.Stop()
	fmt.Println("Closing done channel ...")
	close(doneChan)
	fmt.Println("Exiting ... final tick count is", tickCount)
}