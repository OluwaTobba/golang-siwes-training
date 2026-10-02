package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Student struct {
	Name   string  `json:"name"`
	Score  float64 `json:"score"`
	Grade  string  `json:"grade"`
}

func grade(score float64) string {
	switch {
	case score >= 70: return "A"
	case score >= 60: return "B"
	case score >= 50: return "C"
	default:          return "F"
	}
}

func main() {
	// Write a text file with bufio
	f, _ := os.Create("output.txt")
	w := bufio.NewWriter(f)
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(w, "Line %d: Hello from Go\n", i)
	}
	w.Flush()
	f.Close()

	// Read it back line by line
	f, _ = os.Open("output.txt")
	scanner := bufio.NewScanner(f)
	fmt.Println("=== output.txt ===")
	for scanner.Scan() { fmt.Println(scanner.Text()) }
	f.Close()

	// JSON file: write then read
	students := []Student{
		{"Ada Lovelace", 85.5, ""},
		{"Alan Turing", 92.0, ""},
		{"Grace Hopper", 78.3, ""},
	}
	for i := range students { students[i].Grade = grade(students[i].Score) }

	jf, _ := os.Create("students.json")
	enc := json.NewEncoder(jf)
	enc.SetIndent("", "  ")
	enc.Encode(students)
	jf.Close()

	var loaded []Student
	jf, _ = os.Open("students.json")
	json.NewDecoder(jf).Decode(&loaded)
	jf.Close()
	fmt.Println("\n=== students.json ===")
	for _, s := range loaded { fmt.Printf("%-15s %.1f %s\n", s.Name, s.Score, s.Grade) }

	// CSV: write
	cf, _ := os.Create("students.csv")
	cw := csv.NewWriter(cf)
	cw.Write([]string{"Name", "Score", "Grade"})
	for _, s := range loaded {
		cw.Write([]string{s.Name, fmt.Sprintf("%.1f", s.Score), s.Grade})
	}
	cw.Flush()
	cf.Close()

	// Directory walk
	fmt.Println("\n=== .go files in current dir ===")
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil { return err }
		if filepath.Ext(path) == ".go" { fmt.Println(" ", path) }
		return nil
	})

	os.Remove("output.txt")
	os.Remove("students.json")
	os.Remove("students.csv")
}