// node.go: Node type for TAG core
package core

// Node represents a node in the TAG framework.
type Node struct {
	ID         string
	Function   Vector
	Constraint Vector
	Tolerance  float64
}
