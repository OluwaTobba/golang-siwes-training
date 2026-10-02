package main

import (
	"fmt"
	"time"
)

// Generator — sends numbers into a channel
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Stage 1 — square each number
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Stage 2 — filter even numbers only
func evensOnly(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
	}()
	return out
}

// Fan-out — distribute jobs to N workers
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("  worker %d processing job %d\n", id, j)
		time.Sleep(10 * time.Millisecond) // simulate work
		results <- j * 2
	}
}

func main() {
	fmt.Println("=== Pipeline ===")
	// Pipeline: generate → square → evensOnly → print
	c := generate(1, 2, 3, 4, 5, 6, 7, 8)
	squared := square(c)
	evens := evensOnly(squared)

	for v := range evens {
		fmt.Print(v, " ")
	}
	fmt.Println()

	// Buffered channel
	fmt.Println("\n=== Buffered Channel ===")
	buffered := make(chan string, 3)
	buffered <- "first"
	buffered <- "second"
	buffered <- "third"
	close(buffered)
	for msg := range buffered {
		fmt.Println(msg)
	}

	// Fan-out worker pool
	fmt.Println("\n=== Fan-out (3 workers, 6 jobs) ===")
	jobs    := make(chan int, 6)
	results := make(chan int, 6)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}
	for j := 1; j <= 6; j++ { jobs <- j }
	close(jobs)
	for r := 1; r <= 6; r++ {
		fmt.Println("result:", <-results)
	}
}