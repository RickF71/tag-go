package main

import (
	"fmt"
	"math"
	"strings"
)

// ---------- Receipts (timeline) ----------

type ReceiptType string

const (
	RSpawnChain ReceiptType = "spawn_chain"
	RInject     ReceiptType = "inject"
	RDiffuse    ReceiptType = "diffuse"
	RMetaBirth  ReceiptType = "meta_totelevation"
	RBackfeed   ReceiptType = "backfeed"
	RReconcile  ReceiptType = "reconcile"
	RQuench     ReceiptType = "quench"
)

type Receipt struct {
	Step    int
	Type    ReceiptType
	Subject string
	Note    string
	Value1  float64
	Value2  float64
}

// ---------- Tote / Error / Chaostote ----------

type ToteBubble struct {
	ID        string
	Parent    *ToteBubble
	Child     *ToteBubble
	State     float64 // what it currently has
	Demand    float64 // what parent wants from it
	Tolerance float64 // λ
}

type ErrorBubble struct {
	ID         string
	Origin     *ToteBubble
	Upstream   *ErrorBubble
	Downstream *ErrorBubble
	ErrorValue float64 // mismatch snapshot projected into Ξ
	IsCulprit  bool
	Resolved   bool
}

type Chaostote struct {
	ID        string
	Viscosity float64 // νχ (how fast errors smooth)
	Field     []*ErrorBubble
}

// Spawn a mirrored error chain starting at 'start' and going downstream.
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

func (c *Chaostote) InjectChain(root *ErrorBubble, receipts *[]Receipt, step int) {
	for e := root; e != nil; e = e.Downstream {
		// Snapshot the local error (|state - demand| beyond tolerance)
		err := math.Max(0, math.Abs(e.Origin.State-e.Origin.Demand)-e.Origin.Tolerance)
		e.ErrorValue = err
		*c.FieldPtr() = append(*c.FieldPtr(), e)
		*receipts = append(*receipts, Receipt{
			Step:    step,
			Type:    RInject,
			Subject: e.ID,
			Note:    "inject into chaostote",
			Value1:  err,
		})
	}
}

func (c *Chaostote) FieldPtr() *[]*ErrorBubble { return &c.Field }

func (c *Chaostote) TotalError() float64 {
	sum := 0.0
	for _, e := range c.Field {
		if !e.Resolved {
			sum += e.ErrorValue
		}
	}
	return sum
}

// Simple exponential-like smoothing of error values.
func (c *Chaostote) Diffuse(dt float64, receipts *[]Receipt, step int) {
	decay := math.Exp(-c.Viscosity * dt)
	for _, e := range c.Field {
		if e.Resolved {
			continue
		}
		before := e.ErrorValue
		e.ErrorValue *= decay
		*receipts = append(*receipts, Receipt{
			Step:    step,
			Type:    RDiffuse,
			Subject: e.ID,
			Note:    "chaostote diffusion",
			Value1:  before,
			Value2:  e.ErrorValue,
		})
	}
}

// If total error exceeds limit, the chaostote becomes its own totebubble.
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
			Step:    step,
			Type:    RMetaBirth,
			Subject: meta.ID,
			Note:    "chaostote exceeded limit; spawned meta totebubble",
			Value1:  total,
			Value2:  limit,
		})
		return meta
	}
	return nil
}

// Draw correction from the chaostote back into the functional chain.
// We locate the culprit: first link (downstream) that still violates tolerance.
func BackfeedAndReconcile(root *ErrorBubble, draw float64, receipts *[]Receipt, step int) (culprit *ErrorBubble, used float64) {
	// Find culprit
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
		// fallback: use far downstream node
		candidate = root
		for candidate.Downstream != nil {
			candidate = candidate.Downstream
		}
	}

	// Apply correction: push State toward Demand using available draw
	o := candidate.Origin
	err := o.Demand - o.State
	correction := math.Copysign(math.Min(math.Abs(err), draw), err)
	before := o.State
	o.State += correction
	used = math.Abs(correction)

	*receipts = append(*receipts, Receipt{
		Step:    step,
		Type:    RBackfeed,
		Subject: candidate.ID,
		Note:    "apply correction from chaostote",
		Value1:  before,
		Value2:  o.State,
	})

	// Mark resolved if now within tolerance
	if math.Abs(o.State-o.Demand) <= o.Tolerance {
		candidate.Resolved = true
		*receipts = append(*receipts, Receipt{
			Step:    step,
			Type:    RReconcile,
			Subject: candidate.ID,
			Note:    "culprit within tolerance; local reconciliation",
		})
	}
	return candidate, used
}

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
	used := amount - remaining
	return used
}

// ---------- Demo world ----------

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

func printReceipts(rs []Receipt) {
	const eps = 1e-4
	lastKey := ""
	lastVal := 0.0

	for _, r := range rs {
		// major events always print
		if r.Type == RSpawnChain || r.Type == RMetaBirth || r.Type == RReconcile || r.Type == RQuench {
			fmt.Printf("[%03d] %-14s %-10s %s", r.Step, r.Type, r.Subject, r.Note)
			if r.Value1 != 0 || r.Value2 != 0 {
				fmt.Printf("  (%.4f → %.4f)", r.Value1, r.Value2)
			}
			fmt.Println()
			continue
		}

		// build key by event type + subject (ignore step)
		key := string(r.Type) + ":" + r.Subject
		// only print if numeric change beyond eps or different key
		if key != lastKey || math.Abs(r.Value2-lastVal) > eps {
			fmt.Printf("[%03d] %-14s %-10s %s", r.Step, r.Type, r.Subject, r.Note)
			if r.Value1 != 0 || r.Value2 != 0 {
				fmt.Printf("  (%.4f → %.4f)", r.Value1, r.Value2)
			}
			fmt.Println()
			lastKey = key
			lastVal = r.Value2
		}
	}
}

func main() {
	// Build A -> B -> C -> D
	A := chain("A", "B", "C", "D")

	// Set tolerances and an initial failing scenario:
	// A demands a lot from B; B can't fully meet it; that pressure propagates downstream.
	A.State, A.Demand, A.Tolerance = 0, 0, 0.01
	B := A.Child
	B.State, B.Demand, B.Tolerance = 1.2, 1.6, 0.05 // B under-delivers for A
	C := B.Child
	C.State, C.Demand, C.Tolerance = 1.5, 1.5, 0.05 // C looks fine initially
	D := C.Child
	D.State, D.Demand, D.Tolerance = 1.5, 1.5, 0.05 // D too

	// Chaostote
	chi := &Chaostote{ID: "Χ", Viscosity: 0.05}

	var receipts []Receipt
	step := 0
	dt := 1.0

	// When B fails A, spawn full error mirror A.err->B.err->C.err->D.err (inside chaostote)
	rootErr := SpawnErrorChain(B) // start from B (the failing link relative to A)
	receipts = append(receipts, Receipt{Step: step, Type: RSpawnChain, Subject: "B…D.err", Note: "spawned error mirror chain"})
	chi.InjectChain(rootErr, &receipts, step)

	// Simulation loop
	limit := 0.50 // when total error in chaostote exceeds this, meta totebubble is born
	var meta *ToteBubble

	for step = 1; step <= 40; step++ {
		// keep the failing demand pressure active on B (simulate persistent request from A)
		B.Demand = 1.6 // A still wants more from B

		// re-inject current snapshot mismatch into chaostote (keeps it “lit”)
		chi.InjectChain(rootErr, &receipts, step)

		// diffuse the field (chaostote viscosity)
		chi.Diffuse(dt, &receipts, step)

		// meta birth check
		if meta == nil {
			if m := chi.CheckMetaBirth(limit, &receipts, step); m != nil {
				meta = m
			}
		} else {
			// Meta bubble exists: draw correction energy out of chaostote and apply upstream.
			draw := meta.State * 0.25 // use a fraction each step
			if draw > 0 {
				// pull from chaostote reservoir
				used := DrainChaostote(chi, draw)
				meta.State -= used

				// push correction to the culprit in the error chain
				_, applied := BackfeedAndReconcile(rootErr, used, &receipts, step)

				// If nothing to apply (rare), damp meta a bit to avoid stalling.
				if applied == 0 {
					meta.State *= 0.9
				}
			}

			// If chaostote error mostly gone and meta energy low, we’re reconciled.
			if chi.TotalError() < 1e-3 && meta.State < 1e-3 {
				receipts = append(receipts, Receipt{
					Step:    step,
					Type:    RReconcile,
					Subject: meta.ID,
					Note:    "meta & chaostote reconciled; field collapsed",
				})
				break
			}
		}
	}

	// Pretty print
	title := "=== TAG Demo — Chaostote Meta-Totelevation ==="
	underline := strings.Repeat("=", len(title))
	fmt.Println(title)
	fmt.Println(underline)
	fmt.Println("Chain: A -> B -> C -> D")
	fmt.Printf("Tolerance(B)=%.2f, Demand(B)=%.2f, State(B)=%.2f\n", B.Tolerance, B.Demand, B.State)
	fmt.Printf("Chaostote νχ=%.2f, Limit=%.2f\n\n", chi.Viscosity, limit)
	printReceipts(receipts)
}
