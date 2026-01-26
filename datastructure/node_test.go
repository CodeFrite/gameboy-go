package datastructure

import "testing"

func TestNewNode(t *testing.T) {
	t.Run("New node with int value and next node nil", test_NewNode_int_next_nil)
	t.Run("New node with string value and next node string", test_NewNode_string_next_string)
}

func test_NewNode_int_next_nil(t *testing.T) {
	t.Log("TestNewNode: NewNode should return a new node with the given value and next node")

	x := 1
	newNode := NewNode[int](&x, nil)

	// Check if the value is set correctly
	if *newNode.GetValue() != x {
		t.Errorf("Expected new node value to be %d, got %d", x, *newNode.GetValue())
	}

	// Check if the next node is set correctly
	if newNode.GetNext() != nil {
		t.Errorf("Expected next node to be nil, got %v", newNode.GetNext())
	}
}

func test_NewNode_string_next_string(t *testing.T) {
	t.Log("TestNewNode: NewNode should return a new node with the given value and next node")

	valueNext := "hello"
	nextNode := NewNode[string](&valueNext, nil)
	valueNew := "world"
	newNode := NewNode[string](&valueNew, nextNode)

	// Check if the value is set correctly
	if *newNode.GetValue() != valueNew {
		t.Errorf("Expected new node value to be %s, got %s", valueNew, *newNode.GetValue())
	}

	// Check if the next node is set correctly
	if newNode.GetNext() != nextNode {
		t.Errorf("Expected next node to be %p, got %p", nextNode, newNode.GetNext())
	}

	// Check if the next node value is set correctly
	if newNode.GetNext().GetValue() != &valueNext {
		t.Errorf("Expected next node value to be %v, got %v", &valueNext, newNode.GetNext().GetValue())
	}
}

func TestSetValue(t *testing.T) {
	oldValue := "old value"
	newNode := NewNode[string](&oldValue, nil)
	newValue := "new value"
	newNode.SetValue(&newValue)

	// Check if the value has been correctly set
	if *newNode.GetValue() != newValue {
		t.Errorf("Expected new node value to be %s, got %s", newValue, *newNode.GetValue())
	}
}

func TestSetNext(t *testing.T) {
	nextNodeValue := "next node value"
	nextNode := NewNode[string](&nextNodeValue, nil)

	updatedNextNodeValue := "updated next node value"
	updatedNextNode := NewNode[string](&updatedNextNodeValue, nil)

	newNodeValue := "new node value"
	newNode := NewNode[string](&newNodeValue, nextNode)
	newNode.SetNext(updatedNextNode)

	// Check if the value has been correctly set
	if *newNode.GetNext() != *updatedNextNode {
		t.Errorf("Expected updated next node value to be %p, got %p", updatedNextNode, *newNode.GetNext())
	}
}
