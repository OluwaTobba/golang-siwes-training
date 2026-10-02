package main

import (
	"errors"
	"fmt"
)

// Multiple return values — result + error
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Named return values
func minMax(nums []int) (min, max int) {
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min { min = n }
		if n > max { max = n }
	}
	return // bare return — returns named values
}

// Variadic function
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func average(nums ...int) float64 {
	if len(nums) == 0 { return 0 }
	return float64(sum(nums...)) / float64(len(nums))
}

// Defer — runs LIFO on function exit
func deferDemo() {
	fmt.Println("start")
	defer fmt.Println("first defer  (runs last)")
	defer fmt.Println("second defer (runs first)")
	fmt.Println("end")
}

func main() {
	// Multiple returns
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	_, err = divide(5, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Named returns
	nums := []int{4, 1, 9, 2, 7}
	lo, hi := minMax(nums)
	fmt.Printf("min=%d  max=%d\n", lo, hi)

	// Variadic
	fmt.Println("sum(1,2,3)  =", sum(1, 2, 3))
	fmt.Println("sum(1..5)   =", sum(1, 2, 3, 4, 5))
	fmt.Printf("average     = %.2f\n", average(10, 20, 30, 40))

	// Spread slice into variadic
	s := []int{5, 10, 15}
	fmt.Println("sum(slice)  =", sum(s...))

	// Defer
	deferDemo()
}