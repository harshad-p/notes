package main

import (
	"encoding/json"
	"fmt"
)

type Product struct {
	Name    string
	Price   float64 `json:"price"`
	Website string  `json:"website,omitempty"`
}

func main() {
	iPhoneFold := Product{
		Name:    "iPhone Fold",
		Price:   2500,
		Website: "",
	}

	fmt.Println(iPhoneFold)

	serializedJson, serializationError := json.Marshal(iPhoneFold)
	if serializationError != nil {
		panic(serializationError)
	}

	fmt.Println(string(serializedJson)) // Name, price, Website omitted

	var deserializedProduct Product
	deserializationError := json.Unmarshal(serializedJson, &deserializedProduct)
	if deserializationError != nil {
		panic(deserializationError)
	}

	fmt.Println(deserializedProduct)
}
