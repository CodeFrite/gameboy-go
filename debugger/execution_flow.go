package debugger

import (
	ds "github.com/codefrite/gameboy-go/datastructure"
)

// The execution flow struct is part of the debugger package.
// It is responsible for holding a subset of the state of any group of components of the gameboy during execution.
// The execution flow struct uses a user-defined callback func to record the state of the components.

type StateRecorder[T any] func(<-chan *T)

type ExecutionFlow[T any] struct {
	stateQueue     ds.UpdatableIterator[T]
	stateRecorder  StateRecorder[T]
	maxStatesCount uint64
}

// Create a new execution flow struct
func NewExecutionFlow[T any](stateRecorder StateRecorder[T], maxStatesCount uint64) *ExecutionFlow[T] {
	return &ExecutionFlow[T]{stateQueue: ds.NewFifo[T](STATE_QUEUE_MAX_LENGTH), stateRecorder: stateRecorder, maxStatesCount: maxStatesCount}
}

// Record the state of the components and push it to the state queue
func (ef *ExecutionFlow[T]) Record() {
	state := ef.stateRecorder(cpuState, ppuState, apuState)
	ef.stateQueue.Push(state)
}

// Get the next state
func (ef *ExecutionFlow[T]) Next() <-chan *T {
	// TODO: finish implementing the iterator
	return ef.stateQueue.Iterator()
}

/*
func (ef *ExecutionFlow[T]) Map(fn func(*T)interface{}) ds.UpdatableList[interface{}] {
	var newList ds.UpdatableList[T];

	for t := range(ef.stateQueue.Iterator()) {
		newList.Push(fn(t))
	}

	return newList
}
*/
