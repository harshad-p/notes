package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	product := Product{
		Name:  "iPhone Duo",
		Price: 1999.99,
	}
	fmt.Fprintln(w, product)
}

func loggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, elapsed)
	})
}

func versionEmbedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-API-Version", "1")
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", getProductsHandler)
	handler := loggerHandler(versionEmbedHandler(mux))

	http.ListenAndServe(":8080", handler)
}
