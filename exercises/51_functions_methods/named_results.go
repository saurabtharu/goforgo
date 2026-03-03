package main

import (
	"fmt"
	"math"
)

// FIX: This function returns two float64 values representing latitude and longitude,
// but the caller has no idea which is which from the signature alone.
// Mistake #43: Named result parameters can serve as documentation for ambiguous returns.
// Add named result parameters to make the return values self-documenting.
func geocode(address string) (float64, float64) {
	switch address {
	case "New York":
		return 40.7128, -74.0060
	case "London":
		return 51.5074, -0.1278
	default:
		return 0.0, 0.0
	}
}

// FIX: This function uses naked returns, which make the code confusing.
// Mistake #44: Naked returns with named results hurt readability, especially
// in functions with multiple return paths.
// Remove the naked returns by explicitly returning the named variables.
func calculateCircle(radius float64) (area, circumference float64, err error) {
	if radius < 0 {
		err = fmt.Errorf("radius cannot be negative: %.2f", radius)
		return // naked return - what exactly is being returned here?
	}

	if radius == 0 {
		return // naked return - returning zero values, but it's not obvious
	}

	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	return // naked return - which variables have been set at this point?
}

// FIX: This function converts temperature between Celsius and Fahrenheit.
// The two float64 returns are ambiguous without names.
// Add named results to clarify which is Celsius and which is Fahrenheit.
func convertTemp(celsius float64) (float64, float64) {
	fahrenheit := celsius*9.0/5.0 + 32.0
	return celsius, fahrenheit
}

func main() {
	// geocode returns lat/lng but the signature doesn't tell us which is which
	lat, lng := geocode("New York")
	fmt.Printf("New York: lat=%.4f, lng=%.4f\n", lat, lng)

	lat, lng = geocode("London")
	fmt.Printf("London: lat=%.4f, lng=%.4f\n", lat, lng)

	// calculateCircle has confusing naked returns
	area, circ, err := calculateCircle(5.0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Circle (r=5): area=%.2f, circumference=%.2f\n", area, circ)
	}

	area, circ, err = calculateCircle(-1.0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	area, circ, err = calculateCircle(0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Circle (r=0): area=%.2f, circumference=%.2f\n", area, circ)
	}

	// convertTemp returns two floats - which is which?
	c, f := convertTemp(100.0)
	fmt.Printf("Temperature: %.2f C = %.2f F\n", c, f)
}
