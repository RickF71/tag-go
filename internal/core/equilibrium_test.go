package core

import (
	"math"
	"testing"
)

func TestNewSystem(t *testing.T) {
	system := NewSystem()

	if system == nil {
		t.Fatal("expected non-nil system")
	}

	if len(system.Nodes) != 0 {
		t.Errorf("expected empty nodes, got %d nodes", len(system.Nodes))
	}

	if len(system.Constraints) != 0 {
		t.Errorf("expected empty constraints, got %d constraints", len(system.Constraints))
	}
}

func TestSystemAddNode(t *testing.T) {
	system := NewSystem()
	node := NewNode("A", NewVector(1.0, 2.0))

	system.AddNode(node)

	if len(system.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(system.Nodes))
	}

	if system.Nodes[0].ID != "A" {
		t.Errorf("expected node ID 'A', got '%s'", system.Nodes[0].ID)
	}
}

func TestSystemAddConstraint(t *testing.T) {
	system := NewSystem()
	constraint := Constraint{Type: "test", Value: 1.0}

	system.AddConstraint(constraint)

	if len(system.Constraints) != 1 {
		t.Errorf("expected 1 constraint, got %d", len(system.Constraints))
	}

	if system.Constraints[0].Type != "test" {
		t.Errorf("expected constraint type 'test', got '%s'", system.Constraints[0].Type)
	}
}

func TestSystemGetDimensionality(t *testing.T) {
	system := NewSystem()

	// Empty system
	if dim := system.GetDimensionality(); dim != 0 {
		t.Errorf("expected dimensionality 0, got %d", dim)
	}

	// Add 3D node
	node := NewNode("A", NewVector(1.0, 2.0, 3.0))
	system.AddNode(node)

	if dim := system.GetDimensionality(); dim != 3 {
		t.Errorf("expected dimensionality 3, got %d", dim)
	}
}

func TestSystemComputeEquilibrium(t *testing.T) {
	system := NewSystem()

	// Empty system
	eq := system.ComputeEquilibrium()
	if eq != 0.0 {
		t.Errorf("expected equilibrium 0.0 for empty system, got %f", eq)
	}

	// Add nodes with functions
	node1 := NewNode("A", NewVector(0.0, 0.0))
	node1.SetFunction(1.0)
	node2 := NewNode("B", NewVector(1.0, 0.0))
	node2.SetFunction(0.8)

	system.AddNode(node1)
	system.AddNode(node2)

	// Add constraint
	constraint := Constraint{Type: "boundary", Value: 0.85}
	system.AddConstraint(constraint)

	eq = system.ComputeEquilibrium()
	if eq < 0 || eq > 1 {
		t.Errorf("expected equilibrium in range [0,1], got %f", eq)
	}

	// Verify equilibrium is stored
	if system.Equilibrium != eq {
		t.Errorf("equilibrium not stored correctly: expected %f, got %f", eq, system.Equilibrium)
	}
}

func TestSystemComputeCoherence(t *testing.T) {
	system := NewSystem()

	// System with < 2 nodes
	if coherence := system.ComputeCoherence(); coherence != 1.0 {
		t.Errorf("expected coherence 1.0 for system with < 2 nodes, got %f", coherence)
	}

	// Add nodes
	node1 := NewNode("A", NewVector(0.0, 0.0))
	node2 := NewNode("B", NewVector(1.0, 0.0))
	node3 := NewNode("C", NewVector(0.5, 0.866))

	system.AddNode(node1)
	system.AddNode(node2)
	system.AddNode(node3)

	coherence := system.ComputeCoherence()
	if coherence <= 0 || coherence > 1 || math.IsNaN(coherence) {
		t.Errorf("expected valid coherence value, got %f", coherence)
	}
}

func TestSystemWithCustomConstraint(t *testing.T) {
	system := NewSystem()

	node := NewNode("A", NewVector(0.0, 0.0))
	node.SetFunction(1.0)
	system.AddNode(node)

	// Custom constraint with apply function
	constraint := Constraint{
		Type:  "custom",
		Value: 0.0,
		Apply: func(s *System) float64 {
			return float64(len(s.Nodes)) * 0.5
		},
	}
	system.AddConstraint(constraint)

	eq := system.ComputeEquilibrium()
	if math.IsNaN(eq) || math.IsInf(eq, 0) {
		t.Errorf("expected valid equilibrium, got %f", eq)
	}
}
