package in_memory

import (
	"sync"
	"testing"
	"time"
)

func TestNewInMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	if storage == nil {
		t.Error("Expected to get an InMemoryStorage, got nil")
	}

	if len(storage.WordsStorage) != 0 {
		t.Errorf("Expected empty storage, got %v", storage.WordsStorage)
	}
}

func TestInsert(t *testing.T) {
	storage := NewInMemoryStorage()

	if err := storage.Insert("good"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if storage.WordsStorage["good"] != 1 {
		t.Errorf("Expected 'good' to have count of 1, got %d", storage.WordsStorage["good"])
	}

	if err := storage.Insert("bad"); err == nil {
		t.Error("Expected an error for forbidden word, got nil")
	}
}

func TestConcurrentInsert(t *testing.T) {
	storage := NewInMemoryStorage()

	var wg sync.WaitGroup

	concurrentInserts := 100
	expectedCount := concurrentInserts

	word := "concurrent"

	wg.Add(concurrentInserts)

	for i := 0; i < concurrentInserts; i++ {
		go func() {
			defer wg.Done()
			if err := storage.Insert(word); err != nil {
				t.Errorf("Unexpected error during concurrent insert: %v", err)
			}
		}()
	}

	wg.Wait()

	if storage.WordsStorage[word] != expectedCount {
		t.Errorf("Expected '%s' to have count of %d, got %d", word, expectedCount, storage.WordsStorage[word])
	}
}

func TestFindFrequentByPrefix(t *testing.T) {
	storage := NewInMemoryStorage()

	storage.WordsStorage["apple"] = 2
	storage.WordsStorage["appetite"] = 3
	storage.WordsStorage["banana"] = 1

	word, err := storage.FindFrequentByPrefix("app")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if word != "appetite" {
		t.Errorf("Expected word 'appetite', got %s", word)
	}

	_, err = storage.FindFrequentByPrefix("xyz")
	if err == nil {
		t.Error("Expected an error for prefix not found, got nil")
	}
}

func TestConcurrentFindFrequentByPrefix(t *testing.T) {
	storage := NewInMemoryStorage()
	var wg sync.WaitGroup

	storage.WordsStorage["apple"] = 2
	storage.WordsStorage["appetite"] = 3
	storage.WordsStorage["banana"] = 1

	concurrentReads := 100
	wg.Add(concurrentReads)

	for i := 0; i < concurrentReads; i++ {
		go func() {
			defer wg.Done()

			word, err := storage.FindFrequentByPrefix("app")
			if err != nil {
				t.Errorf("Unexpected error during concurrent read: %v", err)
			}

			if word != "appetite" {
				t.Errorf("During concurrent read, expected word 'appetite', got %s", word)
			}
		}()
	}

	wg.Wait()
}

func TestCleanGarbageCollector(t *testing.T) {
	// Short interval for testing
	gcInterval := 100 * time.Millisecond

	storage := NewInMemoryStorage()

	storage.WordsStorage["apple"] = 3
	storage.WordsStorage["banana"] = 3
	storage.WordsStorage["cherry"] = 10

	go storage.CleanGarbageCollector()

	// Wait for longer than gcInterval to ensure the garbage collector has run
	time.Sleep(2 * gcInterval)

	storage.mu.RLock()
	defer storage.mu.RUnlock()

	if _, exists := storage.WordsStorage["apple"]; exists {
		t.Error("Expected 'apple' to be removed by garbage collector")
	}

	if _, exists := storage.WordsStorage["banana"]; exists {
		t.Error("Expected 'banana' to be removed by garbage collector")
	}

	if _, exists := storage.WordsStorage["cherry"]; !exists {
		t.Error("Expected 'cherry' to remain in storage")
	}
}
