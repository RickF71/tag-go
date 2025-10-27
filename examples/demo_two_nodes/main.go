package main

import (
	"fmt"
	"time"

	"github.com/RickF71/tag-go/internal/core"
)

// Two nodes in sequence: A drives B.
// B's Function is always A's current Constraint.
// A disturbance is injected at step 10 (A.Function changes).
func main() {
	// Top node A
	nodeA := core.Node{
		ID:         "A",
		Function:   core.Vector{X: 1, Y: 1, Z: 0},
		Constraint: core.Vector{X: 0.2, Y: 0.7, Z: 0},
		Tolerance:  0.01,
	}
	// Lower node B
	nodeB := core.Node{
		ID:         "B",
		Function:   nodeA.Constraint,
		Constraint: core.Vector{X: 0.2, Y: 0.6, Z: 0},
		Tolerance:  0.01,
	}

	rate := 0.1
	fmt.Println("Starting Two-Node Equilibrium Demo (with disturbance at step 10)\n")

	for step := 0; step < 40; step++ {
		// Inject a disturbance at step 10: change A's Function vector
		if step == 10 {
			fmt.Printf("\n--- DISTURBANCE @ step %d: A.Function changed from %+v to ", step, nodeA.Function)
			nodeA.Function = core.Vector{X: 1.2, Y: 0.6, Z: 0} // tweak direction & magnitude
			fmt.Printf("%+v ---\n\n", nodeA.Function)
		}

		// Update A first
		nodeA.Step(rate)

		// B follows A's current constraint
		nodeB.Function = nodeA.Constraint
		nodeB.Step(rate)

		fmt.Printf(
			"Step %-2d | A.err %.4f | B.err %.4f | A.clear %-5v | B.clear %-5v | A.C %+v | B.C %+v\n",
			step, nodeA.EquilibriumError(), nodeB.EquilibriumError(),
			nodeA.IsClear(), nodeB.IsClear(),
			nodeA.Constraint, nodeB.Constraint,
		)

		time.Sleep(120 * time.Millisecond)
	}

	fmt.Printf("\nFinal states:\nA.Constraint = %+v\nB.Constraint = %+v\n",
		nodeA.Constraint, nodeB.Constraint)
}
