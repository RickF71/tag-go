package main

import (
	"fmt"
	"os"

	"github.com/RickF71/tag-go/pkg/tag"
)

func main() {
	fmt.Println("TAG - Theory of Asymptotic Geometry")
	fmt.Println("=====================================")
	fmt.Println()

	if len(os.Args) > 1 && os.Args[1] == "demo" {
		runDemo()
		return
	}

	fmt.Println("Usage:")
	fmt.Println("  tag demo    Run a demonstration of TAG equilibrium")
	fmt.Println()
	fmt.Println("TAG models equilibrium between function and constraint across dimensional systems.")
}

func runDemo() {
	fmt.Println("Running TAG Demo...")
	fmt.Println()

	// Create a new system
	system := tag.NewSystem()

	// Create nodes in 3D space
	node1 := tag.NewNode("A", tag.NewVector(0.0, 0.0, 0.0))
	node1.SetFunction(1.0)

	node2 := tag.NewNode("B", tag.NewVector(1.0, 0.0, 0.0))
	node2.SetFunction(0.8)

	node3 := tag.NewNode("C", tag.NewVector(0.5, 0.866, 0.0))
	node3.SetFunction(0.9)

	// Add nodes to system
	system.AddNode(node1)
	system.AddNode(node2)
	system.AddNode(node3)

	fmt.Printf("System dimensionality: %d\n", system.GetDimensionality())
	fmt.Printf("Number of nodes: %d\n", len(system.Nodes))
	fmt.Println()

	// Add constraints
	constraint1 := tag.Constraint{
		Type:  "boundary",
		Value: 0.85,
	}
	system.AddConstraint(constraint1)

	// Compute equilibrium
	equilibrium := system.ComputeEquilibrium()
	fmt.Printf("System equilibrium: %.4f\n", equilibrium)

	// Compute coherence
	coherence := system.ComputeCoherence()
	fmt.Printf("System coherence: %.4f\n", coherence)

	fmt.Println()
	fmt.Println("Demo complete!")
}
