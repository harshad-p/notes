package main

import (
	"errors"
	"fmt"
)

func getUser(id int) (string, error) {
	users := map[int]string{
		1: "harshad",
		2: "arnold",
	}

	user, exists := users[id]

	if exists {
		return user, nil
	} else {
		// For a simple fixed message
		return "", errors.New("User not found")
	}
}

func divide(dividend, divisor float32) (float32, error) {
	if divisor == 0 {
		// fmt.Errorf becomes particularly useful when you want to include values
		return 0, fmt.Errorf("Cannot divide %f by 0.", dividend)
	}

	return dividend / divisor, nil
}

func printDivisionResult(dividend, divisor, result float32, err error) {
	if err == nil {
		fmt.Printf("%f / %f = %f\n", dividend, divisor, result)
	} else {
		fmt.Println(err)
	}
}

func printUserResult(user string, err error) {
	if err == nil {
		fmt.Println(user)
	} else {
		fmt.Println(err)
	}
}

func main() {
	user_1, err_1 := getUser(1)
	printUserResult(user_1, err_1)

	user_100, err_100 := getUser(100)
	printUserResult(user_100, err_100)

	result1, error1 := divide(4, 2)
	printDivisionResult(4, 2, result1, error1)

	result2, error2 := divide(4, 0)
	printDivisionResult(4, 2, result2, error2)
}
