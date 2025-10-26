# TAG - Theory of Asymptotic Geometry

Theory of Asymptotic Geometry (TAG) is a Go implementation for modeling equilibrium between function and constraint across dimensional systems.

## Project Purpose

TAG provides a framework for:
- Modeling equilibrium states in multi-dimensional systems
- Balancing functional values against system constraints
- Computing coherence measures across node distributions
- Simulating asymptotic geometric relationships

## Structure

The project follows idiomatic Go 1.23+ conventions with a clean separation of concerns:

```
tag-go/
├── cmd/tag/           # CLI entry point
│   └── main.go        # Command-line interface
├── internal/core/     # Core implementation (private)
│   ├── vector.go      # Vector mathematics
│   ├── node.go        # Node structures
│   └── equilibrium.go # Equilibrium logic
├── pkg/tag/           # Public API
│   └── tag.go         # Public wrapper interface
├── go.mod             # Module definition
├── LICENSE            # MIT License
└── README.md          # This file
```

## Installation

```bash
go get github.com/RickF71/tag-go/pkg/tag
```

## Usage

### As a Library

```go
package main

import (
    "fmt"
    "github.com/RickF71/tag-go/pkg/tag"
)

func main() {
    // Create a new system
    system := tag.NewSystem()

    // Create nodes in 3D space
    node1 := tag.NewNode("A", tag.NewVector(0.0, 0.0, 0.0))
    node1.SetFunction(1.0)
    
    system.AddNode(node1)

    // Add constraints
    constraint := tag.Constraint{Type: "boundary", Value: 0.85}
    system.AddConstraint(constraint)

    // Compute equilibrium
    equilibrium := system.ComputeEquilibrium()
    fmt.Printf("Equilibrium: %.4f\n", equilibrium)
}
```

### As a CLI Tool

```bash
# Build the CLI
go build -o tag ./cmd/tag

# Run demo
./tag demo
```

## Core Concepts

### Vector
N-dimensional vectors with standard operations (add, subtract, dot product, magnitude, normalization).

### Node
Points in dimensional space with associated functional values.

### System
Collections of nodes with constraints, capable of computing equilibrium and coherence states.

### Equilibrium
The balance between functional values and system constraints, representing the core TAG principle.

## Dependencies

TAG uses only the Go standard library, with no external dependencies.

## License

MIT License - See LICENSE file for details.
