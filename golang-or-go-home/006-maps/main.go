package main

import "fmt"

func main() {
	favoriteLanguages := map[string]string{
		"Harshad": "C#",
		"Alice":   "Go",
	}

	// Printing map as it is
	fmt.Println(favoriteLanguages)

	// Print Harshad's favorite language
	fmt.Println("Harshad's favorite language:", favoriteLanguages["Harshad"])

	value, exists := favoriteLanguages["Bob"]

	// Print Bob's favorite language
	fmt.Println("Bob's favorite language:", value) //empty-string

	if !exists {
		fmt.Println("Bob's favorite language does not exist in the dictionary")
	}

	favoriteLanguages["Bob"] = "Python"
	fmt.Println("After adding Bob's favorite language:", favoriteLanguages["Bob"])

	favoriteLanguages["Alice"] = "Rust"
	fmt.Println("After modifying Alice's favorite language:", favoriteLanguages["Alice"])

	delete(favoriteLanguages, "Bob")
	fmt.Println("After deleting Bob:", favoriteLanguages)
}
