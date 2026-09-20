package main

import (
	"fmt"
	"net/http"
)

func getMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	fmt.Fprintln(w, "Terminator 2: Judgement Day, Avatar")
}

func getMovieByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.PathValue("id") {
	case "1":
		fmt.Fprintln(w, "Terminator 2: Judgement Day")
	case "2":
		fmt.Fprintln(w, "Avatar")
	default:
		fmt.Fprintln(w, "Movie not found")
	}
}

func createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Created a movie.")
}

func deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Deleted a movie")
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", getMoviesHandler)
	mux.HandleFunc("GET /movies/{id}", getMovieByIdHandler)
	mux.HandleFunc("POST /movies", createMovieHandler)
	mux.HandleFunc("DELETE /movies/{id}", deleteMovieHandler)

	http.ListenAndServe(":8080", mux)
}
