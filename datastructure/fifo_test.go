package datastructure

import (
	"fmt"
	"math"
	"testing"
)

func TestNewFifo(t *testing.T) {
	t.Log("TestNewFifo")

	capacity := uint64(10)
	fifo := NewFifo[int](capacity)
	if fifo.capacity != capacity {
		t.Errorf("Expected max node count to be %v, got %v", capacity, fifo.capacity)
	}

	if fifo.head != nil {
		t.Errorf("Expected head to be nil, got %v", fifo.head)
	}

	if fifo.Length() != 0 {
		t.Errorf("Expected fifo length to be 0, got %v", fifo.Length())
	}

}

func TestPush(t *testing.T) {
	t.Log("TestPush")

	capacity := uint64(3)
	fifo := NewFifo[int](capacity)

	// test data
	testData := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// In a loop, push a new element to the fifo and check:
	// - if the head is the new element
	// - if the length is the index + 1 with the capacity as the limit
	for idx, v := range testData {
		fifo.Push(&v)

		// Check the value of the last pushed element
		if *fifo.Peek() != v {
			t.Errorf("Expected fifo head to be %v, got %v", v, *fifo.Peek())
		}

		// Check if the length is idx + 1 with the capacity as the limit
		expectedLength := uint64(math.Min(float64(idx+1), float64(capacity)))
		if fifo.Length() != expectedLength {
			t.Errorf("Expected fifo length to be %v, got %v", expectedLength, fifo.Length())
		}
	}
}

// Test the iterator interface implementation
// Iterator() func returns a channel of *T that sends the fifo values and closes the channel when the fifo is empty
// Therefore, it can be used in a for loop to iterate over the fifo values
func TestIterator(t *testing.T) {
	t.Log("TestIterator")

	// stuff the fifo with some values
	testData := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	capacity := uint64(10)
	fifo := NewFifo[int](capacity)
	for _, v := range testData {
		fifo.Push(&v)
	}

	// loop over the fifo values
	for v := range Iterate[int](fifo) {
		fmt.Println("... v", *v)
	}
}

// Test the iterator map func which apply a func to every single values
// stored in the fifo and returns sends them on a channel
func TestIteratorMap(t *testing.T) {
	t.Log("TestIteratorMapper")

	// define the func to use when mapping fifo content
	fn := func(v *int) *int {
		var val int = *v * 2
		return &val
	}

	// define test data
	testData := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	capacity := uint64(100)

	// instantiate a new fifo
	fifo := NewFifo[int](capacity)

	// fill in the test data into the fifo
	for v := range testData {
		fifo.Push(&v)
	}

	// expected results
	expectedResults := []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	iteratorMapper := Map[int, int](fifo, fn)

	i := 0
	for v := range iteratorMapper {
		fmt.Println("... v", *v)
		if *v != expectedResults[i] {
			t.Errorf("Expected fn(<-fifoChan) to be %v, got %v", expectedResults[i], *v)
		}
		i++
	}

}
