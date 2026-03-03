package main

import (
	"fmt"
	"sync"
)

// BankAccount holds a balance that can be modified concurrently.
// FIX: Some methods use the wrong receiver type. Apply these rules:
// - Use pointer receiver when: method mutates state, struct is large, struct has sync.Mutex
// - Use value receiver when: struct is small and immutable
type BankAccount struct {
	mu      sync.Mutex
	owner   string
	balance float64
}

// Deposit adds money to the account.
// BUG: This uses a value receiver, so the mutation is lost!
// The caller's BankAccount remains unchanged after calling Deposit.
func (a BankAccount) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

// Withdraw removes money from the account.
// BUG: Same problem - value receiver means the balance change is invisible.
func (a BankAccount) Withdraw(amount float64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return false
	}
	a.balance -= amount
	return true
}

// Balance returns the current balance.
// BUG: This also uses a value receiver, but since the struct has a sync.Mutex,
// copying it is dangerous (mutex must never be copied after first use).
func (a BankAccount) Balance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Point represents a 2D coordinate. It's small and immutable in usage.
type Point struct {
	X, Y float64
}

// Distance returns the Manhattan distance from the origin.
// FIX: This uses a pointer receiver unnecessarily. Point is small (16 bytes)
// and Distance doesn't mutate it. A value receiver is more appropriate here.
func (p *Point) Distance() float64 {
	return abs(p.X) + abs(p.Y)
}

// Equal checks if two points are the same.
// FIX: Same issue - pointer receiver is overkill for a small, read-only struct.
func (p *Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Test BankAccount - mutations should be visible
	account := &BankAccount{owner: "Alice", balance: 100.0}
	account.Deposit(50.0)
	account.Withdraw(30.0)

	fmt.Printf("Owner: %s\n", account.owner)
	fmt.Printf("Balance: %.2f\n", account.Balance())
	// Expected: Balance: 120.00

	// Test Point - should work with both value and pointer
	p := Point{3, 4}
	fmt.Printf("Distance: %.2f\n", p.Distance())
	fmt.Printf("Equal: %v\n", p.Equal(Point{3, 4}))
}
