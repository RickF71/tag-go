package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RickF71/tag-go/internal/tag"

	"sync"
)

var latest tag.Frame
var params = tag.Params{
	Viscosity: 0.5,
	Limit:     5.0,
	Dt:        0.05,
}
var mu sync.RWMutex

// Handle slider updates
func handleParams(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var p tag.Params
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		mu.Lock()
		params = p
		mu.Unlock()
		fmt.Printf("[demo] updated params: %+v\n", p)
		w.WriteHeader(http.StatusOK)
		return
	}

	// GET -> return current params
	mu.RLock()
	defer mu.RUnlock()
	json.NewEncoder(w).Encode(params)
}

// Serve starts the HTTP demo server at / and /stream.
func Serve() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("/stream", handleStream)
	go RunLoop(func(f tag.Frame) { latest = f })
	fmt.Println("Demo running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(200 * time.Millisecond) // 5 FPS
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			mu.RLock()
			frame := latest
			mu.RUnlock()

			data, _ := json.Marshal(frame)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func servePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "examples/mechchain/view.html")
}
