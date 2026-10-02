package main

import (
	"fmt"
	"sort"
)

func main() {
	// ── Arrays (fixed size, value type) ──────────────────────
	var arr [5]int
	arr[0] = 10
	arr[4] = 50
	fmt.Println("Array:", arr)

	primes := [5]int{2, 3, 5, 7, 11}
	fmt.Println("Primes:", primes)

	// ── Slices (dynamic, reference type) ─────────────────────
	s := []string{"go", "python", "rust"}
	s = append(s, "zig")
	fmt.Println("Slice:", s, "len:", len(s), "cap:", cap(s))

	// Slice of a slice — shares backing array!
	sub := s[1:3]
	fmt.Println("Sub-slice:", sub)

	// make — specify len and cap
	scores := make([]int, 0, 10)
	for i := 1; i <= 5; i++ {
		scores = append(scores, i*i)
	}
	fmt.Println("Scores:", scores)

	// 2D slice (matrix)
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	for _, row := range matrix {
		fmt.Println(row)
	}

	// ── Maps (key-value, reference type) ─────────────────────
	wordCount := make(map[string]int)
	sentence := []string{"go", "is", "fast", "go", "is", "fun", "go"}
	for _, w := range sentence {
		wordCount[w]++
	}

	// Comma-ok idiom — safe key lookup
	if count, ok := wordCount["go"]; ok {
		fmt.Printf("'go' appears %d times\n", count)
	}
	if _, ok := wordCount["java"]; !ok {
		fmt.Println("'java' not found")
	}

	// Delete a key
	delete(wordCount, "is")
	fmt.Println("After delete:", wordCount)

	// Sort map keys for deterministic output
	keys := make([]string, 0, len(wordCount))
	for k := range wordCount { keys = append(keys, k) }
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, wordCount[k])
	}
}