package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// 100 Go Mistakes #77: JSON handling pitfalls
//
// This exercise covers common JSON mistakes in Go:
// 1. Embedded struct field promotion and duplicate key confusion
// 2. time.Time default marshaling (RFC3339) vs custom format
// 3. json.Number for precise numeric handling vs float64 precision loss
// 4. omitempty behavior with zero values, nil, and empty strings

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

// BUG: When Event embeds both Address and has its own City field,
// the JSON behavior is confusing. If both the embedded and outer struct
// have the same JSON key, the outer field wins but it can cause subtle bugs.
//
// TODO: Fix the struct so that the Address city is properly namespaced
// under an "address" key instead of being promoted to top level.
type Event struct {
	Name string `json:"name"`
	Address
	City string `json:"city"` // This shadows Address.City!
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

	// The address city should be in a nested object, not promoted
	if strings.Contains(jsonStr, `"address"`) && strings.Contains(jsonStr, `"Portland"`) {
		fmt.Println("PASS: address city is properly nested")
	} else {
		fmt.Println("FAIL: address city should be in nested 'address' object")
	}
}

// --- Time Handling ---

type Meeting struct {
	Title string `json:"title"`
	// BUG: time.Time marshals to RFC3339 by default (e.g., "2025-01-15T14:30:00Z").
	// But our API expects "2025-01-15" (date only, no time component).
	//
	// TODO: Change ScheduledDate to a custom type that implements
	// json.Marshaler and json.Unmarshaler to output just the date portion.
	ScheduledDate time.Time `json:"scheduled_date"`
}

// DateOnly is a wrapper type for date-only JSON serialization.
// TODO: Implement MarshalJSON and UnmarshalJSON for this type
// to format as "2006-01-02" (YYYY-MM-DD).
type DateOnly struct {
	time.Time
}

func testTimeHandling() {
	fmt.Println("--- Time Handling ---")

	m := Meeting{
		Title:         "Sprint Review",
		ScheduledDate: time.Date(2025, 1, 15, 14, 30, 0, 0, time.UTC),
	}

	data, _ := json.Marshal(m)
	jsonStr := string(data)
	fmt.Printf("Marshaled: %s\n", jsonStr)

	// We want just the date, not the full RFC3339 timestamp
	if strings.Contains(jsonStr, `"2025-01-15"`) && !strings.Contains(jsonStr, `T14:30`) {
		fmt.Println("PASS: date formatted as YYYY-MM-DD")
	} else {
		fmt.Println("FAIL: date should be '2025-01-15', not full RFC3339")
	}
}

// --- Numeric Precision ---

func testNumericPrecision() {
	fmt.Println("--- Numeric Precision ---")

	// BUG: By default, json.Unmarshal decodes numbers into float64,
	// which can lose precision for large integers.
	// The value 9_999_999_999_999_999 cannot be represented exactly as float64.
	//
	// TODO: Use json.NewDecoder with UseNumber() to preserve numeric precision,
	// then extract the value using json.Number.

	raw := `{"id": 9999999999999999, "amount": 123.456}`

	var result map[string]interface{}
	json.Unmarshal([]byte(raw), &result)

	idVal := fmt.Sprintf("%.0f", result["id"].(float64))
	fmt.Printf("Decoded ID: %s\n", idVal)

	if idVal == "9999999999999999" {
		fmt.Println("PASS: ID preserved with full precision")
	} else {
		fmt.Println("FAIL: ID lost precision (got " + idVal + ")")
	}
}

// --- OmitEmpty Behavior ---

type UserProfile struct {
	Name     string  `json:"name,omitempty"`
	Email    string  `json:"email,omitempty"`
	Age      int     `json:"age,omitempty"`
	Score    float64 `json:"score,omitempty"`
	Verified bool    `json:"verified,omitempty"`
	Bio      *string `json:"bio,omitempty"`
}

func testOmitEmpty() {
	fmt.Println("--- OmitEmpty Behavior ---")

	// BUG: omitempty drops zero-valued fields, but sometimes zero IS a valid value.
	// A user with age=0 or score=0.0 or verified=false might be intentional.
	// The Age, Score, and Verified fields should appear in JSON even when zero.
	//
	// TODO: Fix the struct tags so that Age, Score, and Verified always appear
	// in JSON output. Use pointer types for fields where you still want omitempty
	// to distinguish "not set" from "zero value".

	user := UserProfile{
		Name:     "Alice",
		Age:      0,     // Intentionally zero (e.g., unborn entity in a game)
		Score:    0.0,   // Intentionally zero
		Verified: false, // Intentionally false
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
