package main

import (
	"fmt"
	"net/http"

	"github.com/RickF71/tag-go/internal/tag"
)

func main() {
	http.HandleFunc("/", servePage)
	http.HandleFunc("/stream", handleStream)

	go RunLoop(func(f tag.Frame) { mu.Lock(); latest = f; mu.Unlock() })

	fmt.Println("Demo running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
