package store

import (
	"testing"
	"time"
)

func TestNewInMemoryStore(t *testing.T) {
	store := NewInMemoryStore(1 * time.Second)

	if store == nil {
		t.Error("Expected to get an InMemoryStore, got nil")
	}

	if len(store.wordsStore) != 0 {
		t.Errorf("Expected empty store, got %v", store.wordsStore)
	}
}

func TestInsert(t *testing.T) {
	store := NewInMemoryStore(1 * time.Second)

	if err := store.Insert("good"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if store.wordsStore["good"] != 1 {
		t.Errorf("Expected 'good' to have count of 1, got %d", store.wordsStore["good"])
	}

	if err := store.Insert("bad"); err == nil {
		t.Error("Expected an error for forbidden word, got nil")
	}
}

func TestFindFrequentByPrefix(t *testing.T) {
	store := NewInMemoryStore(1 * time.Second)

	store.wordsStore["apple"] = 2
	store.wordsStore["appetite"] = 3
	store.wordsStore["banana"] = 1

	word, err := store.FindFrequentByPrefix("app")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if word != "appetite" {
		t.Errorf("Expected word 'appetite', got %s", word)
	}

	_, err = store.FindFrequentByPrefix("xyz")
	if err == nil {
		t.Error("Expected an error for prefix not found, got nil")
	}
}

func TestCleanGarbageCollector(t *testing.T) {
	// Short interval for testing
	gcInterval := 100 * time.Millisecond

	store := NewInMemoryStore(gcInterval)
	store.wordsStore["apple"] = 1
	store.wordsStore["banana"] = 1
	store.wordsStore["cherry"] = 2

	// Wait for longer than gcInterval to ensure the garbage collector has run
	time.Sleep(2 * gcInterval)

	store.mux.RLock()
	defer store.mux.RUnlock()

	// "apple" and "banana" should be deleted by garbage collector
	if _, exists := store.wordsStore["apple"]; exists {
		t.Error("Expected 'apple' to be removed by garbage collector")
	}

	if _, exists := store.wordsStore["banana"]; exists {
		t.Error("Expected 'banana' to be removed by garbage collector")
	}

	// "cherry" has count > 1 so should not be removed
	if _, exists := store.wordsStore["cherry"]; !exists {
		t.Error("Expected 'cherry' to remain in store")
	}
}
