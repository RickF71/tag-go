// Package core implements basic vector math for TAG.
// Vectors represent the functional and constraint forces
// that interact within each node to determine equilibrium.
package core

import "math"

// Vector is a simple 3D vector used for geometric operations.
type Vector struct {
	X, Y, Z float64
}

// Add returns the vector sum of v and other.
func (v Vector) Add(other Vector) Vector {
	return Vector{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Sub returns the vector difference v - other.
func (v Vector) Sub(other Vector) Vector {
	return Vector{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Dot returns the dot product of two vectors.
func (v Vector) Dot(other Vector) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Magnitude returns the Euclidean length of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize returns a unit vector in the same direction as v.
func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector{}
	}
	return Vector{v.X / mag, v.Y / mag, v.Z / mag}
}

// Scale multiplies each component of the vector by scalar s.
func (v Vector) Scale(s float64) Vector {
	return Vector{v.X * s, v.Y * s, v.Z * s}
}
