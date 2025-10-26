// Package tag provides a public API for the Theory of Asymptotic Geometry (TAG) system.
// TAG models equilibrium between function and constraint across dimensional systems.
package tag

import "github.com/RickF71/tag-go/internal/core"

// Vector is a public wrapper for core.Vector.
type Vector = core.Vector

// Node is a public wrapper for core.Node.
type Node = core.Node

// System is a public wrapper for core.System.
type System = core.System

// Constraint is a public wrapper for core.Constraint.
type Constraint = core.Constraint

// NewVector creates a new vector with the given components.
func NewVector(components ...float64) Vector {
	return core.NewVector(components...)
}

// NewNode creates a new node with the given ID and position.
func NewNode(id string, position Vector) *Node {
	return core.NewNode(id, position)
}

// NewSystem creates a new dimensional system.
func NewSystem() *System {
	return core.NewSystem()
}
