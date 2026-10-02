// Package main provides the application entry point.
package main

import "errors"

// Round rounds f to d decimal places.
func Round(f float64, d int) float64 {
	pow := 1.0
	for i := 0; i < d; i++ { pow *= 10 }
	return float64(int(f*pow+0.5)) / pow
}

// Percentage returns what percent part is of total.
func Percentage(part, total float64) (float64, error) {
	if total == 0 { return 0, errors.New("total cannot be zero") }
	return (part / total) * 100, nil
}