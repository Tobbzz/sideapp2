package main

import (
	"log"
	"net/http"
)

func main() {
	// Register API endpoints
	http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/update", updateHandler)

	log.Println("Starting Simrail Side (Go) server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
