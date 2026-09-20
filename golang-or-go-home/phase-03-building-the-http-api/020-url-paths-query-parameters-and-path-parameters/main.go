package main

import (
	"fmt"
	"net/http"
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a request on users endpoint.")
	sortResponse := "Query does not have sort parameter"
	if r.URL.Query().Get("sort") != "" {
		sortResponse = "Query has sort parameter"
	}
	fmt.Fprintln(w, "Method:", r.Method, "Path:", r.URL.Path, "Query:", r.URL.Query(), sortResponse)
}

func main() {
	fmt.Println("Starting http server...")
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/harshad", handleUsers) // need to define a handler again
	http.ListenAndServe(":8080", nil)
}
