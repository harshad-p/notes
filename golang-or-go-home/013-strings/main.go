package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "I am learning Golang"

	fmt.Println("Original text:", text)
	fmt.Println("len(): ", len(text))
	fmt.Println("Uppercase:", strings.ToUpper(text))
	fmt.Println("Lowercase:", strings.ToLower(text))
	fmt.Println("Contains Go?: ", strings.Contains(text, "Go"))
}
