package datastructure

// interface to iterate over a list of elements
type Iterable[T any] interface {
	GetHead() *Node[T]
}

type UpdatableList[T any] interface {
	Push(*T) uint64
	Pop() *T
}

type UpdatableIterator[T any] interface {
	Iterable[T]
	UpdatableList[T]
}

// Iterator interface implementation
func Iterate[T any](it Iterable[T]) <-chan *T {
	ch := make(chan *T)
	go func() {
		for node := it.GetHead(); node != nil; node = node.GetNext() {
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
func Map[T, U any](it Iterable[T], fn func(*T) *U) <-chan *U {
	ch := make(chan *U)
	go func() {
		for item := range Iterate(it) {
			ch <- fn(item)
		}
		close(ch)
	}()
	return ch
}
