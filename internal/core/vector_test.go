package core

import (
	"math"
	"testing"
)

func TestNewVector(t *testing.T) {
	v := NewVector(1.0, 2.0, 3.0)
	if v.Dim() != 3 {
		t.Errorf("expected dimension 3, got %d", v.Dim())
	}
	if v.Components[0] != 1.0 || v.Components[1] != 2.0 || v.Components[2] != 3.0 {
		t.Errorf("unexpected vector components: %v", v.Components)
	}
}

func TestVectorAdd(t *testing.T) {
	v1 := NewVector(1.0, 2.0)
	v2 := NewVector(3.0, 4.0)
	result := v1.Add(v2)

	expected := []float64{4.0, 6.0}
	for i, val := range result.Components {
		if val != expected[i] {
			t.Errorf("expected %v, got %v", expected, result.Components)
			break
		}
	}
}

func TestVectorSub(t *testing.T) {
	v1 := NewVector(5.0, 7.0)
	v2 := NewVector(2.0, 3.0)
	result := v1.Sub(v2)

	expected := []float64{3.0, 4.0}
	for i, val := range result.Components {
		if val != expected[i] {
			t.Errorf("expected %v, got %v", expected, result.Components)
			break
		}
	}
}

func TestVectorScale(t *testing.T) {
	v := NewVector(2.0, 4.0)
	result := v.Scale(3.0)

	expected := []float64{6.0, 12.0}
	for i, val := range result.Components {
		if val != expected[i] {
			t.Errorf("expected %v, got %v", expected, result.Components)
			break
		}
	}
}

func TestVectorDot(t *testing.T) {
	v1 := NewVector(1.0, 2.0, 3.0)
	v2 := NewVector(4.0, 5.0, 6.0)
	result := v1.Dot(v2)

	expected := 32.0 // 1*4 + 2*5 + 3*6
	if result != expected {
		t.Errorf("expected %f, got %f", expected, result)
	}
}

func TestVectorMagnitude(t *testing.T) {
	v := NewVector(3.0, 4.0)
	mag := v.Magnitude()

	expected := 5.0
	if mag != expected {
		t.Errorf("expected %f, got %f", expected, mag)
	}
}

func TestVectorNormalize(t *testing.T) {
	v := NewVector(3.0, 4.0)
	normalized := v.Normalize()

	expectedMag := 1.0
	if mag := normalized.Magnitude(); math.Abs(mag-expectedMag) > 1e-10 {
		t.Errorf("expected magnitude %f, got %f", expectedMag, mag)
	}
}

func TestVectorDistance(t *testing.T) {
	v1 := NewVector(0.0, 0.0)
	v2 := NewVector(3.0, 4.0)
	dist := v1.Distance(v2)

	expected := 5.0
	if dist != expected {
		t.Errorf("expected %f, got %f", expected, dist)
	}
}

func TestVectorPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for mismatched dimensions")
		}
	}()

	v1 := NewVector(1.0, 2.0)
	v2 := NewVector(1.0, 2.0, 3.0)
	_ = v1.Add(v2)
}
