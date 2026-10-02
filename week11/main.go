package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	// strings
	s := "  Hello, Go World!  "
	fmt.Println(strings.TrimSpace(s))
	fmt.Println(strings.ToUpper(s))
	fmt.Println(strings.Contains(s, "Go"))
	fmt.Println(strings.ReplaceAll(s, "Go", "Golang"))
	parts := strings.Split("a,b,c,d", ",")
	fmt.Println(parts)
	fmt.Println(strings.Join(parts, " | "))

	// sort
	nums := []int{5, 2, 8, 1, 9, 3}
	sort.Ints(nums)
	fmt.Println("sorted:", nums)

	names := []string{"Zara", "Alice", "Mike", "Bob"}
	sort.Strings(names)
	fmt.Println("sorted:", names)

	// Custom sort — by length
	sort.Slice(names, func(i, j int) bool {
		return len(names[i]) < len(names[j])
	})
	fmt.Println("by length:", names)

	// math
	fmt.Printf("Pi=%.5f  Sqrt(2)=%.5f  Pow(2,10)=%.0f\n",
		math.Pi, math.Sqrt(2), math.Pow(2, 10))

	// time
	now := time.Now()
	fmt.Println("Now:", now.Format("2006-01-02 15:04:05"))
	tomorrow := now.Add(24 * time.Hour)
	fmt.Println("Tomorrow:", tomorrow.Format("Mon, 02 Jan 2006"))
	diff := tomorrow.Sub(now)
	fmt.Printf("Diff: %.1f hours\n", diff.Hours())

	// os
	hostname, _ := os.Hostname()
	fmt.Println("Hostname:", hostname)
	fmt.Println("Env PATH:", os.Getenv("PATH")[:40]+"...")

	// log
	logger := log.New(os.Stdout, "[APP] ", log.LstdFlags)
	logger.Println("Application started")
	logger.Printf("Running on host: %s\n", hostname)
}