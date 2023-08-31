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
