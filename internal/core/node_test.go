package core

import "testing"

func TestNewNode(t *testing.T) {
	pos := NewVector(1.0, 2.0, 3.0)
	node := NewNode("test", pos)

	if node.ID != "test" {
		t.Errorf("expected ID 'test', got '%s'", node.ID)
	}

	if node.Position.Dim() != 3 {
		t.Errorf("expected dimension 3, got %d", node.Position.Dim())
	}

	if node.Function != 0.0 {
		t.Errorf("expected initial function value 0.0, got %f", node.Function)
	}
}

func TestNodeSetFunction(t *testing.T) {
	node := NewNode("test", NewVector(0.0, 0.0))
	node.SetFunction(5.5)

	if node.Function != 5.5 {
		t.Errorf("expected function value 5.5, got %f", node.Function)
	}
}

func TestNodeDistanceTo(t *testing.T) {
	node1 := NewNode("A", NewVector(0.0, 0.0))
	node2 := NewNode("B", NewVector(3.0, 4.0))
	dist := node1.DistanceTo(node2)

	expected := 5.0
	if dist != expected {
		t.Errorf("expected distance %f, got %f", expected, dist)
	}
}
