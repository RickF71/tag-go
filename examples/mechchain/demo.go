package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

// Adjustable anchors
var anchorTopK = 100.0
var anchorBotK = 100.0

// Vec2 = 2D vector
type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Mass = one node in chain
type Mass struct {
	P, V, F Vec2    // position, velocity, accumulated force
	M, K, C float64 // mass, spring stiffness, damping
}

type Anchors struct {
	TopK float64 `json:"topK"`
	BotK float64 `json:"botK"`
}

// 8-node system A–H
var (
	sys     [8]*Mass
	sysMu   sync.Mutex
	dt      = 0.02 // integration timestep
	running bool
	stepCnt int
)

const restLen = 1.0 // rest length between neighbors (must match visual)
const gravity = 0.0 // disable gravity for pure spring behavior

// Initialize system
func initSystem() {
	for i := range sys {
		sys[i] = &Mass{
			P: Vec2{X: float64(i) * restLen, Y: 0},
			V: Vec2{},
			M: 1.0,
			K: 40.0,
			C: 1.5,
		}
	}
}

// Integrate forces one step
func Step() {
	sysMu.Lock()
	defer sysMu.Unlock()

	// clear forces
	for _, m := range sys {
		m.F = Vec2{}
	}

	// Anchors
	left := sys[0]
	right := sys[len(sys)-1]

	dxL := left.P.X - 0
	left.F.X += -anchorTopK * dxL

	target := float64(len(sys)-1) * restLen
	dxR := right.P.X - target
	right.F.X += -anchorBotK * dxR

	// Internal springs
	for i := 0; i < len(sys)-1; i++ {
		left := sys[i]
		right := sys[i+1]
		dx := right.P.X - left.P.X
		dy := right.P.Y - left.P.Y
		dist := math.Hypot(dx, dy)
		if dist == 0 {
			continue
		}
		ux, uy := dx/dist, dy/dist
		fs := left.K * (dist - restLen)
		fv := left.C * ((right.V.X-left.V.X)*ux + (right.V.Y-left.V.Y)*uy)
		f := fs + fv
		left.F.X += f * ux
		left.F.Y += f * uy
		right.F.X -= f * ux
		right.F.Y -= f * uy
	}

	// Integrate
	for _, m := range sys {
		ax := m.F.X/m.M + 0
		ay := m.F.Y/m.M + gravity
		m.V.X += ax * dt
		m.V.Y += ay * dt
		m.P.X += m.V.X * dt
		m.P.Y += m.V.Y * dt
	}
}

// Directly set a node position
func handleSet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/set/"):]
	var id int
	fmt.Sscanf(idStr, "%d", &id)
	var v Vec2
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "invalid vector", 400)
		return
	}
	sysMu.Lock()
	defer sysMu.Unlock()
	if id < 0 || id >= len(sys) {
		http.Error(w, "invalid id", 400)
		return
	}
	sys[id].P = v
	sys[id].V = Vec2{}
	w.Write([]byte("ok"))
}

// Toggle run state
func handleControl(w http.ResponseWriter, _ *http.Request) {
	running = !running
	fmt.Fprintf(w, "running=%v\n", running)
}

// Return current state
func handleState(w http.ResponseWriter, _ *http.Request) {
	sysMu.Lock()
	defer sysMu.Unlock()
	nodes := make([]Vec2, len(sys))
	for i := range sys {
		nodes[i] = sys[i].P
	}
	frame := struct {
		Step  int    `json:"step"`
		Nodes []Vec2 `json:"nodes"`
	}{stepCnt, nodes}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(frame)
}

// Run loop
func runLoop() {
	ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
	defer ticker.Stop()
	for range ticker.C {
		if running {
			Step()
			stepCnt++
		}
	}
}

// Anchor control endpoint
func handleAnchors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(Anchors{TopK: anchorTopK, BotK: anchorBotK})
	case http.MethodPost:
		var a Anchors
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "invalid JSON", 400)
			return
		}
		sysMu.Lock()
		if a.TopK > 0 {
			anchorTopK = a.TopK
		}
		if a.BotK > 0 {
			anchorBotK = a.BotK
		}
		sysMu.Unlock()
		fmt.Printf("Updated anchors: top=%v, bottom=%v\n", anchorTopK, anchorBotK)
		w.Write([]byte("ok"))
	default:
		http.Error(w, "method not allowed", 405)
	}
}

// Main
func main() {
	initSystem()
	running = true
	go runLoop()
	http.HandleFunc("/api/state", handleState)
	http.HandleFunc("/api/set/", handleSet)
	http.HandleFunc("/api/control", handleControl)
	http.HandleFunc("/api/anchors", handleAnchors)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("mechanical chain demo running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
