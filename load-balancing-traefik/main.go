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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from %s!\n", serverName)
	})

	fmt.Printf("Starting %s on port 8080...\n", serverName)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server crashed:", err)
	}
}
