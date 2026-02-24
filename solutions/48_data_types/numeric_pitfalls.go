package main

import (
	"fmt"
	"math"
)

// 100 Go Mistakes #17, #18, #19: Numeric Pitfalls (Solution)
//
// Fixed all three numeric issues:
// 1. Use octal literal 0o644 for file permissions
// 2. Use int (or int64) to avoid overflow
// 3. Use epsilon-based comparison for floats

func setFilePermissions() int {
	// Fixed: Use explicit octal prefix 0o644 (Go 1.13+).
	// This makes it clear we mean octal 644 = decimal 420.
	perm := 0o644
	fmt.Printf("File permission: %d (octal: %o)\n", perm, perm)
	return perm
}

func computeFactorial() {
	// Fixed: Use int instead of int8 to avoid overflow.
	// int is at least 32 bits, which handles 6! = 720 easily.
	var result int = 1
	n := 6
	for i := 1; i <= n; i++ {
		result *= i
	}
	fmt.Printf("Factorial of %d = %d\n", n, result)
}

func compareFloats() {
	// Force runtime computation to prevent Go's constant folding.
	values := []float64{0.1, 0.2, 0.3}
	a := values[0] + values[1]
	b := values[2]

	// Fixed: Use epsilon-based comparison for floating-point values.
	const epsilon = 1e-9
	if math.Abs(a-b) < epsilon {
		fmt.Println("0.1 + 0.2 == 0.3: true (correct)")
	} else {
		fmt.Println("0.1 + 0.2 == 0.3: false (wrong!)")
	}

	// Large float addition: acknowledge the precision limitation
	big := 1_000_000_000.0
	small := 0.000_000_001
	sum := big + small
	if math.Abs(sum-big) < epsilon {
		fmt.Println("Large float absorbed small value (precision lost!)")
	} else {
		fmt.Println("Large float preserved small value")
	}
}

func main() {
	fmt.Println("=== Numeric Pitfalls ===")
	fmt.Println()

	fmt.Println("--- Octal Literals ---")
	perm := setFilePermissions()
	if perm == 420 {
		fmt.Println("PASS: Correct octal permission value")
	} else {
		fmt.Println("FAIL: Permission value is wrong")
	}
	fmt.Println()

	fmt.Println("--- Integer Overflow ---")
	computeFactorial()
	fmt.Println()

	fmt.Println("--- Float Comparison ---")
	compareFloats()
}
