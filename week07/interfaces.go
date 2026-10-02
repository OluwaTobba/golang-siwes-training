package main

import (
	"fmt"
	"math"
)

// Interface definition
type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// Concrete types
type Circle struct{ Radius float64 }
func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }
func (c Circle) Name() string       { return "Circle" }

type Rectangle struct{ Width, Height float64 }
func (r Rectangle) Area() float64      { return r.Width * r.Height }
func (r Rectangle) Perimeter() float64 { return 2 * (r.Width + r.Height) }
func (r Rectangle) Name() string       { return "Rectangle" }

type Triangle struct{ A, B, C float64 }
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }
func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}
func (t Triangle) Name() string { return "Triangle" }

// Polymorphic function — accepts any Shape
func printShape(s Shape) {
	fmt.Printf("%-12s area=%-10.2f perimeter=%.2f\n",
		s.Name(), s.Area(), s.Perimeter())
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes { total += s.Area() }
	return total
}

// Type switch — inspect concrete type at runtime
func describe(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("int: %d\n", v)
	case string:
		fmt.Printf("string: %q\n", v)
	case Circle:
		fmt.Printf("Circle with radius %.2f\n", v.Radius)
	case bool:
		fmt.Printf("bool: %t\n", v)
	default:
		fmt.Printf("unknown type: %T\n", v)
	}
}

func main() {
	shapes := []Shape{
		Circle{Radius: 5},
		Rectangle{Width: 8, Height: 4},
		Triangle{A: 3, B: 4, C: 5},
	}

	for _, s := range shapes {
		printShape(s)
	}
	fmt.Printf("Total area: %.2f\n", totalArea(shapes))

	// Type assertion
	var s Shape = Circle{Radius: 7}
	if c, ok := s.(Circle); ok {
		fmt.Printf("\nRadius of circle: %.2f\n", c.Radius)
	}

	// Type switch
	fmt.Println()
	describe(42)
	describe("hello")
	describe(Circle{Radius: 3})
	describe(true)
}