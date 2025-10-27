# Theory of Asymptotic Geometry (TAG) Go Project

This repository scaffolds the TAG framework in Go.

## Structure

- `cmd/tag/` — CLI entrypoint (hello-world starter)
- `internal/core/` — Core types: `vector.go`, `node.go`, `equilibrium.go`
- `internal/sim/` — Simulation logic (empty)
- `internal/canon/laws/` — Canonical law YAMLs (e.g., `equilibrium.v1.yaml`)
- `pkg/tag/` — Public API (see `api.go`)
- `examples/demo_equilibrium/` — Example usage of the TAG API

## Usage

- Run the CLI: `go run cmd/tag/main.go`
- See example: `go run examples/demo_equilibrium/main.go`

No external dependencies beyond the Go standard library.
