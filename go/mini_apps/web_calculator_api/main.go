package main

import "net/http"

func main() {

	// Lauch server
	port := ":8080"
	println("Starting server on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		println("Error starting server:", err)
	}
}
