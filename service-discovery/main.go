package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	serverName := os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "Unknown Server"
	}

	// 1. User Traffic Route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! You were routed to %s\n", serverName)
	})

	// 2. Active Health Check Route for Consul
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// 3. Start the server
	fmt.Printf("Starting %s on port 8080...\n", serverName)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
