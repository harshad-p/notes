package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	product := Product{
		ID:    1,
		Name:  "iPhone Duo",
		Price: 1999.0,
	}

	sortParam := r.URL.Query().Get("sortBy")
	if sortParam == "" {
		fmt.Println("No sort by query param provided")
	} else {
		fmt.Println("Sorting by:", sortParam)
	}

	writeResponse(w, http.StatusOK, product)
}

func getProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	product := Product{
		ID:    1,
		Name:  "iPhone Duo",
		Price: 1999.0,
	}

	id := r.PathValue("id")
	if id != "1" {
		writeErrorResponse(w, http.StatusBadRequest, "Product with Id "+id+" does not exist")
		return
	}
	path := r.URL.Path
	fmt.Println("Received request on path:", path, "for id:", id)
	writeResponse(w, http.StatusOK, product)
}

func writeResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid Product json body")
		return
	}

	product.ID = 2
	writeResponse(w, http.StatusCreated, product)
}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	errorResponse := ErrorResponse{
		Error: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", getProductsHandler)
	mux.HandleFunc("GET /products/{id}", getProductByIdHandler)
	mux.HandleFunc("POST /products", createProductHandler)
	mux.HandleFunc("DELETE /products", deleteProductHandler)

	http.ListenAndServe(":8080", mux)
}
