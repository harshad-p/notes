package main

import (
	"fmt"
	"net/http"
)

func movieHandler(w http.ResponseWriter, r *http.Request) {
	movies := "Terminator, True Lies, Avatar"

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Received GET Movies request.")
		fmt.Fprintln(w, movies)
	case http.MethodPost:
		fmt.Println("Received POST Movies request.")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, movies+", "+"Avengers")
	default:
		fmt.Println("Received a Movies request with unsupported method.")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Received GET Users request.")
		fmt.Fprintln(w, "Harshad")
	default:
		fmt.Println("Received a Users request with unsupported method.")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	fmt.Println("Starting http server...")
	http.HandleFunc("/movies", movieHandler)
	http.HandleFunc("/users", userHandler)
	http.ListenAndServe(":8080", nil)
}
