package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func getDirectorMovieHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	idPathValue := r.PathValue("id")
	id, err := strconv.Atoi(idPathValue)
	if err != nil {
		http.Error(w, "Director id is invalid", http.StatusBadRequest)
		return
	}
	if id <= 0 {
		http.Error(w, "Director id should be a positive non-zero integer", http.StatusBadRequest)
		return
	}

	var sortByName bool

	sortByNameParam := r.URL.Query().Get("sortByName")
	if sortByNameParam != "" {
		sortByName, err = strconv.ParseBool(sortByNameParam)
		if err != nil {
			http.Error(w, "Sort by name query param has an invalid value", http.StatusBadRequest)
			return
		}
	}

	fmt.Fprintln(w, "Movie Id =", id, "sortByName? =", sortByName)
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /directors/{id}/movies", getDirectorMovieHandler)

	http.ListenAndServe(":8080", mux)
}
