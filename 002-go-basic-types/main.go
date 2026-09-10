package main

import "fmt"

func main() {
	// Variables
	var firstName string = "Harshad" // type declared
	var lastName = "Paradkar"        // type not declared
	age := 110                       // short variable declaration

	fmt.Println(firstName, lastName) // Print in a single call
	fmt.Println(age)                 // Print one value only

	isADeveloper := true
	position := "(Senior Software Engineer)"

	fmt.Println("Is a Developer?:", isADeveloper, position)

	const appName string = "Maintenance App" // constants with explicit type declaration
	const maxRetries = "5"                   // constants with no type declaration

}
