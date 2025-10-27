package main

import (
	"fmt"
	"net/http"

	"github.com/RickF71/tag-go/internal/tag"
)

func main() {
	sim := tag.NewSimulation()
	mux := http.NewServeMux()
	tag.RegisterRoutes(mux, sim)

	// Serve static files from the web/ directory.
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	fmt.Println("TAG Observatory live at http://localhost:8080/observatory.html")
	http.ListenAndServe(":8080", mux)
}
