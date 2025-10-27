package main

import (
	"fmt"
	"math"
	"time"

	"github.com/RickF71/tag-go/internal/core"
)

// Demonstrates a constraint violation: A's new function demands
// motion outside B's allowed constraint region.
func main() {
	nodeA := core.Node{
		ID:         "A",
		Function:   core.Vector{1, 1, 0},
		Constraint: core.Vector{0.2, 0.6, 0},
		Tolerance:  0.01,
	}
	nodeB := core.Node{
		ID:         "B",
		Function:   nodeA.Constraint,
		Constraint: core.Vector{0.3, 0.5, 0},
		Tolerance:  0.01,
	}

	// B's hard physical / policy limit (anything beyond is violation)
	bHardConstraint := core.Vector{X: 0.8, Y: 0.8, Z: 0}

	rate := 0.1
	fmt.Println("=== TAG Constraint-Violation Demo ===\n")

	// Phase 1: equilibrate
	for step := 0; step < 10; step++ {
		nodeA.Step(rate)
		nodeB.Function = nodeA.Constraint
		nodeB.Step(rate)
		time.Sleep(80 * time.Millisecond)
	}
	fmt.Println("Phase 1 complete: both nodes balanced\n")

	// Phase 2: disturb A with a function that demands breaking B's limit
	nodeA.Function = core.Vector{X: 1.5, Y: 1.5, Z: 0}
	fmt.Printf("--- Disturbance: A.Function -> %+v ---\n\n", nodeA.Function)

	for step := 10; step < 40; step++ {
		nodeA.Step(rate)
		nodeB.Function = nodeA.Constraint
		nodeB.Step(rate)

		// Compute how far A’s constraint exceeds B’s hard limit
		violation := distanceBeyond(nodeA.Constraint, bHardConstraint)

		fmt.Printf("Step %-2d | A.err %.4f | B.err %.4f | Violation %.4f\n",
			step, nodeA.EquilibriumError(), nodeB.EquilibriumError(), violation)

		if violation > 0 {
			fmt.Printf("*** Constraint violated at step %d (%.3f beyond limit) ***\n",
				step, violation)
			break
		}
		time.Sleep(80 * time.Millisecond)
	}

	fmt.Printf("\nFinal:\nA.Constraint=%+v\nB.Constraint=%+v\n",
		nodeA.Constraint, nodeB.Constraint)
}

// distanceBeyond returns how far vA exceeds vLimit in magnitude.
func distanceBeyond(vA, vLimit core.Vector) float64 {
	mA := vA.Magnitude()
	mL := vLimit.Magnitude()
	if mA <= mL {
		return 0
	}
	return math.Abs(mA - mL)
}
