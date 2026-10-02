package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WaitGroup
func runWithWaitGroup() {
	var wg sync.WaitGroup
	results := make([]string, 5)
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id*10) * time.Millisecond)
			mu.Lock()
			results[id] = fmt.Sprintf("goroutine-%d done", id)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	for _, r := range results { fmt.Println(" ", r) }
}

// select
func selectDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() { time.Sleep(20 * time.Millisecond); ch1 <- "from ch1" }()
	go func() { time.Sleep(10 * time.Millisecond); ch2 <- "from ch2" }()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("select received:", msg)
		case msg := <-ch2:
			fmt.Println("select received:", msg)
		}
	}
}

// context cancellation
func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("  worker %d: context cancelled (%v)\n", id, ctx.Err())
			return
		case <-time.After(30 * time.Millisecond):
			fmt.Printf("  worker %d: tick\n", id)
		}
	}
}

func runWithContext() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, &wg)
	}
	wg.Wait()
}

// sync.Once
var (
	instance string
	once     sync.Once
)
func getInstance() string {
	once.Do(func() {
		fmt.Println("  initialising singleton...")
		instance = "SingletonService-v1"
	})
	return instance
}

func main() {
	fmt.Println("=== WaitGroup + Mutex ===")
	runWithWaitGroup()

	fmt.Println("\n=== select ===")
	selectDemo()

	fmt.Println("\n=== context.WithTimeout ===")
	runWithContext()

	fmt.Println("\n=== sync.Once ===")
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); fmt.Println(" ", getInstance()) }()
	}
	wg.Wait()
}