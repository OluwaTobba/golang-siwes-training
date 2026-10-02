package mathutil

import "errors"

func Add(a, b int) int { return a + b }
func Sub(a, b int) int { return a - b }
func Mul(a, b int) int { return a * b }

func Div(a, b float64) (float64, error) {
	if b == 0 { return 0, errors.New("division by zero") }
	return a / b, nil
}

func Fibonacci(n int) int {
	if n <= 1 { return n }
	return Fibonacci(n-1) + Fibonacci(n-2)
}