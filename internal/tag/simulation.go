package tag

import (
	"encoding/json"
	"math"
	"sync"
	"time"
)

// Simulation wraps the existing Chaostote + chain structures
// and exposes thread-safe control for the API layer.
type Simulation struct {
	mu        sync.Mutex
	StepNum   int
	Chi       *Chaostote
	Root      *ToteBubble
	ErrRoot   *ErrorBubble
	Meta      *ToteBubble
	Receipts  []Receipt
	ParamsCfg Params
}

// --- construction and setup ---

func NewSimulation() *Simulation {
	A := chain("A", "B", "C", "D")
	A.State, A.Demand, A.Tolerance = 0, 0, 0.01
	B := A.Child
	B.State, B.Demand, B.Tolerance = 1.2, 1.6, 0.05
	C := B.Child
	C.State, C.Demand, C.Tolerance = 1.5, 1.5, 0.05
	D := C.Child
	D.State, D.Demand, D.Tolerance = 1.5, 1.5, 0.05

	chi := &Chaostote{ID: "Χ", Viscosity: 0.05}
	rootErr := SpawnErrorChain(B)

	return &Simulation{
		Chi:       chi,
		Root:      A,
		ErrRoot:   rootErr,
		ParamsCfg: Params{Viscosity: 0.05, Limit: 0.5, Dt: 1.0},
	}
}

// --- control ---

func (s *Simulation) Step() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.StepNum++
	step := s.StepNum
	dt := s.ParamsCfg.Dt

	B := s.Root.Child
	B.Demand = 1.6 // keep A’s demand constant

	s.Chi.InjectChain(s.ErrRoot, &s.Receipts, step)
	s.Chi.Diffuse(dt, &s.Receipts, step)

	if s.Meta == nil {
		if m := s.Chi.CheckMetaBirth(s.ParamsCfg.Limit, &s.Receipts, step); m != nil {
			s.Meta = m
		}
	} else {
		draw := s.Meta.State * 0.25
		if draw > 0 {
			used := DrainChaostote(s.Chi, draw)
			s.Meta.State -= used
			_, _ = BackfeedAndReconcile(s.ErrRoot, used, &s.Receipts, step)
		}
		if s.Chi.TotalError() < 1e-3 && s.Meta.State < 1e-3 {
			s.Receipts = append(s.Receipts, Receipt{
				Step:    step,
				Type:    RReconcile,
				Subject: s.Meta.ID,
				Note:    "meta & chaostote reconciled; field collapsed",
			})
		}
	}
}

func (s *Simulation) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	*s = *NewSimulation()
}

func (s *Simulation) UpdateParams(p Params) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.Viscosity != 0 {
		s.Chi.Viscosity = p.Viscosity
		s.ParamsCfg.Viscosity = p.Viscosity
	}
	if p.Limit != 0 {
		s.ParamsCfg.Limit = p.Limit
	}
	if p.Dt != 0 {
		s.ParamsCfg.Dt = p.Dt
	}
}

func (s *Simulation) Params() Params { return s.ParamsCfg }

func (s *Simulation) Snapshot() SimState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return SimState{
		Step:       s.StepNum,
		TotalError: s.Chi.TotalError(),
		MetaEnergy: func() float64 {
			if s.Meta != nil {
				return s.Meta.State
			}
			return 0
		}(),
		Receipts: s.Receipts,
	}
}

// Stream runs physics at 30 Hz but only sends state to the GUI 5 fps (~every 200 ms).
func (s *Simulation) Stream() <-chan []byte {
	out := make(chan []byte)

	go func() {
		physTicker := time.NewTicker(time.Second / 30)       // 30 Hz physics
		sendTicker := time.NewTicker(200 * time.Millisecond) // 5 Hz output
		defer physTicker.Stop()
		defer sendTicker.Stop()

		var lastStep int
		var lastTotal, lastMeta float64

		for {
			select {
			case <-physTicker.C:
				s.Step()

			case <-sendTicker.C:
				state := s.Snapshot()

				// Only send if something actually changed
				if state.Step != lastStep ||
					math.Abs(state.TotalError-lastTotal) > 1e-5 ||
					math.Abs(state.MetaEnergy-lastMeta) > 1e-5 {

					b, _ := json.Marshal(state)
					select {
					case out <- b:
					default:
					} // drop frame if busy

					lastStep = state.Step
					lastTotal = state.TotalError
					lastMeta = state.MetaEnergy
				}
			}
		}
	}()

	return out
}
