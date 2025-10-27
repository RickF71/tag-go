package tag

import (
	"math"
)

type Chaostote struct {
	ID        string
	Viscosity float64
	Field     []*ErrorBubble
}

// Diffuse smooths error values.
func (c *Chaostote) Diffuse(dt float64, receipts *[]Receipt, step int) {
	decay := math.Exp(-c.Viscosity * dt)
	for _, e := range c.Field {
		if e.Resolved {
			continue
		}
		before := e.ErrorValue
		e.ErrorValue *= decay
		*receipts = append(*receipts, Receipt{
			Step: step, Type: RDiffuse, Subject: e.ID,
			Note: "chaostote diffusion", Value1: before, Value2: e.ErrorValue,
		})
	}
}

// InjectChain snapshots errors from a chain into the chaostote.
func (c *Chaostote) InjectChain(root *ErrorBubble, receipts *[]Receipt, step int) {
	for e := root; e != nil; e = e.Downstream {
		err := math.Max(0, math.Abs(e.Origin.State-e.Origin.Demand)-e.Origin.Tolerance)
		e.ErrorValue = err
		c.Field = append(c.Field, e)
		*receipts = append(*receipts, Receipt{
			Step: step, Type: RInject, Subject: e.ID,
			Note: "inject into chaostote", Value1: err,
		})
	}
}

// TotalError returns sum of unresolved errors.
func (c *Chaostote) TotalError() float64 {
	sum := 0.0
	for _, e := range c.Field {
		if !e.Resolved {
			sum += e.ErrorValue
		}
	}
	return sum
}

// CheckMetaBirth spawns a meta totebubble when chaos saturates.
func (c *Chaostote) CheckMetaBirth(limit float64, receipts *[]Receipt, step int) *ToteBubble {
	total := c.TotalError()
	if total > limit {
		meta := &ToteBubble{
			ID:        c.ID + ".meta",
			State:     total,
			Demand:    0,
			Tolerance: limit * 0.1,
		}
		*receipts = append(*receipts, Receipt{
			Step: step, Type: RMetaBirth, Subject: meta.ID,
			Note:   "chaostote exceeded limit; spawned meta totebubble",
			Value1: total, Value2: limit,
		})
		return meta
	}
	return nil
}

// DrainChaostote removes error energy from the field.
func DrainChaostote(c *Chaostote, amount float64) float64 {
	remaining := amount
	for _, e := range c.Field {
		if e.Resolved || e.ErrorValue <= 0 {
			continue
		}
		if remaining <= 0 {
			break
		}
		take := math.Min(e.ErrorValue, remaining)
		e.ErrorValue -= take
		remaining -= take
	}
	return amount - remaining
}
