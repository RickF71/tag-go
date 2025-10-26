package core

import "math"

// epsilon is a small value used to prevent division by zero.
const epsilon = 1e-10

// System represents a dimensional system with nodes and equilibrium state.
type System struct {
	Nodes       []*Node
	Constraints []Constraint
	Equilibrium float64
}

// Constraint represents a constraint on the system.
type Constraint struct {
	Type  string
	Value float64
	Apply func(*System) float64
}

// NewSystem creates a new dimensional system.
func NewSystem() *System {
	return &System{
		Nodes:       make([]*Node, 0),
		Constraints: make([]Constraint, 0),
		Equilibrium: 0.0,
	}
}

// AddNode adds a node to the system.
func (s *System) AddNode(node *Node) {
	s.Nodes = append(s.Nodes, node)
}

// AddConstraint adds a constraint to the system.
func (s *System) AddConstraint(constraint Constraint) {
	s.Constraints = append(s.Constraints, constraint)
}

// ComputeEquilibrium calculates the equilibrium state between function and constraint.
// This implements the core TAG principle: modeling balance across dimensional systems.
func (s *System) ComputeEquilibrium() float64 {
	if len(s.Nodes) == 0 {
		return 0.0
	}

	// Compute total functional value
	var totalFunction float64
	for _, node := range s.Nodes {
		totalFunction += node.Function
	}
	avgFunction := totalFunction / float64(len(s.Nodes))

	// Compute total constraint impact
	var totalConstraint float64
	for _, constraint := range s.Constraints {
		if constraint.Apply != nil {
			totalConstraint += constraint.Apply(s)
		} else {
			totalConstraint += constraint.Value
		}
	}

	// Equilibrium is the balance between function and constraint
	// Using a normalized difference metric
	if len(s.Constraints) > 0 {
		avgConstraint := totalConstraint / float64(len(s.Constraints))
		s.Equilibrium = 1.0 - math.Abs(avgFunction-avgConstraint)/(math.Abs(avgFunction)+math.Abs(avgConstraint)+epsilon)
	} else {
		s.Equilibrium = avgFunction
	}

	return s.Equilibrium
}

// GetDimensionality returns the dimensionality of the system based on its nodes.
func (s *System) GetDimensionality() int {
	if len(s.Nodes) == 0 {
		return 0
	}
	return s.Nodes[0].Position.Dim()
}

// ComputeCoherence calculates the coherence of the system based on node distribution.
func (s *System) ComputeCoherence() float64 {
	if len(s.Nodes) < 2 {
		return 1.0
	}

	// Compute average distance between nodes
	var totalDistance float64
	var count int
	for i := 0; i < len(s.Nodes); i++ {
		for j := i + 1; j < len(s.Nodes); j++ {
			totalDistance += s.Nodes[i].DistanceTo(s.Nodes[j])
			count++
		}
	}
	avgDistance := totalDistance / float64(count)

	// Compute variance in distances
	var variance float64
	for i := 0; i < len(s.Nodes); i++ {
		for j := i + 1; j < len(s.Nodes); j++ {
			dist := s.Nodes[i].DistanceTo(s.Nodes[j])
			variance += math.Pow(dist-avgDistance, 2)
		}
	}
	variance /= float64(count)

	// Coherence inversely related to variance (normalized)
	coherence := 1.0 / (1.0 + variance)
	return coherence
}
