package main

import (
	"fmt"
	"time"

	"github.com/RickF71/tag-go/internal/core"
)

// Demo: simulate a node rebalancing toward equilibrium
func main() {
	node := core.Node{
		ID:         "node-1",
		Function:   core.Vector{X: 1, Y: 1, Z: 0},
		Constraint: core.Vector{X: 0.3, Y: 0.6, Z: 0},
		Tolerance:  0.02,
	}

	fmt.Println("Starting TAG Equilibrium Demo")
	for step := 0; step < 25; step++ {
		node.Step(0.02)
		clear := node.IsClear()
		fmt.Printf("Step %-2d | Error: %.4f | Clear: %v | Constraint: %+v\n",
			step, node.EquilibriumError(), clear, node.Constraint)

		if clear {
			fmt.Println("Node reached equilibrium — clarity achieved.")
			break
		}
		time.Sleep(150 * time.Millisecond)
	}
}
