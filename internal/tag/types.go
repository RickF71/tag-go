package tag

// Params holds tunable simulation settings.
type Params struct {
	Viscosity float64 `json:"viscosity"`
	Limit     float64 `json:"limit"`
	Dt        float64 `json:"dt"`
}

// SimState is a JSON snapshot returned by /api/tag/state and /api/tag/step.
type SimState struct {
	Step       int       `json:"step"`
	TotalError float64   `json:"total_error"`
	MetaEnergy float64   `json:"meta_energy"`
	Receipts   []Receipt `json:"receipts,omitempty"`
}

// internal/tag/types.go
type Bubble struct {
	Label     string     `json:"label"`      // "A.tote", "A.err", "B.tote", "B.err"
	Center    [2]float64 `json:"center"`     // normalized [-1,1] space
	Radius    float64    `json:"radius"`     // scaled magnitude for drawing
	Excess    float64    `json:"excess"`     // functional excess parked in error bubble
	Energy    float64    `json:"energy"`     // optional detail
	ColorHint string     `json:"color_hint"` // e.g. "primary","error"
}

type Frame struct {
	Step       int `json:"step"`
	Equilibria struct {
		A bool `json:"A"`
		B bool `json:"B"`
		C bool `json:"C"`
		D bool `json:"D"`
	} `json:"equilibria"`
	Bubbles    []Bubble `json:"bubbles"`      // four entries: A.tote, A.err, B.tote, B.err
	ErrScalarA float64  `json:"err_scalar_a"` // recent error metric
	ErrScalarB float64  `json:"err_scalar_b"`
}
