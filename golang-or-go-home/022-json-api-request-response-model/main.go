package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateProductRequest struct {
	Name  *string  `json:"name"`
	Price *float64 `json:"price"`
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var createProductRequest CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&createProductRequest)
	if err != nil {
		http.Error(w, "Could not decode create product request body.", http.StatusBadRequest)
		return
	}

	if createProductRequest.Name == "" {
		http.Error(w, "Product name is required.", http.StatusBadRequest)
		return
	}

	if createProductRequest.Price <= 0 {
		http.Error(w, "Product price should be a positive non-zero value", http.StatusBadRequest)
		return
	}

	createProductResponse := CreateProductResponse{
		ID:    1,
		Name:  createProductRequest.Name,
		Price: createProductRequest.Price,
	}

	w.Header().Set("Content-Type", "application/json")
	createProductResponseJson := json.NewEncoder(w).Encode(createProductResponse)
	fmt.Fprintln(w, createProductResponseJson)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	var updateProductRequest UpdateProductRequest

	err := json.NewDecoder(r.Body).Decode(&updateProductRequest)
	if err != nil {
		http.Error(w, "Could not decode create product request body.", http.StatusBadRequest)
		return
	}

	if updateProductRequest.Name == nil {
		fmt.Fprintln(w, "Will not update Product Name")
	} else {
		fmt.Fprintln(w, "Will update Product Name to", *updateProductRequest.Name)
	}

	if updateProductRequest.Price == nil {
		fmt.Fprintln(w, "Will not update Product Price")
	} else {
		fmt.Fprintln(w, "Will update Product Price to", *updateProductRequest.Price)
	}
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /products", createProductHandler)
	mux.HandleFunc("PUT /products", updateProductHandler)

	http.ListenAndServe(":8080", mux)
}
