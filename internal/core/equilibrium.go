// Package core: equilibrium.go implements the Equilibrium Law.
// Each node reaches clarity when its function and constraint vectors
// are in balance (dot product difference near zero).
package core

import "math"

// EquilibriumError returns the difference in alignment between
// function and constraint vectors.
func (n *Node) EquilibriumError() float64 {
	// Smaller dot product difference = more aligned
	return n.Function.Normalize().Dot(n.Constraint.Normalize())
}

// IsClear returns true if the node is within equilibrium tolerance.
func (n *Node) IsClear() bool {
	return math.Abs(1.0-n.EquilibriumError()) < n.Tolerance
}

// Step adjusts the constraint slightly toward the function,
// simulating rebalancing toward equilibrium.
func (n *Node) Step(rate float64) {
	err := n.EquilibriumError()
	if err < 1.0 {
		diff := n.Function.Sub(n.Constraint).Scale(rate)
		n.Constraint = n.Constraint.Add(diff)
	}
}
