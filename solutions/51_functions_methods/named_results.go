package main

import (
	"fmt"
	"math"
)

// Fixed: Named result parameters document that the first float64 is latitude
// and the second is longitude. This is self-documenting at the call site.
func geocode(address string) (lat, lng float64) {
	switch address {
	case "New York":
		return 40.7128, -74.0060
	case "London":
		return 51.5074, -0.1278
	default:
		return 0.0, 0.0
	}
}

// Fixed: Kept named results for documentation but removed naked returns.
// Each return statement now explicitly shows what's being returned,
// making the control flow clear at every exit point.
func calculateCircle(radius float64) (area, circumference float64, err error) {
	if radius < 0 {
		return 0, 0, fmt.Errorf("radius cannot be negative: %.2f", radius)
	}

	if radius == 0 {
		return 0, 0, nil
	}

	area = math.Pi * radius * radius
	circumference = 2 * math.Pi * radius
	return area, circumference, nil
}

// Fixed: Named results clarify which value is Celsius and which is Fahrenheit.
func convertTemp(celsius float64) (c, f float64) {
	f = celsius*9.0/5.0 + 32.0
	return celsius, f
}

func main() {
	lat, lng := geocode("New York")
	fmt.Printf("New York: lat=%.4f, lng=%.4f\n", lat, lng)

	lat, lng = geocode("London")
	fmt.Printf("London: lat=%.4f, lng=%.4f\n", lat, lng)

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

	c, f := convertTemp(100.0)
	fmt.Printf("Temperature: %.2f C = %.2f F\n", c, f)
}
