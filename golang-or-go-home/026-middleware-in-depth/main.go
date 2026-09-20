package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type MyResponseWriter struct {
	http.ResponseWriter
	statusCode      int
	isStatusCodeSet bool
}

func (w *MyResponseWriter) WriteHeader(statusCode int) {
	if w.isStatusCodeSet {
		// HTTP uses the first status code written for the response.
		// Later calls to WriteHeader do not change the response status.
		//
		// We return here as well so our wrapper keeps the first status code,
		// allowing middleware such as a logger to report the actual status.
		return
	}

	w.statusCode = statusCode
	w.isStatusCodeSet = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *MyResponseWriter) Write(data []byte) (int, error) {
	if !w.isStatusCodeSet {
		// If no status has been explicitly written, the underlying
		// http.ResponseWriter automatically commits the response with 200 OK
		// when the body is written. This does not mean Write() calls the
		// public WriteHeader() method; the underlying implementation handles
		// the implicit 200 internally. We call WriteHeader() explicitly here
		// so our wrapper records the status as well.
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(data)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		myResponseWriter := &MyResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(myResponseWriter, r)

		elapsed := time.Since(start)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, myResponseWriter.statusCode, elapsed)
	})
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Howdy!")
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	product := Product{
		Name:  "iPhone Duo",
		Price: 1999.99,
	}
	json.NewEncoder(w).Encode(product)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid product details", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func main() {
	fmt.Println("Starting http server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", homepageHandler)
	mux.HandleFunc("GET /products", getProductsHandler)
	mux.HandleFunc("POST /products", createProductHandler)
	handler := loggerMiddleware(mux)

	http.ListenAndServe(":8080", handler)
}
