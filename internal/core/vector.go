package core

import "math"

// Vector represents a mathematical vector in n-dimensional space.
type Vector struct {
	Components []float64
}

// NewVector creates a new vector with the given components.
func NewVector(components ...float64) Vector {
	return Vector{Components: components}
}

// Dim returns the dimensionality of the vector.
func (v Vector) Dim() int {
	return len(v.Components)
}

// Add returns the sum of two vectors.
func (v Vector) Add(other Vector) Vector {
	if v.Dim() != other.Dim() {
		panic("vector dimensions must match")
	}
	result := make([]float64, v.Dim())
	for i := range v.Components {
		result[i] = v.Components[i] + other.Components[i]
	}
	return Vector{Components: result}
}

// Sub returns the difference of two vectors.
func (v Vector) Sub(other Vector) Vector {
	if v.Dim() != other.Dim() {
		panic("vector dimensions must match")
	}
	result := make([]float64, v.Dim())
	for i := range v.Components {
		result[i] = v.Components[i] - other.Components[i]
	}
	return Vector{Components: result}
}

// Scale multiplies the vector by a scalar.
func (v Vector) Scale(scalar float64) Vector {
	result := make([]float64, v.Dim())
	for i, c := range v.Components {
		result[i] = c * scalar
	}
	return Vector{Components: result}
}

// Dot returns the dot product of two vectors.
func (v Vector) Dot(other Vector) float64 {
	if v.Dim() != other.Dim() {
		panic("vector dimensions must match")
	}
	var sum float64
	for i := range v.Components {
		sum += v.Components[i] * other.Components[i]
	}
	return sum
}

// Magnitude returns the magnitude (length) of the vector.
func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.Dot(v))
}

// Normalize returns a unit vector in the same direction.
func (v Vector) Normalize() Vector {
	mag := v.Magnitude()
	if mag == 0 {
		return v
	}
	return v.Scale(1.0 / mag)
}

// Distance returns the Euclidean distance to another vector.
func (v Vector) Distance(other Vector) float64 {
	return v.Sub(other).Magnitude()
}
