package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	harshad := User{
		Name: "Harshad Paradkar",
		Age:  110,
	}

	fmt.Println(harshad)

	var pointerWithExplicitType *User = &harshad
	pointerWithExplicitType.Age = 500
	fmt.Println(harshad) // value changed

	pointerWithImplicitType := &harshad  // type is *User
	fmt.Println(pointerWithImplicitType) // Prints with an '&'
	fmt.Println("Name:", (*pointerWithExplicitType).Name,
		"\nAge:", pointerWithExplicitType.Age) // Automatic dereference
}
