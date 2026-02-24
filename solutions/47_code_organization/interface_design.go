// interface_design.go - SOLUTION
// Fixed: Small consumer-side interfaces, concrete return types, no interface pollution.

package main

import (
	"fmt"
	"strings"
)

// Fixed: No bloated DataStore interface. The concrete type stands on its own.

// InMemoryStore is a simple key-value store.
type InMemoryStore struct {
	data map[string]string
}

// Fixed: Returns concrete *InMemoryStore, not an interface.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{data: make(map[string]string)}
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *InMemoryStore) Set(key string, value string) {
	s.data[key] = value
}

func (s *InMemoryStore) Delete(key string) {
	delete(s.data, key)
}

func (s *InMemoryStore) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *InMemoryStore) Values() []string {
	vals := make([]string, 0, len(s.data))
	for _, v := range s.data {
		vals = append(vals, v)
	}
	return vals
}

func (s *InMemoryStore) Len() int {
	return len(s.data)
}

func (s *InMemoryStore) Clear() {
	s.data = make(map[string]string)
}

func (s *InMemoryStore) Contains(key string) bool {
	_, ok := s.data[key]
	return ok
}

// --- Consumer-side interfaces: defined where they're USED, not where they're implemented ---

// Setter is defined by the consumer that needs to write values.
type Setter interface {
	Set(key string, value string)
}

// Fixed: Accepts a small consumer-side interface.
func cacheUserPreferences(store Setter, userID string, prefs map[string]string) {
	for key, val := range prefs {
		fullKey := userID + ":" + key
		store.Set(fullKey, val)
	}
	fmt.Printf("Cached %d preferences for user %s\n", len(prefs), userID)
}

// Getter is defined by the consumer that needs to read a single value.
type Getter interface {
	Get(key string) (string, bool)
}

// Fixed: Accepts a single-method interface.
func lookupPreference(store Getter, userID, key string) string {
	fullKey := userID + ":" + key
	val, ok := store.Get(fullKey)
	if !ok {
		return "(not set)"
	}
	return val
}

// KeyValueReader is defined by the consumer that needs to list and read.
type KeyValueReader interface {
	Keys() []string
	Get(key string) (string, bool)
}

// Fixed: Accepts a small interface with just what it needs.
func printAllEntries(store KeyValueReader, prefix string) {
	keys := store.Keys()
	sortKeys(keys)
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			v, _ := store.Get(k)
			fmt.Printf("  %s = %s\n", k, v)
		}
	}
}

func sortKeys(keys []string) {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
}

func main() {
	fmt.Println("=== Interface Design ===")

	store := NewInMemoryStore()

	fmt.Println("\n--- Caching preferences ---")
	cacheUserPreferences(store, "alice", map[string]string{
		"theme":    "dark",
		"language": "en",
		"timezone": "UTC",
	})

	fmt.Println("\n--- Looking up preferences ---")
	fmt.Println("alice:theme =", lookupPreference(store, "alice", "theme"))
	fmt.Println("alice:language =", lookupPreference(store, "alice", "language"))
	fmt.Println("alice:missing =", lookupPreference(store, "alice", "missing"))

	fmt.Println("\n--- All alice entries ---")
	printAllEntries(store, "alice:")

	fmt.Println("\nInterface design fixed!")
}
