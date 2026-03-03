package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 100 Go Mistakes #77: JSON handling pitfalls
//
// Solution: Properly namespace embedded fields, implement custom marshalers,
// use json.Number for precision, and understand omitempty with zero values.

func main() {
	fmt.Println("=== JSON Mistakes ===")
	fmt.Println()

	testEmbeddedFields()
	fmt.Println()

	testTimeHandling()
	fmt.Println()

	testNumericPrecision()
	fmt.Println()

	testOmitEmpty()
	fmt.Println()

	fmt.Println("All JSON checks passed!")
}

// --- Embedded Fields ---

type Address struct {
	City string `json:"city"`
}

// FIXED: Give the embedded field a JSON tag to nest it under "address"
// instead of promoting its fields to the top level.
// Also removed the shadowing City field.
type Event struct {
	Name    string  `json:"name"`
	Address Address `json:"address"` // FIXED: Named field with json tag, not embedded
	City    string  `json:"city"`
}

func testEmbeddedFields() {
	fmt.Println("--- Embedded Fields ---")

	e := Event{
		Name:    "GoConf",
		Address: Address{City: "Portland"},
		City:    "Seattle",
	}

	data, _ := json.Marshal(e)
	jsonStr := string(data)
	fmt.Printf("Marshaled: %s\n", jsonStr)

	if strings.Contains(jsonStr, `"address"`) && strings.Contains(jsonStr, `"Portland"`) {
		fmt.Println("PASS: address city is properly nested")
	} else {
		fmt.Println("FAIL: address city should be in nested 'address' object")
	}
}

// --- Time Handling ---

// FIXED: Use a custom DateOnly type with MarshalJSON/UnmarshalJSON
// to output just the date portion in YYYY-MM-DD format.
type DateOnly struct {
	time.Time
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Time.Format("2006-01-02"))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// FIXED: ScheduledDate uses DateOnly instead of time.Time
type Meeting struct {
	Title         string   `json:"title"`
	ScheduledDate DateOnly `json:"scheduled_date"`
}

func testTimeHandling() {
	fmt.Println("--- Time Handling ---")

	m := Meeting{
		Title:         "Sprint Review",
		ScheduledDate: DateOnly{time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC)},
	}

	data, _ := json.Marshal(m)
	jsonStr := string(data)
	fmt.Printf("Marshaled: %s\n", jsonStr)

	if strings.Contains(jsonStr, `"2025-01-15"`) && !strings.Contains(jsonStr, `T14:30`) {
		fmt.Println("PASS: date formatted as YYYY-MM-DD")
	} else {
		fmt.Println("FAIL: date should be '2025-01-15', not full RFC3339")
	}
}

// --- Numeric Precision ---

func testNumericPrecision() {
	fmt.Println("--- Numeric Precision ---")

	// FIXED: Use json.NewDecoder with UseNumber() to preserve numeric precision.
	// json.Number keeps the original string representation, avoiding float64 loss.

	raw := `{"id": 9999999999999999, "amount": 123.456}`

	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()

	var result map[string]interface{}
	decoder.Decode(&result)

	idVal := result["id"].(json.Number).String()
	fmt.Printf("Decoded ID: %s\n", idVal)

	if idVal == "9999999999999999" {
		fmt.Println("PASS: ID preserved with full precision")
	} else {
		fmt.Println("FAIL: ID lost precision (got " + idVal + ")")
	}
}

// --- OmitEmpty Behavior ---

// FIXED: Removed omitempty from Age, Score, and Verified so zero values are preserved.
// Fields where "not set" vs "zero" distinction matters should use pointer types.
type UserProfile struct {
	Name     string  `json:"name,omitempty"`
	Email    string  `json:"email,omitempty"`
	Age      int     `json:"age"`               // FIXED: removed omitempty
	Score    float64 `json:"score"`              // FIXED: removed omitempty
	Verified bool    `json:"verified"`           // FIXED: removed omitempty
	Bio      *string `json:"bio,omitempty"`      // Pointer: nil=absent, ""=present
}

func testOmitEmpty() {
	fmt.Println("--- OmitEmpty Behavior ---")

	user := UserProfile{
		Name:     "Alice",
		Age:      0,
		Score:    0.0,
		Verified: false,
	}

	data, _ := json.Marshal(user)
	jsonStr := string(data)
	fmt.Printf("Marshaled: %s\n", jsonStr)

	checks := 0
	if strings.Contains(jsonStr, `"age"`) {
		fmt.Println("PASS: age field present (zero value preserved)")
		checks++
	} else {
		fmt.Println("FAIL: age field missing (omitempty dropped zero)")
	}

	if strings.Contains(jsonStr, `"score"`) {
		fmt.Println("PASS: score field present (zero value preserved)")
		checks++
	} else {
		fmt.Println("FAIL: score field missing (omitempty dropped zero)")
	}

	if strings.Contains(jsonStr, `"verified"`) {
		fmt.Println("PASS: verified field present (false preserved)")
		checks++
	} else {
		fmt.Println("FAIL: verified field missing (omitempty dropped false)")
	}

	if checks == 3 {
		fmt.Println("PASS: all zero-valued fields preserved correctly")
	}
}
