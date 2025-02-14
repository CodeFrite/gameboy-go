package datastructure

type Iterator[T any] interface {
	Iterator() <-chan *T
}

type UpdatableList[T any] interface {
	Push(*T) uint64
	Pop() *T
}

type UpdatableIterator[T any] interface {
	Iterator[T]
	UpdatableList[T]
}
