package main

import (
	"log"
	"net/http"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	// Register API endpoints
	http.HandleFunc("/api", withCORS(apiHandler))
	http.HandleFunc("/update", withCORS(updateHandler))
	// Serve React static files
	http.Handle("/", http.FileServer(http.Dir("./frontend/build")))

	log.Println("Starting Simrail Side (Go) server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
