package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
)

// Custom time format for JSON
type JSONDate time.Time

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}
func (d *JSONDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil { return err }
	t, err := time.Parse("2006-01-02", s)
	if err != nil { return err }
	*d = JSONDate(t)
	return nil
}

type Employee struct {
	Name      string    `json:"name"       xml:"name"`
	Role      string    `json:"role"       xml:"role"`
	StartDate JSONDate  `json:"start_date" xml:"-"`
	Salary    *float64  `json:"salary,omitempty" xml:"salary,attr"`
}

// XML feed structure
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Items   []Item   `xml:"item"`
}
type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func main() {
	// JSON round-trip with custom date and omitempty
	salary := 500000.0
	emp := Employee{
		Name:      "David Ogundipe",
		Role:      "DevOps Engineer",
		StartDate: JSONDate(time.Now()),
		Salary:    &salary,
	}
	data, _ := json.MarshalIndent(emp, "", "  ")
	fmt.Println("=== JSON ===\n" + string(data))

	// omitempty — nil pointer omitted
	emp2 := Employee{Name: "Intern", Role: "Trainee"}
	data2, _ := json.MarshalIndent(emp2, "", "  ")
	fmt.Println("\n=== JSON (omitempty) ===\n" + string(data2))

	// Unmarshal
	raw := `{"name":"Ada","role":"Engineer","start_date":"2024-01-15","salary":650000}`
	var e Employee
	json.Unmarshal([]byte(raw), &e)
	fmt.Printf("\nUnmarshalled: %+v\n", e)
	if e.Salary != nil { fmt.Printf("Salary: %.2f\n", *e.Salary) }

	// XML decode
	xmlData := `
<feed>
  <item><title>Go 1.22 Released</title><link>https://go.dev/blog</link><pubDate>Feb 2024</pubDate></item>
  <item><title>gRPC Best Practices</title><link>https://grpc.io</link><pubDate>Mar 2024</pubDate></item>
</feed>`
	var feed Feed
	xml.Unmarshal([]byte(xmlData), &feed)
	fmt.Println("\n=== XML Feed ===")
	for _, item := range feed.Items {
		fmt.Printf("  [%s] %s\n", item.PubDate, item.Title)
	}
}