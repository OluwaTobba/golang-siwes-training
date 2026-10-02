package main

import (
	"encoding/json"
	"fmt"
	"math"
)

// Base struct
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// Embedding — Employee "has a" Person
type Employee struct {
	Person                    // promoted fields
	Department string `json:"department"`
	Salary     float64 `json:"salary"`
}

// Pointer receiver — modifies the struct
func (e *Employee) Promote(raise float64) {
	e.Salary += raise
}

func (e Employee) String() string {
	return fmt.Sprintf("%s | %s | NGN %.2f",
		e.FullName(), e.Department, e.Salary)
}

// Geometric shapes with methods
type Circle struct {
	Radius float64
}
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

type Rectangle struct {
	Width, Height float64
}
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }

func main() {
	emp := Employee{
		Person:     Person{FirstName: "Michael", LastName: "Ogundipe", Age: 27},
		Department: "DevOps",
		Salary:     450_000,
	}

	fmt.Println(emp)
	emp.Promote(50_000)
	fmt.Println("After promotion:", emp)

	// Access promoted fields directly
	fmt.Println("Full name:", emp.FullName())
	fmt.Println("Age:", emp.Age)

	// JSON marshal
	data, _ := json.MarshalIndent(emp, "", "  ")
	fmt.Println("\nJSON:\n" + string(data))

	// JSON unmarshal
	raw := `{"first_name":"Ada","last_name":"Lovelace","age":36,"department":"Engineering","salary":600000}`
	var emp2 Employee
	json.Unmarshal([]byte(raw), &emp2)
	fmt.Println("\nUnmarshalled:", emp2)

	// Shapes
	c := Circle{Radius: 5}
	r := Rectangle{Width: 8, Height: 4}
	fmt.Printf("\nCircle    area=%.2f  perimeter=%.2f\n", c.Area(), c.Perimeter())
	fmt.Printf("Rectangle area=%.2f  perimeter=%.2f\n", r.Area(), r.Perimeter())
}