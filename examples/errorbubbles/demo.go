package main

import (
	"math"
	"time"

	"github.com/RickF71/tag-go/internal/tag"
)

// GenerateDemoFrame creates 8 bubbles (A–D + error partners).
func GenerateDemoFrame(step int) tag.Frame {
	t := float64(step) / 25.0
	nodes := []string{"A", "B", "C", "D"}
	var bubbles []tag.Bubble

	for i, name := range nodes {
		angle := float64(i) * math.Pi / 2
		tx := math.Cos(angle) * 0.7
		ty := math.Sin(angle) * 0.7

		mainR := 0.1 + 0.05*math.Sin(t+float64(i))
		errMag := 0.05 * math.Abs(math.Sin(t*0.5+float64(i)))

		bubbles = append(bubbles,
			tag.Bubble{
				Label:     name + ".tote",
				Center:    [2]float64{tx, ty},
				Radius:    mainR,
				Energy:    mainR * 8,
				ColorHint: "primary",
			},
			tag.Bubble{
				Label:     name + ".err",
				Center:    [2]float64{tx * 0.6, ty * 0.6},
				Radius:    errMag,
				Excess:    errMag,
				ColorHint: "error",
			},
		)
	}

	f := tag.Frame{
		Step:    step,
		Bubbles: bubbles,
	}
	f.Equilibria.A = math.Abs(bubbles[1].Radius) < 0.005
	f.Equilibria.B = math.Abs(bubbles[3].Radius) < 0.005
	f.Equilibria.C = math.Abs(bubbles[5].Radius) < 0.005
	f.Equilibria.D = math.Abs(bubbles[7].Radius) < 0.005

	return f
}

// RunLoop emits frames at a fixed rate and calls the given callback.
func RunLoop(callback func(tag.Frame)) {
	step := 0
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		callback(GenerateDemoFrame(step))
		step++
	}
}
