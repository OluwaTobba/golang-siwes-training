package main

import "fmt"

func classify(n int) string {
	// if with short initialisation statement
	if rem := n % 2; rem == 0 {
		return "even"
	} else {
		return "odd"
	}
}

func dayName(d int) string {
	switch d {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	default:
		return "Weekend"
	}
}

func main() {
	// Standard for loop
	fmt.Print("Squares: ")
	for i := 1; i <= 5; i++ {
		fmt.Printf("%d ", i*i)
	}
	fmt.Println()

	// While-style for loop
	count := 1
	for count <= 3 {
		fmt.Printf("count = %d\n", count)
		count++
	}

	// Range over a slice
	fruits := []string{"mango", "pineapple", "banana"}
	for idx, fruit := range fruits {
		fmt.Printf("[%d] %s\n", idx, fruit)
	}

	// Range — index only
	for i := range fruits {
		fmt.Printf("index: %d\n", i)
	}

	// if/else
	fmt.Println(classify(7))  // odd
	fmt.Println(classify(10)) // even

	// switch
	for d := 1; d <= 5; d++ {
		fmt.Printf("Day %d: %s\n", d, dayName(d))
	}

	// Multiplication table
	fmt.Println("\n--- 5x Table ---")
	for i := 1; i <= 10; i++ {
		fmt.Printf("5 x %2d = %2d\n", i, 5*i)
	}
}