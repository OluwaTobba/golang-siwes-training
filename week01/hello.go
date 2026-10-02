package main

import (
	"fmt"
	"runtime"
)

func run() {
	// Print a greeting
	fmt.Println("Hello, World!")

	// Explore the runtime package
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
}