// interface_design.go
// Fix interface pollution, producer-side interfaces, and over-engineering
//
// Go interfaces should be:
// - Small (1-3 methods ideally)
// - Defined by the CONSUMER, not the producer
// - Return concrete types, not interfaces
//
// This code has several interface design mistakes from "100 Go Mistakes":
// 1. A bloated interface with too many methods (interface pollution)
// 2. Producer-side interfaces (defined next to the implementation)
// 3. A function returning an interface instead of a concrete type
//
// Fix all the interface design problems.

package main

import (
	"fmt"
	"strings"
)

// BUG 1: This interface is too large. No consumer needs all these methods.
// TODO: Remove this bloated interface. Let consumers define what they need.
type DataStore interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
	Keys() []string
	Values() []string
	Len() int
	Clear()
	Contains(key string) bool
	Merge(other DataStore)
}

// InMemoryStore is a simple key-value store.
// BUG 2: This implements the bloated DataStore interface unnecessarily.
// A concrete type should just have methods - no interface needed at definition site.
type InMemoryStore struct {
	data map[string]string
}

// BUG 3: This returns an interface instead of the concrete type.
// TODO: Return *InMemoryStore instead of DataStore.
func NewInMemoryStore() DataStore {
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

func (s *InMemoryStore) Merge(other DataStore) {
	for _, k := range other.Keys() {
		v, _ := other.Get(k)
		s.data[k] = v
	}
}

// --- Consumer code ---

// BUG 4: This function accepts the bloated DataStore interface but only
// uses Get and Set. It should define its own small interface.
// TODO: Define a small consumer-side interface with just what this function needs.
func cacheUserPreferences(store DataStore, userID string, prefs map[string]string) {
	for key, val := range prefs {
		fullKey := userID + ":" + key
		store.Set(fullKey, val)
	}
	fmt.Printf("Cached %d preferences for user %s\n", len(prefs), userID)
}

// BUG 5: This function accepts DataStore but only calls Get.
// TODO: Define a tiny interface (single method) for this consumer.
func lookupPreference(store DataStore, userID, key string) string {
	fullKey := userID + ":" + key
	val, ok := store.Get(fullKey)
	if !ok {
		return "(not set)"
	}
	return val
}

// BUG 6: This function accepts DataStore but only needs Keys and Get.
// TODO: Define a consumer-side interface for listing/reading.
func printAllEntries(store DataStore, prefix string) {
	keys := store.Keys()
	// Sort for deterministic output
	sortKeys(keys)
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			v, _ := store.Get(k)
			fmt.Printf("  %s = %s\n", k, v)
		}
	}
}

// Simple insertion sort for deterministic output (no sort import needed)
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
