package storage

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	logWordAddINFO = "word %s succesfully inserted to storage"
	logWordGetINFO = "most frequent word with this prefix '%s' is %s"
	logCleanGC     = "cleaning the garbage collector"
	logStorage     = "current storage: %v"
)

// forbiddenWords serves as an example list of words that are disallowed
// for insertion. Given that our in-memory storage simply adds words to a map,
// there's a minimal chance of operational errors. However, to adhere to our
// interface which expects an error return, we use this list to simulate potential
// error scenarios during the Insert() operation.
var forbiddenWords = []string{"bad", "words", "example"}

type InMemoryStore struct {
	WordsStore map[string]int

	mu         sync.RWMutex
	GcInterval time.Duration
}

func NewInMemoryStore(GcInterval time.Duration) *InMemoryStore {
	store := &InMemoryStore{
		WordsStore: make(map[string]int, 0),
		GcInterval: GcInterval,
	}

	go store.cleanGarbageCollector()

	return store
}

func (s *InMemoryStore) Insert(word string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulated error for this challenge
	for _, w := range forbiddenWords {
		if word == w {
			return fmt.Errorf("forbidden word")
		}
	}

	s.WordsStore[word]++
	log.Info().Msgf(logWordAddINFO, word)
	return nil
}

func (s *InMemoryStore) FindFrequentByPrefix(prefix string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var maxWord string
	var maxCount int
	for word, count := range s.WordsStore {
		if strings.HasPrefix(word, prefix) && count > maxCount {
			maxWord = word
			maxCount = count
		}
	}
	if maxWord == "" {
		return "", fmt.Errorf("word notfound")
	}

	log.Info().Msgf(logWordGetINFO, prefix, maxWord)
	return maxWord, nil
}

// cleanGarbageCollector trims infrequent words (e.g., count of 3 for testing challenge, but 500 in real scenarios)
// from InMemoryStore to save memory and enhance performance.
// Especially useful when storage swells, limiting iteration overhead over vast entries like 10,000 words.
func (s *InMemoryStore) cleanGarbageCollector() {
	ticker := time.NewTicker(s.GcInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Info().Msg(logCleanGC)
			log.Info().Msgf(logStorage, s.WordsStore)
			s.mu.Lock()
			for word := range s.WordsStore {
				if s.WordsStore[word] == 3 {
					delete(s.WordsStore, word)
				}
			}
			s.mu.Unlock()
		}
	}
}
