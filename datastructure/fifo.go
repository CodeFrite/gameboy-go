package datastructure

import "math"

// FIFO STRUCT - first-in-first-out data structure:
// - it pushes the new elements at the end of the fifo and returns nil
// - it pops the oldest element located at the head of the fifo and returns it
// - it has a limit capacity and when it is full, it will pop the oldest element and return it
// - when the fifo is empty, the head is nil and the fifo length is 0
// - the fifo keeps track of the number of elements in it to maintain its max capacity

const FIFO_INFINITE_NODE_COUNT uint64 = math.MaxUint64

type Fifo[T any] struct {
	capacity uint64
	count    uint64
	head     *Node[T]
}

func NewFifo[T any](capacity uint64) *Fifo[T] {
	return &Fifo[T]{capacity: capacity}
}

// Push a new node at the end of the fifo (far from the head)
func (f *Fifo[T]) Push(value *T) uint64 {
	// instantiate a new node
	newNode := NewNode(value, nil)

	// locate the last node and the previous one
	if f.count == 0 {
		f.head = newNode
	} else {
		lastNode := f.head
		for lastNode.GetNext() != nil {
			lastNode = lastNode.GetNext()
		}
		// add a new node at the end of the fifo
		lastNode.SetNext(newNode)
	}

	// increment the count
	f.count++

	// check if the fifo is full and pop the head if it is
	if f.count > f.capacity {
		f.head.SetNext(f.head.GetNext())
		f.count--
	}

	// return the number of elements present in the fifo
	return f.count
}

// Pops the oldest node pointed by the head
func (f *Fifo[T]) Pop() *T {
	if f.head == nil {
		return nil
	}
	poppedNode := f.head.GetNext()
	f.head.SetNext(poppedNode.GetNext())
	f.count--
	return poppedNode.GetValue()
}

func (f *Fifo[T]) Peek() *T {
	if f.count == 0 {
		return nil
	} else {
		curr := f.head
		for curr.GetNext() != nil {
			curr = curr.GetNext()
		}
		return curr.GetValue()
	}
}

func (f *Fifo[T]) Length() uint64 {
	return f.count
}

// Iterator interface implementation
func (f *Fifo[T]) Iterator() <-chan *T {
	ch := make(chan *T)
	go func() {
		for node := f.head; node != nil; node = node.GetNext() {
			ch <- node.GetValue()
		}
		close(ch)
	}()
	return ch
}

// FifoMapper:
// - iterates through the elements of the fifo
// - applies a func fn to all *T elements of the fifo
// - and returns the new array of type []U over a channel of *U
// Note: Impossible to achieve via a receiver func since they do not accept types parameter
func FifoMapper[T, U any](f *Fifo[T], fn func(*T) *U) <-chan *U {
	ch := make(chan *U)
	go func() {
		for item := range f.Iterator() {
			ch <- fn(item)
		}
		close(ch)
	}()
	return ch
}
