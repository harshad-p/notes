package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my first web server developed using Golang.")
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {

	product := Product{
		Name:  "iPhone Duo",
		Price: 1999,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	responseText := encoder.Encode(product)
	fmt.Fprintln(w, responseText)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var product Product

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&product)
	if err != nil {
		http.Error(w, "Error decoding payload", http.StatusBadRequest)
		return
	}

	fmt.Println(product)

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	responseText := encoder.Encode(product)
	fmt.Fprintln(w, responseText)
}

func main() {
	fmt.Println("Starting http server...")

	http.HandleFunc("/", homepageHandler)
	http.HandleFunc("/products", getProductsHandler)
	http.HandleFunc("/products/create", createProductHandler)
	http.ListenAndServe(":8080", nil)
}
