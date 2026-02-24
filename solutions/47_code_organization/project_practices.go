// project_practices.go - SOLUTION
// Fixed: Go naming conventions, no stuttering, no utils, no Get prefix.

package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Fixed: Renamed UserUser to User. No stuttering.
// Using unexported fields with Go-style accessors (no Get prefix).
type User struct {
	name  string
	email string
	age   int
}

func NewUser(name, email string, age int) *User {
	return &User{name: name, email: email, age: age}
}

// Fixed: Name() instead of GetName(). Go convention.
func (u *User) Name() string {
	return u.name
}

// Fixed: Email() instead of GetEmail().
func (u *User) Email() string {
	return u.email
}

// Fixed: Age() instead of GetAge().
func (u *User) Age() int {
	return u.age
}

// Fixed: SetName is acceptable when you need controlled mutation.
func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetAge(age int) {
	u.age = age
}

func (u *User) String() string {
	return fmt.Sprintf("%s <%s> (age %d)", u.name, u.email, u.age)
}

// Fixed: Standalone functions instead of a StringUtils struct.
// Each function is named by what it does, not lumped in a "utils" bag.

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func capitalizeWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

func countVowels(s string) int {
	count := 0
	for _, r := range strings.ToLower(s) {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// Fixed: Renamed ConfigConfig to Config. Fields don't stutter.
type Config struct {
	Value string // Was ConfigValue - no need for prefix
	Type  string // Was ConfigType - no need for prefix
}

func NewConfig(value, configType string) *Config {
	return &Config{
		Value: value,
		Type:  configType,
	}
}

func main() {
	fmt.Println("=== Project Practices ===")

	// Fixed: Using Go-style accessors
	fmt.Println("\n--- User (fix Get/Set anti-pattern) ---")
	user := NewUser("Alice", "alice@example.com", 30)
	fmt.Println("Name:", user.Name())
	fmt.Println("Email:", user.Email())
	fmt.Println("Age:", user.Age())
	user.SetName("Bob")
	user.SetEmail("bob@example.com")
	user.SetAge(25)
	fmt.Println("Updated:", user)

	// Fixed: Standalone functions, no utils struct
	fmt.Println("\n--- StringUtils (fix utils anti-pattern) ---")
	fmt.Println("Reverse 'hello':", reverseString("hello"))
	fmt.Println("Capitalize 'hello world':", capitalizeWords("hello world"))
	fmt.Println("Vowels in 'beautiful':", countVowels("beautiful"))

	// Fixed: Clean naming, no stuttering
	fmt.Println("\n--- Config (fix stuttering names) ---")
	cfg := NewConfig("production", "environment")
	fmt.Println("Config value:", cfg.Value)
	fmt.Println("Config type:", cfg.Type)

	fmt.Println("\nProject practices fixed!")
}
