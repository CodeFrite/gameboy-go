package datastructure

// NODE STRUCT
// Node is a data structure that holds a value of any type (alias for interface{}) and a pointer to the next node
// The value and next node are private and can be accessed using the getter/setter methods
type Node[T any] struct {
	value *T
	next  *Node[T]
}

func NewNode[T any](value *T, next *Node[T]) *Node[T] {
	return &Node[T]{value: value, next: next}
}

func (n *Node[T]) GetValue() *T {
	return n.value
}

func (n *Node[T]) SetValue(value *T) {
	n.value = value
}

func (n *Node[T]) GetNext() *Node[T] {
	return n.next
}

func (n *Node[T]) SetNext(next *Node[T]) {
	n.next = next
}
