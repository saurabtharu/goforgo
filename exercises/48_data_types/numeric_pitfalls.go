package main

import (
	"fmt"
	"math"
)

// 100 Go Mistakes #17, #18, #19: Numeric Pitfalls
//
// This exercise covers three common numeric mistakes in Go:
// 1. Octal literal confusion (0-prefixed integers are octal!)
// 2. Integer overflow (silently wraps around)
// 3. Floating-point comparison (direct == fails for computed values)
//
// FIX all three functions so the program prints correct results.

// setFilePermissions demonstrates octal literal confusion.
// The developer wants Unix permission 644 (rw-r--r--) but
// writes the literal incorrectly.
func setFilePermissions() int {
	// FIX: This is meant to represent Unix permission 644 (decimal 420).
	// But 0644 in Go is an octal literal, which equals 420 in decimal.
	// The bug here is the print statement treats it as if it were decimal 644.
	perm := 644 // BUG: This is decimal 644, not octal 644!
	fmt.Printf("File permission: %d (octal: %o)\n", perm, perm)
	return perm
}

// computeFactorial shows how integer overflow silently wraps.
func computeFactorial() {
	// FIX: int8 can only hold values from -128 to 127.
	// Computing 6! = 720 will overflow silently.
	// Use a wider integer type.
	var result int8 = 1
	n := 6
	for i := 1; i <= n; i++ {
		result *= int8(i)
	}
	fmt.Printf("Factorial of %d = %d\n", n, result)
}

// compareFloats shows why direct float comparison is dangerous.
func compareFloats() {
	// Force runtime computation to prevent Go's constant folding.
	// Go's arbitrary-precision constants would hide the bug.
	values := []float64{0.1, 0.2, 0.3}
	a := values[0] + values[1] // 0.1 + 0.2 at runtime
	b := values[2]             // 0.3

	// FIX: Direct floating-point comparison fails due to IEEE 754
	// representation. Use an epsilon-based comparison instead.
	if a == b {
		fmt.Println("0.1 + 0.2 == 0.3: true (correct)")
	} else {
		fmt.Println("0.1 + 0.2 == 0.3: false (wrong!)")
	}

	// Also demonstrate: large float addition losing precision
	big := 1_000_000_000.0
	small := 0.000_000_001
	sum := big + small
	if sum == big {
		fmt.Println("Large float absorbed small value (precision lost!)")
	} else {
		fmt.Println("Large float preserved small value")
	}

	_ = math.Abs(0) // hint: math.Abs is useful for epsilon comparison
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
