package main

import (
	"fmt"
	"packages-and-modules/math"
)

func main() {
	fmt.Println(getGreetingMessage())

	a := 5
	b := 10

	fmt.Println("Max of", a, "and", b, "is", math.Max(a, b))
}
