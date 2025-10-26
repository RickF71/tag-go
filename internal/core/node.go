package core

// Node represents a point in dimensional space with associated properties.
type Node struct {
	ID       string
	Position Vector
	Function float64 // Functional value at this node
}

// NewNode creates a new node with the given ID and position.
func NewNode(id string, position Vector) *Node {
	return &Node{
		ID:       id,
		Position: position,
		Function: 0.0,
	}
}

// SetFunction sets the functional value for this node.
func (n *Node) SetFunction(value float64) {
	n.Function = value
}

// DistanceTo returns the distance from this node to another.
func (n *Node) DistanceTo(other *Node) float64 {
	return n.Position.Distance(other.Position)
}
