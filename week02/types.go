package main

import (
	"fmt"
	"strconv"
)

// Typed constants
const Pi float64 = 3.14159

// Iota enumeration
type Direction int

const (
	North Direction = iota // 0
	East                   // 1
	South                  // 2
	West                   // 3
)

func (d Direction) String() string {
	return [...]string{"North", "East", "South", "West"}[d]
}

func main() {
	// var declaration — explicit type
	var age int = 25
	var name string = "David"
	var isActive bool = true
	var salary float64 = 150_000.50

	// Short declaration — type inferred
	city := "Lagos"
	score := 98.5

	// Zero values (default when not initialised)
	var zeroInt int       // 0
	var zeroStr string    // ""
	var zeroBool bool     // false

	fmt.Printf("Name : %s\n", name)
	fmt.Printf("Age : %d\n", age)
	fmt.Printf("Active : %t\n", isActive)
	fmt.Printf("Salary : %.2f\n", salary)
	fmt.Printf("City : %s\n", city)
	fmt.Printf("Score : %.1f\n", score)
	fmt.Printf("Zeros : %d | %q | %t\n", zeroInt, zeroStr, zeroBool)

	// Type conversion — always explicit in Go
	var x int = 42
	var y float64 = float64(x)
	var z string = strconv.Itoa(x)
	fmt.Printf("int→float64: %f | int→string: %s\n", y, z)

	// Iota directions
	fmt.Printf("Direction: %s (%d)\n", North, North)
	fmt.Printf("Direction: %s (%d)\n", East, East)
}