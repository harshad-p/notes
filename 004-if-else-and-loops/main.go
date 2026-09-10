package main

import "fmt"

func adultOrMinor(age int) string {
	if age >= 18 {
		return "Adult"
	} else {
		return "Minor"
	}
}

func main() {
	var age1 int = 15
	var age2 int = 21
	fmt.Println("Age =", age1, "=>", adultOrMinor(age1))
	fmt.Println("Age =", age2, "=>", adultOrMinor(age2))

	// regular for loop demo
	fmt.Println("\nNumbers from 1 to 5:")
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}

	// while loop demo
	fmt.Println("\nEven numbers between from 1 and 10 (inclusive):")
	num := 1
	for num <= 10 {
		if num%2 == 0 {
			fmt.Println(num)
		}
		num++
	}

	// infinite loop
	// fmt.Print("Running")
	// for {
	// 	fmt.Print(".")
	// }
}
