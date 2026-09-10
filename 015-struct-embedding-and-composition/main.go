package main

import "fmt"

type Person struct {
	Name         string
	Organization // no name field. direct type
}

type Organization struct {
	Name string
	Size int
}

func (org Organization) OrganizationToFriendlyName() string {
	return org.Name + " has org.Size employee(s)."
}

func main() {
	harshad := Person{
		Name: "Harshad Paradkar",
		Organization: Organization{
			Name: "Skynet",
			Size: 1000,
		},
	}

	fmt.Println(harshad.Name)
	fmt.Println(harshad.Organization.Name)            // no promotion
	fmt.Println(harshad.Size)                         // promotion
	fmt.Println(harshad.OrganizationToFriendlyName()) // promotion
}
