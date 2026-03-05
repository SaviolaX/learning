package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	firstNum  int
	secondNum int
	sign      string
}

func main() {

	router := chi.NewRouter()

	router.Get("/", computationHandler)

	// Lauch server
	port := ":8080"
	println("Starting server on port:", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		println("Error starting server:", err)
	}
}

func computationHandler(w http.ResponseWriter, r *http.Request) {

	first, err := strconv.Atoi(r.URL.Query().Get("firstNum"))
	if err != nil {
		http.Error(w, "invalid firstNum", http.StatusBadRequest)
		return
	}
	second, err := strconv.Atoi(r.URL.Query().Get("secondNum"))
	if err != nil {
		http.Error(w, "invalid secondNum", http.StatusBadRequest)
		return
	}

	sign := r.URL.Query().Get("sign")

	switch sign {
	case "plus":
		add(first, second, w)
	case "minus":
		minus(first, second, w)
	case "multiply":
		mltply(first, second, w)
	case "divide":
		dvde(first, second, w)
	default:
		w.Write([]byte("nothing\n"))
	}
}

func add(first, second int, w http.ResponseWriter) {
	result := first + second
	fmt.Fprintf(w, "%d\n", result)
}

func minus(first, second int, w http.ResponseWriter) {
	result := first - second
	fmt.Fprintf(w, "%d\n", result)
}

func dvde(first, second int, w http.ResponseWriter) {
	result := first / second
	fmt.Fprintf(w, "%d\n", result)
}

func mltply(first, second int, w http.ResponseWriter) {
	result := first * second
	fmt.Fprintf(w, "%d\n", result)
}
