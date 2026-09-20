package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Pages  int
}

func (book Book) ToString() {
	fmt.Println("Title\t\t:\t", book.Title)
	fmt.Println("Author\t\t:\t", book.Author)
	fmt.Println("Pages\t\t:\t", book.Pages)
}

func (book Book) NotPossibleToResetPages() {
	book.Pages = 0
}

func (book *Book) ResetPages() {
	book.Pages = 0
}

type Cat struct {
	Name  string
	Color string
}

func (c Cat) speak() {
	fmt.Println("Meow")
}

type Dog struct {
	Name  string
	Color string
}

func (d Dog) speak() {
	fmt.Println("Woof")
}

type Speaker interface {
	speak()
}

func makeSound(s Speaker) {
	s.speak()
}

func main() {
	// STRUCTS

	harryPotterBook := Book{
		Title:  "Harry Potter and the Goblet of Fire",
		Author: "J. K. Rowling",
		Pages:  100,
	}

	// Print a struct directly
	fmt.Println(harryPotterBook)

	// Print a struct field
	fmt.Println(harryPotterBook.Pages)

	// Modify value of a field
	harryPotterBook.Pages = 700
	fmt.Println(harryPotterBook.Pages)

	// Use connected method to print Book
	harryPotterBook.ToString()

	// Try to modify pages using instance copy
	harryPotterBook.NotPossibleToResetPages()
	fmt.Println("After trying to modify pages using instance copy:")
	harryPotterBook.ToString()

	// Try to modify pages using pointer
	harryPotterBook.ResetPages()
	fmt.Println("After trying to modify pages using pointer:")
	harryPotterBook.ToString()

	// INTERFACES
	cat := Cat{
		Name:  "Destiny",
		Color: "Orange and White",
	}

	dog := Dog{
		Name:  "Max",
		Color: "Brown",
	}

	cat.speak()
	dog.speak()
}
