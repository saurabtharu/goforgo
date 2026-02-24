package main

import (
	"fmt"
	"sync"
)

// BankAccount holds a balance that can be modified concurrently.
// All methods use pointer receivers because:
// 1. Deposit/Withdraw mutate state
// 2. The struct contains a sync.Mutex (must not be copied)
// 3. Consistency - if any method needs a pointer receiver, use it for all
type BankAccount struct {
	mu      sync.Mutex
	owner   string
	balance float64
}

// Deposit adds money to the account.
// Fixed: pointer receiver so mutation is visible to the caller,
// and the sync.Mutex is not copied.
func (a *BankAccount) Deposit(amount float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

// Withdraw removes money from the account.
// Fixed: pointer receiver for the same reasons as Deposit.
func (a *BankAccount) Withdraw(amount float64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.balance < amount {
		return false
	}
	a.balance -= amount
	return true
}

// Balance returns the current balance.
// Fixed: pointer receiver because the struct contains a sync.Mutex.
// Even though Balance doesn't mutate, copying a Mutex is a bug.
func (a *BankAccount) Balance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Point represents a 2D coordinate. It's small (16 bytes) and immutable.
type Point struct {
	X, Y float64
}

// Distance returns the Manhattan distance from the origin.
// Fixed: value receiver is appropriate because Point is small and
// Distance doesn't mutate it. Value receivers are simpler and safer.
func (p Point) Distance() float64 {
	return abs(p.X) + abs(p.Y)
}

// Equal checks if two points are the same.
// Fixed: value receiver for consistency and because Point is small/immutable.
func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	account := &BankAccount{owner: "Alice", balance: 100.0}
	account.Deposit(50.0)
	account.Withdraw(30.0)

	fmt.Printf("Owner: %s\n", account.owner)
	fmt.Printf("Balance: %.2f\n", account.Balance())

	p := Point{3, 4}
	fmt.Printf("Distance: %.2f\n", p.Distance())
	fmt.Printf("Equal: %v\n", p.Equal(Point{3, 4}))
}
