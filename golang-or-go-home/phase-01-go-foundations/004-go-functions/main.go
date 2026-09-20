package main

import "fmt"

// Simple function
func sayHello() {
	fmt.Println("Hello!")
}

// Function with parameter
func sayHelloTo(name string) {
	fmt.Println("Hello,", name)
}

// Function with return type
func add(a int, b int) int {
	return a + b
}

// Function with parameter shortening of the same type.
func multiply(a, b int) int {
	return a * b
}

// Multiple Return Values
func getAdditionAndSubtraction(a, b int) (int, int) {
	return a + b, a - b
}

func main() {
	a := 2
	b := 3

	sayHello()
	sayHelloTo("Harshad")
	fmt.Println("Operations on", a, "and", b)
	fmt.Println("Addition:", add(a, b))
	fmt.Println("Multiplication:", multiply(a, b))

	sum, difference := getAdditionAndSubtraction(a, b)
	fmt.Println("Addition and Subtraction:", sum, ",", difference)
}
