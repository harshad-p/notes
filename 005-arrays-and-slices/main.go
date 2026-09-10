package main

import "fmt"

func main() {

	// ARRAYS

	daysOfWeek := [7]string{}

	daysOfWeek[0] = "Monday"

	// Direct printing of array elements
	fmt.Println("Days of week:")
	fmt.Println(daysOfWeek)

	// SLICES

	numbers := []int{1, 2, 3, 4, 5}

	// Direct printing of slice elements
	fmt.Println("Numbers:")
	fmt.Println(numbers)

	movies := []string{"Terminator", "Avengers", "True Lies"}

	// Using for loop with index to print array elements
	fmt.Println("\nMovies:")
	for i := 0; i < len(movies); i++ {
		fmt.Println(movies[i])
	}

	// Add a movie
	movies = append(movies, "Tenet")

	// The "Go way" - Using range
	fmt.Println("\nMovies after appending a value:")
	for index, value := range movies {
		fmt.Println(index, value)
	}
}
