package tag

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterRoutes exposes /api/tag/* endpoints.
func RegisterRoutes(mux *http.ServeMux, sim *Simulation) {
	mux.HandleFunc("/api/tag/state", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(sim.Snapshot())
	})

	mux.HandleFunc("/api/tag/step", func(w http.ResponseWriter, r *http.Request) {
		sim.Step()
		json.NewEncoder(w).Encode(sim.Snapshot())
	})

	mux.HandleFunc("/api/tag/reset", func(w http.ResponseWriter, r *http.Request) {
		sim.Reset()
		json.NewEncoder(w).Encode(sim.Snapshot())
	})

	mux.HandleFunc("/api/tag/params", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var p Params
			json.NewDecoder(r.Body).Decode(&p)
			sim.UpdateParams(p)
		}
		json.NewEncoder(w).Encode(sim.Params())
	})

	mux.HandleFunc("/api/tag/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		for b := range sim.Stream() {
			fmt.Fprintf(w, "data: %s\n\n", b)
			flusher.Flush()
		}
	})
}
