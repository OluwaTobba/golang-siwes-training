package mathutil

import (
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"positive", 2, 3, 5},
		{"negative", -1, -2, -3},
		{"zero",     0, 5, 5},
		{"mixed",   -3, 7, 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Add(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Add(%d,%d) = %d; want %d", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{"normal",     10, 2,  5.0, false},
		{"decimal",    7,  2,  3.5, false},
		{"by zero",    5,  0,  0,   true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Div(tc.a, tc.b)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Div(%v,%v) error=%v; wantErr=%v", tc.a, tc.b, err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Errorf("Div(%v,%v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// Benchmark
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

// Run: go test ./... -v
// Run: go test -bench=. -benchmem