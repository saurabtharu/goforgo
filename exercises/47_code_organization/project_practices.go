// project_practices.go
// Fix common Go project anti-patterns: naming, getters, utility packages, docs
//
// This program demonstrates several Go project organization mistakes:
// 1. Java-style Get/Set methods for exported fields
// 2. "utils" grab-bag package (simulated with a utils struct)
// 3. Poor naming conventions (stuttering, misleading names)
// 4. Missing documentation on exported types
//
// Fix all the anti-patterns to follow Go conventions.

package main

import (
	"fmt"
	"strings"
	"unicode"
)

// BUG 1: Java-style getters and setters on exported fields.
// In Go, if a field is exported (capitalized), just use it directly.
// If you need encapsulation, make the field unexported and use a method
// named without "Get" prefix (e.g., Name() not GetName()).
//
// TODO: Fix this type. Either:
// - Make fields exported and remove Get/Set methods, OR
// - Make fields unexported and rename GetX() to X()
type UserUser struct { // BUG 3: Stuttering name "UserUser" - should just be "User"
	name  string
	email string
	age   int
}

func NewUserUser(name, email string, age int) *UserUser {
	return &UserUser{name: name, email: email, age: age}
}

// BUG: Java-style getter. In Go, this should be Name() not GetName().
func (u *UserUser) GetName() string {
	return u.name
}

// BUG: Unnecessary setter. If you need to set, either export the field
// or provide a method with a meaningful name.
func (u *UserUser) SetName(name string) {
	u.name = name
}

func (u *UserUser) GetEmail() string {
	return u.email
}

func (u *UserUser) SetEmail(email string) {
	u.email = email
}

func (u *UserUser) GetAge() int {
	return u.age
}

func (u *UserUser) SetAge(age int) {
	u.age = age
}

func (u *UserUser) String() string {
	return fmt.Sprintf("%s <%s> (age %d)", u.name, u.email, u.age)
}

// BUG 2: "Utils" grab-bag. In Go, avoid packages/types named "utils",
// "helpers", "common", etc. Instead, name things by what they DO.
//
// TODO: Break this into properly named functions:
// - StringUtils.Reverse -> just reverseString()
// - StringUtils.Capitalize -> just capitalizeWords()
// - StringUtils.CountVowels -> just countVowels()
// Each function should be a standalone function, not a method on a struct.
type StringUtils struct{}

func (StringUtils) Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (StringUtils) Capitalize(s string) string {
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

func (StringUtils) CountVowels(s string) int {
	count := 0
	for _, r := range strings.ToLower(s) {
		switch r {
		case 'a', 'e', 'i', 'o', 'u':
			count++
		}
	}
	return count
}

// BUG 3: More stuttering names and poor conventions.
// ConfigConfig should be Config.
// The ConfigConfig.ConfigValue field stutters: config.ConfigValue
// TODO: Fix the naming.
type ConfigConfig struct {
	ConfigValue string
	ConfigType  string
}

func NewConfigConfig(value, configType string) *ConfigConfig {
	return &ConfigConfig{
		ConfigValue: value,
		ConfigType:  configType,
	}
}

// BUG: GetConfigValue stutters and uses Get prefix.
func (c *ConfigConfig) GetConfigValue() string {
	return c.ConfigValue
}

func main() {
	fmt.Println("=== Project Practices ===")

	// Using the Java-style getters/setters
	fmt.Println("\n--- User (fix Get/Set anti-pattern) ---")
	user := NewUserUser("Alice", "alice@example.com", 30)
	fmt.Println("Name:", user.GetName())
	fmt.Println("Email:", user.GetEmail())
	fmt.Println("Age:", user.GetAge())
	user.SetName("Bob")
	user.SetEmail("bob@example.com")
	user.SetAge(25)
	fmt.Println("Updated:", user)

	// Using the utility grab-bag
	fmt.Println("\n--- StringUtils (fix utils anti-pattern) ---")
	utils := StringUtils{}
	fmt.Println("Reverse 'hello':", utils.Reverse("hello"))
	fmt.Println("Capitalize 'hello world':", utils.Capitalize("hello world"))
	fmt.Println("Vowels in 'beautiful':", utils.CountVowels("beautiful"))

	// Using stuttering names
	fmt.Println("\n--- Config (fix stuttering names) ---")
	cfg := NewConfigConfig("production", "environment")
	fmt.Println("Config value:", cfg.GetConfigValue())
	fmt.Println("Config type:", cfg.ConfigType)

	fmt.Println("\nProject practices fixed!")
}
