package tag

import "math"

type ToteBubble struct {
	ID        string
	Parent    *ToteBubble
	Child     *ToteBubble
	State     float64
	Demand    float64
	Tolerance float64
}

type ErrorBubble struct {
	ID         string
	Origin     *ToteBubble
	Upstream   *ErrorBubble
	Downstream *ErrorBubble
	ErrorValue float64
	IsCulprit  bool
	Resolved   bool
}

// Build chain of totebubbles.
func chain(ids ...string) *ToteBubble {
	if len(ids) == 0 {
		return nil
	}
	head := &ToteBubble{ID: ids[0]}
	curr := head
	for i := 1; i < len(ids); i++ {
		next := &ToteBubble{ID: ids[i], Parent: curr}
		curr.Child = next
		curr = next
	}
	return head
}

// Spawn full error chain for a failing link.
func SpawnErrorChain(start *ToteBubble) *ErrorBubble {
	root := &ErrorBubble{ID: start.ID + ".err", Origin: start}
	prev := root
	for t := start.Child; t != nil; t = t.Child {
		next := &ErrorBubble{ID: t.ID + ".err", Origin: t}
		prev.Downstream = next
		next.Upstream = prev
		prev = next
	}
	return root
}

// Backfeed correction and reconcile if within tolerance.
func BackfeedAndReconcile(root *ErrorBubble, draw float64, receipts *[]Receipt, step int) (*ErrorBubble, float64) {
	candidate := root
	for candidate != nil {
		o := candidate.Origin
		if math.Abs(o.State-o.Demand) > o.Tolerance {
			candidate.IsCulprit = true
			break
		}
		candidate = candidate.Downstream
	}
	if candidate == nil {
		candidate = root
		for candidate.Downstream != nil {
			candidate = candidate.Downstream
		}
	}

	o := candidate.Origin
	err := o.Demand - o.State
	correction := math.Copysign(math.Min(math.Abs(err), draw), err)
	before := o.State
	o.State += correction

	*receipts = append(*receipts, Receipt{
		Step: step, Type: RBackfeed, Subject: candidate.ID,
		Note: "apply correction from chaostote", Value1: before, Value2: o.State,
	})

	if math.Abs(o.State-o.Demand) <= o.Tolerance {
		candidate.Resolved = true
		*receipts = append(*receipts, Receipt{
			Step: step, Type: RReconcile, Subject: candidate.ID,
			Note: "culprit within tolerance; local reconciliation",
		})
	}
	return candidate, math.Abs(correction)
}
