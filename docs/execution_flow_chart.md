# Execution Flow Chart

More then often during the development of our emulator, we are stuck in a situation where we are not able to understand why our emulator is not working as expected. The issue can come from multiple sources:

- some instructions are not implemented correctly
- cpu interrupts are not triggered or handled correctly
- the ppu rendering is not working as expected
- we are trying to access memory location that is not supposed to be accessed
- an issue with the orchestration of the different components of the emulator
- etc...

Due to the nature of the system under development, it is not always possible to simply add debug logs to understand what is happening. Indeed, the cpu & ppu are driven @4.194304 MHz and do output a lot of data (cpu registers, memory changes, ppu image and registers, etc.) that are not easy to follow in a simple log file.

One way to approach the issue is to develop and use a Debugger. This allows us to pause the execution of the emulator at any time, inspect the state of the different components, and step through the execution of the different instructions.

In other situations, we only want to know if the cpu is stuck in a loop, visiting the same pc addresses over and over again. This is where the Execution Flow Chart comes in handy.

## Mermaid Flow Chart

Instead of creating a flow chart viewer in go, we will make use of `mermaid`. `Mermaid` is a `JavaScript based diagramming and charting tool that renders Markdown-inspired text definitions to create and modify diagrams dynamically`. It offers various types of diagrams, such as flow charts, sequence diagrams, class diagrams, etc. We will use the flow chart diagram to represent the execution flow of the cpu.

The approach will be to:

- record the pc addresses visited by the cpu during the execution of the emulator
- generate the corresponding flow chart using `mermaid` syntax, identifying the graph nodes and edges
- integrate the mermaid flow chart into a log file in markdown format

Let's first begin by sketching an example flow chart diagram that we want to generate. This will help us determine how to record the pc addresses visited by the cpu.

## Desired Execution Flow Chart Look & Feel

The graph below shows a simple flow chart diagram representing the execution flow of the cpu as we are going to implement it.

```mermaid
---
title: Node
---
flowchart TD
  subgraph bootrom
    A([0x0000])
    B([0x0001])
    C([0x0005])
    D([0x0007])
    E([0x0009])
    F([0x000A])
    A -->|1| B
    B -->|123| C
    C -->|122| B
    C -->|15| D
    D -->|1| E
    E -->|9| F
    E -->|2| B
  end
  subgraph gamerom
    F -->|1| G([0x0100])
  end
```

It has named nodes which contain the pc address visited by the cpu. The edges are labeled with the number of times the cpu transitioned from one pc address to another. It also contains subgraphs to represent the different memory regions (bootrom, gamerom, etc...), which will be considered as a Nice to Have.

## Implementation

Let's talk about how to implement the recording and the generation of the flow chart.

### Recording the Execution Flow

Now comes the question of how to record the pc addresses visited by the cpu. We will use a simple approach:

- use a `[]uint16` to record the unique pc addresses visited by the cpu
- use a `map[int]uint16` to record the order in which the pc addresses were visited
- use a `map[uint16]map[uint16]int` to record the number of times a transition from one pc address to another has been made
- use some kind of struct to list the different memory regions that we want to emphasize

### Generating the Flow Chart

Now that we have recorded, we can generate the flow chart using the `mermaid` syntax. There are basically three steps to follow:

- adding a reference to all the nodes visited at the beginning of the flow chart definition using the `[]uint16` slice
- adding the edges to represent transitions between one pc address to another using the `map[int]uint16` map
- adding a label to the edges representing the number of times the transition has been exercised using the `map[uint16]map[uint16]int` map
- adding the subgraphs to represent the different memory regions using some user defined struct (Nice to Have)

### Creating a new md log file

Once the execution has ended and the data has been recorded, we can create a new md file named `execution_flow_chart_${date}_${time}.md` and save it to the `logs` directory. This file will have the following structure:

```markdown
# Execution Flow Chart

- game: ${game}
- date: ${date}
- time: ${time}

`Insert the mermaid flow chart here`
```

### Where do we put the code

The first question that we should ask ourselves is where to put the code. Here are our options along with the pros and cons:

| #   | Location    | Pros                                                                                    | Cons                                                                                                          |
| --- | ----------- | --------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| 1   | cpu.go      | Easy access to cpu.pc                                                                   | Shouldn't contain the logic to generate the flow chart                                                        |
| 2   | gameboy.go  | Easy access to cpu.pc                                                                   | Cannot be used in debug mode I guess                                                                          |
|     |             | Can be extended to take care of execution flows for other components (memory, ppu, ...) |                                                                                                               |
| 3   | debugger.go | Does not slow the core gameboy-go functionnalities (images, sounds & joypad generation) | Cannot be used in play mode (where we only care about images & sounds produces by the gameboy)                |
| 4   | server.go   | Separates the core gameboy implementation with other tools                              | Introduces complexity in server which should normally only transmit messages in both ways (client <-> server) |

The most logical place to put the code is in the `debugger`. At this point I think I should refactor my code and make the `debugger` a separate package. This will allow me to add more functionalities to the debugger without polluting the core gameboy implementation.

Now that I have moved the `debugger` package to its own directory, my project structure looks like this:

```plaintext
+ gameboy-go          (module)
  + datastructure　   (package)
  + debugger          (package)
    - debugger.go
  + gameboy           (package)
    - gameboy.go
    - cpu.go
    - memory.go
    - ...
  - go.mod
  - go.sum
```

It has only one go module named `gameboy-go` which contains three packages: `gameboy`, `debugger` & `datastructure`. The whole module is versioned in the same git repository.

To have access to the gameboy package from the debugger package, I have to import it in the `debugger.go` file:

```go
package debugger

import "github.com/codefrite/gameboy-go/gameboy"

const STATE_QUEUE_MAX_LENGTH = 100

// debugger struct: combination of a gameboy, its internal state and a list of breakpoints set by the user
type Debugger struct {
	// state
	gameboy     *gameboy.Gameboy
  ...
```

### Data Structures

Now that I have split my code correctly, I can define a few data structures that will help me record any type of execution flow state that I want to monitor. These structures should allow me to:

- `Record any kind of state data`: visited pc address, changed memory addresses & values, whole cpu state, etc... or any combination of these
- Once the data has been recorded, I should be able to `iterate` over it
- During the iteration, I should be able to `apply some transformation functions` to the data, for example to aggregate memory changes, filter out some data, or even generate the flow chart as a `mermaid` diagram

#### Node

To achieve these objectives, I first defined a data container flexible enough to store any kind of data: a `Node` struct:

```go
type Node[T any] struct {
	value *T
	next  *Node[T]
}
```

Nothing incredible here, just a node with a type parameter `T` that can store any kind of data. Along with this structure, I defined a `constructor` and a few `getter` and `setter` functions:

```go
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
```

#### FIFO

Next, I defined a simple linked list data structure, or more exactly a `FIFO` (First In First Out) data structure. This data structure will allow me to record the states in a chronological order and iterate over them:

```go
type Fifo[T any] struct {
	capacity uint64
	count    uint64
	head     *Node[T]
}
```

As we can see, it has a `capacity` field to limit the number of elements that can be stored in the list, a `count` field to keep track of the number of elements stored in the list, and a `head` field to point to the first element of the list. Along with this structure, I defined a `constructor` and the traditional `Push`, `Pop` & `Peek` functions:

```go
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
```

It is important to note here that the oldest element is always pointed by the `head` field, while the newest one is `pushed` at the end of the list and has no `next` node, i.e. `next` is `nil`. To illustrate this, let's consider the example where the fifo has a capacity of 3 and we push 5 elements into it (0, 1, 2, 3, 4) and then pop 3 elements from it:

```mermaid
---
title: FIFO with capacity 3 - adding 5 elements and popping 3
---
flowchart LR
  H1[head] --> N1[nil]
  H2[head] --> E0[0] --> N2[nil]
  H3[head] --> E1[0] --> E2[1] --> N3[nil]
  H4[head] --> E3[0] --> E4[1] --> E5[2] --> N4[nil]
  H5[head] --> E6[1] --> E7[2] --> E8[3] --> N5[nil]
  H6[head] --> E9[2] --> E10[3] --> E11[4] --> N6[nil]
  H7[head] --> E12[3] --> E13[4] --> E14[5] --> N7[nil]
  H8[head] --> E15[4] --> E16[5] --> N8[nil]
  H9[head] --> E17[5] --> N9[nil]
  H10[head] --> N10[nil]
```

When an element is pushed to the fifo or popped from it, the `count` field is updated accordingly. A getter allows to see how many elements are present in the fifo at any time:

```go
func (f *Fifo[T]) Length() uint64 {
	return f.count
}
```

Finally, we reach the interesting part: how can we iterate over the fifo leveraging go language specificities, and how can we apply some transformation functions to the data stored in the fifo in a clever way?

This is done with a combination of interfaces, channels & types parameters. Before going any further, let's define a getter to the head of the fifo. This will come in handy when we want to iterate over the fifo:

```go
func (f *Fifo[T]) GetHead() *Node[T] {
  return f.head
}
```

### Iterable Interface

The first interface that we need to define is the `Iterable` interface. This interface will allow us to iterate over any data structure that exposes an accessor to its first element, here named `head`:

```go
type Iterable[T any] interface {
	GetHead() *Node[T]
}
```

It doesn't see much but combined with the following function, we can iterate over any data structure that uses a `Node[T]` as its data container and has a `GetHead` function:

```go
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
```

We can now reuse this function to iterate over the fifo and apply some transformation functions to the data stored in it. This is done by defining a `Map` interface:

```go
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
```

In the future, if we decide to define other data structures, we will be able to reuse these functions and extend them to support filtering, reducing, etc ... as long as the data structure implements the `Iterable` interface, which is easy and powerful.

Now we are finally set to record the execution flow of the cpu and generate the flow chart using `mermaid` syntax!
