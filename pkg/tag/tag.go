// Package tag provides a public API for the Theory of Asymptotic Geometry (TAG) system.
// TAG models equilibrium between function and constraint across dimensional systems.
//
// This package provides type aliases to internal/core types, which is appropriate
// for this simple API where we want to maintain a stable public interface while
// keeping implementation details internal. This follows Go's re-export pattern.
package tag

import "github.com/RickF71/tag-go/internal/core"

// Vector is a public wrapper for core.Vector.
// It represents a mathematical vector in n-dimensional space.
type Vector = core.Vector

// Node is a public wrapper for core.Node.
// It represents a point in dimensional space with associated properties.
type Node = core.Node

// System is a public wrapper for core.System.
// It represents a dimensional system with nodes and equilibrium state.
type System = core.System

// Constraint is a public wrapper for core.Constraint.
// It represents a constraint on the system.
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
