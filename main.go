package main

import (
	"fmt"
	"net/http"
	"os"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux() // Create a new ServeMux

	server := &http.Server{ // Create a new server with
		Addr:    ":8080", // Address at :8080
		Handler: mux,     // & Handler at mux
	}
	// Create a file server for the current directory
	fileServer := http.FileServer(http.Dir("."))
	// Add health check handler
	mux.HandleFunc("/healthz", healthCheckHandler)
	// Handle Requests to the /app/ path with the filServer
	mux.Handle("/app/", http.StripPrefix("/app", fileServer)) // Handle requests to the root path with the file server

	fmt.Println("Server is starting on http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
