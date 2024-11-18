package gameboy

/* NODE STRUCT */

type node[T any] struct {
	value *T
	next  *node[T]
}

func newNode[T any](value *T, next *node[T]) *node[T] {
	return &node[T]{value: value, next: next}
}

func (n *node[T]) getValue() *T {
	return n.value
}

func (n *node[T]) getNext() *node[T] {
	return n.next
}

func (n *node[T]) setNext(next *node[T]) {
	n.next = next
}

/* FIFO STRUCT */
const INFINITE_MAX_NODE_COUNT uint64 = 2 ^ 64 - 1

// fifo is a first-in-first-out data structure with a head pointer
// it has a limit capacity and when it is full, it will remove the oldest element
// when the fifo is empty, the head is nil and the fifo length is 0
// the most recent element is always at the begining of the fifo and pointed by the head
// the oldest element points to nothing
type fifo[T any] struct {
	maxNodeCount uint64
	head         *node[T]
}

func newFifo[T any](maxNodeCount uint64) *fifo[T] {
	return &fifo[T]{head: nil, maxNodeCount: maxNodeCount}
}

func (f *fifo[T]) push(value *T) {
	// the new node become the one the head is pointing to
	saveHead := f.head
	f.head = newNode(value, saveHead)
	if uint64(f.len()) > f.maxNodeCount {
		// remove the oldest element
		f.pop()
	}
}

func (f *fifo[T]) pop() *T {
	if f.head == nil {
		return nil
	}

	value := f.head.getValue()
	f.head = f.head.getNext()
	return value
}

func (f *fifo[T]) peek() *T {
	if f.head == nil {
		return nil
	}
	return f.head.getValue()
}

func (f *fifo[T]) len() int {
	count := 0
	current := f.head
	for current != nil {
		count++
		current = current.next
	}
	return count
}
