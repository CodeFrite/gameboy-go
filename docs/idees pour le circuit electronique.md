```go
package main

import (
  "fmt"
  "sync"
  "time"
)

// State structure to hold values and provide notification
type State struct {
  mu sync.RWMutex
  counter int
  notify chan struct{}
}

// NewState initializes the state with an initial value
func NewState(initialValue int) *State {
  return &State{
    counter: initialValue,
    notify: make(chan struct{}),
  }
}

// Set updates the counter and notifies listeners
func (s *State) Set(newValue int) {
  s.mu.Lock()
  defer s.mu.Unlock()
  if s.counter != newValue {
    s.counter = newValue
    close(s.notify) // Notify that a change has occurred <--------------------------------
    s.notify = make(chan struct{}) // Reset notify channel
  }
}

// Get retrieves the current value
func (s *State) Get() int {
  s.mu.RLock()
  defer s.mu.RUnlock()
  return s.counter
}

// Effect function to mimic useEffect behavior
func (s *State) Effect(callback func(int)) {
  previousValue := s.Get()
  for {
    <-s.notify // Wait for notification
    currentValue := s.Get()
    if currentValue != previousValue {
      callback(currentValue) // Call the callback with the new value
      previousValue = currentValue // Update the previous value
    }
  }
}

func main() {
// Initialize state
state := NewState(0)

    // Start effect watching the state in a separate goroutine
    go state.Effect(func(newValue int) {
    	fmt.Printf("State changed: %d\n", newValue)
    })

    // Simulate changing the state
    for i := 1; i <= 5; i++ {
    	time.Sleep(1 * time.Second) // Simulate some work
    	state.Set(i)                // Update the state
    }

    // Prevent the main function from exiting immediately
    time.Sleep(5 * time.Second)

}
```
